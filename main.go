/*This package is a CLI tool that will make requests against a
specified URL and will print out the response body. If if the profile flag is
set with an integer, it will create a profile for multiple requests including
the request-response cycle time, response sizes, and any errors returned.*/
package main

import (
	"crypto/tls"
	"flag"
	"fmt"
	"io"
	"log"
	"net"
	"sort"
	"strconv"
	"strings"
	"time"
)

/* Main handles parsing the flags, including the url and request count for the
profile.  URL is parsed into scheme (for setting TCP port to default), host, and
path.  If profile flag is set with an positive integer, requests will be made
with the makeRequest function and profile statistics will be calculated and
printed in main.  If profile flag is not set, then a single makeRequest call is
made which will also print the response body.*/
func main() {
	stringURL := flag.String("url", "https://linktree.justinlowen.workers.dev/links", "a string representing a full url, Example: https://www.cloudflare.com")
	requestCount := flag.Int("profile", 0, "an integer specifying number of requests for profile")

	flag.Parse()

	parsedURL := strings.SplitN(*stringURL, "//", 2)
	var port string
	if parsedURL[0] == "https:" {
		port = "443"
	} else if parsedURL[0] == "http:" {
		port = "80"
	} else {
		log.Fatal("Please specify scheme as either http: or https: in URL.")
	}
	parsedURL = strings.SplitN(parsedURL[1], "/", 2)
	host := parsedURL[0]
	var path string
	if len(parsedURL) > 1 {
		path = parsedURL[1]
	} else {
		path = ""
	}

	if *requestCount > 0 {
		// Initialize request statistics array
		times := make([]int64, *requestCount)
		responseSizes := make([]int, *requestCount)
		errors := make([]string, *requestCount)
		success := make([]int, *requestCount)

		// Perform requests, then sort times for finding median
		for i := 0; i < *requestCount; i++ {
			times[i], responseSizes[i], errors[i], success[i] = makeRequest(host, port, path, false)
		}
		sort.Slice(times, func(i, j int) bool {
			return times[i] < times[j]
		})

		// Calculate profile statistics
		fastestTime := times[0]
		slowestTime := times[0]
		meanTime := times[0]
		medianTime := times[len(times)/2]
		percentSuccess := float64(success[0])
		smallestResponse := responseSizes[0]
		largestResponse := responseSizes[0]
		for i := 1; i < *requestCount; i++ {
			if times[i] < fastestTime {
				fastestTime = times[i]
			}
			if times[i] > slowestTime {
				slowestTime = times[i]
			}
			meanTime += times[i]

			percentSuccess += float64(success[i])

			if responseSizes[i] < smallestResponse {
				smallestResponse = responseSizes[i]
			}
			if responseSizes[i] > largestResponse {
				largestResponse = responseSizes[i]
			}
		}
		meanTime = meanTime / int64(len(times))
		percentSuccess = percentSuccess / float64(len(success)) * 100

		// Final profile printout
		fmt.Println("Profile")
		fmt.Println(*stringURL)
		fmt.Println("Number of Requests:", *requestCount)
		fmt.Println("Fastest Time: ", fastestTime, "ms")
		fmt.Println("Slowest Time:", slowestTime, "ms")
		fmt.Println("Mean Time:", meanTime, "ms")
		fmt.Println("Median Time:", medianTime, "ms")
		fmt.Println("Smallest Response:", smallestResponse, "bytes")
		fmt.Println("Largest Response:", largestResponse, "bytes")
		fmt.Printf("Request Success: %0.0f %% \n", percentSuccess)
		fmt.Println("Error Codes:", errors)
	} else {
		// If no profile, make a sample request and print request body.
		makeRequest(host, port, path, true)
	}
}

/*makeRequest takes in the parsed URL flag's host, port, and path variables and
makes either an http or https request to the URL.  It also takes the
printResBody argument which is set to true if a profile is not being generated and
will print the response body.  makeRequest returns profile datapoints to be
added to arrays in the main function. Each request's time, size, response error
codes, and a success/failure integer is returned.*/
func makeRequest(host string, port string, path string, printResBody bool) (int64, int, string, int) {
	// Establish connection, if connection fails program will exit and error will be printed
	var conn net.Conn
	var err error
	startTime := time.Now()
	if port == "443" {
		conn, err = tls.Dial("tcp", host+":"+port, nil)
	} else if port == "80" {
		conn, err = net.Dial("tcp", host+":"+port)
	}

	if err != nil {
		log.Fatal(err)
	}

	//  send Request
	defer conn.Close()
	if len(path) > 0 {
		fmt.Fprintf(conn, "GET /"+path+" HTTP/1.1\r\nHost: "+host+"\r\nConnection: Close\r\n\r\n")
	} else {
		fmt.Fprintf(conn, "GET / HTTP/1.1\r\nHost: "+host+"\r\nConnection: Close\r\n\r\n")
	}

	// read Response
	buf := make([]byte, 0, 4096)
	for {
		tmp := make([]byte, 256)
		n, err := conn.Read(tmp)
		if err != nil {
			if err != io.EOF {
				fmt.Println("read error:", err)
				endTime := time.Now()
				timeElapsed := endTime.Sub(startTime).Milliseconds()
				return timeElapsed, len(buf), err.Error(), 0
			}
			break
		}
		buf = append(buf, tmp[:n]...)

	}
	endTime := time.Now()
	timeElapsed := endTime.Sub(startTime).Milliseconds()
	s := string(buf)

	// Split response into component parts
	result := strings.Split(s, "\r\n\r\n")
	headers := strings.Split(result[0], "\r\n")
	status := strings.Split(headers[0], " ")
	statusCode := status[1]

	// Check for status indicating error in request/response
	testCode, err := strconv.Atoi(statusCode)
	var errorCode string
	var successValue int
	if testCode >= 400 {
		errorCode = statusCode
		successValue = 0
	} else {
		errorCode = ""
		successValue = 1
	}

	// If not making a profile, print response's body for inspection.
	if printResBody {
		fmt.Println(":::Response Body:::")
		fmt.Println(result[1])
		fmt.Println(":::End of Response Body:::")
	}

	return timeElapsed, len(buf), errorCode, successValue

}
