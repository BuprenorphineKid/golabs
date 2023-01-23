#!/bin/bash
# Bash script for creating a new empty session.

tmp="$0/../.labs/template"
session="$0/../.labs/session/lab.go"

if [[ ! -e $tmp ]]
	touch $tmp
	echo "package main\n\nimport(\n\n)\n\nfunc main() {\n\n}\n" >> $tmp
fi

rm -f $session
cp $tmp $session
