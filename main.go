package main

import (
	"fmt"
	"github.com/ifchange/botKit/config"
)

func main() {
	fmt.Println(config.GetConfig())
}
