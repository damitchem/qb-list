package main

import (
	"bufio"
	"encoding/csv"
	"flag"
	"fmt"
	"github.com/damitchem/qb-list/pkg/io"
	"github.com/damitchem/qb-list/pkg/logger"
	"github.com/damitchem/qb-list/pkg/parser"
	"os"
	"path/filepath"
)

func main() {
	accountQbFileFlag := flag.String("input", "", "Input file for your exported QBs. Relative or absolute path")
	outputDirectoryFlag := flag.String("output", ".", "Output directory for CSVs, defaults to current directory")
	logLevel := flag.String("loglevel", "Info", "Set minimum logging level, defaults to Info. Available values are Trace, Debug, Info, Warn, Error, Critical")

	flag.Parse()

	outputDirectory := *outputDirectoryFlag

	l := logger.New(logger.MinLogLevel(logger.GetLevel(*logLevel)))

	if accountQbFileFlag == nil || *accountQbFileFlag == "" {
		l.Critical("Please provide either an absolute or relative path containing your QB output!", "formats", []string{
			"C:\\directory\\path\\file.txt",
			"relative/directory/path/file.txt",
			"file.txt",
		})
		return
	}

	i := io.New(*accountQbFileFlag)
	p := parser.New(l, i)
	datalist, err := p.Parse()
	if err != nil {
		fmt.Println(err)
		return
	}
	err = writeOutput(datalist, outputDirectory)
	if err != nil {
		l.Critical("Unable to write to output directory", "error", err)
	}
}

func writeOutput(datalist parser.MasterList, outputDirectory string) error {
	err := ensureOutputDirectory(outputDirectory)
	if err != nil {
		return err
	}
	err = writeQuestFile(datalist.ServerQuestList, outputDirectory)
	if err != nil {
		return err
	}
	err = writeQbFile(datalist.ServerQbList, outputDirectory)
	return nil
}

func ensureOutputDirectory(outputDirectory string) error {
	if _, err := os.Stat(outputDirectory); os.IsNotExist(err) {
		return os.Mkdir(outputDirectory, 0755)
	} else if os.IsPermission(err) {
		return fmt.Errorf("error accessing output directory: %v", err)
	}
	return nil
}

func writeQbFile(list parser.ServerQbList, outputDirectory string) error {
	fileName := "qb-list.csv"
	f, err := os.Create(filepath.Join(outputDirectory, fileName))
	if err != nil {
		return fmt.Errorf("error creating quest file in output directory: %v", err)
	}
	writer := bufio.NewWriter(f)
	for _, line := range list.HeaderLines {
		_, err = writer.WriteString(line)
		if err != nil {
			return fmt.Errorf("error writing %v header: %v", fileName, err)
		}
		writer.Flush()
	}
	csvWriter := csv.NewWriter(writer)
	defer f.Close()

	var unconfirmed []*parser.QB
	// First write all the confirmed QBs
	for _, qb := range list.QBs {
		if qb.Unconfirmed {
			unconfirmed = append(unconfirmed, qb)
			continue
		}
		csvWriter.Write(qb.GetRecord())
	}

	// Write our unconfirmed beneath this line warning
	writer.WriteString(list.UnconfirmedLine)

	// Now write the rest of the unconfirmed QBs
	for _, qb := range unconfirmed {
		err = csvWriter.Write(qb.GetRecord())
		if err != nil {
			return err
		}
	}
	csvWriter.Flush()

	return nil
}

func writeQuestFile(list parser.ServerQuestList, outputDirectory string) error {
	fileName := "quest-list.csv"
	f, err := os.Create(filepath.Join(outputDirectory, fileName))
	if err != nil {
		return fmt.Errorf("error creating qb file in output directory: %v", err)
	}
	writer := bufio.NewWriter(f)
	for _, line := range list.HeaderLines {
		_, err = writer.WriteString(line)
		if err != nil {
			return fmt.Errorf("error writing %v header: %v", fileName, err)
		}
	}
	defer f.Close()
	csvWriter := csv.NewWriter(writer)
	for _, quest := range list.Quests {
		csvWriter.Write(quest.GetRecord())
	}
	csvWriter.Flush()
	return nil
}
