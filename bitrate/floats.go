package bitrate

import (
	"math"
)

func Fill(value float64, length int) []float64 {
	a := make([]float64, length)
	for i := range a {
		a[i] = value
	}
	return a
}

func Add(s1 []float64, s2 []float64) []float64 {
	a := make([]float64, len(s1))
	for i := 0; i < len(a); i++ {
		a[i] = s1[i] + s2[i]
	}
	return a
}

func Abs(s []float64) []float64 {
	a := make([]float64, len(s))
	for i := 0; i < len(a); i++ {
		a[i] = math.Abs(s[i])
	}
	return a
}

func Mul(s1 []float64, s2 []float64) []float64 {
	a := make([]float64, len(s1))
	for i := 0; i < len(a); i++ {
		a[i] = s1[i] * s2[i]
	}
	return a
}

func Divide(s []float64, divider float64) []float64 {
	a := make([]float64, len(s))
	for i := 0; i < len(a); i++ {
		a[i] = s[i] / divider
	}
	return a
}

func Normalize(s []float64) []float64 {
	a := make([]float64, len(s))
	for i := 0; i < len(a); i++ {
		a[i] = math.Log10(s[i])
	}
	return a
}

func Reverse(s []float64) []float64 {
	for i, j := 0, len(s)-1; i < j; i, j = i+1, j-1 {
		s[i], s[j] = s[j], s[i]
	}
	return s
}

func sumProduct(s1 []float64, s2 []float64) float64 {
	var a float64
	for i := range s1 {
		a += s1[i] * s2[i]
	}
	return a
}

func CrossCorrelation(s1 []float64, s2 []float64) []float64 {
	if len(s2) > len(s1) {
		s1, s2 = s2, s1
	}

	a := make([]float64, (len(s1) - len(s2) + 1))

	for i := 0; i < len(a); i++ {
		a[i] = sumProduct(s1[i:len(s2)+i], s2)
	}

	return a
}

func ClosestValue(value float64, s []float64) float64 {
	c := s[0]

	for _, v := range s {
		m := v - value
		m = math.Abs(m)
		if m < c {
			c = m
		}
	}

	return c
}
