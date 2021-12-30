package main

import (
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"time"
)

// program that takes a list of URL's and will fetch them all concurrently
// theyll all start same time rather than one after the other finishes
func main() {
	start := time.Now()

	// a channel allows communication from one goroutine to another
	// 'main' creates a channel of strings to recieve the results
	mainChannel := make(chan string)

	// for each url in the arguments a new goroutine is created that fetches the url asynchronously
	for _, url := range os.Args[1:] { //args[1:]
		go fetch(url, mainChannel) // go-routine starts here
	}

	// recieve results to the main channel and print them
	for range os.Args[1:] {
		fmt.Println(<-mainChannel)
	}
	fmt.Printf("%.2fs taken for all calls to finish. CONCURRENT EXECUTION \n", time.Since(start).Seconds())

	//this is the sequential execution version
	for _, url := range os.Args[1:] {
		start := time.Now()

		response, err := http.Get(url)
		if err != nil {
			fmt.Fprint(os.Stderr, "fetch: %v\n", err)
			os.Exit(1)
		}
		nbytes, err := io.Copy(ioutil.Discard, response.Body)
		response.Body.Close() // prevents resource leaking even though no longer user
		if err != nil {
			fmt.Fprintf(os.Stderr, "Fetch: reading %s: %v\n", url, err)
			os.Exit(1)
		}
		secs := time.Since(start).Seconds()

		fmt.Printf("%.2fs %7d %s\n", secs, nbytes, url)
	}
	fmt.Printf("%.2fs taken for all calls to finish. SEQUENTIAL EXECUTION \n\n", time.Since(start).Seconds())

}

func fetch(url string, ch chan<- string) {
	start := time.Now()
	response, err := http.Get(url)
	if err != nil { //can re-use same err the whole way down
		ch <- fmt.Sprint(err) // send to channel
	}

	nbytes, err := io.Copy(ioutil.Discard, response.Body)
	response.Body.Close() // prevents resource leaking even though no longer user
	if err != nil {
		ch <- fmt.Sprintf("while reading %s: %v", url, err)
		return
	}
	secs := time.Since(start).Seconds()
	ch <- fmt.Sprintf("%.2fs %7d %s", secs, nbytes, url) // once all error points passed we return
}
