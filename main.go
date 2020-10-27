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

func main() {
	stringURL := flag.String("url", "https://linktree.justinlowen.workers.dev/links", "a string representing a full url")
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
		times := make([]int64, *requestCount)
		responseSizes := make([]int, *requestCount)
		errors := make([]string, *requestCount)
		success := make([]int, *requestCount)

		for i := 0; i < *requestCount; i++ {
			times[i], responseSizes[i], errors[i], success[i] = makeRequest(host, port, path, false)
		}
		sort.Slice(times, func(i, j int) bool {
			return times[i] < times[j]
		})

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
		makeRequest(host, port, path, true)
	}
}

func makeRequest(host string, port string, path string, printResBody bool) (int64, int, string, int) {
	startTime := time.Now()
	var conn net.Conn
	var err error
	if port == "443" {
		conn, err = tls.Dial("tcp", host+":"+port, nil)
	} else if port == "80" {
		conn, err = net.Dial("tcp", host+":"+port)
	}

	if err != nil {
		log.Fatal(err)
		endTime := time.Now()
		timeElapsed := endTime.Sub(startTime).Milliseconds()

		return timeElapsed, 0, err.Error(), 0
	}
	defer conn.Close()
	buf := make([]byte, 0, 4096)
	// Define the request string
	if len(path) > 0 {
		fmt.Fprintf(conn, "GET /"+path+" HTTP/1.1\r\nHost: "+host+"\r\nConnection: Close\r\n\r\n")
	} else {
		fmt.Fprintf(conn, "GET / HTTP/1.1\r\nHost: "+host+"\r\nConnection: Close\r\n\r\n")
	}

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
	s := string(buf)
	result := strings.Split(s, "\r\n\r\n")
	headers := strings.Split(result[0], "\r\n")
	status := strings.Split(headers[0], " ")
	statusCode := status[1]
	testCode, err := strconv.Atoi(statusCode)
	var errorCode string
	if testCode >= 400 {
		errorCode = statusCode
	}

	endTime := time.Now()
	timeElapsed := endTime.Sub(startTime).Milliseconds()

	if printResBody {
		fmt.Println(":::Response Body:::")
		fmt.Println(result[1])
		fmt.Println(":::End of Response Body:::")
	}
	return timeElapsed, len(buf), errorCode, 0

}
