package main

import (
	structs "./Structures"
	"bufio"
	"fmt"
	"net/url"
	"os"
	"strings"
	"time"
)

func main() {
	maxLenChannel := 5
	reader := bufio.NewReader(os.Stdin)
	fmt.Println("Paste URL:")

	for {
		fmt.Print("=> ")

		total := 0

		chUrls := make(chan *structs.Request, maxLenChannel)
		quit := make(chan int, 1)

		text, _ := reader.ReadString('\n')
		splitStrings := readUrls(text)
		countUrls := len(splitStrings)

		go func(q chan int, arr []string) {
			for {
				if countUrls <= 0 {
					quit <- 0
					return
				}

				time.Sleep(150 * time.Millisecond)
			}
		}(quit, splitStrings)

		for _, urlString := range splitStrings {
			go initNewRequest(chUrls, urlString)
		}

		allow := true
		for allow {
			select {
			case r := <-chUrls:
				urlString := r.Url()
				countWord := r.CountGo()

				str := "Count for " + urlString + ":"
				fmt.Println(str, countWord)
				total += countWord

				countUrls--
			case <-quit:
				fmt.Println("Total: ", total)
				allow = false
				break
			}
		}
	}
}

func readUrls(input string) []string {
	readString := strings.ReplaceAll(input, "\r", "")
	splitStrings := strings.Split(readString, "\\n")

	var arr []string
	for _, value := range splitStrings {
		value = strings.TrimSpace(value)

		if len(value) > 0 {
			if isValidUrl := isValidUrl(value); isValidUrl {
				arr = append(arr, value)
			}
		}
	}

	return arr
}

func isValidUrl(urlString string) bool {
	_, err := url.ParseRequestURI(urlString)
	if err != nil {
		return false
	}

	u, err := url.Parse(urlString)
	if err != nil || u.Scheme == "" || u.Host == "" {
		return false
	}

	return true
}

func initNewRequest(channel chan *structs.Request, urlString string) {
	request := structs.NewRequest(urlString)

	err := request.Send()
	if err != nil {
		fmt.Println("Error in Send for", urlString)
	}

	request.CountWord()

	channel <- request
}

// https://yourbasic.org/golang/split-string-into-slice/
// https://yourbasic.org/golang/split-string-into-slice/ \n https://www.google.com/search?q=explode+string+golang&oq=explode+string+golang&aqs=chrome..69i57j0l3.4639j1j7&sourceid=chrome&ie=UTF-8
// https://yourbasic.org/golang/split-string-into-slice/ \n https://www.google.com/search?q=explode+string+golang&oq=explode+string+golang&aqs=chrome..69i57j0l3.4639j1j7&sourceid=chrome&ie=UTF-8 \n https://golangcode.com/how-to-check-if-a-string-is-a-url/
// https://yourbasic.org/golang/split-string-into-slice/ \n https://www.google.com/search?q=explode+string+golang&oq=explode+string+golang&aqs=chrome..69i57j0l3.4639j1j7&sourceid=chrome&ie=UTF-8 \n https://golangcode.com/how-to-check-if-a-string-is-a-url/ \n https://tour.golang.org/concurrency/6 \n https://golangbot.com/goroutines/ \n https://stackoverflow.com/questions/29898400/import-struct-from-another-package-and-file-golang
