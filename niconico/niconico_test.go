package niconico

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

	dummyCahnnelID := `co123456`
	dummyURL := `https://com.nicovideo.jp/api/v1/communities/123456/lives.json?limit=1&offset=0`
	testJson, err := ioutil.ReadFile("./test_json/not_exist_live.json")
	if err != nil {
		t.Error(err)
	}

	httpmock.RegisterResponder(
		"GET",
		dummyURL,
		httpmock.NewStringResponder(200, testJson),
	)

	client := New()

	live, err := client.GetLive(dummyCahnnelID)
	if err != nil {
		t.Error(err)
	}

	assert.Equal(t, "lv1234", live.ID)
	assert.Equal(t, "test title", live.Title)
	assert.Equal(t, "test description", live.Description)
	assert.Equal(t, "ON_AIR", live.Status)
	assert.Equal(t, "1234", live.UserID)
	assert.Equal(t, "https://live.nicovideo.jp/watch/lv1234", live.WatchURL)
}

func TestGetNotOnAirLive(t *testing.T) {
	client := New()

	live, err := client.GetLive("co3782975")
	if err != nil {
		t.Error(err)
	}
	fmt.Println(live)
}

func TestGetUser(t *testing.T) {
	client := New()

	// mock用意するのめんどいわ
	user, err := client.GetUser("26578404")
	if err != nil {
		t.Error(err)
	}

	assert.Equal(t, "teguru", user.Name)
	assert.Equal(t, "26578404", user.UserID)
	assert.Equal(t, "https://secure-dcdn.cdn.nimg.jp/nicoaccount/usericon/2657/26578404.jpg", user.IconURL)
}
