package golang

import (
	"fmt"
	"testing"
)

func TestAa(t *testing.T) {
	s := make([]string, 3)
	fmt.Println("example:", s)

	s[0] = "a"
	s[1] = "b"
	s[2] = "c"
	fmt.Println("set:", s)
	fmt.Println("get s[2]:", s[2])

	fmt.Println("len s:", len(s))

	s = append(s, "d")
	s = append(s, "e", "f")

	fmt.Println("new s:", s)

	c := make([]string, len(s))
	copy(c, s)
	fmt.Println("cpy c:", c)

	l := s[2:5]
	fmt.Println("sl1:", l)

	l = s[:5]
	fmt.Println("sl2:", l)

	l = s[2:]
	fmt.Println("sl3:", l)

	m := []string{"c", "d", "e"}
	fmt.Println("sli m:", m)

	td := make([][]int, 3)
	for i := 0; i < 3; i++ {
		tt := i + 1
		td[i] = make([]int, tt)
		for j := 0; j < tt; j++ {
			td[i][j] = j + i
		}
		fmt.Println("td[", i, "]:", td[i])
	}
	fmt.Println("td:", td)
}
