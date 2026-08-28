package cmd

import (
	"fmt"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/socialbu/socialbu-cli/internal/output"
	"github.com/spf13/cobra"
)

type post struct {
	ID        int    `json:"id"`
	Content   string `json:"content"`
	Status    string `json:"status"`
	PublishAt string `json:"publish_at"`
}

func newPostsCmd() *cobra.Command {
	cmd := &cobra.Command{Use: "post", Aliases: []string{"posts"}, Short: "Manage posts"}
	cmd.AddCommand(newPostsListCmd())
	cmd.AddCommand(newPostsGetCmd())
	cmd.AddCommand(newPostsCreateCmd())
	return cmd
}

func newPostsListCmd() *cobra.Command {
	var postType string
	cmd := &cobra.Command{
		Use:   "list",
		Short: "List posts",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			if postType != "" {
				if err := validateChoice(postType, "--type", "scheduled", "awaiting_approval", "scheduled_or_awaiting_approval", "draft", "published"); err != nil {
					return err
				}
			}
			cli, err := apiClient()
			if err != nil {
				return err
			}
			q := url.Values{}
			if postType != "" {
				q.Set("type", postType)
			}
			var resp any
			if err := cli.Request(cmd.Context(), "GET", "/posts", q, nil, &resp); err != nil {
				return err
			}
			if jsonOutput {
				return output.JSON(resp)
			}
			posts := postsFromAny(resp)
			rows := make([][]string, 0, len(posts))
			for _, p := range posts {
				rows = append(rows, []string{strconv.Itoa(p.ID), p.Status, p.PublishAt, truncate(p.Content, 60)})
			}
			output.Table([]string{"ID", "Status", "Publish At", "Content"}, rows)
			return nil
		},
	}
	cmd.Flags().StringVar(&postType, "type", "", "filter by type: scheduled, awaiting_approval, scheduled_or_awaiting_approval, draft, or published")
	return cmd
}

func newPostsGetCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "get <id>",
		Short: "Get post details",
		Args:  exactPositiveIDArg("post ID"),
		RunE: func(cmd *cobra.Command, args []string) error {
			cli, err := apiClient()
			if err != nil {
				return err
			}
			var resp map[string]any
			if err := cli.Request(cmd.Context(), "GET", fmt.Sprintf("/posts/%s", args[0]), nil, nil, &resp); err != nil {
				return err
			}
			return output.JSON(resp)
		},
	}
}

func newPostsCreateCmd() *cobra.Command {
	var accounts []int
	var publishAt, content string
	var draft bool
	cmd := &cobra.Command{
		Use:   "create",
		Short: "Create a post",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			if len(accounts) == 0 {
				return fmt.Errorf("at least one --accounts is required")
			}
			if err := validatePositiveInts(accounts, "--accounts"); err != nil {
				return err
			}
			if err := validatePublishTime(publishAt, draft, time.Now()); err != nil {
				return err
			}
			cli, err := apiClient()
			if err != nil {
				return err
			}
			payload := map[string]any{
				"accounts":   accounts,
				"publish_at": publishAt,
				"content":    content,
				"draft":      draft,
			}
			var resp map[string]any
			if err := cli.Request(cmd.Context(), "POST", "/posts", nil, payload, &resp); err != nil {
				return err
			}
			return output.JSON(resp)
		},
	}
	cmd.Flags().IntSliceVar(&accounts, "accounts", nil, "account IDs")
	cmd.Flags().StringVar(&publishAt, "publish-at", "", "publish time in UTC, format YYYY-MM-DD HH:MM:SS")
	cmd.Flags().StringVar(&content, "content", "", "post content")
	cmd.Flags().BoolVar(&draft, "draft", false, "save as draft")
	return cmd
}

func truncate(s string, n int) string {
	s = strings.TrimSpace(strings.ReplaceAll(s, "\n", " "))
	runes := []rune(s)
	if len(runes) <= n {
		return s
	}
	if n <= 1 {
		return "…"
	}
	return string(runes[:n-1]) + "…"
}

func postsFromAny(resp any) []post {
	var items []map[string]any
	switch value := resp.(type) {
	case []any:
		items = make([]map[string]any, 0, len(value))
		for _, item := range value {
			if entry, ok := item.(map[string]any); ok {
				items = append(items, entry)
			}
		}
	case map[string]any:
		items = output.SliceFromMap(value, "items", "data", "posts")
	}

	posts := make([]post, 0, len(items))
	for _, item := range items {
		posts = append(posts, post{
			ID:        output.IntFromMap(item, "id"),
			Content:   output.StringFromMap(item, "content", "text"),
			Status:    output.StringFromMap(item, "status"),
			PublishAt: output.StringFromMap(item, "publish_at", "published_at"),
		})
	}
	return posts
}
