package main

import (
	"flag"
	"fmt"
	"os"

	"encoding/json"

	"github.com/GitbookIO/go-gitbook-api"
)

func createGitBookAPIObj() *gitbook.API {

	var api *gitbook.API

	// Obtain the GitBook user with which this client will authenticate with
	// the GitBook API -this GitBook user does not need to be the author of
	// the book- from the environment variable GITBOOK_USER

	gitbookUser := os.Getenv("GITBOOK_USER")

	if gitbookUser != "" {
		gitbookPasswd := os.Getenv("GITBOOK_PASSWD")

		// immediately clear the environment variables in this process
		// before even attempting connecting to the GitBook API, so they
		// can't be seen through the "/proc/<pid>/environ" file.
		os.Setenv("GITBOOK_USER", "")
		os.Setenv("GITBOOK_PASSWD", "")

		// try to connect and authenticate to the GitBook API
		api = gitbook.NewAPI(gitbook.APIOptions{

			// Hit API with a specific user
			Username: gitbookUser,
			Password: gitbookPasswd,
		})
	} else {
		// no GitBook user known: connect to the GitBook API anonymously
		api = gitbook.NewAPI(gitbook.APIOptions{})
	}

	return api
}

func reportGitBook(gitbookApi *gitbook.API, gitbookItem string) {

	book, err := gitbookApi.Book.Get(gitbookItem)
	// book, err := gitbookApi.Books.List()

	if err == nil {
		// No error retrieving the book information from the GitBook API

		// Try to print the entire JSON results of the GitBook metadata
		bookPrettyPr, errMarsh := json.MarshalIndent(book, "", "  ")
		if errMarsh == nil {
			fmt.Printf("Book metadata = %s\n", bookPrettyPr)
		} else {
			fmt.Fprintln(os.Stderr, "Error in book's JSON response: ", errMarsh)
		}

		// Print directly the download links for the book in GitBook in EPUB, PDF, and
		// Mobi document formats
		fmt.Printf("EPUB=%s\n", book.Urls.Download.Epub)
		fmt.Printf("PDF=%s\n", book.Urls.Download.Pdf)
		fmt.Printf("Mobi=%s\n", book.Urls.Download.Mobi)
	} else {
		fmt.Fprintf(os.Stderr, "Error querying the GitBook API for book: %q\n", err)
	}
}

func main() {

	var gitbookAuthor string
	flag.StringVar(&gitbookAuthor, "author", "",
		"The GitBook account which authored the book.")
	var gitbookBook string
	flag.StringVar(&gitbookBook, "book", "",
		"The name of the GitBook book belonging to that author.")
	flag.Parse()

	fmt.Printf("Looking for book %q written by author %q\n",
		gitbookBook, gitbookAuthor)
	var item string
	if gitbookBook != "" {
		item = gitbookAuthor + "/" + gitbookBook
	} else {
		// TODO: To-Check: this call to the GitBook /book/ Restful API
		// entry may fail
		item = gitbookAuthor
	}

	api := createGitBookAPIObj()

	reportGitBook(api, item)
}
