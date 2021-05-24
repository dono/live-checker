package youtube

import (
	"fmt"
	"io/ioutil"
	"testing"

	"github.com/jarcoal/httpmock"
	"github.com/stretchr/testify/assert"
)

func TestGetOnAirLive(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	dummyChannelID := `channel/dummy`
	dummyURL := fmt.Sprintf("https://www.youtube.com/%s", dummyChannelID)

	testHtml, err := ioutil.ReadFile("./test_html/live_test.html")
	if err != nil {
		t.Error(err)
	}

	httpmock.RegisterResponder(
		"GET",
		dummyURL,
		httpmock.NewStringResponder(200, string(testHtml)),
	)

	client := New()

	live, err := client.GetLive(dummyChannelID)
	if err != nil {
		t.Error(err)
	}

	assert.Equal(t, "ウェザーニュース", live.ChannelName)
	assert.Equal(t, "https://yt3.ggpht.com/ytc/AAUvwnih5YprsDRqTVkJeBa25c1DG_kIpulgFZPG2nhN=s900-c-k-c0x00ffffff-no-rj", live.ChannelIconURL)
	assert.Equal(t, "【LIVE】 最新地震・気象情報　ウェザーニュースLiVE　2021年5月24日(月) 14時から", live.Title)
	assert.Equal(t, "【最新の天気に関する情報】お天気アプリ「ウェザーニュース」からも随時最新情報をお伝えしていますhttps://weathernews.jp/s/download/weathernewstouch.htm...", live.Description)
	assert.Equal(t, "ON_AIR", live.Status)
	assert.Equal(t, "https://www.youtube.com/channel/dummy/live", live.URL)
}

func TestGetNotOnAirLive(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	dummyChannelID := `channel/dummy`
	dummyURL := fmt.Sprintf("https://www.youtube.com/%s", dummyChannelID)

	testHtml, err := ioutil.ReadFile("./test_html/not_live_test.html")
	if err != nil {
		t.Error(err)
	}

	httpmock.RegisterResponder(
		"GET",
		dummyURL,
		httpmock.NewStringResponder(200, string(testHtml)),
	)

	client := New()

	live, err := client.GetLive(dummyChannelID)
	if err != nil {
		t.Error(err)
	}
	fmt.Println(live.ChannelName)
	fmt.Println(live.ChannelIconURL)
	fmt.Println(live.Title)
	fmt.Println(live.Description)
	fmt.Println(live.Status)
	fmt.Println(live.URL)

	// assert.Equal(t, "", live.ChannelName)
	// assert.Equal(t, "", live.ChannelIconURL)
	// assert.Equal(t, "", live.Title)
	// assert.Equal(t, "", live.Description)
	// assert.Equal(t, "", live.Status)
	// assert.Equal(t, "", live.URL)
}
