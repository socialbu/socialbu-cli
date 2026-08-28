package output

import (
	"bytes"
	"encoding/json"
	"io"
	"os"
	"reflect"
	"strings"
	"testing"
)

func TestStringHandlesDecodedJSONValues(t *testing.T) {
	tests := []struct {
		value any
		want  string
	}{
		{value: nil, want: ""},
		{value: " value ", want: "value"},
		{value: json.Number("9007199254740993"), want: "9007199254740993"},
		{value: float64(12.5), want: "12.5"},
		{value: float64(12), want: "12"},
		{value: true, want: "true"},
	}
	for _, tt := range tests {
		if got := String(tt.value); got != tt.want {
			t.Errorf("String(%#v) = %q, want %q", tt.value, got, tt.want)
		}
	}
}

func TestIntFromMapHandlesDecodedJSONNumber(t *testing.T) {
	item := map[string]any{"id": json.Number("190345")}
	if got := IntFromMap(item, "id"); got != 190345 {
		t.Fatalf("IntFromMap = %d", got)
	}
}

func TestJoinStringsHandlesStringsAndNestedObjects(t *testing.T) {
	if got := JoinStrings(map[string]any{"authors": "Usama Ejaz"}, "authors"); got != "Usama Ejaz" {
		t.Fatalf("string value = %q", got)
	}
	item := map[string]any{
		"topics": []any{
			"social media",
			map[string]any{"topic": map[string]any{"name": "automation"}},
			map[string]any{"attributes": map[string]any{"title": "analytics"}},
		},
	}
	if got := JoinStrings(item, "topics"); got != "social media, automation, analytics" {
		t.Fatalf("nested values = %q", got)
	}
}

func TestMapAndSliceHelpersUseFirstMatchingKey(t *testing.T) {
	nested := map[string]any{"id": json.Number("7")}
	item := map[string]any{
		"data":  nested,
		"items": []any{nested},
	}
	if got := MapFromMap(item, "missing", "data"); !reflect.DeepEqual(got, nested) {
		t.Fatalf("MapFromMap = %#v", got)
	}
	if got := SliceFromMap(item, "missing", "items"); len(got) != 1 || IntFromMap(got[0], "id") != 7 {
		t.Fatalf("SliceFromMap = %#v", got)
	}
}

func TestStringFromNestedMapFallsBackAcrossPaths(t *testing.T) {
	item := map[string]any{"result": map[string]any{"content": "generated"}}
	got := StringFromNestedMap(item, []string{"data", "content"}, []string{"result", "content"})
	if got != "generated" {
		t.Fatalf("StringFromNestedMap = %q", got)
	}
}

func TestKeysAreSorted(t *testing.T) {
	got := Keys(map[string]any{"z": 1, "a": 2, "m": 3})
	want := []string{"a", "m", "z"}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("Keys = %#v", got)
	}
}

func TestOutputWriters(t *testing.T) {
	jsonText := captureOutput(t, func() {
		if err := JSON(map[string]any{"ok": true}); err != nil {
			t.Fatalf("JSON: %v", err)
		}
	})
	if !strings.Contains(jsonText, `"ok": true`) {
		t.Fatalf("JSON output = %q", jsonText)
	}

	table := captureOutput(t, func() {
		Table([]string{"Name", "Status"}, [][]string{{"Main\tPage", "active\nnow"}})
	})
	for _, want := range []string{"NAME", "STATUS", "Main Page", "active now"} {
		if !strings.Contains(table, want) {
			t.Fatalf("table missing %q:\n%s", want, table)
		}
	}

	keyValues := captureOutput(t, func() {
		KeyValue("Config", map[string]string{"z": "last", "a": "first"})
	})
	if !strings.HasPrefix(keyValues, "Config\na: first\nz: last\n") {
		t.Fatalf("key/value output = %q", keyValues)
	}

	list := captureOutput(t, func() {
		PrintSection("Items")
		PrintList([]string{"one", " ", "two"})
	})
	for _, want := range []string{"Items\n-----", "- one", "- two"} {
		if !strings.Contains(list, want) {
			t.Fatalf("list output missing %q:\n%s", want, list)
		}
	}
}

func TestBoolFromMap(t *testing.T) {
	if !BoolFromMap(map[string]any{"active": true}, "active") {
		t.Fatal("boolean true was not detected")
	}
	if !BoolFromMap(map[string]any{"active": " TRUE "}, "active") {
		t.Fatal("string true was not detected")
	}
	if BoolFromMap(map[string]any{"active": false}, "active") {
		t.Fatal("false value was reported as true")
	}
}

func captureOutput(t *testing.T, fn func()) string {
	t.Helper()
	old := os.Stdout
	reader, writer, err := os.Pipe()
	if err != nil {
		t.Fatalf("pipe: %v", err)
	}
	os.Stdout = writer
	defer func() { os.Stdout = old }()

	fn()
	_ = writer.Close()
	var buffer bytes.Buffer
	if _, err := io.Copy(&buffer, reader); err != nil {
		t.Fatalf("copy output: %v", err)
	}
	_ = reader.Close()
	return buffer.String()
}
