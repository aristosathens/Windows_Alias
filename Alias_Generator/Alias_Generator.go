// Aristos Athens
// wiki: https://github.com/aristosathens/Windows_Cmd_Aliases/wiki

package main

import (
	. "Cmd_Commands_Windows"
	"bufio"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"text/tabwriter"
)

const (
	folder = "C:\\Cmd_Aliases\\"
)

var reader *bufio.Reader
var currentAliases []string

// ------------------------------------------- Main ------------------------------------------- //

func main() {
	fmt.Println("")
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
		fmt.Println("Ready to use.")
		fmt.Println("Type 'alias help' for help.")

		return
	}
	checkPath()

	// If alias.cmd is missing, replace it
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

	// Check for keywords, respond as necessary
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

	// If user input passes all tests, generate .cmd file
	addAlias(args)
}

// ------------------------------------------- Private ------------------------------------------- //

// If valid arguments provided, generates .cmd file to persistently assign alias to command
func addAlias(args []string) {

	if !validArguments(args) {
		return
	}

	if isInArray(args[1], currentAliases) {
		fmt.Print("You already have an alias with that name. Replace it? (y/n) ")
		for {
			text := readUserInput()
			if text == "no" || text == "n" {
				return
			} else if text == "yes" || text == "y" {
				break
			} else {
				fmt.Println("Invalid input. Enter yes or no.")
			}
		}
	}
	argsString := strings.Join(args[2:], " ")
	generateCMD(args[1], []string{argsString}, folder)
}

// Add an alias for executing multiple commands.
// Use this for using special cmd syntax, likek %CD%, etc.
// fmt.Println("See here for more details: http://www.robvanderwoude.com/parameters.php, http://www.robvanderwoude.com/batchcommands.php")
func addSpecialAlias() {

	var text string
	var name string
	var commands []string

	fmt.Println("Here you can add multiple commands to an alias and use special cmd syntax for arguments.")
	fmt.Print("Enter alias name: ")
	for {
		text = readUserInput()
		if isNameAvailable(text) {
			name = text
			break
		} else {
			return
		}
	}

	var command string
	for {
		fmt.Print("Enter command: ")
		text = readUserInput()
		if isCommandAvailable(text) {
			command = text
		} else {
			continue
		}
		fmt.Print("Enter arguments for " + command + " command (empty string -> no arguments): ")
		text = readUserInput()
		command += " " + text
		commands = append(commands, command)
		fmt.Print("Add more commands? (y/n) ")
		for {
			text = strings.ToLower(readUserInput())
			if text == "no" || text == "n" {
				generateCMD(name, commands, folder)
				return
			} else if text == "yes" || text == "y" {
				break
			} else {
				fmt.Println("Invalid input. Enter yes or no.")
			}
		}
	}
}

// If alias exists, deletes it
func removeAlias(name string) {
	path := folder + name + ".cmd"
	if fileExists(path) {
		os.Remove(path)
		fmt.Println("Removed '" + name + "' alias.")
	} else {
		fmt.Println("No such alias. Type 'alias list' to see all aliases.")
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

// Checks if command line arguments are valid
func validArguments(args []string) bool {

	if len(args) < 3 {
		fmt.Println("Wrong number of arguments. Expected 2 arguments. Enter commands in the following format")
		fmt.Println("$ alias <yourAlias> <yourCommand>")
		return false
	}

	name := args[1]
	command := args[2]

	if !isNameAvailable(name) {
		return false
	}
	if !isCommandAvailable(command) {
		fmt.Println(command + " is not a valid command.")
		return false
	}
	if name == "alias" {
		fmt.Println("Cannot overwrite 'alias' name.")
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

// ------------------------------------------- Display ------------------------------------------- //

// Prints help
func displayHelp() {
	tabWriter := tabwriter.NewWriter(os.Stdout, 0, 0, 0, ' ', tabwriter.Debug)

	fmt.Fprintln(tabWriter, "COMMAND \t DESCRIPTION")
	fmt.Fprintln(tabWriter, "------- \t -----------")
	fmt.Fprintln(tabWriter, "alias <Alias> <Command> \t Assign <Alias> to <Command>")
	fmt.Fprintln(tabWriter, "alias delete <Alias> \t Remove assigned alias")
	fmt.Fprintln(tabWriter, "alias special \t Add multi command or special alias")
	fmt.Fprintln(tabWriter, "alias list \t List all assigned aliases")
	fmt.Fprintln(tabWriter, "alias help \t Display help")

	tabWriter.Flush()
}

// Prints currently defined aliases
func displayAliases() {
	tabWriter := tabwriter.NewWriter(os.Stdout, 0, 0, 0, ' ', tabwriter.Debug)
	fmt.Fprintln(tabWriter, "NAME \t COMMAND(S)")
	fmt.Fprintln(tabWriter, "---- \t ----------")

	for _, alias := range currentAliases {
		command := getAliasCommand(alias)
		leftCol := alias
		for {
			i := strings.Index(command, "\n")
			if i == -1 {
				i = len(command)
				fmt.Fprintln(tabWriter, leftCol+" \t "+command[:i])
				break
			}
			column := leftCol + " \t " + command[:i]
			fmt.Fprintln(tabWriter, column)
			command = command[i+1:]
			leftCol = ""
		}
	}
	tabWriter.Flush()
}

// ------------------------------------------- Utility ------------------------------------------- //

// Respond to errors
func checkError(err error) {
	if err != nil {
		panic(err)
	}
}

// Get user input
func readUserInput() string {
	if reader == nil {
		reader = bufio.NewReader(os.Stdin)
	}
	input, err := reader.ReadString('\n')
	checkError(err)
	input = strings.TrimSpace(input)
	return input
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

// Checks if command exists
func isCommandAvailable(command string) bool {
	if fileExists(command) || commandExists(command) {
		return true
	} else {
		fmt.Print("The command/file you have entered does not exist. Use anyway? (y/n) ")
		var text string
		for {
			text = readUserInput()
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

// Checks if file exists
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

// Checks if command is available in command line by calling "where" command
func commandExists(name string) bool {
	if isInArray(name, GetAllCmdCommands()) {
		return true
	}
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
	// func Replace(s, old, new string, n int)
	index := strings.Index(string(body), "\n")
	return string(body)[index+1:]
}

// Recursively reads through a string with elements separated by semicolons. Checks if input string is an element
func pathContains(input string, path string) bool {

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

	if !pathContains(folder, path) {
		cwd, err := filepath.Abs(filepath.Dir(os.Args[0]))
		checkError(err)
		generatePathChangerCMD(cwd + "\\")
		c := exec.Command("cmd", "/c", cwd+"\\addToUserPath.cmd", folder)
		c.Run()
		err = os.Remove(cwd + "\\addToUserPath.cmd")
		checkError(err)
		os.Setenv("Path", path+folder+";")
		fmt.Println("Your path has been updated.")
	}
}

// This generates a .cmd file (similar to .bat) which let's us permanently change the User's Path variable
func generatePathChangerCMD(location string) {
	generateCMD(
		"addToUserPath",
		[]string{
			"REM usage: append_user_path \"path\"",
			"SET Key=\"HKCU\\Environment\"",
			"FOR /F \"usebackq tokens=2*\" %%A IN (`REG QUERY %Key% /v PATH`) DO Set CurrPath=%%B",
			"ECHO %CurrPath% > user_path_backup.txt",
			"SETX PATH \"%CurrPath%\"%1;",
		},
		location,
	)
}
