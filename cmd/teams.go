package cmd

import (
	"fmt"
	"net/url"
	"strconv"
	"strings"

	"github.com/socialbu/socialbu-cli/internal/output"
	"github.com/spf13/cobra"
)

type team struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

func newTeamsCmd() *cobra.Command {
	cmd := &cobra.Command{Use: "team", Aliases: []string{"teams"}, Short: "Manage teams"}
	var teamType string
	listCmd := &cobra.Command{
		Use:   "list",
		Short: "List teams",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := validateChoice(teamType, "--type", "created", "joined"); err != nil {
				return err
			}
			cli, err := apiClient()
			if err != nil {
				return err
			}
			query := url.Values{"type": {teamType}}
			var raw any
			if err := cli.Request(cmd.Context(), "GET", "/teams", query, nil, &raw); err != nil {
				return err
			}
			if jsonOutput {
				return output.JSON(raw)
			}
			teams := teamsFromAny(raw)
			rows := make([][]string, 0, len(teams))
			for _, t := range teams {
				rows = append(rows, []string{strconv.Itoa(t.ID), t.Name})
			}
			output.Table([]string{"ID", "Name"}, rows)
			return nil
		},
	}
	listCmd.Flags().StringVar(&teamType, "type", "created", "filter by team type: created or joined")
	cmd.AddCommand(listCmd)
	var accounts []int
	var requiresApproval bool
	createCmd := &cobra.Command{
		Use:   "create <name>",
		Short: "Create a team",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			if strings.TrimSpace(args[0]) == "" {
				return fmt.Errorf("team name cannot be empty")
			}
			if len(accounts) == 0 {
				return fmt.Errorf("at least one --accounts value is required")
			}
			if err := validatePositiveInts(accounts, "--accounts"); err != nil {
				return err
			}
			cli, err := apiClient()
			if err != nil {
				return err
			}
			accountItems := make([]map[string]int, 0, len(accounts))
			for _, id := range accounts {
				accountItems = append(accountItems, map[string]int{"id": id})
			}
			payload := map[string]any{
				"name":                      args[0],
				"accounts":                  accountItems,
				"requires_content_approval": requiresApproval,
			}
			var resp map[string]any
			if err := cli.Request(cmd.Context(), "POST", "/teams", nil, payload, &resp); err != nil {
				return err
			}
			return output.JSON(resp)
		},
	}
	createCmd.Flags().IntSliceVar(&accounts, "accounts", nil, "account IDs to include in the team")
	createCmd.Flags().BoolVar(&requiresApproval, "require-approval", false, "require approval before team content is published")
	cmd.AddCommand(createCmd)
	cmd.AddCommand(&cobra.Command{
		Use:   "delete <id>",
		Short: "Delete a team",
		Args:  exactPositiveIDArg("team ID"),
		RunE: func(cmd *cobra.Command, args []string) error {
			cli, err := apiClient()
			if err != nil {
				return err
			}
			if err := cli.Request(cmd.Context(), "DELETE", fmt.Sprintf("/teams/%s", args[0]), nil, nil, nil); err != nil {
				return err
			}
			fmt.Fprintln(cmd.OutOrStdout(), "Team deleted")
			return nil
		},
	})
	return cmd
}

func teamsFromAny(raw any) []team {
	var items []map[string]any
	switch value := raw.(type) {
	case []any:
		items = make([]map[string]any, 0, len(value))
		for _, item := range value {
			if entry, ok := item.(map[string]any); ok {
				items = append(items, entry)
			}
		}
	case map[string]any:
		items = output.SliceFromMap(value, "items", "data", "teams")
	}

	teams := make([]team, 0, len(items))
	for _, item := range items {
		if nested := output.MapFromMap(item, "team"); nested != nil {
			item = nested
		}
		teams = append(teams, team{ID: output.IntFromMap(item, "id"), Name: output.StringFromMap(item, "name")})
	}
	return teams
}
