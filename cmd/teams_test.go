package cmd

import (
	"encoding/json"
	"testing"
)

func TestTeamsResponseUnmarshalWrappedData(t *testing.T) {
	var resp teamsResponse
	if err := json.Unmarshal([]byte(`{"data":[{"id":1,"name":"Core"}]}`), &resp); err != nil {
		t.Fatalf("unmarshal wrapped response: %v", err)
	}
	if len(resp.Data) != 1 || resp.Data[0].ID != 1 || resp.Data[0].Name != "Core" {
		t.Fatalf("unexpected teams: %#v", resp.Data)
	}
}

func TestTeamsResponseUnmarshalBareArray(t *testing.T) {
	var resp teamsResponse
	if err := json.Unmarshal([]byte(`[{"id":2,"name":"Agency"}]`), &resp); err != nil {
		t.Fatalf("unmarshal bare array response: %v", err)
	}
	if len(resp.Data) != 1 || resp.Data[0].ID != 2 || resp.Data[0].Name != "Agency" {
		t.Fatalf("unexpected teams: %#v", resp.Data)
	}
}
