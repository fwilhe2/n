package main

import (
	"testing"
)

func TestHelloWorld(t *testing.T) {
	actual := formatDate("Fri, 09 Apr 2021 00:21:49 +0000")
	if actual != "2021-04-09 02:21:49" {
		t.Fail()
	}
	//todo: make this work: formatDate("Thu, 27 May 2021 19:30:00 GMT")

}