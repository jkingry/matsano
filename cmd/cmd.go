package cmd

import "flag"
import "fmt"
import "os"
import "strings"
import "io/ioutil"

type SubCommand interface {
	String() string
	Run([]string)
}

type Command struct {
	Flags *flag.FlagSet

	name string
	run  func([]string)
}

func NewCommand(name string, flags *flag.FlagSet, run func([]string)) *Command {
	cmd := new(Command)
	cmd.Flags = flags
	if cmd.Flags == nil {
		cmd.Flags = flag.NewFlagSet(name, flag.ContinueOnError)
	}
	cmd.name = name
	cmd.run = run

	return cmd
}

func (this *Command) String() string {
	return this.name
}

func (this *Command) Run(args []string) {
	if ok := this.Flags.Parse(args); ok != nil {
		return
	}

	this.run(this.Flags.Args())
}

type CommandSet struct {
	Flags *flag.FlagSet

	name     string
	commands []SubCommand
}

func NewCommandSet(name string) *CommandSet {
	cmd := new(CommandSet)
	cmd.name = name
	cmd.Flags = flag.NewFlagSet(name, flag.ContinueOnError)
	cmd.Flags.Usage = func() {
		fmt.Printf("Commands: %v", cmd.commands)
	}
	return cmd
}

func (this *CommandSet) String() string {
	return this.name
}

func (this *CommandSet) Run(args []string) {
	if ok := this.Flags.Parse(args); ok != nil {
		return
	}

	if this.Flags.NArg() == 0 {
		this.Flags.Usage()
		return
	}

	name := this.Flags.Arg(0)
	for _, cmd := range this.commands {
		if strings.HasPrefix(cmd.String(), name) {
			cmd.Run(this.Flags.Args()[1:])
			return
		}
	}

	fmt.Println("Unknown command:", name)
	this.Flags.Usage()
}

func (this *CommandSet) Add(cmd SubCommand) {
	this.commands = append(this.commands, cmd)
}

var rootCommandSet = NewCommandSet("")

func Add(cmd SubCommand) {
	rootCommandSet.Add(cmd)
}

func Run() {
	rootCommandSet.Run(os.Args[1:])
}

func GetInput(args []string) string {
	if len(args) != 0 {
		return args[0]
	}
	data, _ := ioutil.ReadAll(os.Stdin)

	return string(data)
}
