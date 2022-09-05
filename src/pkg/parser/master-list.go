package parser

import (
	"github.com/damitchem/qb-list/pkg/logger"
)

type MasterList struct {
	logger          *logger.Logger
	AccountQbList   AccountQbList
	ServerQbList    ServerQbList
	ServerQuestList ServerQuestList
}

func (m *MasterList) associateQbsToQuests() *MasterList {
	for flag := range m.ServerQbList.QBs {
		qb := m.ServerQbList.QBs[flag]
		if qb.Quest == "" || qb.Quest == "Retired" {
			continue
		}
		quest := m.ServerQuestList.GetQuest(qb.Quest)
		if quest == nil && !qb.Unconfirmed {
			m.logger.Debug("Unable to match Confirmed QB to quest, this is not an error, just a missing quest record", "flag", qb.Flag, "quest", qb.Quest)
			continue
		} else if qb.Unconfirmed {
			continue
		}
		quest.AddQB(qb)
		qb.AssociateQuest(quest)
	}
	return m
}

func (m *MasterList) markQbsAsComplete() *MasterList {
	for _, flag := range m.AccountQbList.QBs {
		serverQb := m.ServerQbList.GetQB(flag)
		if serverQb == nil {
			m.logger.Info("QB marked as completed, but we don't have a record of the flag, report this!", "flag", flag)
			m.ServerQbList.addUndiscoveredQb(flag)
			continue
		}
		serverQb.Completed = true
	}
	return m
}

func (m *MasterList) markQuestsAsComplete() *MasterList {
	for name := range m.ServerQuestList.Quests {
		m.ServerQuestList.Quests[name].Complete()
	}
	return m
}
