[![License: MPL 2.0](https://img.shields.io/badge/License-MPL%202.0-brightgreen.svg)](https://opensource.org/licenses/MPL-2.0)

Chego implements chessboard state management and legal move generation.

Piece positions are stored as bitboards.

Move generation is implemented using the Magic Bitboards method.

It is assigned to use in the web-servers (for example, [justchess.org](https://justchess.org/)),
hence it does not provide any GUI or CLI.

## Usage

To install chego, run `go get`:

```
go get github.com/BelikovArtem/chego
```

Here is a simple example: 

```go
package main

import (
	"fmt"

	"github.com/BelikovArtem/chego"
)

func main() {
	// It is important to call InitAttackTables as close to the program
	// start as possible, otherwise the move generation won't work.
	chego.InitAttackTables()

	g := chego.NewGame()
	// Scholar's mate.
	g.PushMove(chego.NewMove(chego.SF3, chego.SF2, chego.MoveNormal))
	g.PushMove(chego.NewMove(chego.SE5, chego.SE7, chego.MoveNormal))
	g.PushMove(chego.NewMove(chego.SG4, chego.SG2, chego.MoveNormal))
	g.PushMove(chego.NewMove(chego.SH4, chego.SD8, chego.MoveNormal))
 	// Prints "Is checkmate: true"
	fmt.Printf("Is checkmate: %t\n", g.IsCheckmate())
}
```

## Local installation

First install the Go compiler version 1.24.1 or newer (see https://go.dev/dl).

Once the compiler is installed, clone this repository:

```
git clone https://github.com/BelikovArtem/chego
cd chego
```

## Tests and benchmarks

To execute the performance test, run this command in the chego folder:  

```
go run ./internal/perft.go -depth {IntValue}
```	

## License

Copyright (c) 2025 Artem Bielikov

This project is available under the Mozilla Public License, v. 2.0.  
See the [LICENSE](LICENSE) file for details.