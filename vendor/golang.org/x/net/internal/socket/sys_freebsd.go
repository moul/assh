// Copyright 2017 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package socket

import (
	"syscall"
	"unsafe"
)

func probeProtocolStack() int {
	var p uintptr
	align := int(unsafe.Sizeof(p))
	// In the case of kern.supported_archs="amd64 i386", we need
	// to know the underlying kernel's architecture because the
	// alignment for socket facilities are set at the build time
	// of the kernel.
	conf, _ := syscall.Sysctl("kern.conftxt")
	for i, j := 0, 0; j < len(conf); j++ {
		if conf[j] != '\n' {
			continue
		}
		s := conf[i:j]
		i = j + 1
		if len(s) > len("machine") && s[:len("machine")] == "machine" {
			s = s[len("machine"):]
			for k := 0; k < len(s); k++ {
				if s[k] == ' ' || s[k] == '\t' {
					s = s[1:]
				}
				break
			}
			if s == "amd64" {
				align = 8
			}
			break
		}
	}
	return align
}
