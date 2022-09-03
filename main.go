package main

import (
	"fmt"
	"os"
	"io"
	"strings"
	"sync"
	"net/http"
	"time"
	"github.com/go-xmlfmt/xmlfmt"
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing/object"
)

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func main() {
	//Read file for list of web pages to save
	data, err := os.ReadFile("sites.list")
	check(err)

	content := string(data)
	lines := strings.Split(content, "\n")

	r, err := git.PlainOpen("plutarchs-journal")
	check(err)

	w, err := r.Worktree()
	check(err)

	var wg sync.WaitGroup

	getAndSave := func(index int, lines []string) {
		currentTime := time.Now().Format("2006-01-02T15:04:05.000000")

		line := strings.Split(lines[index], ":::")

		if len(line) != 2 {
			panic("URL file formatted incorrectly")
		}

		fmt.Printf("[%v] - %v - %v\n", currentTime, line[0], line[1])

		response, err := http.Get(line[1])
		check(err)
		defer response.Body.Close()

		body, err := io.ReadAll(response.Body)
		check(err)

		bodyString := string(body)

		prettyBody := xmlfmt.FormatXML(bodyString, "\t", " ")

		err = os.WriteFile(fmt.Sprint("plutarchs-journal/", line[0]), []byte(prettyBody), 0666)
		check(err)		

		_, err = w.Add(line[0])
		check(err)

		commit, err := w.Commit("Changes to file", &git.CommitOptions{
			Author: &object.Signature{
				Name: "Mark Ehresman",
				Email: "435089@gmail.com",
				When: time.Now(),
			},
		})
		check(err)

		_, err = r.CommitObject(commit)
		check(err)

		wg.Done()
	}

	wg.Add(len(lines))

	for index := 0; index < len(lines); index++ {
		go getAndSave(index, lines)
	}

	fmt.Println("Waiting for jobs to finish...")
	wg.Wait()
	fmt.Println("Jobs finished")

	fmt.Println("Pushing site changes...")
	err = r.Push(&git.PushOptions{})
	check(err)
	fmt.Println("Site changes pushed")
}