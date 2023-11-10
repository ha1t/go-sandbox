package main

import (
	php "./php"
	"os"
)

func main() {
	result := php.File_get_contents(os.Args[1])
	php.File_put_contents(os.Args[2], result)
}
