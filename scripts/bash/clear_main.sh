#!/bin/bash
# Bash script for creating a new empty session.

dir=$(dirname $0)

tmp="$dir/../../.labs/template"
session="$dir/../../.labs/session/lab.go"

if [[ ! -e $tmp ]]
then
	touch $tmp
	echo "package main\n\nimport(\n\n)\n\nfunc main() {\n\n}\n" >> $tmp
fi

i=1
start=0
end=0

nest=false
while IFS= read -r line
do
	if [[ $line == "func main() {" ]]
	then
		start=$(($i + 1))
		continue
	fi

	if [[ $line =~ "{" ]]
	then
		nest=true
	fi

	if $nest && [[ $line == "}" ]]
	then
		nest=false
		continue
	fi

	if [[ $nest == false ]] && [[ $start != 0 ]] && [[ $line == "}" ]]
	then
		end=$(($i - 1))
		break
	fi

	i=$(($i + 1))

done < $session
if [[ $start == 0 ]] && [[ $end == 0 ]]
then
	return 0
fi

if (($start < $end))
then
	sed -i $start,$end"d" $session
fi

sed -i $start"i\ " $session
