package cmd

import (
	"context"
	"fmt"

	"github.com/spf13/cobra"
	"github.com/usamaejaz/socialbu-cli/internal/output"
)

func newNotificationsCmd() *cobra.Command {
	cmd := &cobra.Command{Use: "notifications", Short: "Manage SocialBu notifications"}
	cmd.AddCommand(newNotificationsListCmd())
	cmd.AddCommand(newNotificationsUnreadCmd())
	cmd.AddCommand(newNotificationsGetCmd())
	cmd.AddCommand(newNotificationsMarkReadCmd())
	cmd.AddCommand(newNotificationsMarkUnreadCmd())
	cmd.AddCommand(newNotificationsMarkAllReadCmd())
	return cmd
}

func newNotificationsListCmd() *cobra.Command {
	return notificationGetCmd("list", "List notifications", "/notifications")
}

func newNotificationsUnreadCmd() *cobra.Command {
	return notificationGetCmd("unread", "List unread notifications", "/notifications/unread")
}

func newNotificationsGetCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "get <id>",
		Short: "Get a notification by ID",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return notificationAction("GET", fmt.Sprintf("/notifications/%s", args[0]), nil)
		},
	}
}

func newNotificationsMarkReadCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "mark-read <id>",
		Short: "Mark a notification as read",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return notificationAction("POST", fmt.Sprintf("/notifications/%s/mark_read", args[0]), map[string]any{})
		},
	}
}

func newNotificationsMarkUnreadCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "mark-unread <id>",
		Short: "Mark a notification as unread",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return notificationAction("POST", fmt.Sprintf("/notifications/%s/mark_unread", args[0]), map[string]any{})
		},
	}
}

func newNotificationsMarkAllReadCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "mark-all-read",
		Short: "Mark all notifications as read",
		RunE: func(cmd *cobra.Command, args []string) error {
			return notificationAction("POST", "/notifications/mark_all_read", map[string]any{})
		},
	}
}

func notificationGetCmd(use, short, endpoint string) *cobra.Command {
	return &cobra.Command{
		Use:   use,
		Short: short,
		RunE: func(cmd *cobra.Command, args []string) error {
			return notificationAction("GET", endpoint, nil)
		},
	}
}

func notificationAction(method, endpoint string, body any) error {
	cli, err := apiClient()
	if err != nil {
		return err
	}
	var resp any
	out := &resp
	if method == "POST" && endpoint == "/notifications/mark_all_read" {
		out = nil
	}
	if err := cli.Request(context.Background(), method, endpoint, nil, body, out); err != nil {
		return err
	}
	if out == nil {
		fmt.Println("OK")
		return nil
	}
	return renderNotifications(resp)
}

func renderNotifications(resp any) error {
	if jsonOutput {
		return output.JSON(resp)
	}
	items := notificationsFromAny(resp)
	if len(items) > 1 {
		rows := make([][]string, 0, len(items))
		for _, item := range items {
			rows = append(rows, []string{
				output.FirstNonEmpty(output.StringFromMap(item, "id"), "-"),
				output.FirstNonEmpty(
					output.StringFromMap(item, "title", "subject", "type", "event"),
					output.StringFromNestedMap(item, []string{"event", "name"}, []string{"type", "name"}),
					"-",
				),
				output.FirstNonEmpty(
					output.StringFromMap(item, "message", "body", "description", "text"),
					output.StringFromNestedMap(item, []string{"data", "message"}, []string{"payload", "message"}, []string{"payload", "text"}),
					"-",
				),
				output.FirstNonEmpty(output.StringFromMap(item, "created_at", "date", "updated_at"), output.StringFromNestedMap(item, []string{"timestamps", "created_at"}), "-"),
			})
		}
		output.Table([]string{"ID", "Title", "Message", "Created"}, rows)
		return nil
	}
	if len(items) == 1 {
		return renderNotificationDetail(items[0])
	}
	return output.JSON(resp)
}

func renderNotificationDetail(item map[string]any) error {
	values := map[string]string{}
	for _, key := range []string{"id", "title", "subject", "message", "body", "description", "text", "created_at", "updated_at", "read_at", "status"} {
		if value := output.StringFromMap(item, key); value != "" {
			values[key] = value
		}
	}
	if createdAt := output.FirstNonEmpty(values["created_at"], output.StringFromMap(item, "date"), output.StringFromNestedMap(item, []string{"timestamps", "created_at"})); createdAt != "" {
		values["created_at"] = createdAt
	}
	delete(values, "event")
	delete(values, "type")
	if event := output.FirstNonEmpty(
		output.StringFromNestedMap(item, []string{"event", "name"}, []string{"type", "name"}),
		output.StringFromMap(item, "type", "event"),
	); event != "" {
		values["event"] = event
	}
	if message := output.FirstNonEmpty(
		values["message"],
		values["body"],
		values["description"],
		values["text"],
		output.StringFromNestedMap(item, []string{"data", "message"}, []string{"payload", "message"}, []string{"payload", "text"}),
	); message != "" {
		values["message"] = message
	}
	if actor := output.MapFromMap(item, "actor", "user", "sender"); actor != nil {
		if name := output.FirstNonEmpty(output.StringFromMap(actor, "name", "full_name", "username"), output.StringFromMap(item, "actor_name", "user_name")); name != "" {
			values["actor"] = name
		}
	}
	if target := output.MapFromMap(item, "target", "resource", "entity"); target != nil {
		if label := output.FirstNonEmpty(output.StringFromMap(target, "name", "title", "type"), output.StringFromMap(item, "target_name")); label != "" {
			values["target"] = label
		}
	}
	if len(values) == 0 {
		return output.JSON(item)
	}
	output.KeyValue("Notification", values)
	return nil
}

func notificationsFromAny(resp any) []map[string]any {
	if m, ok := resp.(map[string]any); ok {
		if data := output.MapFromMap(m, "data"); data != nil {
			if items := output.SliceFromMap(data, "notifications", "results", "data"); len(items) > 0 {
				return items
			}
			return []map[string]any{data}
		}
		if items := output.SliceFromMap(m, "notifications", "results", "data"); len(items) > 0 {
			return items
		}
		return []map[string]any{m}
	}
	if items, ok := resp.([]any); ok {
		out := make([]map[string]any, 0, len(items))
		for _, item := range items {
			if entry, ok := item.(map[string]any); ok {
				out = append(out, entry)
			}
		}
		return out
	}
	return nil
}
