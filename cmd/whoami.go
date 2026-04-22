package cmd

import (
	"context"
	"fmt"

	"github.com/spf13/cobra"
	"github.com/usamaejaz/socialbu-cli/internal/output"
)

type statsResponse struct {
	Data map[string]any `json:"data"`
}

func newWhoamiCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "whoami",
		Short: "Check API identity by fetching user stats",
		RunE: func(cmd *cobra.Command, args []string) error {
			cli, err := apiClient()
			if err != nil {
				return err
			}
			var resp statsResponse
			if err := cli.Request(context.Background(), "GET", "/insights/stats", nil, nil, &resp); err != nil {
				return err
			}
			if jsonOutput {
				return output.JSON(resp)
			}
			fmt.Fprintln(cmd.OutOrStdout(), "Authenticated. Use --json to inspect full stats payload.")
			return nil
		},
	}
}
