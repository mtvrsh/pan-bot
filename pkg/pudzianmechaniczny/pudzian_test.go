package pudzianmechaniczny

import (
	"fmt"
	"math"
	"testing"
)

const (
	maxIter  = 20
	roundNum = 1000 // 3 decimal places
)

type pm func(float64) float64

func iterateN(f pm, expect float64, n int) string {
	for i := 0.0; i < float64(n); i++ {
		got := math.Round(f(i)*roundNum) / roundNum
		want := math.Round(expect*i*roundNum) / roundNum
		// fmt.Printf("PMToW(%v) got: %v, want: %v\n", i, got, want)
		if got != want {
			return fmt.Sprintf("PMToW(%v) got: %v, want: %v\n", i, got, want)
		}
	}
	return ""
}

func TestPMToW(t *testing.T) {
	err := iterateN(PMToW, 20054.092500000002, maxIter)
	if err != "" {
		t.Errorf(err)
	}
}

func TestPMToHP(t *testing.T) {
	err := iterateN(PMToHP, 26.89297639801529, maxIter)
	if err != "" {
		t.Errorf(err)
	}
}

func TestWToPM(t *testing.T) {
	err := iterateN(WToPM, 4.9865133513271664e-05, maxIter)
	if err != "" {
		t.Errorf(err)
	}
}

func TestHPToPM(t *testing.T) {
	err := iterateN(HPToPM, 0.03718443006084668, maxIter)
	if err != "" {
		t.Errorf(err)
	}
}

func TestReverse(t *testing.T) {
	n := PMToW(WToPM(1))
	if n != 1 {
		t.Errorf("PMToW(WToPM(1)) is %v, should be 1", n)
	}

	n = PMToHP(HPToPM(1))
	if n != 0.9999999999999999 { // why?
		t.Errorf("PMToHP(HPToPM(1)) is %v, should be 0.9999999999999999", n)
	}

	n = WToPM(PMToW(1))
	if n != 1 {
		t.Errorf("WToPM(PMToW(1)) is %v, should be 1", n)
	}

	n = HPToPM(PMToHP(1))
	if n != 1 {
		t.Errorf("HPToPM(PMToHP(1)) is %v, should be 1", n)
	}
}
