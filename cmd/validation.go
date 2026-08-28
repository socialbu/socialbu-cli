package cmd

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/spf13/cobra"
)

func exactPositiveIDArg(name string) cobra.PositionalArgs {
	return func(cmd *cobra.Command, args []string) error {
		if err := cobra.ExactArgs(1)(cmd, args); err != nil {
			return err
		}
		id, err := strconv.Atoi(args[0])
		if err != nil || id <= 0 {
			return fmt.Errorf("%s must be a positive integer", name)
		}
		return nil
	}
}

func validateDateRange(start, end string) error {
	if strings.TrimSpace(start) == "" || strings.TrimSpace(end) == "" {
		return fmt.Errorf("--start and --end are required")
	}
	startDate, err := time.Parse("2006-01-02", start)
	if err != nil {
		return fmt.Errorf("--start must use YYYY-MM-DD")
	}
	endDate, err := time.Parse("2006-01-02", end)
	if err != nil {
		return fmt.Errorf("--end must use YYYY-MM-DD")
	}
	if endDate.Before(startDate) {
		return fmt.Errorf("--end cannot be before --start")
	}
	return nil
}

func validateOptionalDateRange(start, end string) error {
	if start == "" && end == "" {
		return nil
	}
	if start == "" || end == "" {
		return fmt.Errorf("--from and --to must be used together")
	}
	startDate, err := time.Parse("2006-01-02", start)
	if err != nil {
		return fmt.Errorf("--from must use YYYY-MM-DD")
	}
	endDate, err := time.Parse("2006-01-02", end)
	if err != nil {
		return fmt.Errorf("--to must use YYYY-MM-DD")
	}
	if endDate.Before(startDate) {
		return fmt.Errorf("--to cannot be before --from")
	}
	return nil
}

func validateDateTime(value, flag string) error {
	if strings.TrimSpace(value) == "" {
		return fmt.Errorf("%s is required", flag)
	}
	if _, err := time.Parse("2006-01-02 15:04:05", value); err != nil {
		return fmt.Errorf("%s must use YYYY-MM-DD HH:MM:SS in UTC", flag)
	}
	return nil
}

func validatePublishTime(value string, draft bool, now time.Time) error {
	if err := validateDateTime(value, "--publish-at"); err != nil {
		return err
	}
	publishAt, _ := time.Parse("2006-01-02 15:04:05", value)
	if !draft && !publishAt.After(now.UTC()) {
		return fmt.Errorf("--publish-at must be in the future unless --draft is set")
	}
	return nil
}

func validateChoice(value, flag string, choices ...string) error {
	for _, choice := range choices {
		if value == choice {
			return nil
		}
	}
	return fmt.Errorf("%s must be one of: %s", flag, strings.Join(choices, ", "))
}

func validatePositiveInts(values []int, flag string) error {
	for _, value := range values {
		if value <= 0 {
			return fmt.Errorf("%s values must be positive integers", flag)
		}
	}
	return nil
}
