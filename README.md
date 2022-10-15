# computercraft-go
go libary that hosts a weboscket server to interact with computers from the computercraft mod in minecraft.

not every tweaked.cc api is implemented

## Getting Started
```go
import (
  "github.com/m4schini/computercraft-go"
  "fmt"
)

func main() {
  // callback for connected turtles
  computercraft.OnTurtleConnected(func(remoteAddr, uuid string, turtle computer.Turtle) {
		fmt.Println(turtle)
	})
  
  // callback for disconnected turtles
  computercraft.OnTurtleDisconnected(func(remoteAddr, uuid string, turtle computer.Turtle) {
    fmt.Println(turtle)
	})
  
  // serve entry point
  computercraft.Serve(0.0.0.0:8080)
}
```
