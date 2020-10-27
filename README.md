# System Request Profiler - CLI Tool

This CLI tool will make requests against a specified URL and will print
out the response body. If specified with a flag, it will create a profile
for multiple requests including the request-response cycle time, response sizes,
and any errors returned.  

## Getting Started

To get a local copy of the System Profiler up and running, follow these steps.

### Prerequisites

[Go Compiler](https://golang.org/doc/install)

This project is written in the Go language.  It can be compiled to be run without
a Go installation, but this tool is best used with the Go installed locally as it
will allow for changing the URL you are testing as well as the number of requests
to be made for profile statistics.

### Installation

1. Clone the repo
`git clone https://github.com/JLowe-N/Systems-Profiler.git`
2. Use Go's CLI to run the program.  With the terminal in the directory with main.go:
`go run main.go -url=https://linktree.justinlowen.workers.dev/links -profile=0`

Replace the url flag value with the URL you wish to make requests to.

The default profile flag of 0 will return the response body for inspection.
If the profile flag integer is set to a positive value, multiple requests of that
amount will be made to the specified URL, and profile statistics will be printed.

For details on the program flags, enter the following command in the terminal.
`go run main.go -h`

For documentation, with the terminal in the same directory as main.go:
`go doc`

### Usage

Using the above commands, a profile for requests to a target site can be
made. The profile will give information about the response time, the size of
responses, any errors or HTTP response status error codes, and the % of requests
that are successful.

![Profiled Cloudflare Workers Linktree Project](/Profiles/linktree-justinlowen-workers-dev.JPG "CF Workers Linktree Page") 

I used this tool to compare my [Cloudflare Workers project](https://github.com/JLowe-N/CF-Workers-Linktree) that implements a
[Linktree](https://linktr.ee/) style website to the [Queensland Australia Linktree page](https://linktr.ee/queensland).
A Linktree is a simple page where brands, media platforms, and influences can 
share links to important content with their audience.  The [Cloudflare Workers](https://workers.cloudflare.com/)
project allows developers to deploy serverless JavaScript applications on their
global CDN.  To create the Linktree style page, an HTML template is retrieved from
another Worker, and Cloudflare's HTMLRewriter class enables my application to transform
the HTML as it streams into my Worker before returning it to the user client.

### Example Output
For the given sample size, the page performance looks comparable, although the response
size for the Queensland Linktree page is almost 8 times larger at around 60000 bytes.
Based on Chrome DevTools, the Queensland page has a larger response size due to
additional scripts that appear to be framework/webpack related.  My Cloudflare Worker's
returned page size is mostly attributable to the background image and my avatar.

Both profiles did not detect any errors during the 1000 requests made.

[Profile 1](/Profiles/Profile_linktree-justinlowen-workers-dev.JPG)
```go run main.go -url=https://linktree.justinlowen.workers.dev -profile=1000```
```
Profile (Cloudflare Workers Project)
https://linktree.justinlowen.workers.dev
Number of Requests: 1000
Fastest Time:  132 ms
Slowest Time: 1433 ms
Mean Time: 177 ms
Median Time: 167 ms
Smallest Response: 8414 bytes
Largest Response: 8688 bytes
Request Success: 100%
Error Codes: []
```

[Profile 2](/Profiles/Profile_linktr-ee-queensland.JPG)
```go run main.go -url=https://linktr.ee/queensland -profile=1000```
```
Profile (Linktr.ee Page)
https://linktr.ee/queensland
Number of Requests: 1000
Fastest Time:  136 ms
Slowest Time: 21348 ms
Mean Time: 591 ms
Median Time: 176 ms
Smallest Response: 60267 bytes
Largest Response: 60390 bytes
Request Success: 100%
Error Codes: []
```

Inspection of https://linktree.justinlowen.workers.dev/links
Profile flag set to 0
Page returns JSON response containing the page links
Note: If profile flag is not set, Systems Profiler will only print response body.
```go run main.go -url=https://linktree.justinlowen.workers.dev/links```
```
:::Response Body:::
[
  {
    "name": "Portfolio",
    "url": "https://JLowe-N.github.io"
  },
  {
    "name": "Project: React-based Netflix",
    "url": "https://jlowen-netflix.netlify.app/"
  },
  {
    "name": "Project: Front End UI - Beer Locator",
    "url": "https://justin-lowen.herokuapp.com/punk-api-beer-app"
  },
  {
    "name": "Download My Resume (PDF)",
    "url": "https://github.com/JLowe-N/MyResume/raw/master/Justin%20Lowen%20-%20Software%20Engineer%20-%20Sept%202020%20-%20Public%20Copy.pdf"
  },
  {
    "name": "Learn Something New: How To Brew",
    "url": "http://www.howtobrew.com/"
  }
]
:::End of Response Body:::
```

## License

Distributed under the MIT License. See `LICENSE` for more information.

## Contact
Justin Lowen: 
- [https://jlowe-n.github.io/](https://jlowe-n.github.io/)
- [Justin.G.Lowen@gmail.com](mailto:Justin.G.Lowen@gmail.com)

Project Link: [https://github.com/JLowe-N/Systems-Profiler](https://github.com/JLowe-N/Systems-Profiler)

## Acknowledgements
- [Cloudflare](https://www.cloudflare.com/)
- [Cloudflare Workers](https://workers.cloudflare.com/)







