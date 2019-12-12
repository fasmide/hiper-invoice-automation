package main

import (
	"log"
	"os"

	"github.com/fasmide/hipermads/hiper"
)

func main() {

	session, err := hiper.Login(os.Getenv("HIPER_USERNAME"), os.Getenv("HIPER_PASSWORD"))
	if err != nil {
		log.Fatalf("could not login: %s", err)
	}

	// If there was no error - we are logged in
	log.Print("Logged in!")

	invoices, err := session.Invoices()
	log.Printf("Found invoices: %s", invoices.String())

}
