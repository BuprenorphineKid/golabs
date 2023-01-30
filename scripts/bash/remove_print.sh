#!/bin/bash
# Bash script for creating a new empty session.

dir=$(dirname $0)

tmp="$dir/../../.labs/template"
session="$dir/../../.labs/session/lab.go"
copy="$dir/../../.labs/session/lab_copy.go"

if [[ ! -e $tmp ]]
then
	touch $tmp
	echo "package main\n\nimport(\n\n)\n\nfunc main() {\n\n}\n" >> $tmp
fi

i=1
while IFS= read -r line
do
	if [[ $line == "func main() {" ]]
	then
		start=$i
	fi

	if [[ $start != 0 ]] && [[ $line == "}" ]]
	then
		end=$i
		break
	fi

	i=$(($i + 1))

done < $session

cp $session $copy

sed -i $start,$end"s/^\s*\(fmt\)\?\.\?[Pp]rint.*$//g" $copy

rm $session

mv $copy $session
