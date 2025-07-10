[![License: MPL 2.0](https://img.shields.io/badge/License-MPL%202.0-brightgreen.svg)](https://opensource.org/licenses/MPL-2.0)

Chego implements chess board state management and legal move generation.

Piece positions are stored as bitboards.

Move generation is implemented using the Magic Bitboards method.

It is assigned to use in the web-servers (for example, [justchess.org](https://justchess.org/)),
hence it does not provide any GUI or CLI.

## Local installation

First install the Go compiler version 1.24.1 or newer (see https://go.dev/dl).

Once the compiler is installed, clone this repository and run `go install`:

```
	$ git clone https://github.com/BelikovArtem/chego
	$ cd chego
	$ go install
```

## Tests and benchmarks

To run tests and benchmarks, run this commands in the chego folder:  

```
	$ go test ./...
	$ go test ./... -bench=. -benchmem
```	

Here are the benchmark results recieved on Intel Core i7-10750H CPU:

![Benchmark results](./doc/benchmarks.png)
