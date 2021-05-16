package youtube

import (
	"testing"
)

func TestGetLive(t *testing.T) {
	client := New()

	_, err := client.GetLive("dummy")
	if err != nil {
		t.Error(err)
	}
}