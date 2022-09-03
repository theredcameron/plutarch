package main

import (
	"fmt"
	"os"
	"strings"
)

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func main() {
	//Read file for list of web pages to save
	data, err := os.ReadFile("sites")
	check(err)

	content := string(data)
	lines := strings.Split(content, "\n")

	for index := 0; index < len(lines); index++ {
		fmt.Println(lines[index])
		
	}
		//ASYNC - For each web page, get contents and save.
			//Get page content
			//If page file exists, overrite.
			//Else, create new file
			//Git add and commit
			//Push value through channel to signify task finished
	//When all tasks are finished, Git push repo
}