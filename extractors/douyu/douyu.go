package douyu

import (
	"encoding/json"
	"errors"

	"github.com/iawia002/annie/downloader"
	"github.com/iawia002/annie/request"
	"github.com/iawia002/annie/utils"
)

type douyuData struct {
	Error int `json:"error"`
	Data  struct {
		VideoURL string `json:"video_url"`
	} `json:"data"`
}

type douyuURLInfo struct {
	URL  string
	Size int64
}

func douyuM3u8(url string) ([]douyuURLInfo, int64, error) {
	var (
		data            []douyuURLInfo
		temp            douyuURLInfo
		size, totalSize int64
		err             error
	)
	urls, err := utils.M3u8URLs(url)
	if err != nil {
		return nil, 0, err
	}
	for _, u := range urls {
		size, err = request.Size(u, url)
		if err != nil {
			return nil, 0, err
		}
		totalSize += size
		temp = douyuURLInfo{
			URL:  u,
			Size: size,
		}
		data = append(data, temp)
	}
	return data, totalSize, nil
}

// Download main download function
func Download(url string) ([]downloader.VideoData, error) {
	var err error
	liveVid := utils.MatchOneOf(url, `https?://www.douyu.com/(\S+)`)
	if liveVid != nil {
		return downloader.EmptyData, errors.New("暂不支持斗鱼直播")
	}

	html, err := request.Get(url, url, nil)
	if err != nil {
		return downloader.EmptyData, err
	}
	title := utils.MatchOneOf(html, `<title>(.*?)</title>`)[1]

	vid := utils.MatchOneOf(url, `https?://v.douyu.com/show/(\S+)`)[1]
	dataString, err := request.Get("http://vmobile.douyu.com/video/getInfo?vid="+vid, url, nil)
	if err != nil {
		return downloader.EmptyData, err
	}
	var dataDict douyuData
	json.Unmarshal([]byte(dataString), &dataDict)

	m3u8URLs, totalSize, err := douyuM3u8(dataDict.Data.VideoURL)
	if err != nil {
		return downloader.EmptyData, err
	}
	urls := make([]downloader.URLData, len(m3u8URLs))
	for index, u := range m3u8URLs {
		urls[index] = downloader.URLData{
			URL:  u.URL,
			Size: u.Size,
			Ext:  "ts",
		}
	}

	format := map[string]downloader.FormatData{
		"default": {
			URLs: urls,
			Size: totalSize,
		},
	}
	return []downloader.VideoData{
		{
			Site:    "斗鱼 douyu.com",
			Title:   title,
			Type:    "video",
			Formats: format,
		},
	}, nil
}
