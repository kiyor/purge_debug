/* -.-.-.-.-.-.-.-.-.-.-.-.-.-.-.-.-.-.-.-.

* File Name : main.go

* Purpose :

* Creation Date : 01-09-2015

* Last Modified : Wed 14 Jun 2017 12:54:22 AM UTC

* Created By : Kiyor

_._._._._._._._._._._._._._._._._._._._._.*/

package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"regexp"
	"runtime"
	"time"
)

var (
	reIp        = regexp.MustCompile(`(\d+\.\d+\.\d+\.\d+)`)
	ch          = make(chan *Resp)
	host        string
	expectLM    *string        = flag.String("lm", "", "Last-Modified expected header")
	expectRatio *float64       = flag.Float64("ratio", 1.00, "good ratio, if good return 0, bad return 1")
	version     *bool          = flag.Bool("v", false, "output version and exit")
	verbose1    *bool          = flag.Bool("vv", false, "verbose output")
	verbose2    *bool          = flag.Bool("vvv", false, "more verbose output")
	ignoreCert  *bool          = flag.Bool("k", false, "ignore cert check")
	curl        *bool          = flag.Bool("curl", false, "show curl")
	requestUrl  *string        = flag.String("u", `http://test.com/test`, "url")
	method      *string        = flag.String("method", "GET", "request method")
	sock        *string        = flag.String("socks5", "", "request using socks5 proxy")
	timeout     *time.Duration = flag.Duration("timeout", 10*time.Second, "request timeout")
	akamai      *bool          = flag.Bool("akamai", false, "add akamai debug header")
	chinacache  *bool          = flag.Bool("chinacache", false, "add Chinacache debug header")
	sum         *bool          = flag.Bool("md5", false, "get body and cal md5")
	grep        *string        = flag.String("grep", "", "grep body, syntax /abc/def/")
	gzip        *bool          = flag.Bool("gzip", false, "add gzip header")
	useragent   *string        = flag.String("useragent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_10_1) AppleWebKit/600.2.5 (KHTML, like Gecko) Version/8.0.2 Safari/600.2.5", "request useragent")
	showHeaders *string        = flag.String("headers", "", "show only headers, syntax -headers a,b")
	worker      *int           = flag.Int("w", 2*runtime.NumCPU(), "worker")
	xcache      *string        = flag.String("xcache", "Powered-By-Chinacache|X-Powered-By-Chinacache", "cache status possible headers")
	jsonOut     *bool          = flag.Bool("json", false, "json output")
	VER                        = "1.0"
	buildtime   string

// 	client *http.Client
)

func init() {
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage of %s:\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "VER: %v.%v\n", VER, buildtime)
		flag.PrintDefaults()
	}
	flag.Var(&flagAddHeader, "H", "add header, -H 'key1:value1', -H 'key2:value2'")
	flag.Parse()
	// 	if *verbose2 {
	// 		fmt.Println(flag.Args())
	// 	}
	if *version {
		fmt.Printf("%v.%v", VER, buildtime)
		os.Exit(0)
	}
	if *requestUrl == "http://test.com/test" {
		// 		if len(flag.Args()) > 0 {
		// 			*requestUrl = flag.Args()[0]
		// 		} else {
		fmt.Println("tsisreq -u http://a.com/b")
		os.Exit(0)
		// 		}
	}
	var err error
	Url, err := url.Parse(*requestUrl)
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(0)
	}
	host = Url.Host
	if len(*grep) > 0 {
		*method = "GET"
	}
	runtime.GOMAXPROCS(runtime.NumCPU())
}
func main() {
	stop := make(chan bool)
	in := make(chan string, 100)
	out := make(chan *Resp, 100)
	w := NewWorkerGroup(in, out)
	for n := 0; n < *worker; n++ {
		go w.Start()
	}
	go func() {
		reader := bufio.NewReader(os.Stdin)
		for {
			l, err := reader.ReadString('\n')
			if err != nil {
				if err == io.EOF {
					w.Wait()
					stop <- true
				} else {
					log.Println(err.Error())
					os.Exit(1)
				}
			} else {
				w.Add(1)
				in <- l
			}
		}
	}()
HERE:
	for {
		select {
		case resp := <-out:
			resp.output()
			w.Done()
		case <-stop:
			break HERE
		}
	}
	if expectReq/totalReq >= *expectRatio {
		os.Exit(0)
	}
	os.Exit(1)
}
func Curl(req *http.Request) {
	s := "curl -I "
	for k, v := range req.Header {
		s += fmt.Sprintf("-H '%v: ", k)
		for _, h := range v {
			s += fmt.Sprintf("%v", h)
		}
		s += fmt.Sprint("' ")
	}
	s += fmt.Sprintf("-H 'Host: %v'", req.Host)
	fmt.Printf("%v '%v://%v%v'\n", s, req.URL.Scheme, req.URL.Host, req.URL.Path)
}
