
# 参数处理
## 1. query参数
> 首先处理query参数，比如：http://xxx.com/user/add?id=1&age=20&username=张三
> 
> 记得将路由的URL匹配改为：uri := ctx.R.URL.Path

### 1.1 map类型参数
类似于 `http://loaclhost:8080/queryMap?user[id]=1&user[name]=张三`

## 2. Post表单参数
> 获取表单参数借助 http.Request.PostForm

Form属性包含了post表单和url后面跟的get参数。
PostForm属性只包含了post表单参数。

## 3. 文件参数
> 借助http.Request.FormFile

# 错误处理
>当程序发生异常的时候，比如panic程序直接就崩溃了，很明显在应用运行的过程中，不允许发生这样的事情，那么我们的框架就需要支持对这样问题的处理

对这样问题的处理，我们很容易想到recover，他会捕获panic错误，同时我们应当在defer进行捕获，defer是一个先进后出的栈结构，在return之前执行

<b>recover函数只在defer中生效</b>

## 1. Recovery中间件

## 2. 打印出错位置
>像这种处理异常的行为，单纯的报error没有意义，我们需要打印出出错位置

这里我们要用到runtime.Caller, Caller()报告当前go调用栈所执行的函数的文件和行号信息
- `skip`  
上溯的栈帧数，0表示Caller的调用者(Caller所在的调用栈) (0-当前函数,1-上一层函数,...)

- `pc`  
调用栈标识符

- `file`  
文件路径

- `line`  
该调用在文件中的行号

- `ok`  
如果无法获得信息，ok会被设为false

Callers 用来返回调用栈的程序计数器，第0个Caller是Callers本身，第1个是上一层trace，第2个是再上一层的`defer func`

# 协程池
>go的优势是高并发，高并发是由go+channel的组合来完成的。  

## 1.GMP模型  

![img](image/fa1499c59281b0e5bceba438676ed58a-16558274359922.png)
- 每个P维护一个G的本地队列
- 当一个G被创建出来，或者变为可执行状态时，优先把它放到P的本地队列中，否则放到全局队列
- 当一个G在M里执行结束后，P会从队列中把该G取出，如果此时P的队列为空，即没有其他G可以执行，如果全局队列为空，它会随机挑选另外一个P，从它的队列里拿走一半G到自己队列中执行

**P的数量在默认情况下，会被设定为CPU的核数。而M虽然需要跟P绑定执行，但数量上并不与P相等。这是因为M会因为系统调用或者其他事情被阻塞，因此随着程序的执行，M的数量可能增长，而P在没有用户干预的情况下，则会保持不变**

### 1.1大量创建go协程的代价
- 内存开销  
go协程大约占2k的内存  
`src/runtime/runtime2.go`  


- 调度开销  
虽然go协程的调度开销非常小，但也有一定的开销。
`runntime.Gosched()`当前协程主动让出CPU去执行另外一个协程 


- gc开销  
协程占用的内存最终需要gc来回收  


- 隐性的CPU开销  
最终协程是要内核线程来执行，我们知道在GMP模型中，G阻塞后，会新创建M来执行，一个M往往对应一个内核线程，当创建大量go协程的时候，内核线程的开销可能也会增大  
```
GO: runtime: program exceeds 10000-thread limit
```
>gmp模型中，本地队列的限制是256  
- 资源开销大的任务  
针对资源开销过大的任务，本身也不应当创建大量的协程，以免对CPU造成过大的任务，影响整体上的单机性能
- 任务堆积  
当创建过多协程，G阻塞增多，本地队列堆积过多，很可能造成内存溢出
- 系统任务影响  
runtime调度、gc等都是运行在go协程上的，当goroutine规模过大，会影响其他任务 

## 2. 协程池
基于以上一些理由，有必要创建一个协程池，将协程有效的管理起来，不要随意的创建过多的协程。
`同时池化的核心在于复用，所以我们可以这么想，一个goroutine是否可以处理多个任务，而不是一个goroutine处理一个任务`


### 2.1 需求
罗列一下需求：  
1. 希望创建固定数量的协程
2. 有一个任务队列，等待协程进行调度执行
3. 协程用完时，其他任务处于等待状态，一旦有协程空余，立即获取任务执行
4. 协程长时间空余，清空，以免占用空间
5. 有超时时间，如果一个任务长时间完成不了，就超时，让出协程

### 2.2 设计
![img](image/image-20220623180017764.png)

### 2.5 引入sync.pool
>将worker的创建放入pool中提取暴露(缓存)，用的时候从pool中获取，用完还回pool中，这样性能更高

### 2.6 引入sync.Cond
sync.Cond 是基于互斥锁/读写锁实现的条件变量，用来协调想要访问共享资源的那些Goroutine。 

当共享资源状态发生变化时，sync.Cond可以用来通知等待条件发生而阻塞的Goroutine。

在上述的场景中，我们可以将其应用在等待worker那里，可以使用sync.Cond阻塞，当worker执行完任务后，通知其继续执行。 

* `Signal方法`，允许调用Caller唤醒一个等待此Cond和goroutine。如果此时没有等待的goroutine，现实无需通知waiter；
如果Cond等待队列中有一个或者多个等待的goroutine，则需要从等待队列中移除第一个goroutine并把它唤醒。
在Java语言中，Signal方法也叫做notify方法。调用Signal方法时，不强求你一定要持有c.L的锁。
* `Broadcast方法`，允许调用者Caller唤醒所有等待此Cond的goroutine。如果此时没有等待的goroutine，显示无需通知waiter；
如果Cond等待队列中有一个或者多个等待的goroutine，则清空所有等待goroutine，并全部唤醒。
在Java语言中，Broadcast方法也叫做notifyAll方法。同样地，调用Broadcast方法时，也不强求你一定持有c.L的锁。
* `Wait方法`，会把调用者Caller放入Cond的等待队列中并阻塞，真到被Signal或者Broadcast的方法从等待队列中移除并唤醒。
调用Wait方法必须要持有c.L的锁

### 2.7 任务超时释放
>针对任务超时，需要使用工具的开发者，在程序中自动处理，及时退出goroutine

### 2.8 异常处理
>当task发生问题时，需要能捕获到，对外提供入口，让开发者自定义错误处理方式
