package main

import (
		"golang.org/x/net/proxy"
		"net/http"
		"net/url"
		"fmt"
		"os"
	  "strconv"
	  "time"
	  "log"
		"sync"
	  "bufio"
	  "flag"
	  "strings"
)

var NotFound, Found, Forbidden, Other int = 0, 0, 0, 0


func fatalf(fmtStr string, args interface{}) {
	fmt.Fprintf(os.Stderr, fmtStr, args)
	os.Exit(-1)
}

func MakeTorRequest(target string, verb bool) {
    //start := time.Now()

    // Create a transport that uses Tor Browser's SocksPort.  If
    // talking to a system tor, this may be an AF_UNIX socket, or
    // 127.0.0.1:9050 instead.
    tbProxyURL, err := url.Parse("socks5://127.0.0.1:9150")
    if err != nil {
        fatalf("Failed to parse proxy URL: %v\n", err)
    }

    // Get a proxy Dialer that will create the connection on our
    // behalf via the SOCKS5 proxy.  Specify the authentication
    // and re-create the dialer/transport/client if tor's
    // IsolateSOCKSAuth is needed.
    tbDialer, err := proxy.FromURL(tbProxyURL, proxy.Direct)
    if err != nil {
        fatalf("Failed to obtain proxy dialer: %v\n", err)
    }

    // Make a http.Transport that uses the proxy dialer, and a
    // http.Client that uses the transport.
    tbTransport := &http.Transport{Dial: tbDialer.Dial}
    client := &http.Client{Transport: tbTransport}

    // Example: Fetch something.  Real code will probably want to use
    // client.Do() so they can change the User-Agent.
    resp, resperr := client.Get(target)
		if resperr != nil {
        fmt.Printf("\033[31m%s \033[0m\n", resperr)
    } else{
		    //secs := time.Since(start).Seconds()
		    //body, _ := ioutil.ReadAll(resp.Body)
		    //ch <- fmt.Sprintf("%.2f elapsed with response length: %d %s", secs, len(body), target)
		    status_code := resp.StatusCode
		    file, err := os.OpenFile("result.txt", os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0600)
		    if err != nil {
		        log.Fatal("Cannot create file", err)
		    }
		    if status_code == 401{
		        Other += 1
		        fmt.Fprintf(file, "%s %s\n", target, strconv.Itoa(status_code) )
		        if verb == true {
		            log.Printf("%s HTTP Response Code %s\n", target, strconv.Itoa(status_code) )
		        }
		    } else if status_code == 200{
		        Found += 1
		        fmt.Fprintf(file, "%s %s\n", target, strconv.Itoa(status_code) )
		        if verb == true {
		            log.Printf("%s HTTP Response Code %s\n", target, strconv.Itoa(status_code) )
		        }
		    } else if status_code == 403{
		  	    Forbidden = Forbidden + 1
		    } else if status_code == 404{
		        NotFound += 1
		    } else{
					  Other += 1
		    }
		    defer file.Close()
				defer resp.Body.Close()
	  }
}

func MakeRequest(target string, verb bool) {
    //start := time.Now()
    resp, resperr := http.Get(target)
		if resperr != nil {
        fmt.Printf("\033[31m%s \033[0m\n", resperr)
    } else{
		    //secs := time.Since(start).Seconds()
		    //body, _ := ioutil.ReadAll(resp.Body)
		    //ch <- fmt.Sprintf("%.2f elapsed with response length: %d %s", secs, len(body), target)
		    status_code := resp.StatusCode
		    file, err := os.OpenFile("result.txt", os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0600)
		    if err != nil {
						fmt.Fprintf(file, "Cannot create file: %s", err)
		        //log.Fatal("Cannot create file", err)
		    }
		    if status_code == 401{
		        Other += 1
		        fmt.Fprintf(file, "%s %s\n", target, strconv.Itoa(status_code) )
		        if verb == true {
		            log.Printf("%s HTTP Response Code %s\n", target, strconv.Itoa(status_code) )
		        }
		    } else if status_code == 200{
		        Found += 1
		        fmt.Fprintf(file, "%s %s\n", target, strconv.Itoa(status_code) )
		        if verb == true {
		            log.Printf("%s HTTP Response Code %s\n", target, strconv.Itoa(status_code) )
		        }
		    } else if status_code == 403{
		  	    Forbidden = Forbidden + 1
		    } else if status_code == 404{
		        NotFound += 1
		    } else{
					  Other += 1
		    }
		    defer file.Close()
		    defer resp.Body.Close()
		}
}

func readLines(baseURL, path string) ([]string, error) {
    file, err := os.Open(path)
    if err != nil {
        return nil, err
    }
    defer file.Close()

    var lines []string
    scanner := bufio.NewScanner(file)
    for scanner.Scan() {
        target := baseURL + "/" + scanner.Text()
        lines = append(lines, target)
    }
    return lines, scanner.Err()
}

func main() {
    start := time.Now()

    URL := flag.String("URL", "https://facebook.com", "Base URL")
    FNAME := flag.String("FNAME", "directorylist.txt", "Path to file of keywords for scanning")
    V := flag.Bool("V", true, "Verbose Mode, Prints 200 and 401 codes to the screen.")
    TOR := flag.Bool("TOR", false, "Make GET request over Tor (or not). Currently set to use Tor Browser (please have it open).")
    flag.Parse()

    if strings.HasPrefix(*URL, "http") != true {
  		  fmt.Printf("Please include the http:// or https:// \n")
  		  os.Exit(0)
    }
		if strings.HasSuffix(*URL, "onion") == true {
				if *TOR == false{
  		  		fmt.Printf("Requested Onion address but you are not using Tor! \n")
  		  		os.Exit(0)
				}
    }

    urls, _ := readLines(*URL, *FNAME)

		const workers = 5
		wg := new(sync.WaitGroup)
    in := make(chan string, 2*workers)

		for i := 0; i < workers; i++ {
        wg.Add(1)
        go func() {
            defer wg.Done()
            for url := range in {
								if *TOR{
                		MakeTorRequest(url, *V)
								} else{
										MakeRequest(url, *V)
								}
            }
        }()
    }

    for _, url := range urls {
        if url != "" {
            in <- url
        }
    }
    close(in)
    wg.Wait()

    fmt.Printf("\033[37mTime Elased: %.2f \033[32mFound(200): %d \033[36mNotFound(404): %d \033[31mUnauthorised(403): %d \033[34mOther(401): %d\033[0m\n", time.Since(start).Seconds(), Found, NotFound, Forbidden, Other)
}
