package buffer

import (
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/parnurzeal/gorequest"
	"github.com/wlcy/tron/explorer/lib/log"
	"github.com/wlcy/tron/explorer/lib/mysql"
	"github.com/wlcy/tron/explorer/web/entity"
)

/*
store all marketBuffer data in memory
load from db every 30 seconds
*/

var _marketBuffer *marketBuffer
var onceMarketOnce sync.Once

//GetMarketBuffer ...
func GetMarketBuffer() *marketBuffer {
	return getMarketBuffer()
}

// getMarketBuffer
func getMarketBuffer() *marketBuffer {
	onceMarketOnce.Do(func() {
		_marketBuffer = &marketBuffer{}
		_marketBuffer.load()

		go func() {
			time.Sleep(5 * time.Second)
			_marketBuffer.load()
		}()
	})
	return _marketBuffer
}

type marketBuffer struct {
	sync.RWMutex

	marketInfoList []*entity.MarketInfo

	updateTime string
}

func (w *marketBuffer) GetMarket() (witness []*entity.MarketInfo) {
	if len(w.marketInfoList) == 0 {
		log.Infof("get market info from buffer nil, data reload at :[%v]", w.updateTime)
		w.load()
	}
	log.Infof("get market info from buffer, buffer data updated at :[%v]", w.updateTime)
	return w.marketInfoList
}

func (w *marketBuffer) load() {
	marketInfos := make([]*entity.MarketInfo, 0)
	marketURL := "https://coinmarketcap.com/currencies/tron/"
	_, body, errs := gorequest.New().Get(marketURL).End()
	if errs != nil && len(errs) > 0 {
		log.Error(errs)
		return
	}
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(body))
	if err != nil {
		log.Error(err)
		return
	}
	doc.Find("#markets-table > tbody > tr").Each(func(i int, s *goquery.Selection) {
		marketInfo := &entity.MarketInfo{}
		node := strconv.Itoa(i + 1)
		rank, _ := s.Find("tr:nth-child(" + node + ") > td:nth-child(1)").Html()
		name, _ := s.Find("tr:nth-child(" + node + ") > td:nth-child(2)").Attr("data-sort")
		pair, _ := s.Find("tr:nth-child(" + node + ") > td:nth-child(3)").Attr("data-sort")
		link, _ := s.Find("tr:nth-child(" + node + ") > td:nth-child(3) > a").Attr("href")
		volume, _ := s.Find("tr:nth-child(" + node + ") > td:nth-child(4) > span[class=volume]").Attr("data-usd")
		volumeNative, _ := s.Find("tr:nth-child(" + node + ") > td:nth-child(4) > span[class=volume]").Attr("data-native")
		price, _ := s.Find("tr:nth-child(" + node + ") > td:nth-child(5) > span[class=price]").Attr("data-usd")
		volumePercentage, _ := s.Find("tr:nth-child(" + node + ") > td:nth-child(6)").Attr("data-sort")
		marketInfo.Rank = mysql.ConvertStringToInt64(rank, 0)
		marketInfo.Name = name
		marketInfo.Pair = pair
		marketInfo.Link = link
		marketInfo.Volume = mysql.ConvertStringToFloat(volume, 0)
		marketInfo.VolumeNative = mysql.ConvertStringToFloat(volumeNative, 0)
		marketInfo.VolumePercentage = mysql.ConvertStringToFloat(volumePercentage, 0)
		marketInfo.Price = mysql.ConvertStringToFloat(price, 0)
		marketInfos = append(marketInfos, marketInfo)
	})

	log.Infof("market in buffer : parse page data done.")
	w.Lock()
	w.marketInfoList = marketInfos
	w.updateTime = time.Now().Local().Format(mysql.DATETIMEFORMAT)
	w.Unlock()
}
