package console

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

//Getenv get env variable with default value
func Getenv(key, defaultValue string) string {
	if val, ok := os.LookupEnv(key); ok {
		return val
	}
	return defaultValue
}

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

//AskUserChoice ask user a choice question
func AskUserChoice(question string, choices ...string) (c string, err error) {
	res, err := AskUser(question)
	if err != nil {
		return
	}
	if contains(choices, res) {
		c = res
	} else {
		err = fmt.Errorf("invalid choice, only one of these value is possible: [%s] ", strings.Join(choices, "/"))
	}
	return
}

func contains(a []string, x string) bool {
	for _, n := range a {
		if x == n {
			return true
		}
	}
	return false
}
