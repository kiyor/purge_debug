/* -.-.-.-.-.-.-.-.-.-.-.-.-.-.-.-.-.-.-.-.

* File Name : wg.go

* Purpose :

* Creation Date : 08-02-2016

* Last Modified : Wed 17 May 2017 07:33:14 PM UTC

* Created By : Kiyor

_._._._._._._._._._._._._._._._._._._._._.*/

package main

import (
	"crypto/tls"
	"errors"
	"golang.org/x/net/proxy"
	"log"
	"net"
	"net/http"
	"strings"
	"sync"
	"time"
)

type WorkerGroup struct {
	in  chan string
	n   int
	out chan *Resp
	sync.WaitGroup
}

func NewWorkerGroup(in chan string, out chan *Resp) *WorkerGroup {
	return &WorkerGroup{
		in:  in,
		out: out,
	}
}

func (wg *WorkerGroup) run(id int) {
	for line := range wg.in {
		wg.Do(id, line)
	}
}

func (wg *WorkerGroup) Start() {
	go wg.run(wg.n)
	wg.n++
}

func (wg *WorkerGroup) Do(id int, line string) {
	// 	defer wg.Done()
	if len(line) > 0 {
		line = strings.Trim(line, "\n")
	}
	for strings.Contains(line, "  ") {
		line = strings.Replace(line, "  ", " ", -1)
	}
	line = strings.Trim(line, " ")
	var hostPart string
	if reIp.MatchString(line) {
		i := reIp.FindStringSubmatch(line)[1]
		if ip := net.ParseIP(i); ip == nil {
			log.Println("error parse ip", i)
			return
		} else {
			hostPart = ip.String()
		}

	}
	// 	ip := reIp.FindAllStringSubmatch(ipPart, -1)
	// 	server := strings.Split(line, " ")[1]

	url := strings.Replace(*requestUrl, host, hostPart, 1)

	req, err := http.NewRequest(*method, url, nil)
	if err != nil {
		ch <- &Resp{req, nil, err, 0, id}
	}

	req.Host = host

	req.Header.Add("User-Agent", *useragent)

	if *akamai {
		req.Header.Add("Pragma", "akamai-x-cache-on, akamai-x-cache-remote-on, akamai-x-check-cacheable, akamai-x-get-cache-key, akamai-x-get-extracted-values, akamai-x-get-nonces, akamai-x-get-ssl-client-session-id, akamai-x-get-true-cache-key, akamai-x-serial-no")
	}
	if *chinacache {
		req.Header.Add("x-c3-debug", "enabled")
	}
	if *gzip {
		req.Header.Add("Accept-Encoding", "gzip")
	}

	t := time.Tick(*timeout)
	r := make(chan *Resp)

	var client *http.Client

	if len(*sock) == 0 {
		client = &http.Client{
			Transport: &http.Transport{
				Proxy: nil,
				TLSClientConfig: &tls.Config{
					ServerName:         host,
					InsecureSkipVerify: *ignoreCert,
				},
				DisableCompression: true,
				DisableKeepAlives:  true,
				DialContext: (&net.Dialer{
					Timeout:   *timeout,
					KeepAlive: 0,
				}).DialContext,
			},
		}
	} else {
		dialer, err := proxy.SOCKS5("tcp", *sock,
			nil,
			&net.Dialer{
				Timeout:   *timeout,
				KeepAlive: 30 * time.Second,
			},
		)
		if err == nil {
			client = &http.Client{
				Transport: &http.Transport{
					Proxy: nil,
					TLSClientConfig: &tls.Config{
						ServerName:         host,
						InsecureSkipVerify: *ignoreCert,
					},
					DisableCompression: true,
					DisableKeepAlives:  true,
					Dial:               dialer.Dial,
				},
			}
		}

	}

	go func() {
		t1 := time.Now()
		var resp *http.Response
		var err error
		if *curl {
			Curl(req)
		}
		resp, err = client.Transport.RoundTrip(req)

		r <- &Resp{req, resp, err, time.Since(t1), id}
	}()

	select {
	case a := <-r:
		wg.out <- a
	case <-t:
		wg.out <- &Resp{req, nil, errors.New("request timeout"), *timeout, id}
	}
}
