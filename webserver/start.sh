#!/bin/bash

# Make sure only root can run our script
if [[ $EUID -ne 0 ]]; then
   echo "This script must be run as root" 1>&2
   exit 1
fi

DIR=`dirname $0`
cd $DIR

go build goserver.go
./goserver &
echo $! > /var/run/goserver.pid
