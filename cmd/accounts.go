package cmd

import (
	"fmt"
	"net/url"
	"strconv"

	"github.com/socialbu/socialbu-cli/internal/output"
	"github.com/spf13/cobra"
)

type account struct {
	ID       int    `json:"id"`
	Name     string `json:"name"`
	Platform string `json:"platform"`
	Type     string `json:"type"`
	Status   string `json:"status"`
}

func newAccountsCmd() *cobra.Command {
	cmd := &cobra.Command{Use: "account", Aliases: []string{"accounts"}, Short: "Manage social accounts"}
	var accountType string
	listCmd := &cobra.Command{
		Use:   "list",
		Short: "List accounts",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := validateChoice(accountType, "--type", "all", "user", "shared"); err != nil {
				return err
			}
			cli, err := apiClient()
			if err != nil {
				return err
			}
			q := url.Values{}
			if accountType != "" {
				q.Add("type", accountType)
			}
			var raw any
			if err := cli.Request(cmd.Context(), "GET", "/accounts", q, nil, &raw); err != nil {
				return err
			}
			if jsonOutput {
				return output.JSON(raw)
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
						Platform: output.StringFromMap(m, "provider", "platform", "_type"),
						Type:     output.StringFromMap(m, "type", "account_type"),
						Status:   accountStatus(m),
					})
				}
			case map[string]any:
				items := output.SliceFromMap(v, "data", "items", "accounts")
				for _, m := range items {
					accounts = append(accounts, account{
						ID:       output.IntFromMap(m, "id"),
						Name:     output.StringFromMap(m, "name"),
						Platform: output.StringFromMap(m, "provider", "platform", "_type"),
						Type:     output.StringFromMap(m, "type", "account_type"),
						Status:   accountStatus(m),
					})
				}
			}
			rows := make([][]string, 0, len(accounts))
			for _, a := range accounts {
				rows = append(rows, []string{strconv.Itoa(a.ID), a.Name, a.Type, a.Platform, a.Status})
			}
			output.Table([]string{"ID", "Name", "Type", "Platform", "Status"}, rows)
			return nil
		},
	}
	listCmd.Flags().StringVar(&accountType, "type", "all", "filter by account type: all, user, or shared")
	cmd.AddCommand(listCmd)
	cmd.AddCommand(&cobra.Command{
		Use:   "get <id>",
		Short: "Get account details",
		Args:  exactPositiveIDArg("account ID"),
		RunE: func(cmd *cobra.Command, args []string) error {
			cli, err := apiClient()
			if err != nil {
				return err
			}
			var resp map[string]any
			if err := cli.Request(cmd.Context(), "GET", fmt.Sprintf("/accounts/%s", args[0]), nil, nil, &resp); err != nil {
				return err
			}
			return output.JSON(resp)
		},
	})
	return cmd
}

func accountStatus(item map[string]any) string {
	if status := output.StringFromMap(item, "status"); status != "" {
		return status
	}
	if active, ok := item["active"].(bool); ok {
		if active {
			return "active"
		}
		return "inactive"
	}
	return ""
}
