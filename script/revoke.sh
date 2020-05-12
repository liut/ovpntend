set -e

addr=$1

if [ ! -n "$addr" ] ; then
	exit 11
fi

# revoke

expect << EOF
spawn $EASYRSA/easyrsa revoke $addr nopass 
expect "with revocation:"
send "yes\n"
expect "/private/ca.key:"
send $PASSWORD\n
expect eof
EOF

if [ $? -ne 0 ] ; then
	exit 12
fi

# update crl

expect << EOF
spawn $EASYRSA/easyrsa gen-crl
expect "/private/ca.key:"
send $PASSWORD\n
expect eof
EOF

