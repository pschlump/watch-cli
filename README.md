watch-cli: Watch file system and run command when files change
==============================================================

Go (golang) tool for watching a path and running a command when files change.

Example.   Lets say you have a directory `~/Projects/tab-server1` and whenever this gets recompiled
you want `/home/pschlump/Projects/rpt-q/rpt-q-init.sh` to run.

``` bash

	$ watch-cli -c "bash -c '/home/pschlump/Projects/rpt-q/rpt-q-init.sh'" ~/Projects/tab-server1/ 

```

By default the watch-cli will read in a ./watch-cli-cfg.json configuration file.
This file can specify the files to watch.  This file is checked for after the `-t` change directory
command line option is applied.

```JSON

{
	"FilesToWatch": [ "abc.def", "./config-dir" ]
}

```
