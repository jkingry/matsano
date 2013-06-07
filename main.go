/**
 * Created with IntelliJ IDEA.
 * User: jkingry
 * Date: 5/30/13
 * Time: 12:17 PM
 * To change this template use File | Settings | File Templates.
 */
package main

import "fmt"
import "flag"
import "bitbucket.org/jkingry/matsano/package1"

func main() {
	commands := map[string]func([]string){
		"pkg1": package1.CommandLine,
	}

	flag.Usage = func() {
		fmt.Println("Commands: ")
		for key, value := range commands {
			fmt.Println(key)
			value([]string{"-help"})
		}
	}

	flag.Parse()

	if flag.NArg() == 0 {
		flag.Usage()
		return
	}

	name := flag.Arg(0)
	args := flag.Args()[1:]
	if cmd, ok := commands[name]; ok {
		cmd(args)
	} else {
		fmt.Printf("Unknown command: %v", name)
		flag.Usage()
	}
}
