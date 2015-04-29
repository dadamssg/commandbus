CommandBus
==========

#### Installation
Make sure that Go is installed on your computer.
Type the following command in your terminal:

    go get github.com/dadamssg/CommandBus

After it the package is ready to use.

#### Import package in your project
Add following line in your `*.go` file:
```go
import "github.com/dadamssg/CommandBus"
```

#### Example
```go
package main

import (
    "fmt"
    "github.com/dadamssg/CommandBus"
)

type RegisterUserCommand struct {
    Username string
}

func main() {
    bus := CommandBus.New()

    bus.RegisterHandler(&RegisterUserCommand{}, func(cmd interface{}) {
        command, _ := cmd.(*RegisterUserCommand)
        fmt.Println("Registered: ", command.Username)
    })

    // add a middleware with a priority of 0
    bus.AddMiddleware(0, func(cmd interface{}, next CommandBus.HandlerFunc) {
        fmt.Println("Enter mock caching middleware")
        next(cmd)
        fmt.Println("Exit mock caching middleware")
    })

    // add a middleware with a priority of 1
    bus.AddMiddleware(1, func(cmd interface{}, next CommandBus.HandlerFunc) {
        fmt.Println("Enter mock logging middleware")
        next(cmd)
        fmt.Println("Exit mock logging middleware")
    })

    bus.Handle(&RegisterUserCommand{
        Username: "David",
    })
}
```

The above will produce the following output

```
Enter mock logging middleware
Enter mock caching middleware
Registered:  David
Exit mock caching middleware
Exit mock logging middleware
```

#### Why? 
I've come to appreciate hexagonal architecture and was missing a package that would allow me to architect an app in such a way. Structuring an app using a command bus allows you to have several different client interfaces(http, cli, etc.) that are only responsible for creating command structs and sending them into the bus to be handled. A command only has a single handler but you can add as many middleware layers as needed, allowing to easily plug in logging, caching, authentication, etc.
