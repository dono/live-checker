package niconico

import (
	"testing"

	"github.com/jarcoal/httpmock"
	"github.com/stretchr/testify/assert"
)

func TestGetLive(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	dummyURL := `https://com.nicovideo.jp/api/v1/communities/888888/lives.json?limit=1&offset=0`

	dummyJSON := `{
					"meta": {"status": 200},
	                "data": {
						"total": 1106,
					    "lives": [
							{"id": "lv1234",
						 	 "title": "test title",
							 "description": "test description",
							 "status": "ON_AIR",
							 "user_id": 1234,
							 "watch_url": "https:\/\/live.nicovideo.jp\/watch\/lv1234",
							 "features": {"is_member_only":false},
							 "timeshift": {"enabled":true,"can_view":false},
							 "started_at": "2021-05-09T20:47:36+0900"}
						]
					}
				}`

	

	httpmock.RegisterResponder(
		"GET",
		dummyURL,
		httpmock.NewStringResponder(200, dummyJSON),
	)

	client := New()

	live, err := client.GetLive("co888888")
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