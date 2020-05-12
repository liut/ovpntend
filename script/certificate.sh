set -e

addr=$1

if [ ! -n "$addr" ] ; then
	exit 11
fi

expect << EOF
spawn $EASYRSA/easyrsa build-client-full $addr nopass 
expect "/private/ca.key:"
send $PASSWORD\n
expect eof
EOF
