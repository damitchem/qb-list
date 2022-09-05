package parser

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
)

type ServerQbList struct {
	HeaderLines     []string
	UnconfirmedLine string
	QBs             map[string]*QB
}

// QB List file is expected in Flag,Quest,Action,Notes,ServerDef format
func (s *ServerQbList) addQb(record []string, unconfirmed bool) error {
	if len(record) == 0 {
		return errors.New("encountered empty quest record")
	}
	flag := strings.ToUpper(strings.TrimSpace(record[0]))
	if _, ok := s.QBs[flag]; ok {
		return errors.New(fmt.Sprintf("encountered duplicate flag: %v", record[0]))
	}
	s.QBs[flag] = &QB{
		Flag:        strings.TrimSpace(record[0]),
		Quest:       strings.TrimSpace(record[1]),
		Action:      record[2],
		Notes:       record[3],
		ServerDef:   record[4],
		Completed:   false,
		Unconfirmed: unconfirmed,
	}
	return nil
}

func (s *ServerQbList) addUndiscoveredQb(flagName string) {
	flag := strings.ToUpper(strings.TrimSpace(flagName))

	if _, ok := s.QBs[flag]; ok {
		// Somehow we're trying to add a QB we already have again
		return
	}
	s.QBs[flag] = &QB{
		Flag:        flagName,
		Quest:       "Unknown",
		Action:      "Unknown",
		Notes:       "This flag was not known in the master list, please report it in Discord",
		ServerDef:   "Unknown",
		Completed:   true,
		Unconfirmed: false,
	}
}

func (s *ServerQbList) GetQB(flagName string) *QB {
	flag := strings.ToUpper(strings.TrimSpace(flagName))
	return s.QBs[flag]
}

type QB struct {
	Flag        string `csv:"Flag"`
	Quest       string `csv:"Quest"`
	Action      string `csv:"Action"`
	Notes       string `csv:"Notes"`
	ServerDef   string `csv:"Server Side Definition"`
	Completed   bool   `csv:"Completed"`
	Unconfirmed bool   `csv:"-"`
	// AssociatedQuest is the referenced Quest -- not every QB has a referenced Quest
	AssociatedQuest *Quest `csv:"-"`
}

func (q *QB) AssociateQuest(quest *Quest) {
	q.AssociatedQuest = quest
}

func (q *QB) GetRecord() []string {
	return []string{
		q.Flag,
		q.Quest,
		q.Action,
		q.Notes,
		q.ServerDef,
		strconv.FormatBool(q.Completed),
	}
}
