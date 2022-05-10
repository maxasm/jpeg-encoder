// The Entropy Coded Segment
package main

import (
	"math"
)

// holds the symbol frequency and corresponding code
type SymbolMapping struct {
	code      uint16 // codes can have a max lenght of 16
	frequency uint64 // the frequency of the symbol
}

// helper function to  get the number of bits required to represent a given +value
func bitLength(val uint8) uint8 {
	v := uint8(0)
	n := uint8(0)
	for {
		// v is the largest n-bit value
		v = (1 << n) - 1
		if v >= val {
			return n
		}
		n += 1
	}
}

func generateSymbolTable(MCUs []MCU) {
	// the array of blocks
	ecsBlocks := []uint16{}
	// the ac symbol table
	var ACFrequencyTable = map[uint8]int64{}

	// the dc symbol table
	var DCFrequencyTable = map[uint8]int64{}

	for a := range MCUs {
		mcu := MCUs[a]

		for ch := 0; ch < 3; ch++ {
			var channel [64]int
			switch ch {
			case 0:
				channel = mcu.ch1
			case 1:
				channel = mcu.ch2
			case 2:
				channel = mcu.ch3
			}
			// iterate through all the coeffecients
			c := 0
			// number of zeroes
			nz := 0
			for {
				if c == 0 {
					// the upper nibble for dc values is always 0
					un := uint8(0x00)
					// the lower nibble contains the coeffecient length in bits
					cValue := channel[c]
					ln := bitLength(uint8(math.Abs(float64(cValue))))
					if cValue < 0 {
						cValue += ((1 << ln) - 1)
					}
					// create the DC symbol
					sym := uint8(0x00)
					sym |= un
					sym <<= 4
					sym |= ln

					// update the frequency table
					if _, ok := DCFrequencyTable[sym]; ok {
						DCFrequencyTable[sym] += 1
					} else {
						DCFrequencyTable[sym] = 1
					}
					// create a new ecs-block
					block := uint16(0x0000)
					block |= 0
					block <<= 8
					block |= uint16(sym)
					block <<= 7
					block |= uint16(cValue)
					// add the ecs-block
					ecsBlocks = append(ecsBlocks, block)
					c += 1
				} else {
					if c == 64 {
						if nz > 0 {
							// Todo: encode the block
							ecsBlocks = append(ecsBlocks, uint16(0x00))
						}
						break
					}

					if channel[c] == 0 {
						nz += 1
						c += 1
						continue
					}

					if channel[c] != 0 {
						// the number of 16-zeroe bands
						_16band := nz / 16
						// the remaining zeroes that aren't part of a 16-zeroe band
						_rem := uint8(nz % 16)

						// 'encode' the 16-zeroe bands
						for k := 0; k < _16band; k++ {
							// update the frequnecy table
							if _, ok := ACFrequencyTable[0xf0]; ok {
								ACFrequencyTable[0xf0] += 1
							} else {
								ACFrequencyTable[0xf0] = 1
							}
							// Todo: encode the block
							ecsBlocks = append(ecsBlocks, uint16(0xf0))
						}

						// 'encode' the remaining zeroes
						cValue := channel[c]
						ln := bitLength(uint8(math.Abs(float64(cValue))))
						// get the 'encode' value of cValue
						if cValue < 0 {
							cValue += ((1 << ln) - 1)
						}
						// create the symbol
						sym := uint8(0x00)
						sym |= _rem
						sym <<= 4
						sym |= ln
						// update the AC Frequency Table
						if _, ok := ACFrequencyTable[sym]; ok {
							ACFrequencyTable[sym] += 1
						} else {
							ACFrequencyTable[sym] = 1
						}
						// create the new block
						block := uint16(0x0000)
						block |= 1
						block <<= 8
						block |= uint16(sym)
						block <<= 7
						block |= uint16(cValue)
						// append the new block
						ecsBlocks = append(ecsBlocks, block)
						nz = 0
						c += 1
					}
				}
			}
		}
	}
}
