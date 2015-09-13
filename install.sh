#!/bin/bash

if [ "$(whoami)" == "root" ] ; then
	:
else
	echo "Usage: !! run as root"
	exit 1
fi

cp w-watch-init.sh /etc/init.d/w-watch
cd /etc
ln -s /etc/init.d/w-watch ./rc0.d/K98w-watch
ln -s /etc/init.d/w-watch ./rc1.d/K98w-watch
ln -s /etc/init.d/w-watch ./rc2.d/S98w-watch
ln -s /etc/init.d/w-watch ./rc3.d/S98w-watch
ln -s /etc/init.d/w-watch ./rc4.d/S98w-watch
ln -s /etc/init.d/w-watch ./rc5.d/S98w-watch
ln -s /etc/init.d/w-watch ./rc6.d/K98w-watch

