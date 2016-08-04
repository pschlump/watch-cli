// code for linux

package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"
)

var IsWindows = false
var IsBSD = false
var IsLinux = true
var IsUnix = true

var opts struct {
	Cmd  string `short:"c" long:"cmd"             description:"Command to run when files change"         default:"echo Hw"`
	CdTo string `short:"t" long:"cdto"            description:"Directory to cd to before running comand" default:"."`
	Cfg  string `short:"C" long:"cfg"             description:"JSON configuration file"                  default:"./watch-cli-cfg.json"`
}

var hookableSignals []os.Signal
var sigChan chan os.Signal

func init() {
	hookableSignals = []os.Signal{
		syscall.SIGHUP,
		syscall.SIGUSR1,
		syscall.SIGUSR2,
		syscall.SIGINT,
		syscall.SIGTERM,
		syscall.SIGTSTP,
	}
	sigChan = make(chan os.Signal)
}

// handleSignals listens for os Signals and ignores some.
func handleSignals() {
	var sig os.Signal

	signal.Notify(
		sigChan,
		hookableSignals...,
	)

	pid := syscall.Getpid()
	for {
		sig = <-sigChan
		switch sig {
		case syscall.SIGHUP:
			log.Printf("Pid=%d Received SIGHUP. ignored.\n", pid)
		case syscall.SIGUSR1:
			log.Printf("Pid=%d Received SIGUSR1. ignored.\n", pid)
		case syscall.SIGUSR2:
			log.Printf("Pid=%d Received SIGUSR2. ignored.\n", pid)
		case syscall.SIGINT:
			log.Printf("Pid=%d Received SIGINT. exit.\n", pid)
			os.Exit(1)
		case syscall.SIGTERM:
			log.Printf("Pid=%d Received SIGTERM. exit.\n", pid)
			os.Exit(1)
		case syscall.SIGTSTP:
			log.Printf("Pid=%d Received SIGTSTP. ignored.\n", pid)
		default:
			log.Printf("Received %v: irrelevant signal...\n", sig)
		}
	}
}
