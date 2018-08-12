// Aristos Athens
// Tool for automating the creation of persistent cmd aliases in Windows
// Run once to set up. Then run again with any of the following commands:
// alias <yourAlias> <yourCommand>
// alias delete <yourAlias>
// alias list
// alist help

package main

import (
	. "Cmd_Commands_Windows"
	"bufio"
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

const (
	folder = "C:\\Cmd_Aliases\\"
)

var allCmdCommands []string
var currentAliases []string

// ------------------------------------------- Main ------------------------------------------- //

func main() {
	args := os.Args
	if _, err := os.Stat(folder); os.IsNotExist(err) {
		fmt.Println("Alias_Generator running for the first time.")
		fmt.Println("Setting up...")
		fmt.Println("Creating " + folder + " directory.")
		// 0700 is the permissions level. Gives user read, write, execute permissions. http://permissions-calculator.org/
		os.Mkdir(folder, 0700)
		fmt.Println("Setting self alias.")
		generateOwnCMD()
		if checkPath() {
			fmt.Println("To use this tool, enter commands in the following format:")
			fmt.Println("$ alias <yourAliasName> <yourCommand>")
		}
		return
	}

	if !checkPath() {
		return
	}

	allCmdCommands = GetAllCmdCommands()
	currentAliases = getCurrentAliases()
	if !isInArray("alias", currentAliases) {
		generateOwnCMD()
	}

	if len(args) > 1 {
		arg := strings.TrimSpace(strings.ToLower(args[1]))
		if arg == "list" || (args[1] == "delete" && len(args) == 2) {
			displayAliases()
			return
		} else if args[1] == "delete" && len(args) == 3 {
			removeAlias(args[2])
			return
		} else if args[1] == "help" && len(args) == 2 {
			displayHelp()
			return
		}
	}

	addAlias(args)
}

// ------------------------------------------- Private ------------------------------------------- //

// Prints help
func displayHelp() {
	fmt.Println("-------------------- Alias Help --------------------")
	fmt.Println("-Add alias: alias <yourName> <yourCommand>")
	fmt.Println("-Remove alias: alias delete <yourName>")
	fmt.Println("-List all aliases: alias list")
	fmt.Println("-Display help: alias help")
}

// If valid arguments provided, generates .cmd file to persistently assign alias to command
func addAlias(args []string) {
	if !validArguments(args) {
		return
	}

	if isInArray(args[1], currentAliases) {
		fmt.Println("You already have an alias with that name. Replace it? (y/n)")
		reader := bufio.NewReader(os.Stdin)
		var text string
		for {
			text, _ = reader.ReadString('\n')
			text = strings.TrimSpace(text)
			if text == "no" || text == "n" {
				return
			} else if text == "yes" || text == "y" {
				break
			} else {
				fmt.Println("Invalid input. Enter yes or no.")
			}
		}
	}
	argsString := concatenateStringsWithSpaces(args[2:])

	generateCMD(args[1], argsString)
}

// If alias exists, deletes it. Else prints current aliases
func removeAlias(name string) {
	path := folder + name + ".cmd"
	if fileExists(path) {
		os.Remove(path)
		fmt.Println("Removed '" + name + "' alias.")
	} else {
		fmt.Println("No such alias. Type 'alias list' to see all aliases")
	}
}

// Generates a CMD to use "alias" as an alias for this program
func generateOwnCMD() {
	alias := "alias"
	command, err := filepath.Abs(filepath.Dir(os.Args[0]))
	command = "\"" + command + "\\" + os.Args[0]
	command += "\" %*"
	checkError(err)
	generateCMD(alias, command)
}

// Given a name and a command, generates the .cmd file necessary to properly assign the alias
func generateCMD(alias string, command string) {
	fmt.Println("Generating .cmd")
	file, err := os.Create(folder + alias + ".cmd")
	checkError(err)

	defer func() {
		err := file.Close()
		checkError(err)
	}()

	body := []byte("@echo off\n" + command)

	_, err = file.Write(body)
	checkError(err)
}

// Prints currently defined aliases
func displayAliases() {
	fmt.Println("-------------------- Current aliases --------------------")
	for _, alias := range currentAliases {
		fmt.Println(alias + " " + getAliasCommand(alias))
	}
}

// Checks if command line arguments are valid
func validArguments(args []string) bool {

	fmt.Println("num args: ", len(args))
	if len(args) < 3 {
		fmt.Println("Wrong number of arguments. Expected 2 arguments. Enter commands in the following format")
		fmt.Println("$ alias <yourAliasName> <yourCommand>")
		return false
	}

	name := args[1]
	command := args[2]

	if len(args) > 3 {
		command = concatenateStringsWithSpaces(args[2:])
		if !isCommandAvailable(command) {
			fmt.Println(command + " is not a valid command.")
			return false
		}
	}
	if name == "alias" {
		fmt.Println("Cannot overwrite 'alias' name.")
		return false
	}
	if !isNameAvailable(name) {
		return false
	}
	if !isCommandAvailable(command) {
		// fmt.Println("Not a valid command.")
		return false
	}

	return true
}

// Gets all of the custom aliases
func getCurrentAliases() []string {
	files, err := ioutil.ReadDir(folder)
	checkError(err)
	aliases := make([]string, len(files))

	for i, f := range files {
		name := f.Name()
		if name[len(name)-4:] == ".cmd" {
			aliases[i] = name[:len(name)-4]
		}
	}
	return aliases
}

// ------------------------------------------- Utilities ------------------------------------------- //

// Respond to errors
func checkError(err error) {
	if err != nil {
		panic(err)
	}
}

// Checks if array contains elem
func isInArray(elem string, arr []string) bool {
	for _, item := range arr {
		if elem == item {
			return true
		}
	}
	return false
}

func fileExists(path string) bool {
	_, err := os.Stat(path)
	if err == nil {
		return true
	}
	if os.IsNotExist(err) {
		return false
	}
	return true
}

// Concatenates strings, inserting a space between each pair of elements
func concatenateStringsWithSpaces(stringArray []string) string {
	var buffer bytes.Buffer

	for _, item := range stringArray {
		buffer.WriteString(item)
		buffer.WriteString(" ")
	}

	// Remove trailing space
	returnString := buffer.String()
	returnString = returnString[:len(returnString)-1]
	returnString = strings.TrimSpace(returnString)

	return returnString
}

// Checks if command exists
func isCommandAvailable(name string) bool {
	firstSpace := strings.Index(name, " ")
	var command string
	if firstSpace != -1 {
		command = name[:firstSpace]
	} else {
		command = name
	}

	if isInArray(command, allCmdCommands) {
		return true
	} else if isInArray(command, currentAliases) {
		return true
	}

	if fileExists(command) {
		return true
	} else {
		fmt.Println("The program you have entered does not exist. Create alias anyway? (y/n)")
		reader := bufio.NewReader(os.Stdin)
		var text string
		for {
			text, _ = reader.ReadString('\n')
			text = strings.TrimSpace(text)
			if text == "no" || text == "n" {
				return false
			} else if text == "yes" || text == "y" {
				return true
			} else {
				fmt.Println("Invalid input. Enter yes or no.")
			}
		}
	}
}

// Checks if name can be used as alias name
func isNameAvailable(name string) bool {
	if isInArray(name, allCmdCommands) {
		fmt.Println("Cannot use existing command as alias name.")
		return false
	} else if fileExists(name) {
		fmt.Println("Cannot use folder/file name as alias name.")
		return false
	} else {
		return true
	}
}

// Given alias name get the assigned command
func getAliasCommand(name string) string {
	body, err := ioutil.ReadFile(folder + name + ".cmd")
	checkError(err)
	index := strings.Index(string(body), "\n")
	return string(body)[index+1:]
}

// Checks that folder is in the System path. If it's not, prints instructions on how to add it
func checkPath() bool {

	path := os.Getenv("Path")

	if !pathContains(folder, path) {
		fmt.Println("")
		fmt.Println("ATTENTION: The folder containing the .cmd files must be added to the System Path. (Takes ~1 minute)")
		fmt.Println("Instructions (for Windows):")
		fmt.Println("-Open start menu. Search for 'System' (not 'System Information'")
		fmt.Println("-Navigate to System -> Advanced system settings -> Environment Variables")
		fmt.Println("-Under 'System variables', find 'Path'. Click on it and click 'Edit...'")
		fmt.Println("-Click 'New' and input '" + folder + "' (without the apostrophes)")
		fmt.Println("-Save all changes. Run this program again.")
		return false
	}
	return true
}

// Recursively reads through a string with elements separated by semicolons. Checks if input string is an element
func pathContains(input string, path string) bool {

	index := strings.Index(path, ";")

	if index == -1 {
		fmt.Println("WARNING: Path not formatted correctly.")
		return false
	}
	if input[:len(input)-1] == path[:index] {
		return true
	} else {
		// If we have reached the end of path, break the recursion
		if index == len(path)-1 {
			return false
		}
		return pathContains(input, path[index+1:])
	}
}
