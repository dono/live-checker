package youtube

import (
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/dono/live-checker/utils"
)

var ErrLiveNotFound = errors.New("live not found")

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

type Live struct {
	ID             string
	ChannelName    string
	ChannelIconURL string
	Title          string
	Description    string
	Status         string
	URL            string
}

func isOnAir(ytInitialData string) bool {
	feature := `"style":"LIVE","icon":{"iconType":"LIVE"}`
	return strings.Contains(ytInitialData, feature)
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

func (c *Client) GetLive(channelID string) (*Live, error) {
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

	if !isOnAir(ytInitialData) {
		return nil, ErrLiveNotFound // ToDo: 配信してないのかチャンネル自体が無いのか区別したほうが良さげ
	}

	title, err := utils.JpToString([]byte(ytInitialData), titleJP)
	if err != nil {
		return nil, err
	}

	descriptionSnippet, err := utils.JpToString([]byte(ytInitialData), descriptionJP)
	if err != nil {
		return nil, err
	}
	description := strings.Replace(descriptionSnippet, "\\n", "", -1)

	channelName, err := utils.JpToString([]byte(ytInitialData), channelNameJP)
	if err != nil {
		return nil, err
	}

	channelIconURL, err := utils.JpToString([]byte(ytInitialData), channelThumbnailURLJP)
	if err != nil {
		return nil, err
	}

	liveURL := fmt.Sprintf("%s/live", channelURL)

	return &Live{
		ID:             channelID,
		ChannelName:    channelName,
		ChannelIconURL: channelIconURL,
		Title:          title,
		Description:    description,
		Status:         "ON_AIR",
		URL:            liveURL,
	}, nil
}
