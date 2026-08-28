package cmd

import (
	"github.com/socialbu/socialbu-cli/internal/output"
	"github.com/spf13/cobra"
)

func newWhoamiCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "whoami",
		Short: "Show the authenticated user",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			cli, err := apiClient()
			if err != nil {
				return err
			}
			var resp map[string]any
			if err := cli.Request(cmd.Context(), "GET", "/user", nil, nil, &resp); err != nil {
				return err
			}
			if jsonOutput {
				return output.JSON(resp)
			}
			output.KeyValue("Authenticated user", map[string]string{
				"company":  output.StringFromMap(resp, "company"),
				"email":    output.StringFromMap(resp, "email"),
				"id":       output.StringFromMap(resp, "id"),
				"name":     output.StringFromMap(resp, "name"),
				"verified": output.StringFromMap(resp, "verified"),
			})
			return nil
		},
	}
}
