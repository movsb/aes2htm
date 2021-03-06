package main

import (
	"bytes"
	"strings"
	"testing"

	"github.com/romiras/aes2htm/pkg/converter"
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
	m["\033[34m111\033[m"] = `<span style="color:` + converter.Palette[4] + `;">111</span>`
	m["\033[44m222\033[m"] = `<span style="background-color:` + converter.Palette[4] + `;">222</span>`
	m["\033[34m\033[43m333\033[39;49m\033[0m"] = `<span style="color:` + converter.Palette[4] + `;"><span style="background-color:` + converter.Palette[3] + `;">333</span></span>`

	// m = map[string]string{
	// 	"\033[34m\033[43m333\033[39;49m\033[0m": `<span style="color:` + Palette[4] + `;"><span style="background-color:` + Palette[3] + `;">333</span></span>`,
	// }

	for k, v := range m {
		sw := bytes.NewBuffer(nil)
		ah, err := converter.NewAes2Htm(sw)
		if err != nil {
			t.Fatal(err)
		}

		sr := strings.NewReader(k)
		er := ah.Input(sr)
		if er != nil {
			t.Fatal(er)
		}
		if sw.String() != v {
			t.Fatalf("%s -> %s -> %s\n", k, sw.String(), v)
			// } else {
			// t.Logf("Pass: %s -> %s", k, sw.String())
		}
	}
}
