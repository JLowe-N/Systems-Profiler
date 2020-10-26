package main

import (
	"crypto/tls"
	"flag"
	"fmt"
	"log"
	"strings"
	"time"
)

func main() {
	stringURL := flag.String("url", "https://linktree.justinlowen.workers.dev/", "a string representing a full url")
	// requestCount := flag.Int("profile", 3, "an integer specifying number of requests for profile")

	flag.Parse()

	var port string
	parsedURL := strings.SplitN(*stringURL, "//", 2)
	if parsedURL[0] == "https:" {
		port = "443"
	} else if parsedURL[0] == "http:" {
		port = "80"
	} else {
		log.Fatal("Please specify scheme as either http: or https: in URL.")
	}
	parsedURL = strings.SplitN(parsedURL[1], "/", 2)
	host := parsedURL[0]
	path := parsedURL[1]
	fmt.Println(parsedURL)
	fmt.Println(port)
	fmt.Println(host)
	fmt.Println(path)
	makeRequest(host, port, path)
}

// 	fmt.Fprintf(conn, "GET /links HTTP/1.1\r\nHost: linktree.justinlowen.workers.dev\r\nConnection: Close\r\n\r\n")
// 	status, err := bufio.NewReader(conn).ReadString('\n')
// 	fmt.Println(status)
// 	fmt.Println(err)

// 	fmt.Println(*url)
// 	fmt.Println("Number of Requests:", *requestCount)
// 	fmt.Println("Fastest Time:")
// 	fmt.Println("Slowest Time:")
// 	fmt.Println("Mean Time:")
// 	fmt.Println("Median Time:")
// 	fmt.Println("Request Success %:")
// 	fmt.Println("Error Codes:")
// 	fmt.Println("Smallest Response:")
// 	fmt.Println("Largest Response:")
// }

func makeRequest(host string, port string, path string) (int, time.Duration, int, string) {
	startTime := time.Now()
	conn, err := tls.Dial("tcp", host+":"+port, nil)
	if err != nil {
		log.Fatal(err)
		endTime := time.Now()
		timeElapsed := endTime.Sub(startTime)

		return 0, timeElapsed, 0, err.Error()
	}
	defer conn.Close()
	// buf := make([]byte, 0, 4096)
	// Define the request string
	if len(path) > 0 {
		fmt.Fprintf(conn, "GET /"+path+" HTTP/1.1\r\nHost: "+host+"\r\nConnection: Close\r\n\r\n")
	} else {
		fmt.Fprintf(conn, "GET / HTTP/1.1\r\nHost: "+host+"\r\nConnection: Close\r\n\r\n")
	}
	fmt.Println("Success")
	return 0, 0, 0, ""

}
