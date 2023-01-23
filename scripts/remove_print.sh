#!/bin/bash
# Bash script for creating a new empty session.

dir=$(dirname $0)

tmp="$dir/../.labs/template"
session="$dir/../.labs/session/lab.go"
copy="$dir/../.labs/session/lab_copy.go"

if [[ ! -e $tmp ]]
then
	touch $tmp
	echo "package main\n\nimport(\n\n)\n\nfunc main() {\n\n}\n" >> $tmp
fi

cp $session $copy
sed -i s/^.*[Pp]rint.*$//g $copy

rm $session

mv $copy $session
