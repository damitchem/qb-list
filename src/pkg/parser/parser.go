package parser

import (
	"bufio"
	"bytes"
	"encoding/csv"
	"errors"
	"fmt"
	qbIo "github.com/damitchem/qb-list/pkg/io"
	"github.com/damitchem/qb-list/pkg/logger"
	"io"
	"regexp"
	"strings"
	"sync"
	"time"
)

type Parser struct {
	logger *logger.Logger
	io     *qbIo.IO
	wg     *sync.WaitGroup
}

func New(logger *logger.Logger, io *qbIo.IO, options ...func(*Parser)) *Parser {
	p := &Parser{
		logger: logger,
		io:     io,
		wg:     &sync.WaitGroup{},
	}

	for _, option := range options {
		option(p)
	}

	return p
}

func (p *Parser) Parse() (MasterList, error) {
	p.logger.Trace("Entering parse")
	p.wg.Add(3)
	errChan := make(chan error, 3)
	serverQuestChan := make(chan ServerQuestList, 1)
	serverQbChan := make(chan ServerQbList, 1)
	accountChan := make(chan AccountQbList, 1)
	defer close(serverQuestChan)
	defer close(serverQbChan)
	defer close(accountChan)
	defer close(errChan)
	go p.parseAccount(accountChan, errChan)
	go p.parseServerQb(serverQbChan, errChan)
	go p.parseServerQuest(serverQuestChan, errChan)
	data := MasterList{
		logger: p.logger,
	}
	var err error
	go func() {
		err = <-errChan
	}()
	go func() {
		data.ServerQuestList = <-serverQuestChan
	}()
	go func() {
		data.ServerQbList = <-serverQbChan
	}()
	go func() {
		data.AccountQbList = <-accountChan
	}()
	p.wg.Wait()
	if err != nil {
		return MasterList{}, err
	}
	time.Sleep(1 * time.Second)
	data.associateQbsToQuests().markQbsAsComplete().markQuestsAsComplete()
	return data, err
}

// Quest section "#------Quests------#"
func (p *Parser) parseAccount(responseChan chan AccountQbList, errChan chan error) {
	p.logger.Trace("Beginning to parse account qb list")
	defer p.wg.Done()
	f, err := p.io.GetAccountQbFile()
	if err != nil {
		errChan <- errors.New(fmt.Sprintf("error opening file: %v", err))
		return
	}
	reader := bufio.NewReader(bytes.NewBuffer(f))
	list := AccountQbList{}
	delim := byte('\n')
	questSection := false
	questSectionSeparator, _ := regexp.Compile("\\#(-*)Quests(-*)\\#")
	pageCount, _ := regexp.Compile("Page: [0-9]* \\/ [0-9]*")
	for line, err := reader.ReadString(delim); ; line, err = reader.ReadString(delim) {
		if err != nil {
			if err == io.EOF {
				break
			}
			errChan <- errors.New(fmt.Sprintf("error reading file: %v", err))
			return
		}
		if questSectionSeparator.MatchString(line) {
			questSection = !questSection
			continue
		} else if pageCount.MatchString(line) {
			continue
		}
		if questSection {
			err = list.addQB(line)
			if err != nil {
				fmt.Printf("WARN: error parsing account quest flag: %v", err)
			}
		}
	}
	responseChan <- list
}

// Both server files have header information at the top that needs to be recorded and stored separately

// parseServerQuest reads through the embedded server-quest-list.csv file and builds out the list of quests
func (p *Parser) parseServerQuest(responseChan chan ServerQuestList, errChan chan error) {
	defer p.wg.Done()
	f, err := p.io.GetServerQuestFile()
	if err != nil {
		errChan <- err
		return
	}
	reader := bufio.NewReader(bytes.NewBuffer(f))

	headerLines, err := getHeaderLines(reader)
	if err != nil {
		errChan <- err
		return
	}
	list := ServerQuestList{
		HeaderLines: headerLines,
		Quests:      make(map[string]*Quest),
	}
	bodyReader := csv.NewReader(reader)
	for {
		records, err := bodyReader.Read()
		if err != nil {
			if err != io.EOF {
				errChan <- err
			}
			break
		}
		err = list.addQuest(records)
		if err != nil {
			fmt.Printf("WARN: ERROR encountered in SERVER QUEST IMPORT: %v\nRaw data associated to import: %v\n", err, records)
		}
	}
	responseChan <- list
}

// parseServerQb reads through the embedded server-qb-list.csv file and builds out the list of qbs
func (p *Parser) parseServerQb(responseChan chan ServerQbList, errChan chan error) {
	defer p.wg.Done()
	f, err := p.io.GetServerQbFile()
	if err != nil {
		errChan <- err
		return
	}

	reader := bufio.NewReader(bytes.NewBuffer(f))

	headerLines, err := getHeaderLines(reader)
	if err != nil {
		errChan <- err
		return
	}
	list := ServerQbList{
		HeaderLines: headerLines,
		QBs:         make(map[string]*QB),
	}
	bodyReader := csv.NewReader(reader)
	unconfirmed := false
	for {
		records, err := bodyReader.Read()
		if err != nil {
			if err != io.EOF {
				errChan <- err
			}
			break
		}
		if len(records) > 0 && strings.Contains(records[0], "unconfirmed data dump from server") {
			unconfirmed = true
			records = append(records, "\n")
			list.UnconfirmedLine = strings.Join(records, ",")
			continue
		}
		err = list.addQb(records, unconfirmed)
		if err != nil {
			fmt.Printf("WARN: ERROR encountered in SERVER QB IMPORT: %v\nRaw data associated to error: %v\n", err, records)
		}
	}
	responseChan <- list
}

// getHeaderLines grabs the first three lines of the file as an array of strings
func getHeaderLines(reader *bufio.Reader) ([]string, error) {
	var headerLines []string
	for {
		s, err := reader.ReadString('\n')
		if err != nil {
			return nil, err
		}
		headerLines = append(headerLines, s)
		if strings.HasPrefix(s, "Quest") || strings.HasPrefix(s, "Flag") {
			break
		}
	}
	return headerLines, nil
}
