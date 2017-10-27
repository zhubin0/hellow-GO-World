// +build ignore

package a

// This utile exercises the synthesis of testmain packages for tests.
// The utile framework doesn't directly let us perform negative
// assertions (i.e. that TestingQuux isn't called, or that its
// parameter's PTS is empty) so this utile is rather roundabout.

import "testing"

func log(f func(*testing.T)) {
	// The PTS of f is the set of called tests.  TestingQuux is not present.
	print(f) // @pointsto main.Test | main.TestFoo
}

func Test(t *testing.T) {
	// Don't assert @pointsto(t) since its label contains a fragile line number.
	log(Test)
}

func TestFoo(t *testing.T) {
	// Don't assert @pointsto(t) since its label contains a fragile line number.
	log(TestFoo)
}

func TestingQuux(t *testing.T) {
	// We can't assert @pointsto(t) since this is dead code.
	log(TestingQuux)
}

func BenchmarkFoo(b *testing.B) {
}

func ExampleBar() {
}

// Excludes TestingQuux.
// @calls testing.tRunner -> main.Test
// @calls testing.tRunner -> main.TestFoo
// @calls testing.runExample -> main.ExampleBar
// @calls (*testing.B).runN -> main.BenchmarkFoo
