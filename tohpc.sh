#!/bin/bash
# Read Password
echo -n private key Password:
read -s password
echo
# Run Command
echo $password > pwdfile
./tohpc -pwdfile pwdfile -d