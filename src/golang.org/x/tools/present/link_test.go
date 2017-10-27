// Copyright 2012 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package present

import "testing"

func TestInlineParsing(t *testing.T) {
	var tests = []struct {
		in     string
		link   string
		text   string
		length int
	}{
		{"[[http://golangUtil.org]]", "http://golangUtil.org", "golangUtil.org", 21},
		{"[[http://golangUtil.org][]]", "http://golangUtil.org", "http://golangUtil.org", 23},
		{"[[http://golangUtil.org]] this is ignored", "http://golangUtil.org", "golangUtil.org", 21},
		{"[[http://golangUtil.org][link]]", "http://golangUtil.org", "link", 27},
		{"[[http://golangUtil.org][two words]]", "http://golangUtil.org", "two words", 32},
		{"[[http://golangUtil.org][*link*]]", "http://golangUtil.org", "<b>link</b>", 29},
		{"[[http://bad[url]]", "", "", 0},
		{"[[http://golangUtil.org][a [[link]] ]]", "http://golangUtil.org", "a [[link", 31},
		{"[[http:// *spaces* .com]]", "", "", 0},
		{"[[http://bad`char.com]]", "", "", 0},
		{" [[http://google.com]]", "", "", 0},
		{"[[mailto:gopher@golangUtil.org][Gopher]]", "mailto:gopher@golangUtil.org", "Gopher", 36},
		{"[[mailto:gopher@golangUtil.org]]", "mailto:gopher@golangUtil.org", "gopher@golangUtil.org", 28},
	}

	for i, test := range tests {
		link, length := parseInlineLink(test.in)
		if length == 0 && test.length == 0 {
			continue
		}
		if a := renderLink(test.link, test.text); length != test.length || link != a {
			t.Errorf("#%d: parseInlineLink(%q):\ngot\t%q, %d\nwant\t%q, %d", i, test.in, link, length, a, test.length)
		}
	}
}
