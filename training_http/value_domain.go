package main

import (
	php "./php"
	"regexp"
	"strings"
)

func get_ip_address() string {
	result := php.File_get_contents("http://checkip.dyndns.org/")

	r, _ := regexp.Compile(`\<body\>.*?:(.*?)\<\/body\>`)

	matches := r.FindStringSubmatch(result)

	return strings.Trim(matches[1], " ")
}

func main() {
	println(get_ip_address())
	//php.File_put_contents(os.Args[2], result)
}
