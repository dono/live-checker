package main

import (
	"testing"

	"github.com/jarcoal/httpmock"
)

func TestGetLiveStatus(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	dummyURL := `https://com.nicovideo.jp/api/v1/communities/888888/lives.json?limit=1&offset=0`

	dummyJSON := `{
					"meta": { "status": 200 },
	                "data": {
						total": 1106,
					    "lives": [
							{"id": "lv1234",
						 	 "title": "たいとる",
							 "description": "ですくりぷしょん",
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
		httpmock.NewJsonResponder(200, dummyJSON),
	)

	client := New()

	live, err := client.GetLiveStatus("co888888")
	if err != nil {
		t.Error(err)
	}

	if live.ID != "lv1234" {
		t.Errorf("unexpected: %s\n", live.ID)
	}

	if live.Title != "たいとる" {
		t.Errorf("unexpected: %s\n", live.Title)
	}

	if live.Description != "ですくりぷしょん" {
		t.Errorf("unexpected: %s\n", live.Description)
	}

	if live.Status != "ON_AIR" {
		t.Errorf("unexpected: %s\n", live.Status)
	}

	if live.UserID != "1234" {
		t.Errorf("unexpected: %s\n", live.UserID)
	}

	if live.WatchURL != `https:\/\/live.nicovideo.jp\/watch\/lv1234` {
		t.Errorf("unexpected: %s\n", live.WatchURL)
	}
}