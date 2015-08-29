package gosm

import (
	"bytes"
	"encoding/json"
	"log"
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

type AssetIntermediate struct {
	Result map[string]json.RawMessage `json:"result"`
}

type AssetInfo struct {
	Assets  map[string]Asset `json:"assets"`
	Success bool             `json:"success"`
	Error   string           `json:"error"`
}

type Asset struct {
	IconUrl           string                 `json:"icon_url,omitempty"`
	IconUrlLarge      string                 `json:"icon_url_large,omitempty"`
	IconDragUrl       string                 `json:"icon_drag_url,omitempty"`
	Name              string                 `json:"name,omitempty"`
	MarketHashName    string                 `json:"market_hash_name,omitempty"`
	MarketName        string                 `json:"market_name,omitempty"`
	NameColor         string                 `json:"name_color,omitempty"`
	BGColor           string                 `json:"background_color,omitempty"`
	Type              string                 `json:"type,omitempty"`
	Tradable          string                 `json:"tradable,omitempty"`
	Marketable        string                 `json:"marketable,omitempty"`
	Commodity         string                 `json:"commodity,omitempty"`
	TradeRestrict     string                 `json:"market_tradeable_restriction,omitempty"`
	FraudWarnings     string                 `json:"fraudwarnings,omitempty"`
	Descriptions      map[string]Description `json:"descriptions,omitempty"`
	OwnerDescriptions string                 `json:"owner_descriptions,omitempty"`
	Tags              map[string]Tag         `json:"tags,omitempty"`
	ClassId           string                 `json:"classid,omitempty"`
}

type Description struct {
	Type    string `json:"type"`
	Value   string `json:"value"`
	Color   string `json:"color,omitempty"`
	AppData string `json:"appdata"`
}

type Tag struct {
	InternalName string `json:"internal_name"`
	Name         string `json:"name"`
	Category     string `json:"category"`
	Color        string `json:"color,omitempty"`
	CategoryName string `json:"category_name"`
}

func (a *Asset) GetPrice() *AssetPrice {
	resp, err := GetSinglePricePrepped(a.MarketHashName)
	if err != nil {
		panic(err)
	}

	return resp

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

func GetSinglePricePrepped(name string) (*AssetPrice, error) {
	var Url *url.URL

	Url, err := url.Parse(baseUrl + "/market/priceoverview")
	if err != nil {
		panic(err)
	}

	parameters := url.Values{}
	parameters.Add("currency", Currency)
	parameters.Add("appid", AppId)
	parameters.Add("market_hash_name", name)
	Url.RawQuery = parameters.Encode()

	client := http.Client{Timeout: time.Duration(60) * time.Second}
	req, _ := http.NewRequest("GET", Url.String(), nil)

	resp, err := client.Do(req)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	defer resp.Body.Close()

	var jsonResp AssetPrice

	err = json.NewDecoder(resp.Body).Decode(&jsonResp)
	if err != nil {
		return nil, err
	}

	return &jsonResp, nil

}

func GetAssetInfo(key, appid, class_count, classid, instanceid string) (*AssetInfo, error) {
	var Url *url.URL
	Url, err := url.Parse("http://api.steampowered.com/ISteamEconomy/GetAssetClassInfo/v0001/")

	parameters := url.Values{}
	parameters.Add("class_count", class_count)
	parameters.Add("classid0", classid)
	parameters.Add("instanceid0", instanceid)
	parameters.Add("appid", appid)
	parameters.Add("key", key)
	Url.RawQuery = parameters.Encode()

	client := http.Client{Timeout: time.Duration(60) * time.Second}
	req, _ := http.NewRequest("GET", Url.String(), nil)

	resp, err := client.Do(req)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	defer resp.Body.Close()

	var intermediate AssetIntermediate

	json.NewDecoder(resp.Body).Decode(&intermediate)

	var ai AssetInfo

	for k, v := range intermediate.Result {
		r := bytes.NewReader(v)
		d := json.NewDecoder(r)
		var err error
		if k == "error" {
			err = d.Decode(&ai.Error)
		} else if k == "success" {
			err = d.Decode(&ai.Success)
		} else {
			var a Asset
			err = d.Decode(&a)
			if err == nil {
				if ai.Assets == nil {
					ai.Assets = map[string]Asset{}
				}
				ai.Assets[k] = a
			}
		}
		if err != nil {
			panic(err)
		}
	}

	return &ai, nil

}
