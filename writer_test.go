// Copyright 2019 Hong-Ping Lo. All rights reserved.
// Use of this source code is governed by a BDS-style
// license that can be found in the LICENSE file.
package tiff

import (
	"encoding/binary"
	"fmt"
	"math"
	"testing"
	"unsafe"
)

func TestFloat32Change(t *testing.T) {
	a := float32(100.0)
	fmt.Println(a)
	b := (*[4]byte)(unsafe.Pointer(&a))[:]
	fmt.Println(b)
	bits := binary.LittleEndian.Uint32(b)
	u := math.Float32frombits(bits)
	fmt.Println(u)
}
