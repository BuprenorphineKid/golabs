#!/bin/bash

echo -e "\033[34mInstalling Dependencies...\033[0m"
echo -e "\033[33mtbox"
go install services/tbox/main.go
echo -e "goimports\033[0m"
go install golang.org/x/tools/cmd/goimports@latest
echo -e "\033[34mInstalling Labs :)\033[0m"
go build cmd/labs/main.go
mv main labs
echo -e "\033[31mDone.\033[0m"
