package parser

import (
	"errors"
	"strconv"
	"strings"
)

type ServerQuestList struct {
	HeaderLines []string
	Quests      map[string]*Quest
}

// Quest List file is expected in Quest,QuestLink,RepeatYN,Qb#,LumYN,FellowYN,AceYN,OptionalQBYN,LevelReq,RecLevel,Notes,CompletedTF format
func (s *ServerQuestList) addQuest(record []string) error {
	if len(record) == 0 {
		return errors.New("encountered empty quest record")
	}

	qbCount, err := strconv.Atoi(record[3])
	if err != nil {
		return err
	}

	quest := strings.ToUpper(strings.TrimSpace(record[0]))

	s.Quests[quest] = &Quest{
		QuestName:  strings.TrimSpace(record[0]),
		QuestLink:  record[1],
		Repeatable: record[2],
		TotalQb:    qbCount,
		Luminance:  record[4],
		Fellow:     record[5],
		ACE:        record[6],
		OptionalQb: record[7],
		LevelReq:   record[8],
		RecLevel:   record[9],
		Notes:      record[10],
		Completed:  false,
	}
	return nil
}

func (s *ServerQuestList) GetQuest(questName string) *Quest {
	name := strings.ToUpper(strings.TrimSpace(questName))
	return s.Quests[name]
}

type Quest struct {
	QuestName  string
	QuestLink  string
	Repeatable string
	TotalQb    int
	Luminance  string
	Fellow     string
	ACE        string
	OptionalQb string
	LevelReq   string
	RecLevel   string
	Notes      string
	Completed  bool
	// AssociatedQBs are any QBs that reference the Quest by name, a Quest may be Completed when all Associated QBs are finished
	AssociatedQBs []*QB
}

func (q *Quest) AddQB(qb *QB) *Quest {
	q.AssociatedQBs = append(q.AssociatedQBs, qb)
	return q
}

func (q *Quest) canComplete() bool {
	if q.AssociatedQBs == nil {
		return false
	}
	allDone := true
	for _, qb := range q.AssociatedQBs {
		if !qb.Completed {
			allDone = false
			break
		}
	}
	return allDone && len(q.AssociatedQBs) == q.TotalQb
}

func (q *Quest) Complete() {
	q.Completed = q.canComplete()
}

func (q *Quest) GetRecord() []string {
	return []string{
		q.QuestName,
		q.QuestLink,
		q.Repeatable,
		strconv.Itoa(q.TotalQb),
		q.Luminance,
		q.Fellow,
		q.ACE,
		q.OptionalQb,
		q.LevelReq,
		q.RecLevel,
		q.Notes,
		strconv.FormatBool(q.Completed),
	}
}
