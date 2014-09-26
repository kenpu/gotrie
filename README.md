Compressed Trie for uint64
==========================

    export GOPATH=$PWD
    cd src
    go run main.go
    go test gotrie
    go test -cpuprofile=prof.out gotrie
    go tool pprof --lines --web ./gotrie.test ./prof.out


