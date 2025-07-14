[![License: MPL 2.0](https://img.shields.io/badge/License-MPL%202.0-brightgreen.svg)](https://opensource.org/licenses/MPL-2.0)

Chego implements chess board state management and legal move generation.

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

	"github.com/BelikovArtem/chego/enum"
	"github.com/BelikovArtem/chego/game"
	"github.com/BelikovArtem/chego/movegen"
)

func main() {
	// It is important to call InitAttackTables as close to the program
	// start as possible, otherwise the move generation won't work.
	movegen.InitAttackTables()

	game := game.NewGame()
	// Scholar's mate.
	game.PushMove(movegen.NewMove(enum.SF3, enum.SF2, 0, enum.MoveNormal))
	game.PushMove(movegen.NewMove(enum.SE5, enum.SE7, 0, enum.MoveNormal))
	game.PushMove(movegen.NewMove(enum.SG4, enum.SG2, 0, enum.MoveNormal))
	game.PushMove(movegen.NewMove(enum.SH4, enum.SD8, 0, enum.MoveNormal))

	fmt.Printf("Is checkmate: %t\n", game.IsCheckmate()) // Prints "Is checkmate: true"
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

To run tests and benchmarks, run this commands in the chego folder:  

```
go test ./...
go test ./... -bench=. -benchmem
```	

Here are the benchmark results recieved on Intel Core i7-10750H CPU:

![Benchmark results](./doc/benchmarks.png)

## License

Copyright (c) 2025 Artem Bielikov

This project is licensed under the Mozilla Public License 2.0.  
See the [LICENSE](LICENSE) file for details.