package cli

import (
	"github.com/jtprogru/thingscli/internal/things"
	"github.com/spf13/cobra"
)

func addReadCommands(root *cobra.Command) {
	root.AddCommand(
		builtinListCmd("today", "Today list", things.ListToday),
		builtinListCmd("inbox", "Inbox", things.ListInbox),
		builtinListCmd("upcoming", "Upcoming", things.ListUpcoming),
		builtinListCmd("anytime", "Anytime", things.ListAnytime),
		builtinListCmd("someday", "Someday", things.ListSomeday),
		builtinListCmd("logbook", "Logbook", things.ListLogbook),

		&cobra.Command{
			Use:   "list <name>",
			Short: "Todos of any list by its (localized) name",
			Args:  cobra.ExactArgs(1),
			RunE: func(_ *cobra.Command, args []string) error {
				c, err := newClient()
				if err != nil {
					return err
				}
				todos, err := c.List(args[0])
				if err != nil {
					return err
				}
				return renderTodos(todos)
			},
		},
		&cobra.Command{
			Use:   "project <name>",
			Short: "Todos of a project",
			Args:  cobra.ExactArgs(1),
			RunE: func(_ *cobra.Command, args []string) error {
				c, err := newClient()
				if err != nil {
					return err
				}
				todos, err := c.Project(args[0])
				if err != nil {
					return err
				}
				return renderTodos(todos)
			},
		},
		&cobra.Command{
			Use:   "search <text>",
			Short: "Search todos by name substring",
			Args:  cobra.ExactArgs(1),
			RunE: func(_ *cobra.Command, args []string) error {
				c, err := newClient()
				if err != nil {
					return err
				}
				todos, err := c.Search(args[0])
				if err != nil {
					return err
				}
				return renderTodos(todos)
			},
		},
		&cobra.Command{
			Use:   "show <id>",
			Short: "Show a single todo by id",
			Args:  cobra.ExactArgs(1),
			RunE: func(_ *cobra.Command, args []string) error {
				c, err := newClient()
				if err != nil {
					return err
				}
				t, err := c.Show(args[0])
				if err != nil {
					return err
				}
				return renderTodo(t)
			},
		},
		&cobra.Command{
			Use:   "projects",
			Short: "List all projects",
			Args:  cobra.NoArgs,
			RunE: func(_ *cobra.Command, _ []string) error {
				c, err := newClient()
				if err != nil {
					return err
				}
				ps, err := c.Projects()
				if err != nil {
					return err
				}
				return renderProjects(ps)
			},
		},
		&cobra.Command{
			Use:   "areas",
			Short: "List all areas",
			Args:  cobra.NoArgs,
			RunE: func(_ *cobra.Command, _ []string) error {
				c, err := newClient()
				if err != nil {
					return err
				}
				as, err := c.Areas()
				if err != nil {
					return err
				}
				return renderAreas(as)
			},
		},
		&cobra.Command{
			Use:   "tags",
			Short: "List all tag names",
			Args:  cobra.NoArgs,
			RunE: func(_ *cobra.Command, _ []string) error {
				c, err := newClient()
				if err != nil {
					return err
				}
				tags, err := c.Tags()
				if err != nil {
					return err
				}
				return renderTags(tags)
			},
		},
		&cobra.Command{
			Use:   "locale",
			Short: "Show the resolved Things locale",
			Args:  cobra.NoArgs,
			RunE: func(_ *cobra.Command, _ []string) error {
				c, err := newClient()
				if err != nil {
					return err
				}
				return renderLocale(c.Locale)
			},
		},
	)
}

func builtinListCmd(name, short string, key things.ListKey) *cobra.Command {
	return &cobra.Command{
		Use:   name,
		Short: short,
		Args:  cobra.NoArgs,
		RunE: func(_ *cobra.Command, _ []string) error {
			c, err := newClient()
			if err != nil {
				return err
			}
			todos, err := c.BuiltinList(key)
			if err != nil {
				return err
			}
			return renderTodos(todos)
		},
	}
}
