#!/bin/bash
# Bash script for creating a new empty session.


tmp="$HOME/.labs/template"
session="$HOME/.labs/session/lab.go"
eval="$HOME/.labs/session/eval.go"

if [[ ! -e $tmp ]]
then
	touch $tmp
	echo "package main\n\n\nimport(\n\n)\n\n\nfunc main() {\n\n\n}\n" >> $tmp
fi

rm -f $session $eval
cp $tmp $session
cp $tmp $eval
