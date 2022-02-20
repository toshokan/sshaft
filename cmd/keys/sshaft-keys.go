package main

import (
	"flag"
	"fmt"
	"github.com/toshokan/sshaft/internal"
	"os"
)

var cfgFile string

func init() {
	flag.StringVar(&cfgFile, "config", "", "Config file path")
	flag.Parse()
}

func main() {
	if cfgFile == "" {
		fmt.Println("Supply --config <path>")
		os.Exit(1)
	}
	cfg, err := internal.LoadCfg(cfgFile)
	if err != nil {
		panic(err)
	}
	token, err := internal.GetToken(cfg)
	if err != nil {
		panic(err)
	}
	keys, err := internal.GetMFAKeys(cfg, token)
	if err != nil {
		panic(err)
	}
	keyLines := internal.GetKeyLines(cfg, keys)
	for _, line := range keyLines {
		fmt.Printf("%s", line)
	}
}
