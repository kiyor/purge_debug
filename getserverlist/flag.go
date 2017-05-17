/* -.-.-.-.-.-.-.-.-.-.-.-.-.-.-.-.-.-.-.-.

* File Name : flag.go

* Purpose :

* Creation Date : 05-17-2017

* Last Modified : Wed 17 May 2017 12:49:46 AM UTC

* Created By : Kiyor

_._._._._._._._._._._._._._._._._._._._._.*/

package main

import (
// 	"flag"
)

var (
	flagLocations flagSliceString
)

type flagSliceString []string

func (i *flagSliceString) String() string {
	return ""
}

func (i *flagSliceString) Set(value string) error {
	*i = append(*i, value)
	return nil
}
