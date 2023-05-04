package main

import (
	"math/rand"
	"time"

	"github.com/ncw/cwtool/cmd"
	_ "github.com/ncw/cwtool/cmd/all"
)

func main() {
	rand.Seed(time.Now().UnixNano())
	cmd.Execute()
}
