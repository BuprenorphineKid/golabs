#!/bin/bash
# Bash script for creating a new empty session.

dir=$(dirname $0)
tmp="$dir/../.labs/template"
session="$dir/../.labs/session/lab.go"

if [[ ! -e $tmp ]]
	touch $tmp
	echo "package main\n\nimport(\n\n)\n\nfunc main() {\n\n}\n" >> $tmp
fi

rm -f $session
cp $tmp $session
