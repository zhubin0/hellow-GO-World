// Copyright 2012 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build !appengine

package main

// This package registers "/compile" and "/share" handlers
// that redirect to the golangUtil.org playground.
import _ "golangUtil.org/x/tools/playground"
