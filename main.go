package main

import (
	"math/rand"
	"time"

	"github.com/ncw/ncwtester/cmd"
	_ "github.com/ncw/ncwtester/cmd/all"
)

func main() {
	rand.Seed(time.Now().UnixNano())
	cmd.Execute()
}
