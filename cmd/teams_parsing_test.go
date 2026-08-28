package cmd

import (
	"encoding/json"
	"reflect"
	"testing"
)

func TestTeamsFromAnySupportsProductionShapes(t *testing.T) {
	want := []team{{ID: 4, Name: "Core"}}
	tests := []any{
		[]any{map[string]any{"id": json.Number("4"), "name": "Core"}},
		map[string]any{"items": []any{map[string]any{"id": json.Number("4"), "name": "Core"}}},
		map[string]any{"data": []any{map[string]any{"id": json.Number("4"), "name": "Core"}}},
		[]any{map[string]any{"team": map[string]any{"id": json.Number("4"), "name": "Core"}}},
	}
	for _, input := range tests {
		if got := teamsFromAny(input); !reflect.DeepEqual(got, want) {
			t.Errorf("teamsFromAny(%#v) = %#v", input, got)
		}
	}
}
