package gold

import (
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"
)

const GOLD_URL = "https://sjc.com.vn/xml/tygiavang.xml"

// Thanks ChatGPT for generating this structure for me
type GoldRoot struct {
	XMLName  xml.Name `xml:"root"`
	Title    string   `xml:"title"`
	URL      string   `xml:"url"`
	RateList RateList `xml:"ratelist"`
}

type RateList struct {
	Updated string `xml:"updated,attr"`
	Unit    string `xml:"unit,attr"`
	Cities  []City `xml:"city"`
}

type City struct {
	Name  string `xml:"name,attr"`
	Items []Item `xml:"item"`
}

type Item struct {
	Buy  string `xml:"buy,attr"`
	Sell string `xml:"sell,attr"`
	Type string `xml:"type,attr"`
}

type goldService struct{}

func (g *goldService) FetchGoldPrice() (int, error) {
	fmt.Println("Fetching gold price!")
	res, err := http.Get(GOLD_URL)
	if err != nil {
		return 0, err
	}
	fmt.Println("Fetching gold price successfully!")
	resBody, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return 0, err
	}

	var goldData *GoldRoot
	err = xml.Unmarshal(resBody, &goldData)
	if err != nil {
		return 0, err
	}

	goldPriceStr := goldData.RateList.Cities[0].Items[0].Buy
	goldPriceStr = strings.Replace(goldPriceStr, ".", "", 1)

	result, err := strconv.Atoi(goldPriceStr)
	if err != nil {
		return 0, err
	}

	return result, nil
}

func NewService() *goldService {
	return &goldService{}
}
