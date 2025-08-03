cpuprof:
	go build -o perft perft.go 
	./perft -depth 6 -cpuprofile="cpu.prof"
	go tool pprof perft cpu.prof

memprof:
	go build -o perft perft.go
	./perft -depth 6 -memprofile="mem.prof"
	go tool pprof perft mem.prof