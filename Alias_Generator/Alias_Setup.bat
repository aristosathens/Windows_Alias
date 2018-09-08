@ECHO OFF
DEL "C:\Cmd_Aliases\alias.cmd"
IF NOT EXIST Alias_Generator.exe go build Alias_Generator.go
START /W Alias_Generator.exe
SET PATH=%PATH%C:\Cmd_Aliases\;