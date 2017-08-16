# restc - simple golang http client to fetch json data

## How to use?

### example

```shell
go get github.com/thbourlove/restc
```

```go
package main

import (
	"log"

	"github.com/thbourlove/restc"
)

func main() {
	var usernames []string
	restc.NewClient().GetJsonDataWithPath("https://api.github.com/users", &usernames, "$[*].login")
	log.Println(usernames)
}
```
