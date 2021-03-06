package facebook

import (
	"fmt"

	"github.com/iawia002/annie/downloader"
	"github.com/iawia002/annie/request"
	"github.com/iawia002/annie/utils"
)

// Download main download function
func Download(url string) ([]downloader.VideoData, error) {
	var err error
	html, err := request.Get(url, url, nil)
	if err != nil {
		return downloader.EmptyData, err
	}
	title := utils.MatchOneOf(html, `<title id="pageTitle">(.+)</title>`)[1]

	format := map[string]downloader.FormatData{}
	for _, quality := range []string{"sd", "hd"} {
		u := utils.MatchOneOf(
			html, fmt.Sprintf(`%s_src:"(.+?)"`, quality),
		)[1]
		size, err := request.Size(u, url)
		if err != nil {
			return downloader.EmptyData, err
		}
		urlData := downloader.URLData{
			URL:  u,
			Size: size,
			Ext:  "mp4",
		}
		format[quality] = downloader.FormatData{
			URLs:    []downloader.URLData{urlData},
			Size:    size,
			Quality: quality,
		}
	}

	return []downloader.VideoData{
		{
			Site:    "Facebook facebook.com",
			Title:   utils.FileName(title),
			Type:    "video",
			Formats: format,
		},
	}, nil
}
