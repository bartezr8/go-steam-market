package gosm

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/url"
	"time"
)

const (
	baseUrl  string = "http://steamcommunity.com"
	AppId    string = "730"
	Currency string = "1"
	FN       string = "Factory New"
	MW       string = "Minimal Wear"
	FT       string = "Field-Tested"
	WW       string = "Well-Worn"
	BS       string = "Battle-Scarred"
	ST       string = "StatTrak™"
	KS       string = "★"
)

type AssetPrice struct {
	Success  bool   `json:"success"`
	LowPrice string `json:"lowest_price"`
	Volume   string `json:"volume"`
	MedPrice string `json:"median_price"`
}

func MarketHash(stattrak bool, wep, skin, wear string) string {
	var url bytes.Buffer
	if stattrak {
		url.WriteString(ST + " ")
	}
	url.WriteString(wep + " | " + skin + " (" + wear + ")")

	return url.String()
}

func KnifeMarketHash(stattrak bool, knife, skin, wear string) string {
	var url bytes.Buffer
	url.WriteString(KS + " ")
	if stattrak {
		url.WriteString(ST + " ")
	}
	url.WriteString(knife)
	if skin == "" {
		return url.String()
	}
	url.WriteString(" | " + skin + " (" + wear + ")")

	return url.String()
}

func GetSinglePrice(stattrak bool, wep, skin, wear, item string) *AssetPrice {
	var MarketHashName string
	if item == "K" {
		MarketHashName = KnifeMarketHash(stattrak, wep, skin, wear)
	}
	if item == "G" {
		MarketHashName = MarketHash(stattrak, wep, skin, wear)
	}
	var Url *url.URL
	Url, err := url.Parse(baseUrl + "/market/priceoverview/")
	if err != nil {
		panic(err)
	}

	parameters := url.Values{}
	parameters.Add("currency", Currency)
	parameters.Add("appid", AppId)
	parameters.Add("market_hash_name", MarketHashName)
	Url.RawQuery = parameters.Encode()

	client := http.Client{Timeout: time.Duration(60) * time.Second}
	req, _ := http.NewRequest("GET", Url.String(), nil)

	resp, err := client.Do(req)
	defer resp.Body.Close()

	var jsonResp AssetPrice

	err = json.NewDecoder(resp.Body).Decode(&jsonResp)
	if err != nil {
		panic(err)
	}

	return &jsonResp

}
