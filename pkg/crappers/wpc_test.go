package crappers //(s)crapper

import (
	"testing"
)

func TestGetWpcPrice(t *testing.T) {
	price, err := GetWpcPrice()
	if err != nil {
		t.Errorf(err.Error())
	}
	if price != "49.99" {
		t.Errorf("got:\n%v\nwant:\n%v", price, "49.99")
	}
}
