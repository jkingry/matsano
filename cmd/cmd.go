package cmd

import "flag"
import "fmt"
import "os"
import "strings"
import "io/ioutil"
import "net/http"

type Command struct {
	Flags   *flag.FlagSet
	Command func([]string)

	name     string
	usage    string
	commands []*Command
}

func NewCommand(name, usage string) *Command {
	cmd := new(Command)
	cmd.Flags = flag.NewFlagSet(name, flag.ContinueOnError)
	cmd.name = name
	cmd.usage = usage

	return cmd
}

func (this *Command) printUsage(level int) {
	indent := strings.Repeat(" ", level*2)
	if this != rootCommands {
		fmt.Println(indent, this.name+": "+this.usage)
	}

	this.Flags.VisitAll(func(flag *flag.Flag) {
		fmt.Printf(indent+"  -%s=%s: %s\n", flag.Name, flag.DefValue, flag.Usage)
	})

	if len(this.commands) == 0 {
		return
	}

	fmt.Println(indent, " commands:")
	for _, c := range this.commands {
		c.printUsage(level + 1)
	}
}

func (this *Command) Run(args []string) {
	if ok := this.Flags.Parse(args); ok != nil {
		return
	}

	if this.Command != nil {
		this.Command(this.Flags.Args())
		return
	}

	if this.Flags.NArg() == 0 {
		this.printUsage(0)
		return
	}

	name := this.Flags.Arg(0)
	for _, c := range this.commands {
		if strings.HasPrefix(c.name, name) {
			c.Run(this.Flags.Args()[1:])
			return
		}
	}

	fmt.Fprintln(os.Stderr, "Unknown command:", name)
	this.printUsage(0)
}

func (this *Command) AddCommand(c *Command) {
	this.commands = append(this.commands, c)
}

func (this *Command) Add(name, usage string) *Command {
	c := NewCommand(name, usage)
	this.AddCommand(c)
	return c
}

var rootCommands = NewCommand("", "")

func Add(name, usage string) *Command {
	return rootCommands.Add(name, usage)
}

func AddCommand(c *Command) {
	rootCommands.AddCommand(c)
}

func Run() {
	rootCommands.Run(os.Args[1:])
}

func GetInput(args []string, index int) string {
	var arg string
	if len(args) < index {
		arg = "stdin"
	} else {
		arg = args[index]
	}

	switch {
	case strings.HasPrefix(arg, "file:"):
		file, _ := os.Open(strings.TrimPrefix(arg, "file:"))
		defer file.Close()
		bytes, _ := ioutil.ReadAll(file)
		return string(bytes)
	case strings.HasPrefix(arg, "http:"):
		resp, _ := http.Get(arg)
		defer resp.Body.Close()
		bytes, _ := ioutil.ReadAll(resp.Body)
		return string(bytes)
	case strings.HasPrefix(arg, "stdin"):
		bytes, _ := ioutil.ReadAll(os.Stdin)
		return string(bytes)
	}

	return arg
}
