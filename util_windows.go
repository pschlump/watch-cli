// code dedicated to windows

// +build windows

package main

var IsWindows = true
var IsBSD = false
var IsLinux = false
var IsUnix = false

var opts struct {
	Cmd  string `short:"c" long:"cmd"             description:"Command to run when files change"         default:"echo Hw"`
	CdTo string `short:"t" long:"cdto"            description:"Directory to cd to before running comand" default:"."`
	Cfg  string `short:"C" long:"cfg"             description:"JSON configuration file"                  default:"./watch-cli-cfg.json"`
}

func handleSignals() {
}
