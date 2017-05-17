/* -.-.-.-.-.-.-.-.-.-.-.-.-.-.-.-.-.-.-.-.

* File Name : random.go

* Purpose :

* Creation Date : 05-17-2017

* Last Modified : Wed 17 May 2017 12:13:43 AM UTC

* Created By : Kiyor

_._._._._._._._._._._._._._._._._._._._._.*/

package main

import (
	"math/rand"
	"time"
)

var letters = []rune("abcdefghijklmnopqrstuvwxyz0123456789")

func RandomString(size int) string {
	rand.Seed(int64(time.Now().Nanosecond()))
	str := make([]rune, size)
	for n := range str {
		str[n] = letters[rand.Intn(len(letters))]
	}
	return string(str)
}
