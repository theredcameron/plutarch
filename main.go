package main

import (
	"fmt"
	"os"
	"io"
	"strings"
	"sync"
	"flag"
	"time"
	"net/http"
	"strconv"
	"hash/fnv"
	"github.com/go-xmlfmt/xmlfmt"
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing/object"
	"gitlab.com/435089/go-logger"
)

var logger *golang_logger.Logger

func LogError(e error) {
	if logger == nil {
		panic(e)
	}
	
	if e != nil{
		logger.Log(e.Error())
		panic("Error found when running program. See logs.")
	}
}

func PassiveLogError(e error) {
	if logger == nil {
		panic(e)
	}

	if e != nil{
		logger.Log(e.Error())
	}
}

func main() {
	var err error
	logger, err = golang_logger.CreateLogger("logs")
	LogError(err)

	inputInterval := flag.Int("interval", 10, "The time interval in seconds for saving pages")
	authorEmail := flag.String("author_email", "", "The email address of the author")
	authorName := flag.String("author_name", "", "The name of the author")
	plutarchJournalPath := flag.String("journal_path", "", "The path of the repository that will hold the journal")
	flag.Parse()
	
	if *inputInterval < 10 {
		LogError(fmt.Errorf("Cannot use interval that is less than 10 seconds"));
	}

	intervalString := fmt.Sprintf("%vs", *inputInterval)
	interval, err := time.ParseDuration(intervalString)
	LogError(err)

	executeTask := func(authorEmail, authorName string) {
		logger.Log("Starting save sequence")
		//Read file for list of web pages to save
		data, err := os.ReadFile("sites-list")
		LogError(err)

		content := string(data)
		lines := strings.Split(content, "\n")

		var r *git.Repository

		if _, err = os.Stat("./plutarchs-journal/.git"); os.IsNotExist(err) {
			r, err = git.PlainClone("./plutarchs-journal", false, &git.CloneOptions{
				URL:	*plutarchJournalPath,
			})
			LogError(err)
		} else {
			r, err = git.PlainOpen("./plutarchs-journal")
			LogError(err)
		}

		w, err := r.Worktree()
		LogError(err)

		var wg sync.WaitGroup

		getAndSave := func(index int, lines []string) {
			logger.Log(lines[index])

			response, err := http.Get(lines[index])
			if err != nil {
				PassiveLogError(err)
				return
			}
			defer response.Body.Close()

			body, err := io.ReadAll(response.Body)
			PassiveLogError(err)

			bodyString := string(body)

			prettyBody := lines[index] + "\n\n\n" + xmlfmt.FormatXML(bodyString, "\t", " ")

			hasher := fnv.New64a()
			hasher.Write([]byte(lines[index]))
			fileHash := hasher.Sum64()
			fileName := strconv.FormatUint(fileHash, 10)

			err = os.WriteFile(fmt.Sprint("plutarchs-journal/", fileName), []byte(prettyBody), 0666)
			PassiveLogError(err)		

			_, err = w.Add(".")
			PassiveLogError(err)

			commit, err := w.Commit("Changes to file", &git.CommitOptions{
				Author: &object.Signature{
					Name: authorName,
					Email: authorEmail,
					When: time.Now(),
				},
			})
			PassiveLogError(err)

			_, err = r.CommitObject(commit)
			PassiveLogError(err)

			wg.Done()
		}

		wg.Add(len(lines))

		for index := 0; index < len(lines); index++ {
			go getAndSave(index, lines)
		}

		logger.Log("Waiting for jobs to finish...")
		wg.Wait()
		logger.Log("Jobs finished")

		logger.Log("Pushing site changes...")
		err = r.Push(&git.PushOptions{})
		PassiveLogError(err)
		logger.Log("Site changes pushed")
		logger.Log("Finished save sequence")
	}

	executeTask(*authorEmail, *authorName)

	for _ = range time.Tick(interval) {
		executeTask(*authorEmail, *authorName)
	}
}