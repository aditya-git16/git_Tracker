package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"os"
	"os/user"
	"strings"
)

// openFile opens a file located at 'filePath' and creates it if it doesn't exist
func openFile(filepath string) *os.File{
	f , err := os.OpenFile(filepath , os.O_APPEND | os.O_WRONLY , 0755)
	if err != nil {
		if os.IsNotExist(err){
			// file does not exist
			_ , err = os.Create(filepath)
			if err != nil {
				panic(err)
			}
		}else{
			panic(err)
		}
	}
	return f	
}

// parseFileLinesToSlice gets the content of a file line by line and returns a slice of it
func parseFileLinesToSlice(filePath string) []string {
	f := openFile(filePath)
	defer f.Close()

	
	var lines []string
	// NewScanner returns a scan object to scan the file passed into it
	scanner := bufio.NewScanner(f)
	for scanner.Scan(){
		lines = append(lines, scanner.Text())
	}
	if err := scanner.Err(); err != nil {
		if err != io.EOF{
			panic(err)
		}
	}

	return lines
}

// joinSlice adds the element of the 'new' slice into existing slice , only if its not already there
func joinSlice(new []string , exisitng []string) []string {
	for _ , i := range new {
		if !sliceContains(exisitng, i){
			exisitng = append(exisitng, i)
		}
	}
	return exisitng
}

// sliceContains to check if 'slice' contains 'value'
func sliceContains(slice []string , value string) bool{
	for _, v := range slice{
		if v == value{
			return true
		}
	}
	return false
}

// dumpStringToFile writes the content to the file in 'filepath'
func dumpStringSliceToFile(repos []string, filePath string) {
	content := strings.Join(repos,"\n")
	os.WriteFile(filePath, []byte(content) , 0755)
}

// addNewSliceElementsToFile given a slice of string paths stores them in a file
func addNewSliceElementsToFile(filePath string , newRepos []string){
	existingRepos := parseFileLinesToSlice(filePath)
	repos := joinSlice(newRepos, existingRepos)
	dumpStringSliceToFile(repos , filePath)
}


// Gives the path of the dotfile which has thew database of the repos path
func getDotFiles() string{
	usr ,err := user.Current()
	if err != nil {
		log.Fatal(err)
	}

	dotFile := usr.HomeDir + "/.gogitlocalstats"
	return dotFile
}

func recursiveScanFolder(folder string) []string {
	return scanGitFolders(make([]string, 0) , folder)
}

// Scans a list of subfolders of a folder ending with '.git'
// Searches recursively
func scanGitFolders(folders []string, folder string) []string{
	// trim the last '/' form file name
	folder = strings.TrimSuffix(folder,"/")

	// os.Open opens the directory?
	f, err := os.Open(folder)
	if err != nil {
		log.Fatal(err)
	}

	// -1 in Readdir means we want to read all the contents of the directory
	files , err := f.Readdir(-1)
	f.Close()
	if err != nil {
		log.Fatal(err)
	}

	var path string

	for _, file := range files{
		if file.IsDir(){
			path = folder + "/" + file.Name()
			if file.Name() == ".git"{
				path = strings.TrimSuffix(path, "/.git")
				fmt.Println(path)
				folders =append(folders, path)
				continue
			}
			if file.Name() == "vendor" || file.Name() == "node_modules"{
				continue
			}
			folders = scanGitFolders(folders, path)
		}
	}
	return folders
}

// Scans a new folder for git repository
func scan(folder string){
	fmt.Printf("Found folders:\n\n")
	repositories := recursiveScanFolder(folder)
	filePath := getDotFiles()
	addNewSliceElementsToFile(filePath, repositories)
	fmt.Printf("\n\nSuccessfully added\n\n")
}