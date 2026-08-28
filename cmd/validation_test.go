package cmd

import (
	"strings"
	"testing"
	"time"

	"github.com/spf13/cobra"
)

func TestExactPositiveIDArg(t *testing.T) {
	validator := exactPositiveIDArg("post ID")
	cmd := &cobra.Command{Use: "get <id>"}
	if err := validator(cmd, []string{"42"}); err != nil {
		t.Fatalf("valid ID: %v", err)
	}
	for _, args := range [][]string{nil, {"0"}, {"-1"}, {"abc"}, {"1", "2"}} {
		if err := validator(cmd, args); err == nil {
			t.Errorf("args %#v were accepted", args)
		}
	}
}

func TestValidateDateRange(t *testing.T) {
	if err := validateDateRange("2026-08-01", "2026-08-31"); err != nil {
		t.Fatalf("valid range: %v", err)
	}
	if err := validateDateRange("2026-08-01", "2026-08-01"); err != nil {
		t.Fatalf("same-day range: %v", err)
	}
	for _, dates := range [][2]string{{"", "2026-08-01"}, {"2026/08/01", "2026-08-02"}, {"2026-08-02", "bad"}, {"2026-08-02", "2026-08-01"}} {
		if err := validateDateRange(dates[0], dates[1]); err == nil {
			t.Errorf("range %#v was accepted", dates)
		}
	}
}

func TestValidateDateTime(t *testing.T) {
	if err := validateDateTime("2026-09-01 12:30:00", "--publish-at"); err != nil {
		t.Fatalf("valid date-time: %v", err)
	}
	for _, value := range []string{"", "2026-09-01T12:30:00Z", "2026-02-30 12:30:00"} {
		if err := validateDateTime(value, "--publish-at"); err == nil {
			t.Errorf("date-time %q was accepted", value)
		}
	}
}

func TestValidatePublishTime(t *testing.T) {
	now := time.Date(2026, 8, 28, 12, 0, 0, 0, time.UTC)
	if err := validatePublishTime("2026-08-28 12:00:01", false, now); err != nil {
		t.Fatalf("future publish time: %v", err)
	}
	if err := validatePublishTime("2026-08-28 11:59:59", true, now); err != nil {
		t.Fatalf("draft publish time: %v", err)
	}
	if err := validatePublishTime("2026-08-28 11:59:59", false, now); err == nil {
		t.Fatal("past non-draft publish time was accepted")
	}
}

func TestValidateChoice(t *testing.T) {
	if err := validateChoice("all", "--type", "all", "user", "shared"); err != nil {
		t.Fatalf("valid choice: %v", err)
	}
	err := validateChoice("invalid", "--type", "all", "user", "shared")
	if err == nil || !strings.Contains(err.Error(), "all, user, shared") {
		t.Fatalf("error = %v", err)
	}
}

func TestValidatePositiveInts(t *testing.T) {
	if err := validatePositiveInts([]int{1, 2}, "--accounts"); err != nil {
		t.Fatalf("valid IDs: %v", err)
	}
	if err := validatePositiveInts([]int{1, 0}, "--accounts"); err == nil {
		t.Fatal("zero ID was accepted")
	}
}
