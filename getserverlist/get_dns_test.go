/* -.-.-.-.-.-.-.-.-.-.-.-.-.-.-.-.-.-.-.-.

* File Name : get_dns_test.go

* Purpose :

* Creation Date : 05-16-2017

* Last Modified : Wed 17 May 2017 12:47:06 AM UTC

* Created By : Kiyor

_._._._._._._._._._._._._._._._._._._._._.*/

package main

import (
	"log"
	"testing"
)

func Test_GetDns(t *testing.T) {
	c := new(Client)

	nameservers, err := c.GetDns("us", 10)
	if err != nil {
		t.Fatal(err)
	}

	records, err := c.GetRecordA("www.google.com", nameservers)
	if err != nil {
		t.Fatal(err)
	}
	log.Println(records)
}
