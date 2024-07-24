package main

import (
	"flag"
	"fmt"
	"github.com/kaspanet/kaspad/domain/dagconfig"
	"github.com/pkg/errors"
	"github.com/svarogg/kaspaper/kaspaperlib"
	"io/ioutil"
	"os"
)

func main() {
	// 定义命令行参数
	password := flag.String("p", "", "password")
	outputFile := flag.String("o", "", "output file")

	// 解析命令行参数
	flag.Parse()

	if *password == "" || *outputFile == "" {
		fmt.Println("Error: Both -p and -o parameters are required. Example: kaspaper -p pass -o index.html")
		os.Exit(1)
	}

	wallet, err := kaspaperlib.NewAPI(&dagconfig.MainnetParams).GenerateWallet()
	if err != nil {
		printErrorAndExit(err)
	}

	walletString, err := renderWallet(wallet, *password)
	if err != nil {
		printErrorAndExit(err)
	}

	err = ioutil.WriteFile(*outputFile, []byte(walletString), 0600)
	if err != nil {
		printErrorAndExit(errors.WithStack(err))
	}
	fmt.Printf("Paperwallet written to %s\n", *outputFile)
}

func printErrorAndExit(err error) {
	fmt.Fprintf(os.Stderr, "A critical error occured:\n%+v\n", err)
	os.Exit(1)
}
