# API_GitBook

List the URL addresses to download the documentation in epub, pdf or mobi format of book projects hosted by GitBooks.com, using the GitBooks API's client-library in Go.

This client is somehow similar to the [API_ReadTheDocs project](https://github.com/je-nunez/API_ReadTheDocs) which lists (and downloads) books hosted by [ReadTheDocs.org](https://ReadTheDocs.org/), but in this case, for [GitBook.com](https://www.gitbook.com/).

# WIP

This is the very first draft, which is a *work in progress*. The implementation is *incomplete* and subject to change. The documentation can be inaccurate.

# How to Use it

After you have [compiled the source code](#how-to-compile), you may need to setup your [GitBook account](https://www.gitbook.com/) to authenticate (some books require this). (**Note**: the GitBook account is **not** the same as the GitHub account, so far, and you may need to create your free GitBook account.)

      # some books require this
      export GITBOOK_USER=<your-gitbook-username>
      export GITBOOK_PASSWD=<your-gitbook-password>

Call the client program:

      ./api_client_gitbook -author "<put-author-here>" [-book "<put-optional-book-entry-here>"]

The `-book "<put-optional-book-entry-here>"` is optional: if omitted, this client will report all books written in GitBook by the given author.

For example:

      ./api_client_gitbook -author 0xax -book linux-insides
       
      ./api_client_gitbook -author unbug -book react-native-training
       
      # to report all books written by a given author, for example:
       
      ./api_client_gitbook -author wxdublin

# How to Compile

This project is in GoLang and requires the [GitBook API Go client library](https://github.com/GitbookIO/go-gitbook-api).

1. [Install GoLang for your architecture](https://golang.org/doc/install)

2. Set the `GOPATH` environment variable to the base directory where you downloaded this project (or another subdirectory):

         GOPATH=<put/here/base/directory/of/git/clone/of/this/project>
         export GOPATH

3. Install the GitBook API Go client library:

         go get github.com/GitbookIO/go-gitbook-api

4. For further installing the compiled executable, set the `GOBIN` to the proper directory where to write the executable. As a very simple example:

         GOBIN="$GOPATH/bin"
         export GOBIN

5. Compile the source code:

         go build src/api_client_gitbook.go
         go install src/api_client_gitbook.go

# Inspiration

This GitBook command-line client is inspired from the GitBook command-line example at [https://github.com/GitbookIO/go-gitbook-api](https://github.com/GitbookIO/go-gitbook-api).
