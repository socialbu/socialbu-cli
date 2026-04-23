package cmd

import (
	"context"
	"fmt"
	"strconv"

	"github.com/spf13/cobra"
	"github.com/usamaejaz/socialbu-cli/internal/output"
)

type team struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

type teamsResponse struct {
	Data []team `json:"data"`
}

func newTeamsCmd() *cobra.Command {
	cmd := &cobra.Command{Use: "team", Aliases: []string{"teams"}, Short: "Manage teams"}
	cmd.AddCommand(&cobra.Command{
		Use:   "list",
		Short: "List teams",
		RunE: func(cmd *cobra.Command, args []string) error {
			cli, err := apiClient()
			if err != nil {
				return err
			}
			var resp teamsResponse
			if err := cli.Request(context.Background(), "GET", "/teams", nil, nil, &resp); err != nil {
				return err
			}
			if jsonOutput {
				return output.JSON(resp)
			}
			rows := make([][]string, 0, len(resp.Data))
			for _, t := range resp.Data {
				rows = append(rows, []string{strconv.Itoa(t.ID), t.Name})
			}
			output.Table([]string{"ID", "Name"}, rows)
			return nil
		},
	})
	cmd.AddCommand(&cobra.Command{
		Use:   "create <name>",
		Short: "Create a team",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			cli, err := apiClient()
			if err != nil {
				return err
			}
			payload := map[string]any{"name": args[0]}
			var resp map[string]any
			if err := cli.Request(context.Background(), "POST", "/teams", nil, payload, &resp); err != nil {
				return err
			}
			return output.JSON(resp)
		},
	})
	cmd.AddCommand(&cobra.Command{
		Use:   "delete <id>",
		Short: "Delete a team",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			cli, err := apiClient()
			if err != nil {
				return err
			}
			if err := cli.Request(context.Background(), "DELETE", fmt.Sprintf("/teams/%s", args[0]), nil, nil, nil); err != nil {
				return err
			}
			fmt.Fprintln(cmd.OutOrStdout(), "Team deleted")
			return nil
		},
	})
	return cmd
}
