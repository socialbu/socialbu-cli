package cmd

import (
	"fmt"
	"strings"

	"github.com/socialbu/socialbu-cli/internal/output"
	"github.com/spf13/cobra"
)

func newAICmd() *cobra.Command {
	cmd := &cobra.Command{Use: "ai", Short: "Generate or autocomplete content with AI"}
	cmd.AddCommand(newAIGenerateCmd())
	cmd.AddCommand(newAIGenerateFromPostCmd())
	cmd.AddCommand(newAIAutocompleteCmd())
	return cmd
}

func newAIGenerateCmd() *cobra.Command {
	var contentType, topic string
	var accountID, teamID int
	cmd := &cobra.Command{
		Use:   "generate",
		Short: "Generate new AI content from a topic",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			if strings.TrimSpace(topic) == "" || accountID <= 0 {
				return fmt.Errorf("--type, --topic, and --account are required")
			}
			if err := validateChoice(contentType, "--type", "tweet", "linkedin_post", "instagram_caption", "generic"); err != nil {
				return err
			}
			cli, err := apiClient()
			if err != nil {
				return err
			}
			payload := map[string]any{"type": contentType, "topic": topic, "account_id": accountID}
			if teamID > 0 {
				payload["team_id"] = teamID
			}
			var resp map[string]any
			if err := cli.Request(cmd.Context(), "POST", "/generated_content", nil, payload, &resp); err != nil {
				return err
			}
			return renderAIResponse(resp)
		},
	}
	cmd.Flags().StringVar(&contentType, "type", "generic", "content type: tweet, linkedin_post, instagram_caption, generic")
	cmd.Flags().StringVar(&topic, "topic", "", "topic or prompt seed")
	cmd.Flags().IntVar(&accountID, "account", 0, "account ID")
	cmd.Flags().IntVar(&teamID, "team", 0, "team ID")
	return cmd
}

func newAIGenerateFromPostCmd() *cobra.Command {
	var postID int
	cmd := &cobra.Command{
		Use:   "from-post",
		Short: "Generate AI content from an existing post",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			if postID <= 0 {
				return fmt.Errorf("--post is required")
			}
			cli, err := apiClient()
			if err != nil {
				return err
			}
			payload := map[string]any{"post_id": postID}
			var resp map[string]any
			if err := cli.Request(cmd.Context(), "POST", "/generated_content/generate_by_post", nil, payload, &resp); err != nil {
				return err
			}
			return renderAIResponse(resp)
		},
	}
	cmd.Flags().IntVar(&postID, "post", 0, "post ID")
	return cmd
}

func newAIAutocompleteCmd() *cobra.Command {
	var content string
	var accountID, teamID int
	cmd := &cobra.Command{
		Use:   "autocomplete",
		Short: "Autocomplete a partial caption or post body",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			if accountID <= 0 {
				return fmt.Errorf("--account is required")
			}
			cli, err := apiClient()
			if err != nil {
				return err
			}
			payload := map[string]any{"account_id": accountID}
			if content != "" {
				payload["content"] = content
			}
			if teamID > 0 {
				payload["team_id"] = teamID
			}
			var resp map[string]any
			if err := cli.Request(cmd.Context(), "POST", "/generated_content/autocomplete_post", nil, payload, &resp); err != nil {
				return err
			}
			return renderAIResponse(resp)
		},
	}
	cmd.Flags().StringVar(&content, "content", "", "partial content to extend")
	cmd.Flags().IntVar(&accountID, "account", 0, "account ID")
	cmd.Flags().IntVar(&teamID, "team", 0, "team ID")
	return cmd
}

func renderAIResponse(resp map[string]any) error {
	if jsonOutput {
		return output.JSON(resp)
	}
	data := resp
	if nested, ok := resp["data"].(map[string]any); ok {
		data = nested
	}
	content := output.FirstNonEmpty(
		output.StringFromMap(data, "content", "text", "caption", "generated_content", "suggestion", "result", "completion"),
		output.StringFromMap(resp, "content", "text", "caption", "generated_content", "suggestion", "result", "completion"),
		output.StringFromNestedMap(
			data,
			[]string{"data", "content"},
			[]string{"data", "text"},
			[]string{"result", "content"},
			[]string{"result", "text"},
			[]string{"result", "suggestion"},
			[]string{"completion", "text"},
			[]string{"suggestion", "text"},
		),
		output.StringFromNestedMap(
			resp,
			[]string{"data", "content"},
			[]string{"data", "text"},
			[]string{"result", "content"},
			[]string{"result", "text"},
			[]string{"result", "suggestion"},
			[]string{"completion", "text"},
			[]string{"suggestion", "text"},
		),
	)
	if content == "" {
		return output.JSON(resp)
	}
	output.PrintSection("AI output")
	fmt.Println(content)
	meta := map[string]string{}
	for _, key := range []string{"type", "tone", "language", "model"} {
		if value := output.FirstNonEmpty(
			output.StringFromMap(data, key),
			output.StringFromNestedMap(data, []string{"meta", key}, []string{"result", key}),
			output.StringFromMap(resp, key),
			output.StringFromNestedMap(resp, []string{"meta", key}, []string{"result", key}),
		); value != "" {
			meta[key] = value
		}
	}
	if len(meta) > 0 {
		fmt.Println()
		output.KeyValue("Details", meta)
	}
	return nil
}
