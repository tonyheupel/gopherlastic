# gopherlastic

A minimal Elasticsearch client


## Using the API

First you create a client, specifying a host and then use it.

```go
package main

import (
        "fmt"
        "github.com/infospace/gopherlastic"
)

type Something struct {
        SomeValue string `json:"someValue"`
        OtherThing string `json:"otherThing"`
}

func main() {
        doc := &Soemthing{
                SomeValue: "some",
                OtherThing: "thingy",
        }

        client := gopherlastic.NewClient(r.Host)
        putResponse, err := client.PutDocument("some-index", "something", "1", doc)

        if err != nil {
                fmt.Println("Problem storing document:", err)
                return
        }

        fmt.Println("Saved document:", putResponse)
}
```


There's also a WIP Search interface.
