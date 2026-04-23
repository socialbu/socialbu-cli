package cmd

import (
	"context"
	"fmt"
	"strconv"

	"github.com/spf13/cobra"
	"github.com/usamaejaz/socialbu-cli/internal/output"
)

type account struct {
	ID       int    `json:"id"`
	Name     string `json:"name"`
	Username string `json:"username"`
	Type     string `json:"type"`
	Status   string `json:"status"`
}

func newAccountsCmd() *cobra.Command {
	cmd := &cobra.Command{Use: "account", Aliases: []string{"accounts"}, Short: "Manage social accounts"}
	cmd.AddCommand(&cobra.Command{
		Use:   "list",
		Short: "List accounts",
		RunE: func(cmd *cobra.Command, args []string) error {
			cli, err := apiClient()
			if err != nil {
				return err
			}
			var raw any
			if err := cli.Request(context.Background(), "GET", "/accounts", nil, nil, &raw); err != nil {
				return err
			}
			accounts := make([]account, 0)
			switch v := raw.(type) {
			case []any:
				for _, item := range v {
					m, ok := item.(map[string]any)
					if !ok {
						continue
					}
					accounts = append(accounts, account{
						ID:       output.IntFromMap(m, "id"),
						Name:     output.StringFromMap(m, "name"),
						Username: output.StringFromMap(m, "username", "handle"),
						Type:     output.StringFromMap(m, "type", "account_type"),
						Status:   output.StringFromMap(m, "status"),
					})
				}
			case map[string]any:
				items := output.SliceFromMap(v, "data", "items", "accounts")
				for _, m := range items {
					accounts = append(accounts, account{
						ID:       output.IntFromMap(m, "id"),
						Name:     output.StringFromMap(m, "name"),
						Username: output.StringFromMap(m, "username", "handle"),
						Type:     output.StringFromMap(m, "type", "account_type"),
						Status:   output.StringFromMap(m, "status"),
					})
				}
			}
			if jsonOutput {
				return output.JSON(map[string]any{"data": accounts})
			}
			rows := make([][]string, 0, len(accounts))
			for _, a := range accounts {
				rows = append(rows, []string{strconv.Itoa(a.ID), a.Name, a.Username, a.Type, a.Status})
			}
			output.Table([]string{"ID", "Name", "Username", "Type", "Status"}, rows)
			return nil
		},
	})
	cmd.AddCommand(&cobra.Command{
		Use:   "get <id>",
		Short: "Get account details",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			cli, err := apiClient()
			if err != nil {
				return err
			}
			var resp map[string]any
			if err := cli.Request(context.Background(), "GET", fmt.Sprintf("/accounts/%s", args[0]), nil, nil, &resp); err != nil {
				return err
			}
			return output.JSON(resp)
		},
	})
	return cmd
}
