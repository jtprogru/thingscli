package cli

import (
	"errors"
	"fmt"

	"github.com/jtprogru/thingscli/internal/things"
	"github.com/spf13/cobra"
)

func addWriteCommands(root *cobra.Command) {
	root.AddCommand(
		addCmd(),
		simpleIDCmd("done", "Mark to-do completed", func(c *things.Client, id string) error { return c.Done(id) }, "DONE"),
		simpleIDCmd("cancel", "Mark to-do canceled", func(c *things.Client, id string) error { return c.Cancel(id) }, "CANCELED"),
		simpleIDCmd("reopen", "Reopen a completed/canceled to-do", func(c *things.Client, id string) error { return c.Reopen(id) }, "REOPENED"),
		simpleIDCmd("trash", "Move to-do to Trash (reversible)", func(c *things.Client, id string) error { return c.Trash(id) }, "TRASHED"),

		&cobra.Command{
			Use:   "rename <id> <new-title>",
			Short: "Rename a to-do",
			Args:  cobra.ExactArgs(2),
			RunE: func(_ *cobra.Command, args []string) error {
				c, err := newClient()
				if err != nil {
					return err
				}
				if err := c.Rename(args[0], args[1]); err != nil {
					return err
				}
				fmt.Println("RENAMED", args[0])
				return nil
			},
		},
		&cobra.Command{
			Use:   "note <id> <text>",
			Short: "Overwrite a to-do's notes",
			Args:  cobra.ExactArgs(2),
			RunE: func(_ *cobra.Command, args []string) error {
				c, err := newClient()
				if err != nil {
					return err
				}
				if err := c.Note(args[0], args[1]); err != nil {
					return err
				}
				fmt.Println("NOTE-SET", args[0])
				return nil
			},
		},
		moveCmd(),
		scheduleCmd(),
		dueCmd(),
		&cobra.Command{
			Use:   "tag <id> <comma,tags>",
			Short: "Replace tag list on a to-do",
			Args:  cobra.ExactArgs(2),
			RunE: func(_ *cobra.Command, args []string) error {
				c, err := newClient()
				if err != nil {
					return err
				}
				if err := c.Tag(args[0], args[1]); err != nil {
					return err
				}
				fmt.Println("TAGGED", args[0])
				return nil
			},
		},
	)
}

func addCmd() *cobra.Command {
	var opt things.AddOptions
	cmd := &cobra.Command{
		Use:   "add <title>",
		Short: "Create a new to-do",
		Args:  cobra.ExactArgs(1),
		RunE: func(_ *cobra.Command, args []string) error {
			// Mutual exclusion check for project/area/list.
			set := 0
			for _, v := range []string{opt.Project, opt.Area, opt.List} {
				if v != "" {
					set++
				}
			}
			if set > 1 {
				return errors.New("--project, --area, --list are mutually exclusive")
			}
			c, err := newClient()
			if err != nil {
				return err
			}
			id, err := c.Add(args[0], opt)
			if err != nil {
				return err
			}
			fmt.Printf("ADDED %s :: %s\n", id, args[0])
			return nil
		},
	}
	cmd.Flags().StringVar(&opt.Notes, "notes", "", "to-do notes")
	cmd.Flags().StringVar(&opt.Project, "project", "", "place in this project")
	cmd.Flags().StringVar(&opt.Area, "area", "", "place in this area")
	cmd.Flags().StringVar(&opt.List, "list", "", "place in this built-in list (localized name)")
	cmd.Flags().StringVar(&opt.When, "when", "", "schedule: today | tomorrow | YYYY-MM-DD")
	cmd.Flags().StringVar(&opt.Tags, "tags", "", "comma-separated tag list")
	return cmd
}

func simpleIDCmd(name, short string, fn func(*things.Client, string) error, verb string) *cobra.Command {
	return &cobra.Command{
		Use:   name + " <id>",
		Short: short,
		Args:  cobra.ExactArgs(1),
		RunE: func(_ *cobra.Command, args []string) error {
			c, err := newClient()
			if err != nil {
				return err
			}
			if err := fn(c, args[0]); err != nil {
				return err
			}
			fmt.Println(verb, args[0])
			return nil
		},
	}
}

func moveCmd() *cobra.Command {
	var project, area, list string
	cmd := &cobra.Command{
		Use:   "move <id>",
		Short: "Move a to-do to a project, area, or list",
		Args:  cobra.ExactArgs(1),
		RunE: func(_ *cobra.Command, args []string) error {
			set := 0
			var t things.MoveTarget
			if project != "" {
				t = things.MoveTarget{Kind: "project", Name: project}
				set++
			}
			if area != "" {
				t = things.MoveTarget{Kind: "area", Name: area}
				set++
			}
			if list != "" {
				t = things.MoveTarget{Kind: "list", Name: list}
				set++
			}
			if set != 1 {
				return errors.New("specify exactly one of --project, --area, --list")
			}
			c, err := newClient()
			if err != nil {
				return err
			}
			if err := c.Move(args[0], t); err != nil {
				return err
			}
			fmt.Printf("MOVED %s -> %s %s\n", args[0], t.Kind, t.Name)
			return nil
		},
	}
	cmd.Flags().StringVar(&project, "project", "", "target project name")
	cmd.Flags().StringVar(&area, "area", "", "target area name")
	cmd.Flags().StringVar(&list, "list", "", "target list name (localized)")
	return cmd
}

func scheduleCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "schedule <id> <today|tomorrow|YYYY-MM-DD|someday>",
		Short: "Schedule a to-do",
		Args:  cobra.ExactArgs(2),
		RunE: func(_ *cobra.Command, args []string) error {
			c, err := newClient()
			if err != nil {
				return err
			}
			if err := c.Schedule(args[0], args[1]); err != nil {
				return err
			}
			fmt.Printf("SCHEDULED %s -> %s\n", args[0], args[1])
			return nil
		},
	}
}

func dueCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "due <id> <YYYY-MM-DD|clear>",
		Short: "Set or clear a to-do's due date",
		Args:  cobra.ExactArgs(2),
		RunE: func(_ *cobra.Command, args []string) error {
			c, err := newClient()
			if err != nil {
				return err
			}
			if err := c.Due(args[0], args[1]); err != nil {
				return err
			}
			if args[1] == "clear" {
				fmt.Println("DUE-CLEARED", args[0])
			} else {
				fmt.Printf("DUE-SET %s -> %s\n", args[0], args[1])
			}
			return nil
		},
	}
}
