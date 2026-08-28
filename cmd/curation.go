package cmd

import (
	"context"
	"fmt"
	"net/url"

	"github.com/spf13/cobra"
	"github.com/usamaejaz/socialbu-cli/internal/output"
)

func newCurationCmd() *cobra.Command {
	cmd := &cobra.Command{Use: "curation", Short: "Browse curated content and topics"}
	cmd.AddCommand(newCurationTopicsCmd())
	cmd.AddCommand(newCurationItemsCmd())
	cmd.AddCommand(newCurationGetItemCmd())
	return cmd
}

func newCurationTopicsCmd() *cobra.Command {
	var query string
	cmd := &cobra.Command{
		Use:   "topics",
		Short: "List curation topics",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			q := url.Values{}
			if query != "" {
				q.Set("q", query)
			}
			return curationGetTopics(cmd.Context(), q)
		},
	}
	cmd.Flags().StringVar(&query, "query", "", "topic search query")
	return cmd
}

func newCurationItemsCmd() *cobra.Command {
	var page, perPage, feedID int
	var search, from, to string
	var authors, topics []string
	cmd := &cobra.Command{
		Use:   "items",
		Short: "List curated content items",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			if page < 0 || perPage < 0 || feedID < 0 {
				return fmt.Errorf("--page, --per-page, and --feed-id cannot be negative")
			}
			if err := validateOptionalDateRange(from, to); err != nil {
				return err
			}
			q := url.Values{}
			if page > 0 {
				q.Set("page", fmt.Sprintf("%d", page))
			}
			if perPage > 0 {
				q.Set("per_page", fmt.Sprintf("%d", perPage))
			}
			if feedID > 0 {
				q.Set("feed_id", fmt.Sprintf("%d", feedID))
			}
			if search != "" {
				q.Set("search", search)
			}
			if from != "" {
				q.Set("from", from)
			}
			if to != "" {
				q.Set("to", to)
			}
			for _, author := range authors {
				q.Add("authors[]", author)
			}
			for _, topic := range topics {
				q.Add("topics[]", topic)
			}
			return curationGet(cmd.Context(), "/curation/items", q)
		},
	}
	cmd.Flags().IntVar(&page, "page", 0, "page number")
	cmd.Flags().IntVar(&perPage, "per-page", 0, "results per page")
	cmd.Flags().IntVar(&feedID, "feed-id", 0, "feed ID")
	cmd.Flags().StringVar(&search, "search", "", "search query")
	cmd.Flags().StringVar(&from, "from", "", "start date in YYYY-MM-DD")
	cmd.Flags().StringVar(&to, "to", "", "end date in YYYY-MM-DD")
	cmd.Flags().StringSliceVar(&authors, "authors", nil, "author filters")
	cmd.Flags().StringSliceVar(&topics, "topics", nil, "topic filters")
	return cmd
}

func newCurationGetItemCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "get <id>",
		Short: "Get one curation item",
		Args:  exactPositiveIDArg("curation item ID"),
		RunE: func(cmd *cobra.Command, args []string) error {
			return curationGet(cmd.Context(), fmt.Sprintf("/curation/items/%s", args[0]), nil)
		},
	}
}

func curationGet(ctx context.Context, endpoint string, query url.Values) error {
	cli, err := apiClient()
	if err != nil {
		return err
	}
	var resp map[string]any
	if err := cli.Request(ctx, "GET", endpoint, query, nil, &resp); err != nil {
		return err
	}
	return renderCuration(resp)
}

func curationGetTopics(ctx context.Context, query url.Values) error {
	cli, err := apiClient()
	if err != nil {
		return err
	}
	var resp any
	if err := cli.Request(ctx, "GET", "/curation/topics", query, nil, &resp); err != nil {
		return err
	}
	return renderCurationAny(resp)
}

func renderCurationAny(resp any) error {
	if m, ok := resp.(map[string]any); ok {
		return renderCuration(m)
	}
	if items, ok := resp.([]any); ok {
		wrapped := map[string]any{"topics": items}
		return renderCuration(wrapped)
	}
	return output.JSON(resp)
}

func renderCuration(resp map[string]any) error {
	if jsonOutput {
		return output.JSON(resp)
	}
	items := output.SliceFromMap(resp, "data", "items", "topics", "results")
	if len(items) > 1 {
		rows := make([][]string, 0, len(items))
		for _, item := range items {
			rows = append(rows, []string{
				output.FirstNonEmpty(output.StringFromMap(item, "id", "slug"), "-"),
				output.FirstNonEmpty(output.StringFromMap(item, "title", "name"), "-"),
				output.FirstNonEmpty(curationSource(item), "-"),
				output.FirstNonEmpty(output.StringFromMap(item, "published_at", "created_at", "date"), "-"),
			})
		}
		output.Table([]string{"ID", "Title", "Source", "Published"}, rows)
		return nil
	}
	if len(items) == 1 {
		return renderCurationDetail(items[0])
	}
	if data := output.MapFromMap(resp, "data"); data != nil {
		return renderCurationDetail(data)
	}
	for _, key := range []string{"items", "topics", "results"} {
		if _, exists := resp[key]; exists {
			fmt.Println("No curation items.")
			return nil
		}
	}
	return output.JSON(resp)
}

func renderCurationDetail(item map[string]any) error {
	values := map[string]string{}
	for _, key := range []string{"id", "slug", "title", "name", "author", "published_at", "created_at", "date", "url", "link", "summary", "description", "excerpt"} {
		if value := output.StringFromMap(item, key); value != "" {
			values[key] = value
		}
	}
	if source := curationSource(item); source != "" {
		values["source"] = source
	}
	if topics := output.JoinStrings(item, "topics", "topic_names", "tags", "keywords"); topics != "" {
		values["topics"] = topics
	}
	if len(values) == 0 {
		return output.JSON(item)
	}
	output.KeyValue("Curation item", values)
	return nil
}

func curationSource(item map[string]any) string {
	return output.FirstNonEmpty(
		output.StringFromMap(item, "author", "authors", "feed_name", "publication"),
		output.StringFromNestedMap(item, []string{"feed", "name"}, []string{"source", "name"}, []string{"author", "name"}),
	)
}
