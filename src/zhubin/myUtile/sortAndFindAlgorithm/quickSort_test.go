package sortAndFindAlgorithm

import (
	"testing"
	"fmt"
)
//快速排序
//1.先从数列中取出一个数作为基准数。
//2.分区过程，将比这个数大的数全放到它的右边，小于或等于它的数全放到它的左边。
//3.再对左右区间重复第二步，直到各区间只有一个数。

func TestQuickSort(t *testing.T) {
	var a = []int{72, 6, 57, 88, 60, 42, 83, 73, 48, 85}
	fmt.Println("front:", a)
	l := 0
	r := len(a) - 1
	//quickSort(a, l, r)
	quickSortUpdate(a, l, r)
	fmt.Println("hou:", a)
	fmt.Println("ok.")
}

func quickSort(a []int, l, r int) {
	if l<r {
		i := getMiddleSort(a, l, r)
		quickSort(a, l, i-1)
		quickSort(a, i+1, r)
	}

}

//72 6 57 88 60 42 83 73 48 85
func getMiddleSort(a []int, l, r int) int{
	x := a[l]
	i,j := l,r
	for i<j {
		for i<j && a[j] >= x {
			j--
		}
		if i<j {
			a[i] = a[j]
			i++
		}

		for i<j && a[i] <= x {
			i++
		}
		if i<j {
			a[j] = a[i]
			j--
		}
	}
	a[i] = x
	return i
}

func quickSortUpdate(a []int, l, r int) {
	if l < r {
		x := a[l]
		i,j := l,r
		for i<j && a[j] >= x {
			j--
		}
		if i<j {
			a[i] = a[j]
			i++
		}

		for i<j && a[i] <= x {
			i++
		}
		if i<j {
			a[j] = a[i]
			j--
		}
		a[i] = x
		quickSort(a, l, i - 1)
		quickSort(a, i + 1, r)
	}
}