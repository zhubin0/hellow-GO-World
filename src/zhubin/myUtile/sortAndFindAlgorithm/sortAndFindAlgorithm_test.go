package sortAndFindAlgorithm

import (
	"fmt"
	"testing"
)

func TestBubbling1(t *testing.T) {
	arr := []int{10, 2, 5, 2, 6, 9, 10, 11, 23, 34, 45, 65, 76, 76, 87}
	fmt.Println("qian: ", arr)
	Bubbling(arr)
	fmt.Println("hou: ", arr)
}

func Bubbling1(a []int) {
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

func TestInertionSort1(t *testing.T) {
	arr := []int{10, 2, 5, 2, 6, 9, 10, 11, 23, 34, 45, 65, 76, 76, 87}
	fmt.Println("qian: ", arr)
	//InsertSort(arr)
	//InsertSortUpdate(arr)
	InsertSortUpdateUpdate(arr)
	fmt.Println("hou: ", arr)
}

func InsertSort1(a []int) {
	len := len(a)
	count := 0
	var i, j, k int
	for i = 1; i < len; i++ {
		for j = i - 1; j >= 0; j-- {
			count++
			if a[i] > a[j] {
				break
			}
		}
		//fmt.Println("j = ", j)
		if j != i-1 {
			temp := a[i]
			for k = i - 1; k > j; k-- {
				a[k+1] = a[k]
			}
			a[k+1] = temp
		}
	}
	fmt.Println("total: ", count)
}

// 进行一下改写，将搜索和数据后移这二个步 骤合并。
// 即每次a[i]先和前面一个数据a[i-1]比较，如果a[i] > a[i-1]说明a[0...i]也 是有序的，无须调整。
// 否则就令j=i-1,temp=a[i]。然后一边将数据a[j]向后移动一 边向前搜索，当有数据a[j]<a[i]时停止并将temp放到a[j + 1]处。
func InsertSortUpdate1(a []int) {
	len := len(a)
	var i, j int
	for i = 1; i < len; i++ {
		if a[i] < a[i-1] {
			temp := a[i]
			for j = i - 1; j >= 0 && temp < a[j]; j-- {
				a[j+1] = a[j]
			}
			a[j+1] = temp
		}
	}
}

// 再对将a[j]插入到前面a[0...j-1]的有序区间所用的方法进行改写，用数据交换代 替数据后移。
// 如果a[j]前一个数据a[j-1] > a[j]，就交换a[j]和a[j-1]，再j--直到a[j-1] <= a[j]。这样也可以实现将一个新数据新并入到有序区间。
func InsertSortUpdateUpdate1(a []int) {
	len := len(a)
	var i, j int
	for i = 1; i < len; i++ {
		if a[i] < a[i-1] {
			for j = i - 1; j >= 0 && a[i] < a[j]; j-- {
				a[j+1], a[j] = a[j], a[j+1]
			}
		}
	}
}
