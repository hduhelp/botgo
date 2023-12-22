package command

import (
	"context"
	"strings"
)

func New(name, key, desc string, run func(context.Context, []string) any) *Command {
	return &Command{
		Name:        name,
		Key:         key,
		Description: desc,
		Run:         run,
	}
}

func (c *Command) AddSubCommands(subCommand ...*Command) *Command {
	if c.SubCommand == nil {
		c.SubCommand = make(map[string]*Command)
	}
	for _, command := range subCommand {
		if _, exists := c.SubCommand[command.Key]; exists {
			return c
		}
		c.SubCommand[command.Key] = command
	}
	return c
}

// parse 解析命令行
func (c *Command) parse(cmd string) []string {
	arr := strings.Split(cmd, " ")
	// remove empty string
	var result []string
	for _, v := range arr {
		if v != "" {
			result = append(result, v)
		}
	}
	return result
}

func (c *Command) Help(indent int) string {
	var result string
	if indent == 0 {
		result += "\n【帮助】\n"
		result += c.Key + " " + c.Name + ":" + c.Description + "\n"
	} else {
		result += strings.Repeat("  ", indent) + c.Key + " " + c.Name + ":" + c.Description + "\n"
	}

	if c.SubCommand != nil {
		for _, v := range c.SubCommand {
			result += v.Help(indent + 2)
		}
	}
	return result
}

func (c *Command) Exec(ctx context.Context, cmd string) any {
	args := c.parse(cmd)
	// 如果没有子命令，直接执行
	if len(c.SubCommand) == 0 {
		return c.Run(ctx, args)
	}
	// 如果有子命令，执行子命令
	if len(args) > 0 {
		if subCommand, exists := c.SubCommand[args[0]]; exists {
			return subCommand.Exec(ctx, strings.Join(args[1:], " "))
		}
	}

	return c.Help(0)
}
