#!/bin/bash

set -e

NUM=$1

if [ ! -n "$NUM" ] ; then
	NUM=2
fi

expect << EOF
spawn telnet 127.0.0.1 7505
expect "for more info"
send "status ${NUM}\n"
expect "END"
send "quit\n"
expect eof
EOF

