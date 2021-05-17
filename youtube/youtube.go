package youtube

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/mattn/go-jsonpointer"
)



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

func jpToString(jsonBytes []byte, jp string) (string, error) {
	var obj interface{}
	json.Unmarshal(jsonBytes, &obj)

	v, err := jsonpointer.Get(obj, jp)
	if err != nil {
		return "", err
	}

	b, err := json.Marshal(v)
	if err != nil {
		return "", err
	}

	return string(b), nil
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
	url := `https://www.youtube.com/channel/UCNsidkYpIAQ4QaufptQBPHQ`
	// url := `https://www.youtube.com/channel/UCoSrY_IQQVpmIRZ9Xf-y93g`
	// url := `https://www.youtube.com/channel/UCXteDRy5qB0IjA8WPusCJ7w`
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

	titleJp := strings.Join([]string{
		"",
		"contents",
		"twoColumnBrowseResultsRenderer",
		"tabs",
		"0",
		"tabRenderer",
		"content",
		"sectionListRenderer",
		"contents",
		"0",
		"itemSectionRenderer",
		"contents",
		"0",
		"channelFeaturedContentRenderer",
		"items",
		"0",
		"videoRenderer",
		"title",
		"runs",
		"0",
		"text",
	}, "/")

	channelNameJp := "/metadata/channelMetadataRenderer/title"
	channelThumbnailURLJp := "/metadata/channelMetadataRenderer/avatar/thumbnails/0/url"


	title, err := jpToString([]byte(ytInitialData), titleJp)
	if err != nil {
		return nil, err
	}

	channelName, err := jpToString([]byte(ytInitialData), channelNameJp)
	if err != nil {
		return nil, err
	}

	channelThumbnailURL, err := jpToString([]byte(ytInitialData), channelThumbnailURLJp)
	if err != nil {
		return nil, err
	}


	fmt.Println(title)
	fmt.Println(channelName)
	fmt.Println(channelThumbnailURL)


	return nil, nil
}