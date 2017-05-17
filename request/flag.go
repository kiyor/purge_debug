/* -.-.-.-.-.-.-.-.-.-.-.-.-.-.-.-.-.-.-.-.

* File Name : flag.go

* Purpose :

* Creation Date : 01-24-2017

* Last Modified : Wed 17 May 2017 07:21:26 PM UTC

* Created By : Kiyor

_._._._._._._._._._._._._._._._._._._._._.*/

package main

import ()

var (
	flagAddHeader flagSliceString
)

type flagSliceString []string

func (i *flagSliceString) String() string {
	return ""
}

func (i *flagSliceString) Set(value string) error {
	*i = append(*i, value)
	return nil
}
