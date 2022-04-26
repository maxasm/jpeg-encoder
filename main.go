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

type Bitmap struct {
	filename string
	width    int // the image width
	height   int // the image height
	size     int // the size of the pixel array in bytes
	MCUs     []MCU
}

func padInt(a int) string {
	aStr := fmt.Sprintf("%d", a)
	rem := 3 - len(aStr)
	for k := 0; k < rem; k++ {
		aStr = "0" + aStr
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

func decodeBitmap(fn string) *Bitmap {
	bmp := Bitmap{
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
		return &bmp
	}
	// the filelength
	fl := bf.read4()
	bmp.size = fl
	// skip the next 4 bytes
	bf.read4()
	// pixel offset. Always 26
	bf.read4()
	// DIB header size. Always 26
	bf.read4()
	bmp.width = bf.read2()
	bmp.height = bf.read2()
	// number of planes
	bf.read2()
	// number of bits per pixel
	bf.read2()

	// create the needed MCUs for the whole image
	mcuWidth := (bmp.width + 7) / 8
	mcuHeight := (bmp.height + 7) / 8
	mcuCount := mcuWidth * mcuHeight
	mcuArray := make([]MCU, mcuCount)

	// Read the RGB Data
	for y := bmp.height - 1; y >= 0; y-- {
		_mcuHeight := y / 8
		_pxHeight := y % 8
		for x := 0; x < bmp.width; x++ {
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
			r := float64(rb)
			g := float64(gb)
			b := float64(bb)
			// conver the RGB data into YCbCr
			mcu.ch1[_pixelIndex] = int(16.0 + 65.481*r + 128.553*g + 24.966*b)
			mcu.ch2[_pixelIndex] = int(128.0 - 37.797*r - 74.203*g + 112.0*b)
			mcu.ch3[_pixelIndex] = int(128.0 + 112.0*r - 93.786*g - 18.214*b)
		}
	}

	writeMCU(mcuArray[0])
	bf.f.Close()
	bmp.MCUs = mcuArray
	return &bmp
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
