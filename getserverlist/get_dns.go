/* -.-.-.-.-.-.-.-.-.-.-.-.-.-.-.-.-.-.-.-.

* File Name : get_dns.go

* Purpose :

* Creation Date : 05-16-2017

* Last Modified : Wed 17 May 2017 07:09:42 PM UTC

* Created By : Kiyor

_._._._._._._._._._._._._._._._._._._._._.*/

package main

import (
	"bytes"
	"fmt"
	"github.com/miekg/unbound"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"os"
	"sort"
	"strings"
	"sync"
)

type Nameserver string

type Client struct{}

type IPS []string

func (a IPS) Len() int      { return len(a) }
func (a IPS) Swap(i, j int) { a[i], a[j] = a[j], a[i] }
func (a IPS) Less(i, j int) bool {
	i1, i2 := net.ParseIP(a[i]), net.ParseIP(a[j])
	return bytes.Compare(i1, i2) < 0
}

const DNSBASELINK = `https://public-dns.info/nameserver/%s.txt`

func (Client) GetDns(location string, limit int) ([]Nameserver, error) {
	u := fmt.Sprintf(DNSBASELINK, location)
	resp, err := http.Get(u)
	if err != nil {
		return []Nameserver{}, err
	}
	defer resp.Body.Close()
	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return []Nameserver{}, err
	}
	fulllist := strings.Split(string(b), "\n")
	if len(fulllist) < limit {
		limit = len(fulllist)
	}
	var nameservers []Nameserver
	for _, v := range fulllist[:limit] {
		nameservers = append(nameservers, Nameserver(v))
	}
	return nameservers, nil
}

type Job struct {
	domain     string
	nameserver Nameserver
}

type Record struct {
	job *Job
	ips []net.IP
	err error
}

func (c *Client) Work(i int, in chan *Job, out chan *Record) {
	for j := range in {
		if j == nil {
			return
		}
		ips, err := c.Dig(j.domain, j.nameserver)
		out <- &Record{j, ips, err}
	}
}

func (c *Client) GetRecordA(domain string, nameservers []Nameserver) (IPS, error) {

	var ips IPS
	wg := new(sync.WaitGroup)

	job := make(chan *Job)
	record := make(chan *Record)

	process := 10

	for i := 0; i < process; i++ {
		go c.Work(i, job, record)
	}

	go func() {
		for res := range record {
			for _, ip := range res.ips {
				ips = append(ips, ip.String())
			}
			wg.Done()
		}
	}()

	wg.Add(len(nameservers))
	for _, nameserver := range nameservers {
		job <- &Job{domain, nameserver}
	}
	for i := 0; i < process; i++ {
		job <- nil
	}
	wg.Wait()

	return c.Clean(ips), nil
}

func (c *Client) Clean(a IPS) IPS {
	a = c.removeIPv6(a)
	a = c.uniqIP(a)
	sort.Sort(a)
	return a
}

func (Client) removeIPv6(a IPS) IPS {
	var l IPS
	for _, v := range a {
		if i := net.ParseIP(v).To4(); i != nil {
			l = append(l, v)
		}
	}
	return l
}

func (Client) uniqIP(a IPS) IPS {
	m := make(map[string]struct{})
	for _, v := range a {
		m[v] = struct{}{}
	}
	var l IPS
	for k := range m {
		l = append(l, k)
	}
	return l
}

func (Client) Dig(domain string, nameserver Nameserver) ([]net.IP, error) {
	u := unbound.New()
	defer u.Destroy()
	tmpFile := "./resolv_" + RandomString(4) + ".conf"

	err := ioutil.WriteFile(tmpFile, []byte("nameserver "+nameserver), 0644)
	if err != nil {
		log.Fatalln(err.Error())
	}
	defer os.Remove(tmpFile)

	u.ResolvConf(tmpFile)

	return u.LookupIP(domain)

	// not able to use timeout since unbound will broken if use timeout
	/*
		type result struct {
			ips []net.IP
			err error
		}

		ch := make(chan result)

		go func() {
			ips, err := u.LookupIP(domain)
			ch <- result{ips, err}
		}()

		t := time.Tick(2 * time.Second)
		select {
		case res := <-ch:
			return res.ips, res.err
		case <-t:
			u.Destroy()
			os.Remove(tmpFile)
			return []net.IP{}, fmt.Errorf("timeout")
		}
	*/
}
