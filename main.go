// This is a stanard encoder
package main

import (
	"errors"
	"fmt"
	"os"
)

// The buffer for reading the files
type Buffer struct {
	buffer []byte
	nByte  int // the next byte to read from buffer
	f      *os.File
	tBytes int // the total number of bytes that have been read so far
	fSize  int // the size of the file in bytes
}

// helper function to get a new Buffer object
func GetBuffer(fn string) *Buffer {
	f, err := os.Open(fn)
	if err != nil {
		fmt.Printf("Error! %s\n", err.Error())
		return nil
	}
	fst, _ := os.Stat(fn)
	bf := Buffer{
		buffer: make([]byte, 8192),
		f:      f,
		nByte:  0,
		fSize:  int(fst.Size()),
		tBytes: 0,
	}
	bf.ref()
	return &bf
}

// helper function to read new data into the buffer
func (bf *Buffer) ref() error {
	_, err := bf.f.Read(bf.buffer)
	bf.nByte = 0
	return err
}

func (bf *Buffer) read() (byte, error) {
	var err error = nil
	if bf.tBytes == bf.fSize {
		return 0, errors.New("End Of File\n")
	}
	if bf.nByte == len(bf.buffer) {
		err = bf.ref()
	}
	bt := bf.buffer[bf.nByte]
	bf.nByte += 1
	bf.tBytes += 1
	return bt, err
}

// helper function to read a byte integer in little endian
func (bf *Buffer) read4() int {
	b1, _ := bf.read()
	b2, _ := bf.read()
	b3, _ := bf.read()
	b4, _ := bf.read()
	bt := 0
	bt |= int(b4)
	bt <<= 8
	bt |= int(b3)
	bt <<= 8
	bt |= int(b2)
	bt <<= 8
	bt |= int(b1)
	return bt
}

// helper function to read a 2 byte integer in little endian
func (bf *Buffer) read2() int {
	b1, _ := bf.read()
	b2, _ := bf.read()
	bt := 0
	bt |= int(b2)
	bt <<= 8
	bt |= int(b1)
	return bt
}

// MCU containing the 3 channel data
type MCU struct {
	ch1 [64]int
	ch2 [64]int
	ch3 [64]int
}

type ImageData struct {
	filename    string
	width       int   // the image width
	height      int   // the image height
	size        int   // the size of the pixel array in bytes
	MCUs        []MCU // The array of MCUs containing the coeffecients
	blockWidth  int   // The number of MCUs on the x-axis
	blockHeight int   // The number of MCUs on the y-axis
	blockCount  int   // The total number of MCUs
}

func padInt(a int) string {
	aStr := fmt.Sprintf("%d", a)
	rem := 3 - len(aStr)
	for k := 0; k < rem; k++ {
		aStr = " " + aStr
	}
	return aStr
}

func writeMCU(mcu MCU) {
	fmt.Printf("Ch1\n")
	for a := 1; a < 65; a++ {
		fmt.Printf("%s ", padInt(mcu.ch1[a-1]))
		if a%8 == 0 {
			fmt.Printf("\n")
		}
	}
	fmt.Printf("\n")
	fmt.Printf("Ch2\n")
	for a := 1; a < 65; a++ {
		fmt.Printf("%s ", padInt(mcu.ch2[a-1]))
		if a%8 == 0 {
			fmt.Printf("\n")
		}
	}
	fmt.Printf("\n")
	fmt.Printf("CH3\n")
	for a := 1; a < 65; a++ {
		fmt.Printf("%s ", padInt(mcu.ch3[a-1]))
		if a%8 == 0 {
			fmt.Printf("\n")
		}
	}
	fmt.Printf("\n")
}

func getImageData(fn string) *ImageData {
	// the image data
	idt := ImageData{
		filename: fn,
	}
	fmt.Printf("** Decoding the bitmap file '%s' **\n", fn)

	bf := GetBuffer(fn)
	if bf == nil {
		return nil
	}
	b1, _ := bf.read()
	b2, _ := bf.read()
	// check for the 'BM' bytes
	if b1 != 66 && b2 != 77 {
		fmt.Printf("Error. The file '%s' is not a valid bitmap file\n", fn)
		return &idt
	}
	// the filelength
	idt.size = bf.read4()
	// skip the next 4 bytes
	bf.read4()
	// pixel offset. Always 26
	bf.read4()
	// DIB header size. Always 26
	bf.read4()
	idt.width = bf.read2()
	idt.height = bf.read2()
	// number of planes
	bf.read2()
	// number of bits per pixel
	bf.read2()

	// create the needed MCUs for the whole image
	mcuWidth := (idt.width + 7) / 8
	mcuHeight := (idt.height + 7) / 8
	mcuCount := mcuWidth * mcuHeight
	mcuArray := make([]MCU, mcuCount)

	idt.blockWidth = mcuWidth
	idt.blockHeight = mcuHeight
	idt.blockCount = mcuWidth * mcuHeight
	// Read the RGB Data
	for y := idt.height - 1; y >= 0; y-- {
		_mcuHeight := y / 8
		_pxHeight := y % 8
		for x := 0; x < idt.width; x++ {
			_mcuWidth := x / 8
			_pxWidth := x % 8

			_mcuIndex := (_mcuHeight * mcuWidth) + _mcuWidth
			_pixelIndex := (_pxHeight * 8) + _pxWidth
			mcu := &mcuArray[_mcuIndex]
			// Get the RGB data
			bb, _ := bf.read()
			gb, _ := bf.read()
			rb, _ := bf.read()
			// conver them into floats
			r := float32(rb)
			g := float32(gb)
			b := float32(bb)
			// convert the RGB to YCbCr
			y := 0.2990*r + 0.5870*g + 0.1140*b - 128
			cb := -0.1687*r - 0.3313*g + 0.5000*b
			cr := 0.5000*r - 0.4187*g - 0.0813*b
			// save the values as coeffecients
			mcu.ch1[_pixelIndex] = int(y)
			mcu.ch2[_pixelIndex] = int(cb)
			mcu.ch3[_pixelIndex] = int(cr)
		}
	}
	//writeMCU(mcuArray[0])
	bf.f.Close()
	idt.MCUs = mcuArray
	return &idt
}

func decodeBitmap(f string) {
	idt := getImageData(f)
	if idt == nil {
		return
	}
	forwardDCT(idt)
	quantize(idt, stb1, stb2)
	//writeMCU(idt.MCUs[0])
	generateSymbolTable(idt.MCUs)
}

func forwardDCT(idt *ImageData) {
	for y := 0; y < idt.blockHeight; y++ {
		for x := 0; x < idt.blockWidth; x++ {
			block := &idt.MCUs[(y*idt.blockWidth)+x]
			componentForwardDCT(&((*block).ch1))
			componentForwardDCT(&((*block).ch2))
			componentForwardDCT(&((*block).ch3))
		}
	}
}

func componentForwardDCT(comp *[64]int) {
	// 1-dimesnional FDCT on the rows
	for a := 0; a < 8; a++ {
		var a0 float64 = float64((*comp)[0*8+a])
		var a1 float64 = float64((*comp)[1*8+a])
		var a2 float64 = float64((*comp)[2*8+a])
		var a3 float64 = float64((*comp)[3*8+a])
		var a4 float64 = float64((*comp)[4*8+a])
		var a5 float64 = float64((*comp)[5*8+a])
		var a6 float64 = float64((*comp)[6*8+a])
		var a7 float64 = float64((*comp)[7*8+a])

		var b0 float64 = a0 + a7
		var b1 float64 = a1 + a6
		var b2 float64 = a2 + a5
		var b3 float64 = a3 + a4
		var b4 float64 = a3 - a4
		var b5 float64 = a2 - a5
		var b6 float64 = a1 - a6
		var b7 float64 = a0 - a7

		var c0 float64 = b0 + b3
		var c1 float64 = b1 + b2
		var c2 float64 = b1 - b2
		var c3 float64 = b0 - b3
		var c4 float64 = b4
		var c5 float64 = b5 - b4
		var c6 float64 = b6 - c5
		var c7 float64 = b7 - b6

		var d0 float64 = c0 + c1
		var d1 float64 = c0 - c1
		var d2 float64 = c2
		var d3 float64 = c3 - c2
		var d4 float64 = c4
		var d5 float64 = c5
		var d6 float64 = c6
		var d7 float64 = c5 + c7
		var d8 float64 = c4 - c6

		var e0 float64 = d0
		var e1 float64 = d1
		var e2 float64 = d2 * m1
		var e3 float64 = d3
		var e4 float64 = d4 * m2
		var e5 float64 = d5 * m3
		var e6 float64 = d6 * m4
		var e7 float64 = d7
		var e8 float64 = d8 * m5

		var f0 float64 = e0
		var f1 float64 = e1
		var f2 float64 = e2 + e3
		var f3 float64 = e3 - e2
		var f4 float64 = e4 + e8
		var f5 float64 = e5 + e7
		var f6 float64 = e6 + e8
		var f7 float64 = e7 - e5

		var g0 float64 = f0
		var g1 float64 = f1
		var g2 float64 = f2
		var g3 float64 = f3
		var g4 float64 = f4 + f7
		var g5 float64 = f5 + f6
		var g6 float64 = f5 - f6
		var g7 float64 = f7 - f4

		(*comp)[0*8+a] = int(g0 * s0)
		(*comp)[4*8+a] = int(g1 * s4)
		(*comp)[2*8+a] = int(g2 * s2)
		(*comp)[6*8+a] = int(g3 * s6)
		(*comp)[5*8+a] = int(g4 * s5)
		(*comp)[1*8+a] = int(g5 * s1)
		(*comp)[7*8+a] = int(g6 * s7)
		(*comp)[3*8+a] = int(g7 * s3)
	}
	// 1-dimensional FDCT on the columns
	for a := 0; a < 8; a++ {
		var a0 float64 = float64((*comp)[a*8+0])
		var a1 float64 = float64((*comp)[a*8+1])
		var a2 float64 = float64((*comp)[a*8+2])
		var a3 float64 = float64((*comp)[a*8+3])
		var a4 float64 = float64((*comp)[a*8+4])
		var a5 float64 = float64((*comp)[a*8+5])
		var a6 float64 = float64((*comp)[a*8+6])
		var a7 float64 = float64((*comp)[a*8+7])

		var b0 float64 = a0 + a7
		var b1 float64 = a1 + a6
		var b2 float64 = a2 + a5
		var b3 float64 = a3 + a4
		var b4 float64 = a3 - a4
		var b5 float64 = a2 - a5
		var b6 float64 = a1 - a6
		var b7 float64 = a0 - a7

		var c0 float64 = b0 + b3
		var c1 float64 = b1 + b2
		var c2 float64 = b1 - b2
		var c3 float64 = b0 - b3
		var c4 float64 = b4
		var c5 float64 = b5 - b4
		var c6 float64 = b6 - c5
		var c7 float64 = b7 - b6

		var d0 float64 = c0 + c1
		var d1 float64 = c0 - c1
		var d2 float64 = c2
		var d3 float64 = c3 - c2
		var d4 float64 = c4
		var d5 float64 = c5
		var d6 float64 = c6
		var d7 float64 = c5 + c7
		var d8 float64 = c4 - c6

		var e0 float64 = d0
		var e1 float64 = d1
		var e2 float64 = d2 * m1
		var e3 float64 = d3
		var e4 float64 = d4 * m2
		var e5 float64 = d5 * m3
		var e6 float64 = d6 * m4
		var e7 float64 = d7
		var e8 float64 = d8 * m5

		var f0 float64 = e0
		var f1 float64 = e1
		var f2 float64 = e2 + e3
		var f3 float64 = e3 - e2
		var f4 float64 = e4 + e8
		var f5 float64 = e5 + e7
		var f6 float64 = e6 + e8
		var f7 float64 = e7 - e5

		var g0 float64 = f0
		var g1 float64 = f1
		var g2 float64 = f2
		var g3 float64 = f3
		var g4 float64 = f4 + f7
		var g5 float64 = f5 + f6
		var g6 float64 = f5 - f6
		var g7 float64 = f7 - f4

		(*comp)[a*8+0] = int(g0 * s0)
		(*comp)[a*8+4] = int(g1 * s4)
		(*comp)[a*8+2] = int(g2 * s2)
		(*comp)[a*8+6] = int(g3 * s6)
		(*comp)[a*8+5] = int(g4 * s5)
		(*comp)[a*8+1] = int(g5 * s1)
		(*comp)[a*8+7] = int(g6 * s7)
		(*comp)[a*8+3] = int(g7 * s3)
	}
}

// performs quantization on the Y channel using qt1
// and on the Cb and Cr channel using qt2
func quantize(idt *ImageData, qt1 [64]int, qt2 [64]int) {
	for y := 0; y < idt.blockHeight; y++ {
		for x := 0; x < idt.blockWidth; x++ {
			mcu := &idt.MCUs[(y*idt.blockWidth)+x]
			// quantize the Y channel
			for a := 0; a < 64; a++ {
				(*mcu).ch1[a] /= qt1[a]
			}
			// quantize the Cb and Cr channels
			for a := 0; a < 64; a++ {
				(*mcu).ch2[a] /= qt2[a]
				(*mcu).ch3[a] /= qt2[a]
			}
		}
	}
}

func main() {
	if len(os.Args) < 2 {
		fmt.Printf("Error! You have not provided any filenames\n")
		os.Exit(1)
	}
	filenames := os.Args[1:]
	for f := range filenames {
		// decode bitmap file
		decodeBitmap(filenames[f])
	}
}
