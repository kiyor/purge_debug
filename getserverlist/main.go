/* -.-.-.-.-.-.-.-.-.-.-.-.-.-.-.-.-.-.-.-.

* File Name : main.go

* Purpose :

* Creation Date : 05-16-2017

* Last Modified : Wed 17 May 2017 10:01:19 PM UTC

* Created By : Kiyor

_._._._._._._._._._._._._._._._._._._._._.*/

package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"runtime"
	"sync"
)

var (
	topDns *int = flag.Int("top", 20, "get head # of dns server")
)

func init() {
	flag.Var(&flagLocations, "a", "add location")
	flag.Parse()
	if len(flagLocations) == 0 {
		flagLocations = append(flagLocations, "us")
	}
	if len(flag.Args()) == 0 {
		fmt.Fprintf(os.Stderr, "Usage:\n  %s [options] $domain\n", os.Args[0])
		os.Exit(1)
	}
}

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())
	c := new(Client)

	var ips []string

	wg := new(sync.WaitGroup)
	wg.Add(len(flagLocations))

	for _, loc := range flagLocations {
		go func(loc string) {
			nameservers, err := c.GetDns(loc, *topDns)
			if err != nil {
				log.Println(err.Error())
			}
			_ips, err := c.GetRecordA(flag.Args()[0], nameservers)
			if err != nil {
				log.Println(err.Error())
			}
			ips = append(ips, _ips...)
			wg.Done()
		}(loc)
	}

	wg.Wait()

	ips = c.Clean(ips)

	for _, v := range ips {
		fmt.Println(v)
	}
}
