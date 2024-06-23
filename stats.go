package main

import(
	//"github.com/go-git/go-git/v5"
	//"gopkg.in/src-d/go-git.v4/plumbing"
	"gopkg.in/src-d/go-git.v4"
	"gopkg.in/src-d/go-git.v4/plumbing/object"

	"fmt"
)

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

