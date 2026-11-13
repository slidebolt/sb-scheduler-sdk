package scheduler

import (
	"encoding/json"
	"testing"

	domain "github.com/slidebolt/sb-domain"
)

func TestCreateRequestRoundTrip(t *testing.T) {
	want := CreateTimerRequest{ID: "news", Name: "News", Timer: domain.Timer{Enabled: true, Target: "plugin.test.recording"}}
	b, err := json.Marshal(want)
	if err != nil {
		t.Fatal(err)
	}
	var got CreateTimerRequest
	if err := json.Unmarshal(b, &got); err != nil {
		t.Fatal(err)
	}
	if got.ID != want.ID || got.Timer.Target != want.Timer.Target {
		t.Fatalf("got %+v", got)
	}
}
