// code for linux

package main

var IsWindows = false
var IsBSD = false
var IsLinux = true
var IsUnix = true

var opts struct {
	Cmd  string `short:"c" long:"cmd"             description:"Command to run when files change"         default:"echo Hw"`
	CdTo string `short:"t" long:"cdto"            description:"Directory to cd to before running comand" default:"."`
}
