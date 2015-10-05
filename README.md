# watch-cli

Go (golang) tool for watching a path and running a command when files change.

Example.   Lets say you have a directory `~/Projects/tab-server1` and whenever this gets recompiled
you want `/home/pschlump/Projects/rpt-q/rpt-q-init.sh` to run.

``` bash

	$ watch-cli -c "bash -c '/home/pschlump/Projects/rpt-q/rpt-q-init.sh" ~/Projects/tab-server1/ 

```
