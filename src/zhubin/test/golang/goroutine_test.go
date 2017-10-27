package golang

import (
	"fmt"
	"testing"
)

func TestAAA(t *testing.T) {
	go sayHello("zhubin")
}

func sayHello(name string) {
	fmt.Println("hello " + name)
}

func TestBBB(t *testing.T) {
	go func() {
		fmt.Println("hello ")
	}()
}
