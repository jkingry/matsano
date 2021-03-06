/**
 * Created with IntelliJ IDEA.
 * User: jkingry
 * Date: 5/30/13
 * Time: 12:17 PM
 * To change this template use File | Settings | File Templates.
 */
package main

import (
	"bitbucket.org/jkingry/matsano/cmd"
	"bitbucket.org/jkingry/matsano/encoding"
	"bitbucket.org/jkingry/matsano/histogram"
	"bitbucket.org/jkingry/matsano/package1"
	"bitbucket.org/jkingry/matsano/package2"
)

func main() {
	encoding.Init(cmd.Flags())

	cmd.AddCommand(package1.Commands)
	cmd.AddCommand(package2.Commands)
	cmd.AddCommand(histogram.Commands)
	cmd.Run()
}
