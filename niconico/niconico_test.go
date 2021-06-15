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

	dummyChannelID := `co123456`
	dummyURL := `https://com.nicovideo.jp/api/v1/communities/123456/lives.json`
	testJson, err := ioutil.ReadFile("./test_json/on_air_live_test.json")
	if err != nil {
		t.Error(err)
	}

	dummyUserID := `12345`
	dummyUserURL := fmt.Sprintf("https://public.api.nicovideo.jp/v1/users.json?userIds=%s", dummyUserID)
	testUserJson, err := ioutil.ReadFile("./test_json/user_test.json")
	if err != nil {
		t.Error(err)
	}

	httpmock.RegisterResponder(
		"GET",
		dummyURL,
		httpmock.NewStringResponder(200, string(testJson)),
	)

	httpmock.RegisterResponder(
		"GET",
		dummyUserURL,
		httpmock.NewStringResponder(200, string(testUserJson)),
	)

	client := New()

	live, err := client.GetLive(dummyChannelID)
	if err != nil {
		t.Error(err)
	}

	assert.Equal(t, "co123456", live.ID)
	assert.Equal(t, "test title", live.Title)
	assert.Equal(t, "test description", live.Description)
	assert.Equal(t, "ON_AIR", live.Status)
	assert.Equal(t, "https://live.nicovideo.jp/watch/lv222222", live.WatchURL)
}

func TestGetNotOnAirLive(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	dummyCahnnelID := `co123456`
	dummyURL := `https://com.nicovideo.jp/api/v1/communities/123456/lives.json`
	testJson, err := ioutil.ReadFile("./test_json/not_on_air_live_test.json")
	if err != nil {
		t.Error(err)
	}

	httpmock.RegisterResponder(
		"GET",
		dummyURL,
		httpmock.NewStringResponder(200, string(testJson)),
	)

	client := New()

	live, err := client.GetLive(dummyCahnnelID)
	if err != nil {
		t.Error(err)
	}

	assert.Equal(t, "NOT_ON_AIR", live.Status)
}

func TestGetNotExistCommunityLive(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	dummyCahnnelID := `co123456`
	dummyURL := `https://com.nicovideo.jp/api/v1/communities/123456/lives.json`
	testJson, err := ioutil.ReadFile("./test_json/not_exist_community_test.json")
	if err != nil {
		t.Error(err)
	}

	httpmock.RegisterResponder(
		"GET",
		dummyURL,
		httpmock.NewStringResponder(403, string(testJson)),
	)

	client := New()

	live, err := client.GetLive(dummyCahnnelID)
	if err != nil {
		t.Error(err)
	}

	assert.Equal(t, "CHANNEL_NOT_FOUND", live.Status)
}

func TestGetUser(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	dummyUserID := `123456`
	dummyURL := fmt.Sprintf("https://public.api.nicovideo.jp/v1/users.json?userIds=%s", dummyUserID)
	testJson, err := ioutil.ReadFile("./test_json/user_test.json")
	if err != nil {
		t.Error(err)
	}

	httpmock.RegisterResponder(
		"GET",
		dummyURL,
		httpmock.NewStringResponder(403, string(testJson)),
	)

	client := New()

	user, err := client.GetUser(dummyUserID)
	if err != nil {
		t.Error(err)
	}

	assert.Equal(t, "hoge", user.Name)
	assert.Equal(t, "123456", user.UserID)
}
