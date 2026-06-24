package main

import (
	"fmt"
	"io"
	"math"
	"net/http"
	"strconv"
	"strings"
	"time"
	"unicode"

	"github.com/mozillazg/go-pinyin"
	"golang.org/x/text/encoding/simplifiedchinese"
	"golang.org/x/text/transform"
)

var stockNameMapping = map[string]string{
	"sh601138": "Foxconn Industrial Internet",
	"sz300985": "Zhiyuan New Energy",
	"sz300059": "East Money Information",
	"sz300601": "Kangtai Biologics",
	"sz000158": "Changshan Beiming",
	"sz000066": "China Great Wall",
	"sz300592": "HuaKai Ebuy",
	"sh600199": "Jin Seed Wine",
	"sh600398": "Hailan Home",
	"sh601006": "Daqin Railway",
	"sz300527": "ST Yingji",
	"sz002424": "ST Bailing",
	"sz301058": "COFCO Technology and Engineering",
	"sh563230": "Fuguo Satellite ETF",
	"sh588780": "Guolian Sci-Tech Chip Design ETF",
	"sh688521": "VeriSilicon",
	"sh688110": "Dosilicon",
	"sh688047": "Loongson",
	"sz000858": "wuliangye",
	"sh601668": "zhong guo jianzhu",
	"sh601166": "xingye Bank",
	"sh601818": "Everbright Bank",
	"sz300474": "Jingjia Micro",
	"sh688256": "Cambricon",
	"sh688041": "Hygon Information",
	"sh688981": "SMIC",
	"sz300498": "Wens Foodstuff",
	"sz300097": "Zhiyun",
	"sh603029": "Swan",
	"sz000564": "Gongxiao Daji",
	"sz003042": "Sino-Agri United",
	"sh600693": "Dongbai Group",
	"sz002251": "Bubugao",
	"sz000759": "Zhongbai Group",
	"sh601933": "Yonghui Superstores",
	"sh601012": "LONGi Green Energy",
	"sh600340": "China Fortune Land",
	"sz300476": "Shenghong Technology",
	"sh000001": "Shanghai Composite Index",
	"sz399001": "Shenzhen Component Index",
	"sz399006": "ChiNext Index",
	"hk01347": "Hua Hong Semiconductor (HK)",
}

var stockCodes []string

func init() {
	stockCodes = make([]string, 0, len(stockNameMapping))
	for code := range stockNameMapping {
		stockCodes = append(stockCodes, code)
	}
}

type Stock struct {
	Name    string  `json:"name"`
	Pinyin  string  `json:"pinyin"`
	Code    string  `json:"code"`
	Price   float64 `json:"price"`
	Open    float64 `json:"open"`
	Change  float64 `json:"change"`
	Ratio   float64 `json:"ratio"`
	High    string  `json:"high"`
	Low     string  `json:"low"`
	Volume  int64   `json:"volume"`
	Amount  string  `json:"amount"`
	Buy1    string  `json:"buy1"`
	Sell1   string  `json:"sell1"`
	Lowest  string  `json:"lowest"`
	Highest string  `json:"highest"`
}

func fetchStockData(codes []string) (string, error) {
	url := "http://qt.gtimg.cn/q=" + strings.Join(codes, ",")

	client := &http.Client{Timeout: 5 * time.Second}
	resp, err := client.Get(url)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("unexpected status: %s", resp.Status)
	}

	reader := transform.NewReader(resp.Body, simplifiedchinese.GBK.NewDecoder())
	body, err := io.ReadAll(reader)
	if err != nil {
		return "", err
	}
	return string(body), nil
}

func parseFloat(s string) float64 {
	v, _ := strconv.ParseFloat(s, 64)
	return v
}

func parseInt64(s string) int64 {
	v, _ := strconv.ParseInt(s, 10, 64)
	return v
}

func round2(v float64) float64 {
	return math.Round(v*100) / 100
}

func lookupName(stockCode, fallback string) string {
	lookupKey := stockCode
	if len(stockCode) == 5 && isDigits(stockCode) {
		lookupKey = "hk" + stockCode
	}
	if name, ok := stockNameMapping[lookupKey]; ok {
		return name
	}
	if name, ok := stockNameMapping[stockCode]; ok {
		return name
	}
	return fallback
}

func isDigits(s string) bool {
	for _, r := range s {
		if !unicode.IsDigit(r) {
			return false
		}
	}
	return len(s) > 0
}

func toPinyin(name string) string {
	args := pinyin.NewArgs()
	py := pinyin.Pinyin(name, args)

	parts := make([]string, 0, len(py))
	for _, p := range py {
		if len(p) > 0 {
			parts = append(parts, p[0])
		}
	}
	if len(parts) > 0 {
		return strings.Join(parts, " ")
	}
	return strings.ToLower(name)
}

func parseStockData(raw string) []Stock {
	if raw == "" {
		return nil
	}

	stocks := make([]Stock, 0)
	for _, line := range strings.Split(raw, ";") {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}

		parts := strings.SplitN(line, "=", 2)
		if len(parts) < 2 {
			continue
		}

		dataStr := strings.Trim(parts[1], "\"")
		fields := strings.Split(dataStr, "~")
		if len(fields) < 11 {
			continue
		}

		stockCode := fields[2]
		name := lookupName(stockCode, fields[1])
		price := parseFloat(fields[3])
		openPrice := parseFloat(fields[4])
		change := price - openPrice
		ratio := 0.0
		if openPrice != 0 {
			ratio = change / openPrice * 100
		}

		lowest := ""
		highest := ""
		if len(fields) > 34 {
			lowest = fields[34]
		}
		if len(fields) > 33 {
			highest = fields[33]
		}

		stocks = append(stocks, Stock{
			Name:    name,
			Pinyin:  toPinyin(name),
			Code:    stockCode,
			Price:   price,
			Open:    openPrice,
			Change:  round2(change),
			Ratio:   round2(ratio),
			High:    fields[5],
			Low:     fields[6],
			Volume:  parseInt64(fields[7]),
			Amount:  fields[8],
			Buy1:    fields[9],
			Sell1:   fields[10],
			Lowest:  lowest,
			Highest: highest,
		})
	}
	return stocks
}

func getStocks() ([]Stock, error) {
	raw, err := fetchStockData(stockCodes)
	if err != nil {
		return nil, err
	}
	return parseStockData(raw), nil
}
