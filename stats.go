package main

import (
	//"github.com/go-git/go-git/v5"
	//"gopkg.in/src-d/go-git.v4/plumbing"
	"time"

	"gopkg.in/src-d/go-git.v4"
	"gopkg.in/src-d/go-git.v4/plumbing/object"

	"fmt"
)

const outOfRange= 99999
const daysInLastSixMonths= 183

// calcOffset determines and returns the amount of days missign to fill
// the last row of the stats graph

func calcOffset() int {
	var offset int
	weekday := time.Now().Weekday()
	switch weekday{
	case time.Sunday:
		offset = 7
	case time.Monday:
		offset = 6
	case time.Tuesday:
		offset = 5
	case time.Wednesday:
		offset = 4
	case time.Thursday:
		offset = 3
	case time.Friday:
		offset = 2
	case time.Saturday:
		offset = 1
	}
	return offset
}

// getBeginningOfDay given a time. Time calculates the start time of that day
func getBeginningOfDay(t time.Time) time.Time{
	year , month , day := t.Date()
	startOfDay := time.Date(year,month,day,0,0,0,0,t.Location())
	return startOfDay
}

// countDaysSinceDate counts how many days passed since the passed 'date'
func countDaysSinceDate(date time.Time) int{
	days :=0
	now := getBeginningOfDay(time.Now())
	for date.Before(now){
		date = date.Add(time.Hour *24)
		days ++
		if days > daysInLastSixMonths{
			return outOfRange
		}
	}
	return days
}

// fillCommits given a repository found in 'path' , gets the commits and
// puts them in to the 'commits' map , returning it when completed
func fillCommits(email string , path string , commits map[int]int) map[int]int{
	// instantiate a git repo object from path
	// repo is of Repository type
	repo , err := git.PlainOpen(path)
	if err != nil {
		panic(err)
	}

	// get the HEAD ref
	ref , err := repo.Head()
	if err != nil {
		panic(err)
	}

	// get the commits history starting from HEAD
	// Log method uses a ref to LogOptions to return the commit history
	iterator , err := repo.Log(&git.LogOptions{From: ref.Hash()})
	if err != nil {
		panic(err)
	}
	
	// iterate the commits
	offset := calcOffset()
	err = iterator.ForEach(func(c *object.Commit) error{
		daysAgo := countDaysSinceDate(c.Author.When) + offset

		if c.Author.Email != email {
			return nil
		}

		if daysAgo != outOfRange {
			commits[daysAgo]++
		}

		return nil
	})

	if err != nil {
		panic(err)
	}
	return commits
}
// processRepositories given a user email , returns the commits in the last 6 months 
func processRepositories(email string) map[int]int{
	filePath := getDotFiles()
	repos := parseFileLinesToSlice(filePath)
	daysInMap := daysInLastSixMonths

	commits := make(map[int]int , daysInMap)
	for i := daysInMap ; i>0 ; i--{
		commits[i] = 0
	}

	for _ , path := range repos{
		commits = fillCommits(email , path , commits)
	}
	return commits
}

// calculates and prints the stats
func stats(email string)  {
	commits := processRepositories(email)
	printCommitStats(commits)
}

