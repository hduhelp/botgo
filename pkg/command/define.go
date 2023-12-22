package command

import "context"

type Command struct {
	Name        string
	Key         string
	Description string
	Run         func(ctx context.Context, args []string) any
	SubCommand  map[string]*Command
}
