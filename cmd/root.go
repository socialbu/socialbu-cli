package cmd

import (
	"fmt"
	"runtime/debug"
	"strings"

	"github.com/socialbu/socialbu-cli/internal/client"
	"github.com/socialbu/socialbu-cli/internal/config"
	"github.com/spf13/cobra"
)

var (
	jsonOutput bool
	cfg        config.Config
	buildInfo  = BuildInfo{Version: "dev", Commit: "none", Date: "unknown"}
)

type BuildInfo struct {
	Version string
	Commit  string
	Date    string
}

func SetBuildInfo(version, commit, date string) {
	if version == "dev" {
		if info, ok := debug.ReadBuildInfo(); ok && info.Main.Version != "" && info.Main.Version != "(devel)" {
			version = strings.TrimPrefix(info.Main.Version, "v")
		}
	}
	buildInfo = BuildInfo{Version: version, Commit: commit, Date: date}
}

func Execute() error {
	return newRootCmd().Execute()
}

func newRootCmd() *cobra.Command {
	rootCmd := &cobra.Command{
		Use:           "socialbu",
		Short:         "SocialBu CLI",
		Long:          "A Go-based SocialBu CLI for accounts, posts, analytics, and more.",
		SilenceUsage:  true,
		SilenceErrors: true,
	}
	rootCmd.PersistentFlags().BoolVar(&jsonOutput, "json", false, "output raw JSON")
	versionCmd := &cobra.Command{
		Use:   "version",
		Short: "Show build metadata",
		Args:  cobra.NoArgs,
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Fprintf(cmd.OutOrStdout(), "socialbu %s\ncommit: %s\nbuilt: %s\n", buildInfo.Version, buildInfo.Commit, buildInfo.Date)
		},
	}
	helpCmd := &cobra.Command{
		Use:   "help [command]",
		Short: "Show help for a command",
		Args:  cobra.ArbitraryArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			if len(args) == 0 {
				return rootCmd.Help()
			}
			target, _, err := rootCmd.Find(args)
			if err != nil {
				return err
			}
			if target == nil || target == rootCmd && len(args) > 0 {
				return fmt.Errorf("unknown help topic %q", strings.Join(args, " "))
			}
			return target.Help()
		},
	}
	completionCmd := &cobra.Command{
		Use:   "completion [bash|zsh|fish|powershell]",
		Short: "Generate shell completion scripts",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			out := cmd.OutOrStdout()
			switch args[0] {
			case "bash":
				return rootCmd.GenBashCompletionV2(out, true)
			case "zsh":
				return rootCmd.GenZshCompletion(out)
			case "fish":
				return rootCmd.GenFishCompletion(out, true)
			case "powershell":
				return rootCmd.GenPowerShellCompletionWithDesc(out)
			default:
				return fmt.Errorf("unsupported shell %q", args[0])
			}
		},
		ValidArgs: []string{"bash", "zsh", "fish", "powershell"},
	}
	rootCmd.SetHelpCommand(helpCmd)
	rootCmd.AddCommand(versionCmd)
	rootCmd.AddCommand(completionCmd)
	rootCmd.AddCommand(newConfigCmd())
	rootCmd.AddCommand(newWhoamiCmd())
	rootCmd.AddCommand(newAccountsCmd())
	rootCmd.AddCommand(newPostsCmd())
	rootCmd.AddCommand(newTeamsCmd())
	rootCmd.AddCommand(newAnalyticsCmd())
	rootCmd.AddCommand(newAICmd())
	rootCmd.AddCommand(newNotificationsCmd())
	rootCmd.AddCommand(newCurationCmd())
	rootCmd.AddCommand(newMediaCmd())
	return rootCmd
}

func apiClient() (*client.Client, error) {
	if err := ensureConfig(); err != nil {
		return nil, err
	}
	if strings.TrimSpace(cfg.APIKey) == "" {
		return nil, fmt.Errorf("missing API key, run `socialbu config set-key <key>` or set SOCIALBU_API_KEY")
	}
	return client.New(cfg.BaseURL, cfg.APIKey), nil
}

func ensureConfig() error {
	if strings.TrimSpace(cfg.APIKey) != "" || strings.TrimSpace(cfg.BaseURL) != "" {
		return nil
	}
	if err := config.Init(); err != nil {
		return err
	}
	cfg = config.Current()
	return nil
}
