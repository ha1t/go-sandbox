# glob

```golang
package main

import (
	fmt "fmt"
	"path/filepath"
)

func main() {

	items, err := filepath.Glob("/home/*")
	if err != nil {
		panic(err)
	}
	for _, item := range items {
		fmt.Printf(item + "\n")
	}
}
```
