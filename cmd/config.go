package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/usamaejaz/socialbu-cli/internal/config"
	"github.com/usamaejaz/socialbu-cli/internal/output"
)

func newConfigCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "config",
		Short: "Manage CLI config",
		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
			return ensureConfig()
		},
	}

	cmd.AddCommand(&cobra.Command{
		Use:   "set-key <api-key>",
		Short: "Store API key",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := config.SetAPIKey(args[0]); err != nil {
				return err
			}
			cfg = config.Current()
			fmt.Fprintln(cmd.OutOrStdout(), "API key saved")
			return nil
		},
	})

	cmd.AddCommand(&cobra.Command{
		Use:   "set-base-url <url>",
		Short: "Store API base URL",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := config.SetBaseURL(args[0]); err != nil {
				return err
			}
			cfg = config.Current()
			fmt.Fprintln(cmd.OutOrStdout(), "Base URL saved")
			return nil
		},
	})

	cmd.AddCommand(&cobra.Command{
		Use:   "show",
		Short: "Show effective config",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			values := map[string]string{
				"api_key_set": fmt.Sprintf("%t", cfg.APIKey != ""),
				"base_url":    cfg.BaseURL,
			}
			if jsonOutput {
				return output.JSON(values)
			}
			output.KeyValue("config", values)
			return nil
		},
	})

	cmd.AddCommand(&cobra.Command{
		Use:   "reset",
		Short: "Reset config",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := config.Reset(); err != nil {
				return err
			}
			cfg = config.Current()
			fmt.Fprintln(cmd.OutOrStdout(), "Config reset")
			return nil
		},
	})

	return cmd
}
