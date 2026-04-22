package cmd

import (
	"context"
	"fmt"
	"net/url"
	"strconv"
	"strings"

	"github.com/spf13/cobra"
	"github.com/usamaejaz/socialbu-cli/internal/output"
)

type post struct {
	ID        int    `json:"id"`
	Content   string `json:"content"`
	Status    string `json:"status"`
	PublishAt string `json:"publish_at"`
}

type postsResponse struct {
	Data []post `json:"data"`
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
		RunE: func(cmd *cobra.Command, args []string) error {
			cli, err := apiClient()
			if err != nil {
				return err
			}
			q := url.Values{}
			if postType != "" {
				q.Set("type", postType)
			}
			var resp postsResponse
			if err := cli.Request(context.Background(), "GET", "/posts", q, nil, &resp); err != nil {
				return err
			}
			if jsonOutput {
				return output.JSON(resp)
			}
			rows := make([][]string, 0, len(resp.Data))
			for _, p := range resp.Data {
				rows = append(rows, []string{strconv.Itoa(p.ID), p.Status, p.PublishAt, truncate(p.Content, 60)})
			}
			output.Table([]string{"ID", "Status", "Publish At", "Content"}, rows)
			return nil
		},
	}
	cmd.Flags().StringVar(&postType, "type", "", "filter by type (scheduled, published, failed, drafts)")
	return cmd
}

func newPostsGetCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "get <id>",
		Short: "Get post details",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			cli, err := apiClient()
			if err != nil {
				return err
			}
			var resp map[string]any
			if err := cli.Request(context.Background(), "GET", fmt.Sprintf("/posts/%s", args[0]), nil, nil, &resp); err != nil {
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
		RunE: func(cmd *cobra.Command, args []string) error {
			cli, err := apiClient()
			if err != nil {
				return err
			}
			if len(accounts) == 0 {
				return fmt.Errorf("at least one --accounts is required")
			}
			if publishAt == "" {
				return fmt.Errorf("--publish-at is required")
			}
			payload := map[string]any{
				"accounts":   accounts,
				"publish_at": publishAt,
				"content":    content,
				"draft":      draft,
			}
			var resp map[string]any
			if err := cli.Request(context.Background(), "POST", "/posts", nil, payload, &resp); err != nil {
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
	if len(s) <= n {
		return s
	}
	return s[:n-1] + "…"
}
