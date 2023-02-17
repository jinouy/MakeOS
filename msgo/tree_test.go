package msgo

import (
	"fmt"
	"testing"
)

func TestTreeNode(t *testing.T) {
	root := &treeNode{name: "/", children: make([]*treeNode, 0)}
	root.Put("/user/:id/get")
	root.Put("/user/**/hello")
	root.Put("/user/create/aaa")
	root.Put("/order/get/aaa")

	node := root.Get("/user/1/get")
	fmt.Println(node)
	node = root.Get("/user/create/hello")
	fmt.Println(node)
	node = root.Get("/user/create/aaa")
	fmt.Println(node)
	node = root.Get("/order/get/aaa")
	fmt.Println(node)

}
