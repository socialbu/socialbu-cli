package cmd

import (
	"context"
	"fmt"
	"net/url"
	"strings"

	"github.com/socialbu/socialbu-cli/internal/output"
	"github.com/spf13/cobra"
)

func newAnalyticsCmd() *cobra.Command {
	cmd := &cobra.Command{Use: "analytics", Short: "Inspect SocialBu insights and analytics"}
	cmd.AddCommand(newAnalyticsPostsCountCmd())
	cmd.AddCommand(newAnalyticsPostsMetricsCmd())
	cmd.AddCommand(newAnalyticsTopPostsCmd())
	cmd.AddCommand(newAnalyticsAccountsMetricsCmd())
	cmd.AddCommand(newAnalyticsFollowersCmd())
	cmd.AddCommand(newAnalyticsEngagementTrendCmd())
	cmd.AddCommand(newAnalyticsEngagementRateCmd())
	cmd.AddCommand(newAnalyticsInboxUnreadCountCmd())
	cmd.AddCommand(newAnalyticsAutomationLogsCmd())
	cmd.AddCommand(newAnalyticsFollowersGrowthCmd())
	cmd.AddCommand(newAnalyticsTeamMetricsCmd())
	cmd.AddCommand(newAnalyticsTeamActivityCmd())
	cmd.AddCommand(newAnalyticsStatsCmd())
	return cmd
}

func newAnalyticsPostsCountCmd() *cobra.Command {
	var start, end, postType string
	var accounts []int
	var team int
	cmd := &cobra.Command{
		Use:   "posts-count",
		Short: "Get date-wise published post counts",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := validateDateRange(start, end); err != nil {
				return err
			}
			if postType != "" {
				if err := validateChoice(postType, "--post-type", "image", "video", "text"); err != nil {
					return err
				}
			}
			if err := validateAnalyticsFilters(accounts, team); err != nil {
				return err
			}
			q := dateRangeQuery(start, end)
			q = addOptionalString(q, "post_type", postType)
			q = addOptionalInt(q, "team", team)
			q = addIntSlice(q, "accounts", accounts)
			var resp map[string]any
			if err := analyticsRequest(cmd.Context(), "/insights/posts/counts", q, &resp); err != nil {
				return err
			}
			return renderAnalyticsPostsCount(resp)
		},
	}
	addDateRangeFlags(cmd, &start, &end)
	cmd.Flags().StringVar(&postType, "post-type", "", "filter by post type: image, video, text")
	cmd.Flags().IntSliceVar(&accounts, "accounts", nil, "account IDs")
	cmd.Flags().IntVar(&team, "team", 0, "team ID")
	return cmd
}

func newAnalyticsPostsMetricsCmd() *cobra.Command {
	var start, end, postType, metrics string
	var accounts []int
	var team int
	cmd := &cobra.Command{
		Use:   "posts-metrics",
		Short: "Get date-wise post metrics",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			if strings.TrimSpace(postType) == "" || strings.TrimSpace(metrics) == "" {
				return fmt.Errorf("--post-type and --metrics are required")
			}
			if err := validateDateRange(start, end); err != nil {
				return err
			}
			if err := validateChoice(postType, "--post-type", "image", "video", "text"); err != nil {
				return err
			}
			if err := validateAnalyticsFilters(accounts, team); err != nil {
				return err
			}
			q := dateRangeQuery(start, end)
			q.Set("post_type", postType)
			q.Set("metrics", metrics)
			q = addOptionalInt(q, "team", team)
			q = addIntSlice(q, "accounts", accounts)
			var resp map[string]any
			if err := analyticsRequest(cmd.Context(), "/insights/posts/metrics", q, &resp); err != nil {
				return err
			}
			return output.JSON(resp)
		},
	}
	addDateRangeFlags(cmd, &start, &end)
	cmd.Flags().StringVar(&postType, "post-type", "", "post type: image, video, text")
	cmd.Flags().StringVar(&metrics, "metrics", "", "comma-separated metric names")
	cmd.Flags().IntSliceVar(&accounts, "accounts", nil, "account IDs")
	cmd.Flags().IntVar(&team, "team", 0, "team ID")
	return cmd
}

func newAnalyticsTopPostsCmd() *cobra.Command {
	var start, end, metrics string
	var accounts []int
	var team int
	cmd := &cobra.Command{
		Use:   "top-posts",
		Short: "Get top posts ranked by metrics",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			if strings.TrimSpace(metrics) == "" {
				return fmt.Errorf("--metrics is required")
			}
			if err := validateDateRange(start, end); err != nil {
				return err
			}
			if err := validateAnalyticsFilters(accounts, team); err != nil {
				return err
			}
			q := dateRangeQuery(start, end)
			q.Set("metrics", metrics)
			q = addOptionalInt(q, "team", team)
			q = addIntSlice(q, "accounts", accounts)
			var resp map[string]any
			if err := analyticsRequest(cmd.Context(), "/insights/posts/top_posts", q, &resp); err != nil {
				return err
			}
			return output.JSON(resp)
		},
	}
	addDateRangeFlags(cmd, &start, &end)
	cmd.Flags().StringVar(&metrics, "metrics", "", "comma-separated metric names used for ranking")
	cmd.Flags().IntSliceVar(&accounts, "accounts", nil, "account IDs")
	cmd.Flags().IntVar(&team, "team", 0, "team ID")
	return cmd
}

func newAnalyticsAccountsMetricsCmd() *cobra.Command {
	var start, end, metrics string
	var accounts []int
	var team int
	var calculateGrowth bool
	cmd := &cobra.Command{
		Use:   "accounts-metrics",
		Short: "Get date-wise account metrics",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			if strings.TrimSpace(metrics) == "" {
				return fmt.Errorf("--metrics is required")
			}
			if err := validateDateRange(start, end); err != nil {
				return err
			}
			if err := validateAnalyticsFilters(accounts, team); err != nil {
				return err
			}
			q := dateRangeQuery(start, end)
			q.Set("metrics", metrics)
			q.Set("calculate_growth", fmt.Sprintf("%t", calculateGrowth))
			q = addOptionalInt(q, "team", team)
			q = addIntSlice(q, "accounts", accounts)
			var resp map[string]any
			if err := analyticsRequest(cmd.Context(), "/insights/accounts/metrics", q, &resp); err != nil {
				return err
			}
			return output.JSON(resp)
		},
	}
	addDateRangeFlags(cmd, &start, &end)
	cmd.Flags().StringVar(&metrics, "metrics", "", "comma-separated metric names")
	cmd.Flags().BoolVar(&calculateGrowth, "calculate-growth", true, "include growth calculations")
	cmd.Flags().IntSliceVar(&accounts, "accounts", nil, "account IDs")
	cmd.Flags().IntVar(&team, "team", 0, "team ID")
	return cmd
}

func newAnalyticsFollowersCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "followers",
		Short: "Get total followers across accessible accounts",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			var resp map[string]any
			if err := analyticsRequest(cmd.Context(), "/insights/accounts/followers", nil, &resp); err != nil {
				return err
			}
			return renderAnalyticsStats(resp)
		},
	}
}

func newAnalyticsEngagementTrendCmd() *cobra.Command {
	var start, end string
	cmd := &cobra.Command{
		Use:   "engagement-trend",
		Short: "Get daily total engagement trend",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := validateDateRange(start, end); err != nil {
				return err
			}
			var resp map[string]any
			if err := analyticsRequest(cmd.Context(), "/insights/accounts/engagement/trend", dateRangeQuery(start, end), &resp); err != nil {
				return err
			}
			return renderAnalyticsEngagementTrend(resp)
		},
	}
	addDateRangeFlags(cmd, &start, &end)
	return cmd
}

func newAnalyticsEngagementRateCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "engagement-rate",
		Short: "Get total engagement rate across accessible accounts",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			var resp map[string]any
			if err := analyticsRequest(cmd.Context(), "/insights/accounts/engagement/rate", nil, &resp); err != nil {
				return err
			}
			return renderAnalyticsStats(resp)
		},
	}
}

func newAnalyticsInboxUnreadCountCmd() *cobra.Command {
	return &cobra.Command{
		Use:    "inbox-unread-count",
		Short:  "Get open social inbox message count",
		Hidden: true, // The documented endpoint still returns 404 in production.
		Args:   cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			var resp map[string]any
			if err := analyticsRequest(cmd.Context(), "/insights/inbox/unread-count", nil, &resp); err != nil {
				return err
			}
			return renderAnalyticsStats(resp)
		},
	}
}

func newAnalyticsAutomationLogsCmd() *cobra.Command {
	var limit int
	cmd := &cobra.Command{
		Use:   "automation-logs",
		Short: "Get recent automation activity logs",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			if limit < 0 {
				return fmt.Errorf("--limit cannot be negative")
			}
			q := url.Values{}
			if limit > 0 {
				q.Set("limit", fmt.Sprintf("%d", limit))
			}
			var resp map[string]any
			if err := analyticsRequest(cmd.Context(), "/insights/automations/logs", q, &resp); err != nil {
				return err
			}
			return output.JSON(resp)
		},
	}
	cmd.Flags().IntVar(&limit, "limit", 5, "maximum number of automation logs")
	return cmd
}

func newAnalyticsFollowersGrowthCmd() *cobra.Command {
	var start, end string
	cmd := &cobra.Command{
		Use:   "followers-growth",
		Short: "Get daily followers growth across accounts",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := validateDateRange(start, end); err != nil {
				return err
			}
			var resp map[string]any
			if err := analyticsRequest(cmd.Context(), "/insights/accounts/followers/growth", dateRangeQuery(start, end), &resp); err != nil {
				return err
			}
			return renderAnalyticsFollowersGrowth(resp)
		},
	}
	addDateRangeFlags(cmd, &start, &end)
	return cmd
}

func newAnalyticsTeamMetricsCmd() *cobra.Command {
	var start, end, metrics string
	var accounts []int
	var team int
	cmd := &cobra.Command{
		Use:   "team-metrics",
		Short: "Get team member performance metrics",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			if strings.TrimSpace(metrics) == "" {
				return fmt.Errorf("--metrics is required")
			}
			if team <= 0 {
				return fmt.Errorf("--team is required")
			}
			if err := validateDateRange(start, end); err != nil {
				return err
			}
			if err := validateAnalyticsFilters(accounts, team); err != nil {
				return err
			}
			q := dateRangeQuery(start, end)
			q.Set("metrics", metrics)
			q = addOptionalInt(q, "team", team)
			q = addIntSlice(q, "accounts", accounts)
			var resp map[string]any
			if err := analyticsRequest(cmd.Context(), "/insights/teams/metrics", q, &resp); err != nil {
				return err
			}
			return renderAnalyticsTeamMetrics(resp)
		},
	}
	addDateRangeFlags(cmd, &start, &end)
	cmd.Flags().StringVar(&metrics, "metrics", "", "comma-separated metric names")
	cmd.Flags().IntSliceVar(&accounts, "accounts", nil, "account IDs")
	cmd.Flags().IntVar(&team, "team", 0, "team ID")
	return cmd
}

func newAnalyticsTeamActivityCmd() *cobra.Command {
	var limit int
	cmd := &cobra.Command{
		Use:   "team-activity",
		Short: "Get recent team activity logs",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			if limit < 0 {
				return fmt.Errorf("--limit cannot be negative")
			}
			q := url.Values{}
			if limit > 0 {
				q.Set("limit", fmt.Sprintf("%d", limit))
			}
			var resp map[string]any
			if err := analyticsRequest(cmd.Context(), "/insights/teams/activity", q, &resp); err != nil {
				return err
			}
			return renderAnalyticsTeamActivity(resp)
		},
	}
	cmd.Flags().IntVar(&limit, "limit", 5, "maximum number of activity logs")
	return cmd
}

func newAnalyticsStatsCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "stats",
		Short: "Get high-level user stats",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			var resp map[string]any
			if err := analyticsRequest(cmd.Context(), "/insights/stats", nil, &resp); err != nil {
				return err
			}
			return renderAnalyticsStats(resp)
		},
	}
}

func analyticsRequest(ctx context.Context, endpoint string, query url.Values, out any) error {
	cli, err := apiClient()
	if err != nil {
		return err
	}
	return cli.Request(ctx, "GET", endpoint, query, nil, out)
}

func addDateRangeFlags(cmd *cobra.Command, start, end *string) {
	cmd.Flags().StringVar(start, "start", "", "start date in YYYY-MM-DD")
	cmd.Flags().StringVar(end, "end", "", "end date in YYYY-MM-DD")
}

func dateRangeQuery(start, end string) url.Values {
	q := url.Values{}
	q.Set("start", start)
	q.Set("end", end)
	return q
}

func addOptionalString(q url.Values, key, value string) url.Values {
	if strings.TrimSpace(value) != "" {
		q.Set(key, value)
	}
	return q
}

func addOptionalInt(q url.Values, key string, value int) url.Values {
	if value > 0 {
		q.Set(key, fmt.Sprintf("%d", value))
	}
	return q
}

func addIntSlice(q url.Values, key string, values []int) url.Values {
	for _, value := range values {
		q.Add(key+"[]", fmt.Sprintf("%d", value))
	}
	return q
}

func validateAnalyticsFilters(accounts []int, team int) error {
	if err := validatePositiveInts(accounts, "--accounts"); err != nil {
		return err
	}
	if team < 0 {
		return fmt.Errorf("--team cannot be negative")
	}
	return nil
}

func renderAnalyticsPostsCount(resp map[string]any) error {
	if jsonOutput {
		return output.JSON(resp)
	}
	items := output.SliceFromMap(resp, "data", "counts", "results")
	if len(items) == 0 {
		return output.JSON(resp)
	}
	rows := make([][]string, 0, len(items))
	for _, item := range items {
		rows = append(rows, []string{
			output.FirstNonEmpty(output.StringFromMap(item, "date", "day", "label"), "-"),
			output.FirstNonEmpty(output.StringFromMap(item, "count", "posts_count", "total"), "0"),
		})
	}
	output.Table([]string{"Date", "Posts"}, rows)
	return nil
}

func renderAnalyticsEngagementTrend(resp map[string]any) error {
	if jsonOutput {
		return output.JSON(resp)
	}
	items := output.SliceFromMap(resp, "data", "trend", "results")
	if len(items) == 0 {
		if data := output.MapFromMap(resp, "data"); data != nil {
			items = output.SliceFromMap(data, "trend", "series", "results")
		}
	}
	if len(items) == 0 {
		return output.JSON(resp)
	}
	rows := make([][]string, 0, len(items))
	for _, item := range items {
		rows = append(rows, []string{
			output.FirstNonEmpty(output.StringFromMap(item, "date", "day", "label"), "-"),
			output.FirstNonEmpty(output.StringFromMap(item, "engagement", "value", "count", "total"), "0"),
		})
	}
	output.Table([]string{"Date", "Engagement"}, rows)
	return nil
}

func renderAnalyticsFollowersGrowth(resp map[string]any) error {
	if jsonOutput {
		return output.JSON(resp)
	}
	items := output.SliceFromMap(resp, "data", "growth", "results")
	if len(items) == 0 {
		if data := output.MapFromMap(resp, "data"); data != nil {
			items = output.SliceFromMap(data, "growth", "series", "results")
		}
	}
	if len(items) == 0 {
		return output.JSON(resp)
	}
	rows := make([][]string, 0, len(items))
	for _, item := range items {
		rows = append(rows, []string{
			output.FirstNonEmpty(output.StringFromMap(item, "date", "day", "label"), "-"),
			output.FirstNonEmpty(output.StringFromMap(item, "total_followers", "followers", "growth", "value", "count"), "0"),
		})
	}
	output.Table([]string{"Date", "Followers"}, rows)
	return nil
}

func renderAnalyticsTeamMetrics(resp map[string]any) error {
	if jsonOutput {
		return output.JSON(resp)
	}
	items := output.SliceFromMap(resp, "data", "metrics", "results")
	if len(items) == 0 {
		if data := output.MapFromMap(resp, "data"); data != nil {
			items = output.SliceFromMap(data, "metrics", "rows", "results")
		}
	}
	if len(items) == 0 {
		return output.JSON(resp)
	}
	if _, productionShape := items[0]["member_id"]; productionShape {
		rows := make([][]string, 0, len(items))
		for _, item := range items {
			rows = append(rows, []string{
				output.FirstNonEmpty(output.StringFromMap(item, "member_name", "name"), "-"),
				output.FirstNonEmpty(output.StringFromMap(item, "total_engagements"), "0"),
				fmt.Sprintf("%d", sumMetricCounts(item["scheduled_posts"])),
				fmt.Sprintf("%d", sumMetricCounts(item["published_posts"])),
				fmt.Sprintf("%d", sumMetricCounts(item["rejected_posts"])),
			})
		}
		output.Table([]string{"Member", "Engagements", "Scheduled", "Published", "Rejected"}, rows)
		return nil
	}
	rows := make([][]string, 0, len(items))
	for _, item := range items {
		member := output.FirstNonEmpty(
			output.StringFromMap(item, "name", "member_name", "user", "label"),
			output.StringFromNestedMap(item, []string{"member", "name"}, []string{"user", "name"}, []string{"actor", "name"}),
		)
		metric := output.FirstNonEmpty(
			output.StringFromMap(item, "metric", "metrics", "value", "count"),
			output.StringFromNestedMap(item, []string{"metric", "name"}, []string{"metrics", "name"}),
		)
		rows = append(rows, []string{
			output.FirstNonEmpty(member, "-"),
			output.FirstNonEmpty(metric, "-"),
		})
	}
	output.Table([]string{"Member", "Metric"}, rows)
	return nil
}

func sumMetricCounts(value any) int {
	items, ok := value.([]any)
	if !ok {
		return 0
	}
	total := 0
	for _, item := range items {
		if entry, ok := item.(map[string]any); ok {
			total += output.IntFromMap(entry, "count", "value")
		}
	}
	return total
}

func renderAnalyticsTeamActivity(resp map[string]any) error {
	if jsonOutput {
		return output.JSON(resp)
	}
	items := output.SliceFromMap(resp, "data", "activity", "logs", "results")
	if len(items) == 0 {
		if data := output.MapFromMap(resp, "data"); data != nil {
			items = output.SliceFromMap(data, "activity", "logs", "results")
		}
	}
	if len(items) == 0 {
		return output.JSON(resp)
	}
	lines := make([]string, 0, len(items))
	for _, item := range items {
		when := output.FirstNonEmpty(output.StringFromMap(item, "created_at", "date", "time"), "unknown time")
		who := output.FirstNonEmpty(
			output.StringFromMap(item, "user", "name", "actor"),
			output.StringFromNestedMap(item, []string{"user", "name"}, []string{"actor", "name"}, []string{"member", "name"}),
			"Someone",
		)
		what := output.FirstNonEmpty(
			output.StringFromMap(item, "action", "description", "message"),
			output.StringFromNestedMap(item, []string{"action", "name"}, []string{"event", "name"}),
			"did something",
		)
		lines = append(lines, when+" - "+who+" - "+what)
	}
	output.PrintSection("Recent team activity")
	output.PrintList(lines)
	return nil
}

func renderAnalyticsStats(resp map[string]any) error {
	if jsonOutput {
		return output.JSON(resp)
	}
	data := resp
	if nested, ok := resp["data"].(map[string]any); ok {
		data = nested
	}

	values := map[string]string{}
	for _, key := range output.Keys(data) {
		value := data[key]
		if text := output.String(value); text != "" && !strings.HasPrefix(text, "map[") && !strings.HasPrefix(text, "[") {
			values[key] = text
			continue
		}
		if nested, ok := value.(map[string]any); ok {
			if text := output.FirstNonEmpty(
				output.StringFromMap(nested, "value", "count", "total", "label", "name", "title"),
				output.StringFromNestedMap(nested,
					[]string{"current", "value"},
					[]string{"current", "count"},
					[]string{"totals", "value"},
					[]string{"totals", "count"},
					[]string{"series", "label"},
				),
			); text != "" {
				values[key] = text
			}
		}
	}
	if len(values) == 0 {
		return output.JSON(resp)
	}
	output.KeyValue("Stats", values)
	return nil
}
