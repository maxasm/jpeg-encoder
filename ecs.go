// The Entropy Coded Segment
package main

import (
	"fmt"
	"math"
)

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

// helper function to print the frequency tables
func printFrequencyTable(ft map[uint8]int64) {
	for k, v := range ft {
		fmt.Printf("0x%x -> %d\n", k, v)
	}
	fmt.Printf("\n")
}

func generateSymbolTable(MCUs []MCU) (map[uint8]int64, map[uint8]int64) {
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
				// Todo: DC Values should be encoded as being relative to one another
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
							// Reaching at this branch means that all the remaining coeffecients are
							// all zero. So the symbol to use is 0x00
							sym := uint8(0x00)
							if _, ok := ACFrequencyTable[sym]; ok {
								ACFrequencyTable[sym] += 1
							} else {
								ACFrequencyTable[sym] = 1
							}
							block := uint16(0x0000)
							block |= 1
							block <<= 15
							ecsBlocks = append(ecsBlocks, block)
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
	return DCFrequencyTable, ACFrequencyTable
}
