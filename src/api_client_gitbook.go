package main

import (
	"bytes"
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

func reportErrorFromGitBook(errMsgPreffix string, errCode error) {
	fmt.Fprintf(os.Stderr, "Error: %s: %q. "+
		"You might need to verify the values of your GITBOOK_USER and GITBOOK_PASSWD environment variables.\n",
		errMsgPreffix, errCode)
}

func reportGitBooksOfAuthor(gitbookApi *gitbook.API, gitbookAuthor string) {

	fmt.Fprintf(os.Stderr, "Looking up for all books written by author %q...\n",
		gitbookAuthor)

	// The GitBook client library in Go doesn't seem to have the option to
	// retrieve all books authored by a given user, so we'll do the request
	// here without the struct/array typing of the proper structures
	// returned in the JSON by the GitBook API (see json.RawMessage below,
	// otherwise we would have issues unmarshalling the JSON result if the
	// types and names of the fields in Go don't match what is given in
	// the JSON Restful result from the API -GitBook doesn't give what is
	// the formal declaration of this complex data type of all the "Books"
	// of a given author, so possibly "json.RawMessage" is the fastest way
	// around, avoiding the formal type declaration of the complex data
	// type.)

	type Books struct {
		List json.RawMessage // TODO: convert this json.RawMessage
		//       into an array of structs
		Total int
		Limit int
	}
	var books Books

	// Make the lower level call to the GitBook Restful API to get all the
	// books written by an author:
	//      http://developer.gitbook.com/books/index.html

	_, errRequest := gitbookApi.Client.Get(
		fmt.Sprintf("/author/%s/books", gitbookAuthor),
		nil,
		&books,
	)

	if errRequest != nil {
		reportErrorFromGitBook("querying the GitBook API for author", errRequest)
	} else {
		// fmt.Printf("Raw JSON for Books = %s\n", string(books.List))

		var booksPrettyPr bytes.Buffer
		errIndent := json.Indent(&booksPrettyPr, books.List, "", "  ")

		if errIndent == nil {
			fmt.Printf("Books = %s\n", string(booksPrettyPr.Bytes()))
		} else {
			fmt.Fprintln(os.Stderr, "Error in author's JSON response: ", errIndent)
		}
	}
}

func reportGitBook(gitbookApi *gitbook.API, author string, bookName string, dumpJson bool) {

	fmt.Fprintf(os.Stderr, "Looking up for book %q written by author %q...\n",
		bookName, author)

	book, err := gitbookApi.Book.Get(author + "/" + bookName)
	// book, err := gitbookApi.Books.List()

	if err == nil {
		// No error retrieving the book information from the GitBook API

		if dumpJson {
			// Try to print the entire JSON results of the GitBook metadata
			bookPrettyPr, errMarsh := json.MarshalIndent(book, "", "  ")
			if errMarsh == nil {
				fmt.Printf("Book full metadata = %s\n", bookPrettyPr)
			} else {
				fmt.Fprintln(os.Stderr, "Error in book's JSON response: ", errMarsh)
			}
		} else {
			// Print directly the download links for the book in GitBook in EPUB, PDF, and
			// Mobi document formats
			fmt.Printf("EPUB=%s\n", book.Urls.Download.Epub)
			fmt.Printf("PDF=%s\n", book.Urls.Download.Pdf)
			fmt.Printf("Mobi=%s\n", book.Urls.Download.Mobi)
		}
	} else {
		reportErrorFromGitBook("querying the GitBook API for book", err)
	}
}

func main() {

	var gitbookAuthor string
	var gitbookBook string
	var dumpAllResults bool

	flag.StringVar(&gitbookAuthor, "author", "",
		"The GitBook account which authored the book.")
	flag.StringVar(&gitbookBook, "book", "",
		"The name of the GitBook book belonging to that author.")
	flag.BoolVar(&dumpAllResults, "dump", false,
		"Dump all details of the GitBook book besides URLs to download book. (default: false)")
	flag.Parse()

	if gitbookAuthor == "" {
		fmt.Fprintln(os.Stderr,
			"The GitBook 'author' argument must be provided in the command-line.")
		os.Exit(1)
	}

	api := createGitBookAPIObj()

	if gitbookBook != "" {
		reportGitBook(api, gitbookAuthor, gitbookBook, dumpAllResults)
	} else {
		// only the GitBook author was given, without a specific book:
		// report all the books written by the given author
		reportGitBooksOfAuthor(api, gitbookAuthor)
	}

}
