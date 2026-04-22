package cmd

import (
	"fmt"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
)

func newFixturesCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "fixtures",
		Short: "Capture real JSON fixtures for renderer work",
		Long:  "Capture real SocialBu API responses into artifacts/samples so renderer changes are backed by actual endpoint shapes instead of guesses.",
	}
	cmd.AddCommand(newFixturesCaptureCmd())
	return cmd
}

func newFixturesCaptureCmd() *cobra.Command {
	var outDir string
	cmd := &cobra.Command{
		Use:   "capture",
		Short: "Print a ready-to-run fixture capture script",
		Long:  "Prints a deterministic shell script that captures the current fixture set into artifacts/samples. Requires a valid SocialBu API key.",
		RunE: func(cmd *cobra.Command, args []string) error {
			if strings.TrimSpace(outDir) == "" {
				outDir = "artifacts/samples"
			}
			absOutDir, err := filepath.Abs(outDir)
			if err != nil {
				return err
			}
			script := buildFixtureCaptureScript(absOutDir)
			fmt.Fprintln(cmd.OutOrStdout(), script)
			return nil
		},
	}
	cmd.Flags().StringVar(&outDir, "out-dir", "artifacts/samples", "output directory for captured fixture JSON")
	return cmd
}

func buildFixtureCaptureScript(outDir string) string {
	cliCmd := "go run ."
	commands := []struct {
		name string
		line string
	}{
		{"notifications-list", fmt.Sprintf("%s notifications list --json > %s/notifications-list.json", cliCmd, shellQuote(outDir))},
		{"notifications-unread", fmt.Sprintf("%s notifications unread --json > %s/notifications-unread.json", cliCmd, shellQuote(outDir))},
		{"notifications-get", fmt.Sprintf("# %s notifications get <id> --json > %s/notifications-get-<id>.json", cliCmd, shellQuote(outDir))},
		{"curation-topics", fmt.Sprintf("%s curation topics --json > %s/curation-topics.json", cliCmd, shellQuote(outDir))},
		{"curation-items", fmt.Sprintf("%s curation items --json > %s/curation-items.json", cliCmd, shellQuote(outDir))},
		{"curation-get", fmt.Sprintf("# %s curation get <id> --json > %s/curation-get-<id>.json", cliCmd, shellQuote(outDir))},
		{"media-status", fmt.Sprintf("# %s media status --key <key> --json > %s/media-status-<key>.json", cliCmd, shellQuote(outDir))},
		{"analytics-stats", fmt.Sprintf("%s analytics stats --json > %s/analytics-stats.json", cliCmd, shellQuote(outDir))},
		{"analytics-team-activity", fmt.Sprintf("%s analytics team-activity --json > %s/analytics-team-activity.json", cliCmd, shellQuote(outDir))},
		{"analytics-team-metrics", fmt.Sprintf("# %s analytics team-metrics --start YYYY-MM-DD --end YYYY-MM-DD --metrics posts,approvals --json > %s/analytics-team-metrics.json", cliCmd, shellQuote(outDir))},
		{"ai-autocomplete", fmt.Sprintf("%s ai autocomplete --content 'Draft social caption about AI scheduling' --json > %s/ai-autocomplete.json", cliCmd, shellQuote(outDir))},
		{"ai-generate", fmt.Sprintf("# %s ai generate --type generic --topic 'Productivity tips for social media managers' --account <id> --json > %s/ai-generate.json", cliCmd, shellQuote(outDir))},
	}

	var b strings.Builder
	b.WriteString("#!/usr/bin/env bash\n")
	b.WriteString("set -euo pipefail\n\n")
	b.WriteString("cd /root/.openclaw/workspace/worktrees/socialbu-cli-go-cobra-recover\n")
	b.WriteString(fmt.Sprintf("mkdir -p %s\n", shellQuote(outDir)))
	b.WriteString("\n")
	b.WriteString("if [[ -z \"${SOCIALBU_API_KEY:-}\" && -f \"$HOME/.socialbu/config.json\" ]]; then\n")
	b.WriteString("  SOCIALBU_API_KEY=$(python3 - <<'PY'\n")
	b.WriteString("import json, os\n")
	b.WriteString("path = os.path.expanduser('~/.socialbu/config.json')\n")
	b.WriteString("try:\n")
	b.WriteString("    data = json.load(open(path))\n")
	b.WriteString("except Exception:\n")
	b.WriteString("    data = {}\n")
	b.WriteString("print((data.get('api_key') or '').strip())\n")
	b.WriteString("PY\n")
	b.WriteString("  )\n")
	b.WriteString("  export SOCIALBU_API_KEY\n")
	b.WriteString("fi\n")
	b.WriteString("if [[ -z \"${SOCIALBU_API_KEY:-}\" ]]; then\n")
	b.WriteString("  echo 'SOCIALBU_API_KEY is not set in the environment or ~/.socialbu/config.json. Run socialbu config set-key <key> or export SOCIALBU_API_KEY first.' >&2\n")
	b.WriteString("  exit 1\n")
	b.WriteString("fi\n\n")
	for _, command := range commands {
		b.WriteString(fmt.Sprintf("echo 'Capturing %s'\n", command.name))
		b.WriteString(command.line)
		b.WriteString("\n\n")
	}
	return strings.TrimSpace(b.String())
}

func shellQuote(value string) string {
	return "'" + strings.ReplaceAll(value, "'", "'\\''") + "'"
}
