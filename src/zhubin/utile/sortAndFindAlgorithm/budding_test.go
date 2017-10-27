package sortAndFindAlgorithm

import (
	"fmt"
	"testing"
)

func TestBubbling(t *testing.T) {
	arr := []int{10, 2, 5, 2, 6, 9, 10, 11, 23, 34, 45, 65, 76, 76, 87}
	fmt.Println("qian: ", arr)
	Bubbling(arr)
	fmt.Println("hou: ", arr)
}

func Bubbling(a []int) {
	arr := a
	len := len(arr)
	sign := len
	count := 0
	for sign > 0 {
		j := sign
		sign = 0
		for i := 1; i < j; i++ {
			count = count + 1
			if arr[i-1] > arr[i] {
				arr[i-1], arr[i] = arr[i], arr[i-1]
				sign = i
			}
		}

	}
	fmt.Println("total: ", count)

}
