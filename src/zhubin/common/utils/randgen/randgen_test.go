package randgen

import (
	"testing"

	"fmt"
	"github.com/deckarep/golang-set"
	"github.com/stretchr/testify/assert"
)

func TestGenRandString(t *testing.T) {
	round := 1000000
	set := mapset.NewSet()
	for i := 0; i < round; i++ {
		s := GenRandString(8)
		//fmt.Println(s)
		set.Add(s)
	}
	assert.Equal(t, round, set.Cardinality(), "random string generation dectects collision!")
}

func TestGenRandNumString(t *testing.T) {
	round := 100000
	set := mapset.NewSet()
	for i := 0; i < round; i++ {
		s := GenRandNumString(20)
		fmt.Println(s)
		set.Add(s)
	}
	assert.Equal(t, round, set.Cardinality(), "random string generation dectects collision!")
}

func TestGenMongoId(t *testing.T) {
	round := 10000
	set := mapset.NewSet()
	for i := 0; i < round; i++ {
		s := GenMongoId()
		fmt.Println(s)
		set.Add(s)
	}
	assert.Equal(t, round, set.Cardinality(), "random string generation dectects collision!")
}

func TestGenUniqueString(t *testing.T) {
	round := 10000
	set := mapset.NewSet()
	for i := 0; i < round; i++ {
		s := GenUniqueString(0)
		fmt.Println(s)
		set.Add(s)
	}
	assert.Equal(t, round, set.Cardinality(), "random string generation dectects collision!")
}
