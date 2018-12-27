## Description
This tool lets you assign persistent aliases to commands in the Windows cmd console. It is similar to the alias tool on *nix systems. It offers some advantages over using DOSKEY/other methods:
* **Automatic** - Don't need to manually add .bat/.cmd files or manually update system path.
* **Persistent** - Aliases remain after closing out of cmd console.
* **Organized** - Let's you view/add/delete aliases all in one place.

## Setup
Run `Alias_Setup`: This will run the Alias_Setup.bat file.

## Usage
`alias <yourAlias> <yourCommand> <yourArguments>`: Adds a new alias for yourCommand with alias name yourAlias. Will warn you if yourCommand is not a valid command or file, but it is up to the user to properly format arguments. Will ask you before overwriting already existing aliases. yourAlias cannot be the name of an existing command. yourCommand _can_ be an existing alias.

***

`alias special`: Use this tool to add a series of commands to a single alias. Will prompt first for an yourAlias, then will continuously prompt for commands and arguments. Also use this for commands that take special arguments or characters.

For example, typing `alias exampleAlias exampleCommand %CD%` will statically get the current directory, at time of typing. If instead you want the argument to actually be %CD% (e.g. if you want the command to dynamically take the current directory as an argument) then add it as an argument to a command assigned via `alias special`.

***

`alias delete <yourAlias>`: Removes alias.

***

`alias list`: Display a list of all assigned aliases.

***

`alias help`: Display help.

## Example
`alias ls dir /A %*`: Assign alias name 'ls' as an alias for the Windows `dir` command. `/A` is an argument to show hidden files. `%*` is an argument which ensures this alias can take arguments and pass it to the underlying command. Once assigned, type `ls` or `ls path/to/your/directory` to use.
