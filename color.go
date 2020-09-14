// Copyright 2019 Hong-Ping Lo. All rights reserved.
// Use of this source code is governed by a BDS-style
// license that can be found in the LICENSE file.

package tiff

import "image/color"

// Gray32Color represents a 32-bit grayscale color.
type Gray32Color struct {
	Y uint32
}

func (c Gray32Color) RGBA() (r, g, b, a uint32) {
	return c.Y, c.Y, c.Y, c.Y
}

var Gray32Model color.Model = color.ModelFunc(gray32Model)

func gray32Model(c color.Color) color.Color {
	if _, ok := c.(Gray32Color); ok {
		return c
	}
	r, g, b, _ := c.RGBA()

	// These coefficients (the fractions 0.299, 0.587 and 0.114) are the same
	// as those given by the JFIF specification and used by func RGBToYCbCr in
	// ycbcr.go.
	//
	// Note that 19595 + 38470 + 7471 equals 65536.
	//Need to fix let total = 4294967295,1284195221, 2521145802,489626272
	y := (1284195221*uint32(r) + 2521145802*uint32(g) + 489626272*uint32(b) + 1<<31) >> 32

	return Gray32Color{Y: uint32(y)}
}

// GrayFloat32Color represents a 32-bit float grayscale color.
type GrayFloat32Color struct {
	Y uint32
}

func (c GrayFloat32Color) RGBA() (r, g, b, a uint32) {
	return c.Y, c.Y, c.Y, c.Y
}

var Gray32FloatModel color.Model = color.ModelFunc(gray32FloatModel)

func gray32FloatModel(c color.Color) color.Color {
	if _, ok := c.(Gray32Color); ok {
		return c
	}
	r, g, b, _ := c.RGBA()

	// These coefficients (the fractions 0.299, 0.587 and 0.114) are the same
	// as those given by the JFIF specification and used by func RGBToYCbCr in
	// ycbcr.go.
	//
	// Note that 19595 + 38470 + 7471 equals 65536.
	//Need to fix let total = 4294967295,1284195221, 2521145802,489626272
	y := (1284195221*uint32(r) + 2521145802*uint32(g) + 489626272*uint32(b) + 1<<31) >> 32

	return Gray32Color{Y: uint32(y)}
}
