package main

import "flag"

func main() {
	genCert := flag.Bool("cert", false, "Generate client certificate")
	flag.Parse()

	if *genCert {
		genCertificate()
	} else {
		runServer()
	}
}
