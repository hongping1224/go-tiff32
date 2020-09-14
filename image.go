// Copyright 2019 Hong-Ping Lo. All rights reserved.
// Use of this source code is governed by a BDS-style
// license that can be found in the LICENSE file.

package tiff

import (
	"image"
	"image/color"
)

// Gray32 is an in-memory image whose At method returns color.Gray32 values.
type Gray32 struct {
	// Pix holds the image's pixels, as gray values in big-endian format. The pixel at
	// (x, y) starts at Pix[(y-Rect.Min.Y)*Stride + (x-Rect.Min.X)].
	Pix []uint32
	// Stride is the Pix stride (in bytes) between vertically adjacent pixels.
	Stride int
	// Rect is the image's bounds.
	Rect image.Rectangle
}

func (p *Gray32) ColorModel() color.Model { return Gray32Model }

func (p *Gray32) Bounds() image.Rectangle { return p.Rect }

func (p *Gray32) At(x, y int) color.Color {
	return p.Gray32At(x, y)
}

func (p *Gray32) Gray32At(x, y int) Gray32Color {
	if !(image.Point{x, y}.In(p.Rect)) {
		return Gray32Color{}
	}
	i := p.PixOffset(x, y)
	return Gray32Color{uint32(p.Pix[i])}
}

// PixOffset returns the index of the first element of Pix that corresponds to
// the pixel at (x, y).
func (p *Gray32) PixOffset(x, y int) int {
	return (y-p.Rect.Min.Y)*p.Stride + (x - p.Rect.Min.X)
}

func (p *Gray32) SetGray32(x, y int, c Gray32Color) {
	if !(image.Point{x, y}.In(p.Rect)) {
		return
	}
	i := p.PixOffset(x, y)
	p.Pix[i] = c.Y
}

// SubImage returns an image representing the portion of the image p visible
// through r. The returned value shares pixels with the original image.
func (p *Gray32) SubImage(r image.Rectangle) image.Image {
	r = r.Intersect(p.Rect)
	// If r1 and r2 are Rectangles, r1.Intersect(r2) is not guaranteed to be inside
	// either r1 or r2 if the intersection is empty. Without explicitly checking for
	// this, the Pix[i:] expression below can panic.
	if r.Empty() {
		return &Gray32{}
	}
	i := p.PixOffset(r.Min.X, r.Min.Y)
	return &Gray32{
		Pix:    p.Pix[i:],
		Stride: p.Stride,
		Rect:   r,
	}
}

// Opaque scans the entire image and reports whether it is fully opaque.
func (p *Gray32) Opaque() bool {
	return true
}

// NewGray32 returns a new Gray16 image with the given bounds.
func NewGray32(r image.Rectangle) *Gray32 {
	w, h := r.Dx(), r.Dy()
	pix := make([]uint32, w*h)
	return &Gray32{pix, w, r}
}

// GrayFloat32 is an in-memory image whose At method returns color.Gray32 values.
type GrayFloat32 struct {
	// Pix holds the image's pixels, as gray values in big-endian format. The pixel at
	// (x, y) starts at Pix[(y-Rect.Min.Y)*Stride + (x-Rect.Min.X)].
	Pix []uint32
	// Stride is the Pix stride (in bytes) between vertically adjacent pixels.
	Stride int
	// Rect is the image's bounds.
	Rect image.Rectangle
}

func (p *GrayFloat32) ColorModel() color.Model { return Gray32FloatModel }

func (p *GrayFloat32) Bounds() image.Rectangle { return p.Rect }

func (p *GrayFloat32) At(x, y int) color.Color {
	return p.Gray32At(x, y)
}

func (p *GrayFloat32) Gray32At(x, y int) Gray32Color {
	if !(image.Point{x, y}.In(p.Rect)) {
		return Gray32Color{}
	}
	i := p.PixOffset(x, y)
	return Gray32Color{uint32(p.Pix[i])}
}

// PixOffset returns the index of the first element of Pix that corresponds to
// the pixel at (x, y).
func (p *GrayFloat32) PixOffset(x, y int) int {
	return (y-p.Rect.Min.Y)*p.Stride + (x - p.Rect.Min.X)
}

func (p *GrayFloat32) SetGray32(x, y int, c GrayFloat32Color) {
	if !(image.Point{x, y}.In(p.Rect)) {
		return
	}
	i := p.PixOffset(x, y)
	p.Pix[i] = c.Y
}

// SubImage returns an image representing the portion of the image p visible
// through r. The returned value shares pixels with the original image.
func (p *GrayFloat32) SubImage(r image.Rectangle) image.Image {
	r = r.Intersect(p.Rect)
	// If r1 and r2 are Rectangles, r1.Intersect(r2) is not guaranteed to be inside
	// either r1 or r2 if the intersection is empty. Without explicitly checking for
	// this, the Pix[i:] expression below can panic.
	if r.Empty() {
		return &GrayFloat32{}
	}
	i := p.PixOffset(r.Min.X, r.Min.Y)
	return &GrayFloat32{
		Pix:    p.Pix[i:],
		Stride: p.Stride,
		Rect:   r,
	}
}

// Opaque scans the entire image and reports whether it is fully opaque.
func (p *GrayFloat32) Opaque() bool {
	return true
}

// NewGrayFloat32 returns a new Gray16 image with the given bounds.
func NewGrayFloat32(r image.Rectangle) *GrayFloat32 {
	w, h := r.Dx(), r.Dy()
	pix := make([]uint32, w*h)
	return &GrayFloat32{pix, w, r}
}
