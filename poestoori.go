package main

import (
	"fmt"
	"net/http"
	"poegen"
)

var (
	html string = `<html>
<head>
	<style>
		body {
			padding: 4em;
			max-width: 800px;
			font-family: sans-serif;
		}
		footer {
			font-family: sans-serif;
			font-size: 10px;
		}
		div {
			margin-bottom: 3em;
		}
	</style>
</head>
<body>
	<h1>Poestoori</h1>
	<div>%s</div>
	<a href="%s">Uus stoori :D</a>
</body>
<footer>
  <p>HTML from Tomi Björckin äijästoori</p>
	</footer> 
</html>`
	apiroot string = "/poegen"
)

func serverecipe(w http.ResponseWriter, req *http.Request) {
	text := poegen.Generate()
	fmt.Fprintf(w, fmt.Sprintf(html, text, apiroot))
}

func main() {
	http.HandleFunc("/", serverecipe)
	http.ListenAndServe(":8090", nil)
}
