package main

import (
	"bytes"
	"strings"
	"testing"
)

func Test(t *testing.T) {
	m := make(map[string]string)

	m["abc"] = "abcf"

	for k, v := range m {
		ah := &Aes2Htm{}
		sw := bytes.NewBuffer(nil)
		sr := strings.NewReader(k)
		er := ah.Input(sw, sr)
		if er != nil {
			t.Fatal(er)
		}
		if sw.String() != v {
			t.Fatalf("%s -> %s\n", k, v)
		}
	}
}
