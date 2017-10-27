package golangUtil

import (
	"fmt"
	"testing"
)

func TestAA(t *testing.T) {
	c := make(chan int, 1)
	c <- 10
	v := <-c
	fmt.Println("v = ", v)
}

func TestBB(t *testing.T) {
	c := make(chan int, 1)
	c <- 10
	v, ok := <-c
	fmt.Println("v = ", v)
	fmt.Println("ok : ", ok)
	close(c)
	v, ok = <-c
	fmt.Println("v = ", v)
	fmt.Println("ok : ", ok)
}

func TestCC(t *testing.T) {
	c := make(chan int, 1)
	select {
	case c <- 10:
	default:
	}
	v := <-c
	fmt.Println("c = ", c)
	fmt.Println("v = ", v)
}

func TestDD(t *testing.T) {
	c := make(chan int, 1)
	select {
	case c <- 10:
	default:
	}

	select {
	case c <- 11:
	default:
	}
	v := <-c
	fmt.Println("c = ", c)
	fmt.Println("v = ", v)
}

func TestEE(t *testing.T) {
	c := make(chan int, 1)
	select {
	case c <- 10:
	default:
	}

	select {
	case c <- 11:
	default:
	}

	select {
	case v, ok := <-c: // 读出来一个，v=10, ok=true
		fmt.Println("v = ", v)
		fmt.Println("ok : ", ok)
	default:
	}

	select {
	case v, ok := <-c:
		fmt.Println("v = ", v)
		fmt.Println("ok : ", ok)
	default: // 没有可读的，走这个分支
	}
}
