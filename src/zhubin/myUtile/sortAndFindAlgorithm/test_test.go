package sortAndFindAlgorithm

import (
	"testing"
	"fmt"
)

func TestTest(t *testing.T) {
	var a = []int{10, 2, 5, 2, 6, 9, 10, 11, 23, 34, 45, 65, 76, 76, 87}
	fmt.Println("front：",a)
	Bubbing(a)
	fmt.Println("back：",a)
}

func Bubbing(a []int) {
	var j,k,flag int
	flag = len(a)
	for flag > 0 {
		k = flag
		flag = 0
		for j=1; j<k; j++{
			if a[j] < a[j-1] {
				a[j-1],a[j] = a[j],a[j-1]
				k = j
			}
		}
	}
}

func SertSort (a []int) {

}

func QuickSort(a []int) {

}