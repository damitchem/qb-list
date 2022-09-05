package parser

import (
	"errors"
	"fmt"
	"regexp"
	"strings"
)

var questFlagRegex *regexp.Regexp

func init() {
	questFlagRegex, _ = regexp.Compile("([a-zA-Z0-9_()']*\\ -\\ Completed\\!)")
}

type AccountQbList struct {
	QBs []string
}

func (a *AccountQbList) addQB(line string) error {
	if !strings.Contains(line, "SYSTEM") || !strings.Contains(line, "Completed!") {
		return errors.New(fmt.Sprintf("probably not a QB: %v", line))
	}

	completedFlag := questFlagRegex.FindString(line)

	if completedFlag == "" {
		return errors.New(fmt.Sprintf("failed to match quest flag: %v", line))
	}

	split := strings.Split(completedFlag, " - ")
	if split[0] == "" {
		return errors.New(fmt.Sprintf("failed to find flag name: %v", line))
	}
	a.QBs = append(a.QBs, strings.TrimSpace(split[0]))

	return nil
}
