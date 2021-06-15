package youtube

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/dono/live-checker/entity"
	"github.com/dono/live-checker/status"
	"github.com/dono/live-checker/utils"
)

const (
	titleJP string = "/contents" + "/twoColumnBrowseResultsRenderer" + "/tabs" + "/0" + "/tabRenderer" +
		"/content" + "/sectionListRenderer" + "/contents" + "/0" + "/itemSectionRenderer" +
		"/contents" + "/0" + "/channelFeaturedContentRenderer" + "/items" + "/0" +
		"/videoRenderer" + "/title" + "/runs" + "/0" + "/text"

	descriptionJP string = "/contents" + "/twoColumnBrowseResultsRenderer" + "/tabs" + "/0" + "/tabRenderer" +
		"/content" + "/sectionListRenderer" + "/contents" + "/0" + "/itemSectionRenderer" +
		"/contents" + "/0" + "/channelFeaturedContentRenderer" + "/items" + "/0" +
		"/videoRenderer" + "/descriptionSnippet" + "/runs" + "/0" + "/text"

	channelNameJP string = "/metadata" + "/channelMetadataRenderer" + "/title"

	channelThumbnailURLJP string = "/metadata" + "/channelMetadataRenderer" + "/avatar" + "/thumbnails" + "/0" + "/url"
)

type Client struct {
	HTTPClient *http.Client
	Header     *http.Header
}

func isOnAir(ytInitialData string) bool {
	feature := `"style":"LIVE","icon":{"iconType":"LIVE"}`
	return strings.Contains(ytInitialData, feature)
}

func isExistChannel(ytInitialData string) bool {
	feature := `"type":"ERROR"`
	if ytInitialData == "" {
		return false
	} else if strings.Contains(ytInitialData, feature) {
		return false
	}
	return true
}

func New() *Client {
	return &Client{
		HTTPClient: http.DefaultClient,
	}
}

func (c *Client) Get(url string) (*http.Response, error) {
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

func (c *Client) GetLive(channelID string) (*entity.Live, error) {
	channelURL := fmt.Sprintf("https://www.youtube.com/%s", channelID)

	resp, err := c.Get(channelURL)
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
			t := strings.Replace(s.Text(), feature, "", 1)
			ytInitialData = strings.Replace(t, ";", "", 1)
			return false
		}
		return true
	})

	if !isExistChannel(ytInitialData) {
		return &entity.Live{
			Status: status.CHANNEL_NOT_FOUND,
		}, nil
	}

	if !isOnAir(ytInitialData) {
		return &entity.Live{
			Status: status.NOT_ON_AIR,
		}, nil
	}

	var obj interface{}
	err = json.Unmarshal([]byte(ytInitialData), &obj)
	if err != nil {
		return nil, err
	}

	title, err := utils.JpToString(obj, titleJP)
	if err != nil {
		return nil, err
	}

	descriptionSnippet, err := utils.JpToString(obj, descriptionJP)
	if err != nil {
		return nil, err
	}
	description := strings.Replace(descriptionSnippet, "\\n", "", -1)

	channelName, err := utils.JpToString(obj, channelNameJP)
	if err != nil {
		return nil, err
	}

	channelIconURL, err := utils.JpToString(obj, channelThumbnailURLJP)
	if err != nil {
		return nil, err
	}

	liveURL := fmt.Sprintf("%s/live", channelURL)

	return &entity.Live{
		Platform:    "youtube",
		ID:          channelID,
		Name:        channelName,
		Title:       title,
		Description: description,
		Status:      status.ON_AIR,
		IconURL:     channelIconURL,
		WatchURL:    liveURL,
	}, nil
}

func (c *Client) GetLives(channelIDs []string) ([]*entity.Live, error) {
	lives := []*entity.Live{}
	for _, channelID := range channelIDs {
		live, err := c.GetLive(channelID)
		if err != nil {
			return nil, err
		}
		lives = append(lives, live)
	}

	return lives, nil
}
