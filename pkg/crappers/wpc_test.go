package crappers //(s)crappers

import (
	"testing"
)

func TestGetWpcPrice(t *testing.T) {
	const expectedPrice = "49.99"
	t.Parallel()
	price, err := GetWpcPrice()
	if err != nil {
		t.Error(err.Error())
	}
	if price != expectedPrice {
		t.Errorf("got:\n%v\nwant:\n%v", price, expectedPrice)
	}
}
