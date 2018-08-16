// Aristos Athens
// Tool for automating the creation of persistent cmd aliases in Windows
// Run once to set up. Then run again with any of the following commands:
// alias <yourAlias> <yourCommand>
// alias delete <yourAlias>
// alias list
// alist help

package main

import (
	// . "Alias_Path_Helper"
	"bufio"
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

const (
	folder = "C:\\Cmd_Aliases\\"
)

var currentAliases []string

// ------------------------------------------- Main ------------------------------------------- //

func main() {
	// test()
	// checkPath()
	// return
	args := os.Args
	if _, err := os.Stat(folder); os.IsNotExist(err) {
		fmt.Println("Alias_Generator running for the first time.")
		fmt.Println("Setting up...")
		fmt.Println("Creating " + folder + " directory.")

		// 0700 is the permissions level. Gives user read, write, execute permissions. http://permissions-calculator.org/
		os.Mkdir(folder, 0700)
		fmt.Println("Setting self alias.")
		generateOwnCMD()
		fmt.Println("Adding " + folder + " to system path.")
		checkPath()
		fmt.Println("To use this tool, enter commands in the following format:")
		fmt.Println("$ alias <yourAliasName> <yourCommand>")
		return
	}
	checkPath()

	// allCmdCommands = GetAllCmdCommands()
	currentAliases = getCurrentAliases()
	if !isInArray("alias", currentAliases) {
		generateOwnCMD()
		if len(args) == 1 {
			return
		}
	}

	if len(args) == 1 {
		fmt.Println("Type 'alias help' for help.")
		return
	}

	if len(args) > 1 {
		arg := strings.TrimSpace(strings.ToLower(args[1]))
		if arg == "list" || (arg == "delete" && len(args) == 2) {
			displayAliases()
			return
		} else if arg == "delete" && len(args) == 3 {
			if args[2] == "alias" {
				fmt.Println("Cannot delete the \"alias\" alias.")
				return
			}
			removeAlias(args[2])
			return
		} else if arg == "help" {
			displayHelp()
			return
		} else if arg == "special" {
			addSpecialAlias()
			return
		}
	}

	addAlias(args)
}

// ------------------------------------------- Debug ------------------------------------------- //

// func test() {
// 	path := os.Getenv("Path")
// 	fmt.Println("Path length: ", len(path))

// 	generatePathChangerCMD()
// }

// ------------------------------------------- Private ------------------------------------------- //

// Prints help
func displayHelp() {
	fmt.Println("-------------------- Alias Help --------------------")
	fmt.Println("-Add alias: alias <yourName> <yourCommand>")
	fmt.Println("-Remove alias: alias delete <yourName>")
	fmt.Println("-Add multi command or special alias: alias special")
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

	generateCMD(args[1], []string{argsString}, folder)
}

func addSpecialAlias() {
	fmt.Println("Here you can add multiple commands to an alias. You can also use special cmd syntax.")
	fmt.Println("For example, if you want a command that requires the path of the calling location:")
	fmt.Println("alias name: <yourName>")
	fmt.Println("command: <yourCommand>")
	fmt.Println("arguments: %CD%")
	// fmt.Println("See here for more details: http://www.robvanderwoude.com/parameters.php, http://www.robvanderwoude.com/batchcommands.php")

	reader := bufio.NewReader(os.Stdin)
	var text string
	var name string
	var commands []string

	fmt.Println("Enter alias name: ")
	for {
		text, _ = reader.ReadString('\n')
		text = strings.TrimSpace(text)
		if isNameAvailable(text) {
			name = text
			break
		}
	}

	var command string
	for {
		fmt.Println("Enter command: ")
		text, _ = reader.ReadString('\n')
		text = strings.TrimSpace(text)
		if isCommandAvailable(text) {
			command = text
			break
		}
	}
	for {
		fmt.Println("Enter arguments for " + command + " command (empty string -> no arguments: ")
		text, _ = reader.ReadString('\n')
		text = strings.TrimSpace(text)
		command += " " + text
		commands = append(commands, command)
		fmt.Println("Add more commands? (yes -> add more, no -> exit)")
		text, _ = reader.ReadString('\n')
		text = strings.ToLower(strings.TrimSpace(text))
		if text == "no" || text == "n" {
			break
		}
	}
	generateCMD(name, commands, folder)

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
	checkError(err)
	command = "\"" + command + "\\" + os.Args[0]
	command += "\" %*"
	generateCMD(alias, []string{command}, folder)
}

// Given a name and a command, generates the .cmd file necessary to properly assign the alias
func generateCMD(alias string, commands []string, location string) {
	fmt.Println("Generating .cmd")
	file, err := os.Create(location + alias + ".cmd")
	checkError(err)

	defer func() {
		err := file.Close()
		checkError(err)
	}()

	body := "@echo off\n"
	for _, command := range commands {
		body += command + "\n"
	}
	// remove trailing '\n'
	body = body[:len(body)-1]

	_, err = file.Write([]byte(body))
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

	if !isCommandAvailable(command) {
		fmt.Println(command + " is not a valid command.")
		return false
	}
	// }
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
func isCommandAvailable(command string) bool {
	if fileExists(command) || commandExists(command) {
		return true
	} else {
		fmt.Println("The command/file you have entered does not exist. Use anyway? (y/n)")
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

func commandExists(name string) bool {
	_, err := exec.Command("cmd", "/c", "where", name).Output()
	if err == nil {
		return true
	} else {
		return false
	}
}

// Checks if name can be used as alias name
func isNameAvailable(name string) bool {
	if fileExists(name) {
		fmt.Println("Cannot use folder/file name as alias name.")
		return false
	} else if strings.Index(name, " ") != -1 {
		fmt.Println("Cannot use spaces in alias name.")
		return false
	} else if name == "" {
		fmt.Println("Cannot use empty string as alias name. ")
		return false
	} else if commandExists(name) && !isInArray(name, currentAliases) {
		fmt.Println("Cannot use existing command as alias name.")
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

// Recursively reads through a string with elements separated by semicolons. Checks if input string is an element
func pathContains(input string, path string) bool {

	// fmt.Println("Pair: ", input, path)

	index := strings.Index(path, ";")

	if index == -1 {
		fmt.Println("WARNING: Path not formatted correctly.")
		return false
	}
	if input == path {
		return true
	}

	if input == path[:index] || input[:len(input)-1] == path[:index] {
		return true
	} else {
		// If we have reached the end of path, break the recursion
		if index == len(path)-1 {
			return false
		}
		return pathContains(input, path[index+1:])
	}
}

// Checks that folder is in the System path. If it's not, prints instructions on how to add it
func checkPath() {

	path := os.Getenv("Path")
	// fmt.Println(path)

	if !pathContains(folder, path) {
		cwd, err := filepath.Abs(filepath.Dir(os.Args[0]))
		// fmt.Println("CWD: ", cwd)
		checkError(err)
		generatePathChangerCMD(cwd + "\\")
		c := exec.Command("cmd", "/c", cwd+"\\addToUserPath.cmd", folder)
		c.Run()
		err = os.Remove(cwd + "\\addToUserPath.cmd")
		checkError(err)
		os.Setenv("Path", path+folder+";")
		fmt.Println("Your path has been updated.")
	}

	// fmt.Println(" ")
	// fmt.Println("New path: ", os.Getenv("Path"))
}

func generatePathChangerCMD(location string) {
	generateCMD(
		"addToUserPath",
		[]string{
			"REM usage: append_user_path \"path\"",
			"SET Key=\"HKCU\\Environment\"",
			"FOR /F \"usebackq tokens=2*\" %%A IN (`REG QUERY %Key% /v PATH`) DO Set CurrPath=%%B",
			"ECHO %CurrPath% > user_path_bak.txt",
			"SETX PATH \"%CurrPath%\"%1;",
		},
		location,
		// "C:/Users/arist/Desktop/Aristos Documents/Projects/Go/src/Alias_Generator/",
	)
}

// func generateVBS() {
// "Set oShell = WScript.CreateObject(\"WScript.Shell\")
// filename = oShell.ExpandEnvironmentStrings(\"%TEMP%\resetvars.bat\")
// Set objFileSystem = CreateObject(\"Scripting.fileSystemObject\")
// Set oFile = objFileSystem.CreateTextFile(filename, TRUE)

// set oEnv=oShell.Environment(\"System\")
// for each sitem in oEnv
//     oFile.WriteLine(\"SET \" & sitem)
// next
// path = oEnv(\"PATH\")

// set oEnv=oShell.Environment(\"User\")
// for each sitem in oEnv
//     oFile.WriteLine(\"SET \" & sitem)
// next

// path = path & \";\" & oEnv(\"PATH\")
// oFile.WriteLine(\"SET PATH=\" & path)
// oFile.Close"
// }
