package crappers //(s)crapper

import (
	"io"
	"net/http"
	"regexp"
)

const wpcURL = "https://sklep.kfd.pl/kfd-pure-wpc-82-instant-700-g-bialko-serwatkowe-naturalne-p-6867.html"

// GetWpcPrice returns price of item found on wpcURL
func GetWpcPrice() (string, error) {
	resp, err := http.Get(wpcURL)
	if err != nil {
		return "", err
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	priceRegex1 := regexp.MustCompile(`"price": ".*",`)
	price := priceRegex1.Find(body)
	priceRegex2 := regexp.MustCompile(`[0-9]+\.?[0-9]*`)
	price = priceRegex2.Find(price)

	return string(price), nil
}
