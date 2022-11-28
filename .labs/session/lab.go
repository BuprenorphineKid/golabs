package main

import(
"flag"
"fmt"
"strings"

)
var test string

func main() {
flag.Parse()
fmt.Println(strings.Split(test, " "))

}

func init() {
flag.StringVar(&test, "g", "", "")
flag.StringVar(&test, "golang", "", "")
}
