#! /bin/sh
### BEGIN INIT INFO
# Provides:		w-watch
# Required-Start:	$syslog $remote_fs
# Required-Stop:	$syslog $remote_fs
# Should-Start:		$local_fs
# Should-Stop:		$local_fs
# Default-Start:	2 3 4 5
# Default-Stop:		0 1 6
# Short-Description:	w-watch - remote command execution for dev/testing
# Description:		w-watch - remote command execution for dev/testing
### END INIT INFO

PATH=/usr/local/sbin:/usr/local/bin:/sbin:/bin:/usr/sbin:/usr/bin
DAEMON=/home/pschlump/Projects/w-watch/w-watch
NAME=w-watch
DESC=w-watch

RUNDIR=/home/pschlump/Projects/w-watch
PIDFILE=$RUNDIR/w-watch.pid

test -x $DAEMON || exit 0

if [ -r /etc/default/$NAME ]
then
	. /etc/default/$NAME
fi

. /lib/lsb/init-functions

set -e

case "$1" in
  start)
	echo -n "Starting $DESC: "
	mkdir -p $RUNDIR
	touch $PIDFILE

	cd $RUNDIR 
	$DAEMON > ,log 2>&1 & 
	THE_PID=$! 
	echo "$THE_PID" >$PIDFILE

	;;
  stop)
	echo -n "Stopping $DESC: "
	if [ -f $PIDFILE ] ; then
		kill $( cat $PIDFILE )
		rm -f $PIDFILE
	fi
	sleep 1
	;;

  restart|force-reload)
	${0} stop
	${0} start
	;;

  status)
	echo "Unknown:TBD"
	# wget -O - 'http://localhost:20000/api/status?fmt=text'
	;;

  *)
	echo "Usage: /etc/init.d/$NAME {start|stop|restart|force-reload|status}" >&2
	exit 1
	;;
esac

exit 0
