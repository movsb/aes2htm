package main

import (
	"bytes"
	"strings"
	"testing"
)

func TestAes2Htm(t *testing.T) {
	m := make(map[string]string)

	// 无效果
	m["abc"] = "abc"

	// 关闭
	m["\033[m"] = ""
	m["\033[m123"] = "123"
	m["123\033[m"] = "123"
	m["123\033[m123"] = "123123"

	// 颜色
	m["\033[34m123\033[m"] = `<span style="color:blue;">123</span>`
	m["\033[44m123\033[m"] = `<span style="background-color:blue;">123</span>`
	m["\033[34m\033[43m123\033[39;49m\033[0m"] = `<span style="color:blue;"><span style="background-color:yellow;">123</span></span>`

	for k, v := range m {
		ah := &Aes2Htm{}
		sw := bytes.NewBuffer(nil)
		sr := strings.NewReader(k)
		er := ah.Input(sw, sr)
		if er != nil {
			t.Fatal(er)
		}
		if sw.String() != v {
			t.Fatalf("%s -> %s -> %s\n", k, sw.String(), v)
		} else {
			t.Logf("Pass: %s -> %s", k, sw.String())
		}
	}
}
