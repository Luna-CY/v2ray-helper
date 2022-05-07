package main

import (
	"fmt"
	"github.com/Luna-CY/v2ray-helper/command/http"
	"os"
)

func main() {
	if err := http.Execute(); nil != err {
		fmt.Println(err)

		os.Exit(1)
	}
}
