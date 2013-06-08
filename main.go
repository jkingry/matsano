/**
 * Created with IntelliJ IDEA.
 * User: jkingry
 * Date: 5/30/13
 * Time: 12:17 PM
 * To change this template use File | Settings | File Templates.
 */
package main

import "bitbucket.org/jkingry/matsano/cmd"
import "bitbucket.org/jkingry/matsano/package1"

func main() {
	cmd.AddCommand(package1.Commands)
	cmd.Run()
}
