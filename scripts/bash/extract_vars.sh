#!/bin/bash
# Extract variables from current session

dir=$(dirname $0)
varfile="$dir/../../.labs/session/vars"
eval="$dir/../../.labs/session/eval.go"

rm -f $varfile
touch $varfile

function extract() {

	var=$(echo "$@" | grep "var" | sed s/"^\s*"//g | sed s/"$"/"\n"/g)
	inf=$(echo "$@" | grep ":=" | sed s/"^\s*"//g | sed s/"$"/"\n"/g)


	if [[ $var != "" ]]
	then
		echo "$var" | sed s/^\s*\B//g >> $varfile
	fi
	if [[ $inf != "" ]]
	then
		echo "$inf" | sed s/^\s*\B//g >> $varfile
	fi

}

extract "$(cat $eval)"


