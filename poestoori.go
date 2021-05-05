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
			padding: 2em;
			max-width: 600px;
			font-family: sans-serif;
		}

		div {
			margin-bottom: 3em;
		}
	</style>
</head>
<body>
	<h1>Poestoori</h1>
	<div>%s</div>
	<a href="/">Uus stoori :D</a>
</body>
<footer>
  <p>HTML from Tomi Björckin äijästoori</p>
	</footer> 
</html>`
)

func serverecipe(w http.ResponseWriter, req *http.Request) {

	fmt.Fprintf(w, fmt.Sprintf(html, poegen.Generate()))
}

func main() {
	http.HandleFunc("/", serverecipe)
	http.ListenAndServe(":8090", nil)
}
