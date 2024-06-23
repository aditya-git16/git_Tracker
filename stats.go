package main

import(
	"github.com/go-git/go-git/v5"
	"gopkg.in/src-d/go-git.v4/plumbing"

	"fmt"
)

// processRepositories given a user email , returns the commits in the last 6 months 
func processRepositories(email string) map[int]int{
	
}

// calculates and prints the stats
func stats(email string)  {
	commits := processRepositories(email)
	printCommitStats(commits)
}

