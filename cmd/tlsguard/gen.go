package main

import (
	"fmt"
	"os"

	"tlsguard/pkg/cert"

	"github.com/fatih/color"
)

func genCertificate() {
	cert, key, err := cert.GenClient()
	if err != nil {
		report(err)
		return
	}

	err = os.WriteFile("client.crt", cert, 0644)
	if err != nil {
		report(err)
		return
	}

	err = os.WriteFile("client.key", key, 0644)
	if err != nil {
		report(err)
		return
	}

	color.Green("Certificate generated")
}

func report(a ...any) {
	fmt.Fprintln(os.Stderr, color.RedString("ERROR ")+fmt.Sprint(a...))
	os.Exit(1)
}
