package main

import (
	"fmt"
	"os"
)

func render() {
	ah := NewAes2Htm(os.Stdout)
	err := ah.Input(os.Stdin)
	if err != nil {
		panic(err)
	}
}

func main() {
	tohtml := len(os.Args) == 2 && os.Args[1] == "--html"
	if tohtml {
		fmt.Fprint(os.Stdout,
			`<!doctype html>
<head>
<meta charset="utf-8" />
<link rel="stylesheet" href="aes2htm.css" />
</head>
<body>
<pre>
`,
		)
	}
	render()
	if tohtml {
		fmt.Fprint(os.Stdout,
			`
</pre>
</body>
</html>
`)
	}
}
