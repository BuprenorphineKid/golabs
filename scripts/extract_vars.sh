#!/bin/bash
# Extract variables from current session

dir=$(dirname $0)
varfile="$dir/../.labs/session/vars"
eval="$dir/../.labs/session/eval.go"

touch $varfile

function extract() {

	var=$(echo "$@" | grep "var" | sed s/"^\s*"//g | sed s/"$"/"\n"/g)
	inf=$(echo "$@" | grep ":=" | sed s/"^\s*"//g | sed s/"$"/"\n"/g)


	echo "$var" >> $varfile
	echo "$inf" >> $varfile

}

extract "$(cat $eval)"


