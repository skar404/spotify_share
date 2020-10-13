package bot

import (
	"errors"
	"strings"
)

var WarningMessageNotCommand = errors.New("message not command")

func isCommand(m string) bool {
	if m[0] == '/' && len(m[1:]) > 0 {
		return true
	}
	return false
}

func getCommand(m string) (Command, error) {
	c := Command{}
	if !isCommand(m) {
		return c, WarningMessageNotCommand
	}

	mSplit := strings.Split(m, " ")
	c.Name = mSplit[0][1:]
	c.Args = mSplit[1:]

	return c, nil
}
