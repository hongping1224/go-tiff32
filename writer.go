// Copyright 2019 Hong-Ping Lo. All rights reserved.
// Use of this source code is governed by a BDS-style
// license that can be found in the LICENSE file.

package tiff

import (
	"encoding/binary"
	"image"
	"io"
	"sort"

	"golang.org/x/image/tiff"
)

// The length of one instance of each data type in bytes.
var lengths = [...]uint32{0, 1, 1, 2, 4, 8}

const (
	dtByte     = 1
	dtASCII    = 2
	dtShort    = 3
	dtLong     = 4
	dtRational = 5
)

// Tags (see p. 28-41 of the spec).
const (
	tImageWidth                = 256
	tImageLength               = 257
	tBitsPerSample             = 258
	tCompression               = 259
	tPhotometricInterpretation = 262

	tStripOffsets    = 273
	tSamplesPerPixel = 277
	tRowsPerStrip    = 278
	tStripByteCounts = 279

	tTileWidth      = 322
	tTileLength     = 323
	tTileOffsets    = 324
	tTileByteCounts = 325

	tXResolution    = 282
	tYResolution    = 283
	tResolutionUnit = 296

	tPredictor    = 317
	tColorMap     = 320
	tExtraSamples = 338
	tSampleFormat = 339
)

const (
	cNone  = 1
	ifdLen = 12 // Length of an IFD entry in bytes.

	prNone       = 1
	pRGB         = 2
	prHorizontal = 2
	pPaletted    = 3
)
const (
	sampleFormat_UINT   = 1
	sampleFormat_INT    = 2
	sampleFormat_IEEEFP = 3
	sampleFormat_VOID   = 4
)

type ifdEntry struct {
	tag      int
	datatype int
	data     []uint32
}

// Encode writes the image m to w. opt determines the options used for
// encoding, such as the compression type. If opt is nil, an uncompressed
// image is written.
func Encode(w io.Writer, m image.Image, opt *tiff.Options) error {
	d := m.Bounds().Size()

	compression := uint32(cNone)
	predictor := false
	_, err := io.WriteString(w, "II\x2A\x00")
	if err != nil {
		return err
	}

	// Compressed data is written into a buffer first, so that we
	// know the compressed size.
	//var buf bytes.Buffer
	// dst holds the destination for the pixel data of the image --
	// either w or a writer to buf.
	var dst io.Writer
	// imageLen is the length of the pixel data in bytes.
	// The offset of the IFD is imageLen + 8 header bytes.
	var imageLen int

	switch compression {
	case cNone:
		dst = w
		// Write IFD offset before outputting pixel data.
		switch m.(type) {
		case *Gray32:
			imageLen = d.X * d.Y * 4
		case *GrayFloat32:
			imageLen = d.X * d.Y * 4
		default:
			imageLen = d.X * d.Y * 4
		}
		err = binary.Write(w, binary.LittleEndian, uint32(imageLen+8))
		if err != nil {
			return err
		}
	}

	pr := uint32(prNone)
	photometricInterpretation := uint32(pRGB)
	samplesPerPixel := uint32(4)
	bitsPerSample := []uint32{8, 8, 8, 8}
	extraSamples := uint32(0)
	colorMap := []uint32{}
	SampleFormat := sampleFormat_UINT
	if predictor {
		pr = prHorizontal
	}
	switch m := m.(type) {
	case *Gray32:
		photometricInterpretation = 1
		samplesPerPixel = 1
		bitsPerSample = []uint32{32}
		err = encodeGray32(dst, m.Pix, d.X, d.Y, m.Stride, predictor)
	case *GrayFloat32:
		photometricInterpretation = 1
		samplesPerPixel = 1
		bitsPerSample = []uint32{32}
		SampleFormat = sampleFormat_IEEEFP
		err = encodeGrayFloat32(dst, m.Pix, d.X, d.Y, m.Stride, predictor)
	default:
		extraSamples = 1 // Associated alpha.
		//	err = encode(dst, m, predictor)
	}
	if err != nil {
		return err
	}

	ifd := []ifdEntry{
		{tImageWidth, dtShort, []uint32{uint32(d.X)}},
		{tImageLength, dtShort, []uint32{uint32(d.Y)}},
		{tBitsPerSample, dtShort, bitsPerSample},
		{tCompression, dtShort, []uint32{compression}},
		{tPhotometricInterpretation, dtShort, []uint32{photometricInterpretation}},
		{tStripOffsets, dtLong, []uint32{8}},
		{tSamplesPerPixel, dtShort, []uint32{samplesPerPixel}},
		{tRowsPerStrip, dtShort, []uint32{uint32(d.Y)}},
		{tStripByteCounts, dtLong, []uint32{uint32(imageLen)}},
		{tSampleFormat, dtShort, []uint32{uint32(SampleFormat)}},
		// There is currently no support for storing the image
		// resolution, so give a bogus value of 72x72 dpi.
		{tXResolution, dtRational, []uint32{72, 1}},
		{tYResolution, dtRational, []uint32{72, 1}},
		{tResolutionUnit, dtShort, []uint32{2}},
	}
	if pr != prNone {
		ifd = append(ifd, ifdEntry{tPredictor, dtShort, []uint32{pr}})
	}
	if len(colorMap) != 0 {
		ifd = append(ifd, ifdEntry{tColorMap, dtShort, colorMap})
	}
	if extraSamples > 0 {
		ifd = append(ifd, ifdEntry{tExtraSamples, dtShort, []uint32{extraSamples}})
	}

	return writeIFD(w, imageLen+8, ifd)
}

type byTag []ifdEntry

func (d byTag) Len() int           { return len(d) }
func (d byTag) Less(i, j int) bool { return d[i].tag < d[j].tag }
func (d byTag) Swap(i, j int)      { d[i], d[j] = d[j], d[i] }

var enc = binary.LittleEndian

func (e ifdEntry) putData(p []byte) {
	for _, d := range e.data {
		switch e.datatype {
		case dtByte, dtASCII:
			p[0] = byte(d)
			p = p[1:]
		case dtShort:
			enc.PutUint16(p, uint16(d))
			p = p[2:]
		case dtLong, dtRational:
			enc.PutUint32(p, uint32(d))
			p = p[4:]
		}
	}
}

func writeIFD(w io.Writer, ifdOffset int, d []ifdEntry) error {
	var buf [ifdLen]byte
	// Make space for "pointer area" containing IFD entry data
	// longer than 4 bytes.
	parea := make([]byte, 1024)
	pstart := ifdOffset + ifdLen*len(d) + 6
	var o int // Current offset in parea.

	// The IFD has to be written with the tags in ascending order.
	sort.Sort(byTag(d))

	// Write the number of entries in this IFD.
	if err := binary.Write(w, enc, uint16(len(d))); err != nil {
		return err
	}
	for _, ent := range d {
		enc.PutUint16(buf[0:2], uint16(ent.tag))
		enc.PutUint16(buf[2:4], uint16(ent.datatype))
		count := uint32(len(ent.data))
		if ent.datatype == dtRational {
			count /= 2
		}
		enc.PutUint32(buf[4:8], count)
		datalen := int(count * lengths[ent.datatype])
		if datalen <= 4 {
			ent.putData(buf[8:12])
		} else {
			if (o + datalen) > len(parea) {
				newlen := len(parea) + 1024
				for (o + datalen) > newlen {
					newlen += 1024
				}
				newarea := make([]byte, newlen)
				copy(newarea, parea)
				parea = newarea
			}
			ent.putData(parea[o : o+datalen])
			enc.PutUint32(buf[8:12], uint32(pstart+o))
			o += datalen
		}
		if _, err := w.Write(buf[:]); err != nil {
			return err
		}
	}
	// The IFD ends with the offset of the next IFD in the file,
	// or zero if it is the last one (page 14).
	if err := binary.Write(w, enc, uint32(0)); err != nil {
		return err
	}
	_, err := w.Write(parea[:o])
	return err
}

func encodeGray32(w io.Writer, pix []uint32, dx, dy, stride int, predictor bool) error {
	buf := make([]byte, dx*4)
	for y := 0; y < dy; y++ {
		min := y*stride + 0
		max := y*stride + dx
		off := 0
		var v0 uint32
		for i := min; i < max; i++ {
			// An image.Gray16's Pix is in big-endian order.
			v1 := pix[i]
			if predictor {
				v0, v1 = v1, v1-v0
			}
			// We only write little-endian TIFF files.
			buf[off+0] = byte(v1)
			buf[off+1] = byte(v1 >> 8)
			buf[off+2] = byte(v1 >> 16)
			buf[off+3] = byte(v1 >> 24)
			off += 4
		}
		if _, err := w.Write(buf); err != nil {
			return err
		}
	}
	return nil
}

func encodeGrayFloat32(w io.Writer, pix []uint32, dx, dy, stride int, predictor bool) error {
	buf := make([]byte, dx*4)
	for y := 0; y < dy; y++ {
		min := y*stride + 0
		max := y*stride + dx
		off := 0
		var v0 uint32
		for i := min; i < max; i++ {
			// An image.Gray16's Pix is in big-endian order.
			v1 := pix[i]
			if predictor {
				v0, v1 = v1, v1-v0
			}
			// We only write little-endian TIFF files.
			buf[off+0] = byte(v1)
			buf[off+1] = byte(v1 >> 8)
			buf[off+2] = byte(v1 >> 16)
			buf[off+3] = byte(v1 >> 24)
			off += 4
		}
		if _, err := w.Write(buf); err != nil {
			return err
		}
	}
	return nil
}
