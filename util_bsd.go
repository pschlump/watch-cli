// code for BSD

// +build darwin dragonfly freebsd netbsd openbsd

package main

var IsWindows = false
var IsBSD = true
var IsLinux = false
var IsUnix = true

var opts struct {
	Cmd  string `short:"c" long:"cmd"             description:"Command to run when files change"         default:"echo Hw"`
	CdTo string `short:"t" long:"cdto"            description:"Directory to cd to before running comand" default:"."`
	Cfg  string `short:"C" long:"cfg"             description:"JSON configuration file"                  default:"./watch-cli-cfg.json"`
}
