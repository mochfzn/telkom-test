package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/spf13/cobra"
)

type Log struct {
	Type        string    `json:"type"`
	DateTime    time.Time `json:"date_time"`
	LineCode    string    `json:"line_code"`
	Description string    `json:"description"`
}

func main() {
	var typeFile string
	var pathDest string

	rootCmd := &cobra.Command{
		Use:   "Exporting Error Log",
		Short: "A simple CLI for exporting error log to specified file",
		Long:  "A simple CLI for exporting error log to plain text file or JSON file.",
		Run: func(cmd *cobra.Command, args []string) {

		},
	}

	exportCmd := &cobra.Command{
		Use:     "mytools",
		Aliases: []string{"tool"},
		Short:   "Exporting error log to specified file",
		Args:    cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			if typeFile == "text" {
				exportToPlainText(args[0], pathDest)
			} else if typeFile == "json" {
				exportToJson(args[0], pathDest)
			} else {
				return fmt.Errorf("TIPE YANG DIPILIH TIDAK TERSEDIA")
			}

			return nil
		},
	}

	exportCmd.Flags().StringVarP(&typeFile, "type", "t", "text", "Only for text of json")
	exportCmd.Flags().StringVarP(&pathDest, "open", "o", "", "Path destination for file name")

	rootCmd.AddCommand(exportCmd)

	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "Terdapat error ketika menjalankan program CLI '%s'", err)
		os.Exit(1)
	}
}

func exportToPlainText(pathSource string, pathDest string) {
	if pathDest == "" {
		pathDest = "plaintext.txt"
	}

	openFile, err := os.Open(pathSource)
	if err != nil {
		log.Fatal(err)
	}

	defer openFile.Close()

	writeFile, err := os.Create(pathDest)
	if err != nil {
		log.Fatal(err)
	}

	scanner := bufio.NewScanner(openFile)

	for scanner.Scan() {
		_, err := writeFile.WriteString(scanner.Text())
		if err != nil {
			log.Fatal(err)
		}
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}
}

func exportToJson(pathSource string, pathDest string) {
	if pathDest == "" {
		pathDest = "plainjson.json"
	}

	openFile, err := os.Open(pathSource)
	if err != nil {
		log.Fatal(err)
	}

	defer openFile.Close()

	writeFile, err := os.Create(pathDest)
	if err != nil {
		log.Fatal(err)
	}

	defer writeFile.Close()

	scanner := bufio.NewScanner(openFile)

	var dataLoggings []Log

	for scanner.Scan() {
		errorLogging := scanner.Text()

		tipe := "[error]"
		logging := strings.TrimPrefix(errorLogging, tipe)
		strSlice := strings.Split(logging, " ")

		timeString := strSlice[0] + " " + strSlice[1]

		waktu, err := time.Parse("2006/01/02 15:04:05", timeString)
		if err != nil {
			log.Fatal(err)
		}

		var description string

		for i := 3; i < len(strSlice); i++ {
			if i != 3 {
				description += " "
			}

			description += strSlice[i]
		}

		dataLogging := Log{
			Type:        "[error]",
			DateTime:    waktu,
			LineCode:    strings.TrimRight(strSlice[2], ":"),
			Description: description,
		}

		dataLoggings = append(dataLoggings, dataLogging)
	}

	byteJson, err := json.Marshal(&dataLoggings)
	if err != nil {
		log.Fatal(err)
	}

	_, err = writeFile.Write(byteJson)
	if err != nil {
		log.Fatal(err)
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}
}

func createErrorLog() {
	// create error log
	fileError, err := os.OpenFile("./myerror.log", os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0644)
	if err != nil {
		log.Fatal(err)
	}

	errorLog := log.New(fileError, "[error]", log.LstdFlags|log.Lshortfile|log.Lmicroseconds)
	errorLog.Println("this is error")
}

func createInfoLog() {
	// create info log
	fileInfo, err := os.OpenFile("./myinfo.log", os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0644)
	if err != nil {
		log.Fatal(err)
	}

	infoLog := log.New(fileInfo, "[info]", log.LstdFlags|log.Lshortfile|log.Lmicroseconds)
	infoLog.Println("this is info")
}
