package main

import "math"

var m0 float64 = 2.0 * math.Cos(1.0/(16.0*2.0*math.Pi))
var m1 float64 = 2.0 * math.Cos(2.0/(16.0*2.0*math.Pi))
var m3 float64 = 2.0 * math.Cos(2.0/(16.0*2.0*math.Pi))
var m5 float64 = 2.0 * math.Cos(3.0/(16.0*2.0*math.Pi))
var m2 = m0 - m5
var m4 = m0 + m5

var s0 float64 = math.Cos(0.0/(16.0*math.Pi)) / math.Sqrt(8.0)
var s1 float64 = math.Cos(1.0/(16.0*math.Pi)) / 2.0
var s2 float64 = math.Cos(2.0/(16.0*math.Pi)) / 2.0
var s3 float64 = math.Cos(3.0/(16.0*math.Pi)) / 2.0
var s4 float64 = math.Cos(4.0/(16.0*math.Pi)) / 2.0
var s5 float64 = math.Cos(5.0/(16.0*math.Pi)) / 2.0
var s6 float64 = math.Cos(6.0/(16.0*math.Pi)) / 2.0
var s7 float64 = math.Cos(7.0/(16.0*math.Pi)) / 2.0

// standard qunatization tables
var stb1 = [64]int{
	16, 11, 10, 16, 24, 40, 51, 61,
	12, 12, 14, 19, 26, 58, 60, 55,
	14, 13, 16, 24, 40, 57, 69, 56,
	14, 17, 22, 29, 51, 87, 80, 62,
	18, 22, 37, 56, 68, 109, 103, 77,
	24, 35, 55, 64, 81, 104, 113, 92,
	49, 64, 78, 87, 103, 121, 120, 101,
	72, 92, 95, 98, 112, 100, 103, 99,
}

var stb2 = [64]int{
	17, 18, 24, 47, 99, 99, 99, 99,
	18, 21, 26, 66, 99, 99, 99, 99,
	24, 26, 56, 99, 99, 99, 99, 99,
	47, 66, 99, 99, 99, 99, 99, 99,
	99, 99, 99, 99, 99, 99, 99, 99,
	99, 99, 99, 99, 99, 99, 99, 99,
	99, 99, 99, 99, 99, 99, 99, 99,
	99, 99, 99, 99, 99, 99, 99, 99,
}
