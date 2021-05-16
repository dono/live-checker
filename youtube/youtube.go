package youtube

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/mattn/go-jsonpointer"
)

type Channel struct {
    Contents         struct {
		TwoColumnBrowseResultsRenderer struct {
			Tabs []struct {
				TabRenderer struct {
					Endpoint       interface{} `json:"endpoint"`
					Title          interface{} `json:"title"`
					Selected       interface{} `json:"selected"`
					Content        interface{} `json:"contents"`
					TrackingParams interface{} `json:"trackingParams"`
				} `json:"tabRenderer"`
			} `json:"tabs"`
		} `json:"twoColumnBrowseResultsRenderer"`
	} `json:"contents"` // 使う
	Metadata         interface{} `json:"metadata"` // 使う
	ResponseContext  interface{} `json:"responseContext"`
	Header           interface{} `json:"header"`
	TrackingParams   interface{} `json:"trackingParams"`
	Topbar           interface{} `json:"topbar"`
	Microformat      interface{} `json:"microformat"`
}


type Client struct {
	HTTPClient *http.Client
	Header *http.Header
}

type Live struct {
	ID          string
	Title       string
	Description string
	Status      string
	UserID      string
	WatchURL    string
}

type User struct {
	UserID  string
	Name    string
	IconURL string
}

func New() *Client {
    return &Client{
        HTTPClient: http.DefaultClient,
    }
}

func (c *Client)Get(url string) (*http.Response, error) {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

    resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return nil, err
	}

	return resp, nil
}


func (c *Client)GetLive(channelID string) (*Live, error) {
	url := `https://www.youtube.com/channel/UCshX7abyGGG2WqCuYC5En7w`
	resp, err := c.Get(url)

	// ytInitialData を抜き出す
	if err != nil {
		return nil, err
	}
    defer resp.Body.Close()

	doc, err := goquery.NewDocumentFromReader(resp.Body)
  	if err != nil {
    	return nil, err
  	}

	feature := "var ytInitialData = "
	ytInitialData := ""

	doc.Find("script").EachWithBreak(func(_ int, s *goquery.Selection) bool {
		if strings.Contains(s.Text(), feature) {
			t1 := strings.Replace(s.Text(), feature, "", 1)
			t2 := strings.Replace(t1, ";", "", 1)
			ytInitialData = t2
			return false
		}
		return true
	})

	
	// js, err := simplejson.NewJson(byte[](ytInitialData))
	// if err != nil {
	// 	return nil, err
	// }

	// js.GetPath("contents", "twoColumnBrowseResultsRenderer", "tabs").GetIndex(0).GetPath("tabRenderer", "contents").GetIndex(0).
	var obj interface{}
	json.Unmarshal([]byte(ytInitialData), &obj)

	// v, err := jsonpointer.Get(obj, "/contents/twoColumnBrowseResultsRenderer/tabs/0/tabRenderer/content/sectionListRenderer/contents/0/itemSectionRenderer/contents/0/channelFeaturedContentRenderer/items/0/videoRenderer")
	v, err := jsonpointer.Get(obj, "/contents/twoColumnBrowseResultsRenderer/tabs/0/tabRenderer/content/sectionListRenderer/contents/0/itemSectionRenderer/contents/0") // なぜかここより先参照できんが
	if err != nil {
		return nil, err
	}

	out, err := json.Marshal(v)
	if err != nil {
		return nil, err
	}

	fmt.Println(string(out))

	return nil, nil
}