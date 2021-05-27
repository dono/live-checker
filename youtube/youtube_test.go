package youtube

import (
	"fmt"
	"io/ioutil"
	"log"
	"testing"

	"github.com/jarcoal/httpmock"
	"github.com/stretchr/testify/assert"
)

func TestGetNotExistLive(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	dummyChannelID := `channel/dummy`
	dummyURL := fmt.Sprintf("https://www.youtube.com/%s", dummyChannelID)

	testHtml, err := ioutil.ReadFile("./test_html/not_exist_channel_test.html")
	if err != nil {
		t.Error(err)
	}

	httpmock.RegisterResponder(
		"GET",
		dummyURL,
		httpmock.NewStringResponder(200, string(testHtml)),
	)

	client := New()

	_, err = client.GetLive(dummyChannelID)
	if err == ErrLiveNotFound {
		return
	}
	log.Fatal(err)
}

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

	assert.Equal(t, "ON_AIR", live.Status)
	assert.Equal(t, "Marine Ch. 宝鐘マリン", live.ChannelName)
	assert.Equal(t, "【Subnautica】初見の深海…恐怖が待ってるらしい【ホロライブ/宝鐘マリン】", live.Title)
	// ToDo: Equal追加
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

	_, err = client.GetLive(dummyChannelID)
	if err == ErrLiveNotFound {
		return
	}

	t.Error(err)

	// assert.Equal(t, "ENDED", live.Status)
}
