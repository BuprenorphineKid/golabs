#!/bin/bash
if [ "$(uname -o)" == "Android" ]
   then is_termux=true;
fi

echo -e "\033[34mInstalling Dependencies...\033[0m"

echo -e "\033[33mtbox"
go install services/tbox/main.go

echo -e "goimports\033[0m"
go install golang.org/x/tools/cmd/goimports@latest

echo -e "\033[34mInstalling Labs :)\033[0m"
go build cmd/labs/main.go
mv main labs

if [[ $is_termux == true ]];
  then 
  echo -e "\033[termux-elf-cleaner0m"
  echo -e "\033[adding dependencies.0m"
  apt install git make automake -y
  
  cur_dir=$(pwd)
  
  cd ~
  echo -e "\033[cloning repo...0m"
  git clone https://github.com/termux/termux-elf-cleaner/ && cd termux-elf-cleaner/
  
  echo -e "\033[automake0m"
  aclocal && autoconf
  automake --add-missing
  automake
  
  bash "./configure"
  bash "./config.status"

  echo -e "\033[make0m"
  make
  
  echo -e "\033[make install0m"
  make install

  echo -e "\033[add to path.0m"
  mv termux-elf-cleaner $PREFIX/bin
  termux-reload-settings

  cd $cur_dir
  
  echo -e "\033[Cleaning elf files...0m"

  android_version=$(termux-info | grep -e "Android version" -A 1 | tail -n 1)

  case "$android_version" in
  "7.0*") 
    api=24
    ;;
  "7.1") 
    api=25
    ;;
  "8.0") 
    api=26
    ;;
  "8.1*") 
    api=27
    ;;
  "9.0") 
    api=28
    ;;
  "9.1") 
    api=29
    ;;
  "10.0*") 
    api=30
    ;;
  "10.1") 
    api=31
    ;;
  "11.0") 
    api=32
    ;;
  esac
  

  termux-elf-cleaner --api-level $api labs
  termux-elf-cleaner --api-level $api $GOPATH/bin/tbox

  echo -e "\033[31mDone.\033[0m"
fi



