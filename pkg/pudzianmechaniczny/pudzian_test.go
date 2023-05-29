package pudzianmechaniczny

import (
	"fmt"
	"math"
	"testing"
)

const (
	maxIter  = 20
	roundNum = 10000 // 4 decimal places precision
)

type pm func(float64) float64

func iterationHelper(t *testing.T, f pm, expect float64) {
	for i := 0.0; i < float64(maxIter); i++ {
		t.Run(fmt.Sprintf("%v", maxIter), func(t *testing.T) {
			got := math.Round(f(i)*roundNum) / roundNum
			want := math.Round(expect*i*roundNum) / roundNum
			if got != want {
				t.Errorf("PMToW(%v) got: %v, want: %v\n", i, got, want)
			}
		})
	}
}

func TestPMToW(t *testing.T) {
	t.Parallel()
	iterationHelper(t, PMToW, 20054.092500000002)
}

func TestPMToHP(t *testing.T) {
	t.Parallel()
	iterationHelper(t, PMToHP, 26.89297639801529)
}

func TestWToPM(t *testing.T) {
	t.Parallel()
	iterationHelper(t, WToPM, 4.9865133513271664e-05)
}

func TestHPToPM(t *testing.T) {
	t.Parallel()
	iterationHelper(t, HPToPM, 0.03718443006084668)
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
