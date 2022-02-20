package main

import (
	"flag"
	"fmt"
	"os"
	"github.com/toshokan/sshaft/internal"
)

var (
	cfgFile string
	user string
)

func init() {
	flag.StringVar(&cfgFile, "config", "", "Config file path")
	flag.StringVar(&user, "user", "", "Username of the user that authenticated")
	flag.Parse()
}

func main() {
	if user == "" {
		fmt.Println("Supply --user <user>")
		os.Exit(1);
	}
	if cfgFile == "" {
		fmt.Println("Supply --config <path>")
		os.Exit(1);
	}
	cfg, err := internal.LoadCfg(cfgFile)
	if err != nil {
		panic(err)
	}
	token, err := internal.GetToken(cfg)
	if err != nil {
		panic(err)
	}
	if err := internal.MFAAccept(cfg, token, user); err != nil {
		panic(err)
	}
	fmt.Printf("Hello %s!\n", user)
}
