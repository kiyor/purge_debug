/* -.-.-.-.-.-.-.-.-.-.-.-.-.-.-.-.-.-.-.-.

* File Name : resp.go

* Purpose :

* Creation Date : 10-03-2016

* Last Modified : Wed 14 Jun 2017 01:56:38 AM UTC

* Created By : Kiyor

_._._._._._._._._._._._._._._._._._._._._.*/

package main

import (
	"crypto/md5"
	"encoding/json"
	"fmt"
	"github.com/kiyor/golib"
	"github.com/kiyor/terminal/color"
	"io/ioutil"
	"log"
	"net/http"
	"sort"
	"strings"
	"time"
)

var (
	totalReq  float64
	expectReq float64
)

func toJson(intf interface{}) string {
	j, err := json.Marshal(intf)
	if err != nil {
		log.Println(err.Error())
	}
	return string(j)
}

type Resp struct {
	Req      *http.Request
	Resp     *http.Response
	Err      error
	Dur      time.Duration
	WorkerId int
}

func (resp *Resp) output() {
	if *jsonOut {
		m := make(map[string]interface{})
		m["Request"] = golib.CopyRequest(resp.Req)
		m["Response"] = golib.CopyResponse(resp.Resp)
		m["Error"] = resp.Err
		m["Duration"] = resp.Dur.String()
		fmt.Println(toJson(m))
		return
	}
	if resp.Err == nil {
		var msg string
		var lastmodGood bool

		if val, ok := resp.Resp.Header["Etag"]; ok {
			if len(val[0]) > 6 {
				msg = color.Sprintf("Etag: %s", val[0][1:6])
			}
		}
		if val, ok := resp.Resp.Header["Last-Modified"]; ok {
			totalReq++
			if val[0] == *expectLM {
				expectReq++
				lastmodGood = true
				msg = color.Sprintf("%s @{g}LM: %s@{|}", msg, val[0])
			} else {
				msg = color.Sprintf("%s @{r}LM: %s@{|}", msg, val[0])
			}
		}
		if val, ok := resp.Resp.Header["Content-Length"]; ok {
			msg = color.Sprintf("%s Len: %s", msg, val[0])
		}
		if val, ok := resp.Resp.Header["Age"]; ok {
			msg = color.Sprintf("%s Age: %s", msg, val[0])
		}

		if *method == "GET" {
			if *sum || *grep != "" {
				b, _ := ioutil.ReadAll(resp.Resp.Body)
				defer resp.Resp.Body.Close()
				if *sum {
					m := fmt.Sprintf("%x", md5.Sum(b))
					msg = color.Sprintf("%s MD5: %s", msg, m[:5])
				}
				if *grep != "" {
					g := *grep
					ident := g[:1]
					strs := strings.Split(strings.Trim(g, ident), ident)
					for _, v := range strs {
						if strings.Contains(string(b), v) {
							msg = color.Sprintf("%s @{g}[✓]%s@{|}", msg, v)
						} else {
							msg = color.Sprintf("%s @{r}[✗]%s@{|}", msg, v)
						}
					}
				}
			}
		}
		if resp.Dur > 1*time.Second {
			msg += color.Sprintf(" @{r}%s@{|}", resp.Dur)
		} else if resp.Dur > 500*time.Millisecond {
			msg += color.Sprintf(" @{y}%s@{|}", resp.Dur)
		} else {
			msg += color.Sprintf(" @{g}%s@{|}", resp.Dur)
		}
		if lastmodGood {
			msg += " success"
		} else {
			msg += " fail"
		}

		var res string
		res += color.Sprintf("%-17s ", resp.Req.URL.Host)
		if resp.Resp.StatusCode > 302 {
			res += color.Sprintf("@{r}%-5d@{|}", resp.Resp.StatusCode)
		} else {
			res += color.Sprintf("@{g}%-5d@{|}", resp.Resp.StatusCode)
		}
		fmt.Println(res + msg)
		if *verbose2 || len(*showHeaders) > 0 {
			switch resp.Resp.StatusCode {
			case http.StatusMovedPermanently, http.StatusFound, http.StatusSeeOther, http.StatusTemporaryRedirect:
			default:
				msg = ""
				msg += color.Sprintf("@{g}%v %v@{|}\n", resp.Resp.Proto, resp.Resp.Status)
				if len(*showHeaders) == 0 {
					hs := []string{}
					for k := range resp.Resp.Header {
						hs = append(hs, k)
					}
					sort.Strings(hs)
					for _, v := range hs {
						for _, h := range resp.Resp.Header[v] {
							msg += color.Sprintf("@{g}%v: %v@{|}\n", v, h)
						}
					}
				} else {
					for _, v := range strings.Split(*showHeaders, ",") {
						for _, h := range resp.Resp.Header[v] {
							msg += color.Sprintf("@{g}%v: %v@{|}\n", v, h)
						}
					}
				}
				fmt.Print(msg)
			}
		}
	} else {
		color.Printf("@{r}%s %-17s %s(%v)@{|}\n", time.Now(), resp.Req.URL.Host, resp.Err.Error(), *timeout)
	}
}
