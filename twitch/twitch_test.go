package twitch

import (
	"fmt"
	"io/ioutil"
	"log"
	"testing"

	"github.com/jarcoal/httpmock"
	"github.com/stretchr/testify/assert"
)

func TestGetOnAirLive(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	dummyChannelID := `dummyID`
	dummyURL := fmt.Sprintf("https://www.twitch.tv/%s", dummyChannelID)

	testHtml, err := ioutil.ReadFile("./test_html/on_air_live_test.html")
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

	assert.Equal(t, "ON_AIR", live.Status)
	assert.Equal(t, "dummyID", live.ID)
	assert.Equal(t, "Jasper7se", live.ChannelName)
	assert.Equal(t, "Jasper7se - Twitch", live.Title)
	assert.Equal(t, "dummy description", live.Description)
}

func TestGetNotOnAirLive(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	dummyChannelID := `dummy`
	dummyURL := fmt.Sprintf("https://www.twitch.tv/%s", dummyChannelID)

	testHtml, err := ioutil.ReadFile("./test_html/not_on_air_live_test.html")
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
	if live.Status != "NOT_ON_AIR" {
		log.Fatal(err)
	}
}

func TestGetNotExistChannel(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	dummyChannelID := `dummy`
	dummyURL := fmt.Sprintf("https://www.twitch.tv/%s", dummyChannelID)

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
	live, err := client.GetLive(dummyChannelID)
	if live.Status != "CHANNEL_NOT_FOUND" {
		log.Fatal(err)
	}
}
