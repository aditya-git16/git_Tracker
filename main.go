package main

import(
	"flag"
)

func main() {
    var folder string
    var email string
	// & in &folder refers to the memory address of the folder variable
	// First input requires pointer to variable to be updated
	// Pointer to storage var , flag name , default val , description of flag
    flag.StringVar(&folder, "add", "", "add a new folder to scan for Git repositories")
    flag.StringVar(&email, "email", "your@email.com", "the email to scan")
	// Parses the flag passed in by the user and allocates values to the variables
	// Defined after all the flags are defined 
    flag.Parse()

    if folder != "" {
        scan(folder)
        return
    }

    stats(email)
}
