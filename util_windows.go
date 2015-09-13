// code dedicated to windows

// +build windows

package main

var IsWindows = true
var IsBSD = false
var IsLinux = false
var IsUnix = false

var opts struct {
	AppName string `short:"A" long:"application"     description:"Application to run"                       default:"watch-cli"`
	Cmd     string `short:"c" long:"cmd"             description:"Command to run when files change"         default:"echo Hw"`
	CdTo    string `short:"t" long:"cdto"            description:"Directory to cd to before running comand" default:"."`
}
