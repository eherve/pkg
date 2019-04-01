package cliutil

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

//AskUser ask user a question
func AskUser(question string) (res string, err error) {
	reader := bufio.NewReader(os.Stdin)
	fmt.Printf(question)
	if res, err = reader.ReadString('\n'); err != nil {
		return
	}
	res = strings.TrimRight(strings.TrimRight(res, "\n"), "\r")
	return
}

//AskUserBool ask user a bool question
// retruns false if response is equal to npTmpl case insensitive
func AskUserBool(question string, noTmpl string) (b bool, err error) {
	res, err := AskUser(question)
	if err != nil {
		return
	}
	b = strings.ToLower(strings.Trim(res, " ")) != strings.ToLower(noTmpl)
	return
}

//AskUserInt ask user a int question
func AskUserInt(question string) (i int, err error) {
	res, err := AskUser(question)
	if err != nil {
		return
	}
	i, err = strconv.Atoi(strings.Trim(res, " "))
	return
}
