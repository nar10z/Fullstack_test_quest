/*
 * @author Nikita Terentyev (nekit.nar10z@gmail.com)
 */

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

		if countUrls <= 0 {
			fmt.Print("No url found")
			continue
		}

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

/*
 * Чтение введенного текста, и разбор его на url
 */
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
