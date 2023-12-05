#!/bin/bash

re="labs"
err=false

ls ./cmd/ || err=true
if [[ err == true ]] then
  echo -e "\033[31mError: \033[33mCould not find cmd/ path\033[0m"
  exit
fi

if [[ ! $(ls ./cmd/) =~ $re ]] then
  echo -e "\033[31mError: \033[33mCould not find path to main.go\033[0m"
  exit
fi

echo -e "\033[33mInstalling Dependencies...\033[0m"

echo -e "\033[35mtbox"
go install "services/tbox/main.go"

echo -e "goimports\033[0m"
go install "golang.org/x/tools/cmd/goimports@latest"

echo -e "\033[32mInstalling Labs :)\033[0m"

go build "cmd/labs/main.go" || err=true
  if [[ err == true ]] then
    echo -e "\033[31mError: \033[33mCould not build labs\033[0m"
    exit
  fi
  

mv ./main ./labs || err=true
  if [[ err == true ]] then
    echo -e "\033[31mError: \033[33mCould not rename binary\033[0m" &&
    exit
  fi

echo -e "\033[32mDone.\033[0m"


