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
