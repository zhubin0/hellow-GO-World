package compress

import (
	"fmt"
	"testing"
)

func TestGZEncode(t *testing.T) {
	out, err := GZEncode([]byte("123"))
	fmt.Println(err, out)

	o, err := GZDecode(out)

	fmt.Println(string(o), err)
}
