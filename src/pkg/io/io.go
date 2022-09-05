package io

import (
	"embed"
	"errors"
	"fmt"
	"io"
	"os"
)

//go:embed server-quest-list.csv
var serverQuestFile embed.FS

const serverQuestFileName = "server-quest-list.csv"

//go:embed server-qb-list.csv
var serverQbFile embed.FS

const serverQbFileName = "server-qb-list.csv"

type IO struct {
	accountQbFileName   string
	outputQbFileName    string
	outputQuestFileName string
}

func New(fileName string) *IO {
	return &IO{
		accountQbFileName: fileName,
	}
}

func (i *IO) GetServerQuestFile() ([]byte, error) {
	f, err := serverQuestFile.Open(serverQuestFileName)
	if err != nil {
		return nil, errors.New(fmt.Sprintf("error reading embedded server quest file: %v\n", err))
	}
	fileBuffer, err := io.ReadAll(f)
	if err != nil {
		return nil, errors.New(fmt.Sprintf("error reading quest file into buffer: %v\n", err))
	}
	return fileBuffer, nil
}

func (i *IO) GetServerQbFile() ([]byte, error) {
	f, err := serverQbFile.Open(serverQbFileName)
	if err != nil {
		return nil, errors.New(fmt.Sprintf("error reading embedded server qb file: %v\n", err))
	}
	fileBuffer, err := io.ReadAll(f)
	if err != nil {
		return nil, errors.New(fmt.Sprintf("error reading server qb file into buffer: %v\n", err))
	}
	return fileBuffer, nil
}

func (i *IO) GetAccountQbFile() ([]byte, error) {
	f, err := os.Open(i.accountQbFileName)
	if err != nil {
		return nil, errors.New(fmt.Sprintf("error reading qb file: %v\n", err))
	}
	fileBuffer, err := io.ReadAll(f)
	if err != nil {
		return nil, errors.New(fmt.Sprintf("error reading server qb file into buffer: %v\n", err))
	}
	return fileBuffer, nil
}
