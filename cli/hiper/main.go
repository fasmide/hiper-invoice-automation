package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"

	"github.com/fasmide/hipermads/hiper"
	"gopkg.in/gomail.v2"
)

// SentIDs represents invoice IDs which we already emailed
type SentIDs map[int]bool

func main() {

	// before communicating with hiper, lets get our own data stright
	fd, err := os.Open(os.Getenv("HIPER_ID_DB_PATH"))
	if err != nil {
		log.Fatalf("could not open db: %s", err)
	}
	defer fd.Close()

	sentIDs := make(SentIDs)
	decoder := json.NewDecoder(fd)
	err = decoder.Decode(&sentIDs)
	if err != nil {
		log.Printf("unable to decode db: %s", err)
	}

	// communicate with hiper - lookup first page of invoices
	session, err := hiper.Login(os.Getenv("HIPER_USERNAME"), os.Getenv("HIPER_PASSWORD"))
	if err != nil {
		log.Fatalf("could not login: %s", err)
	}

	log.Print("Logged in!")

	invoices, err := session.Invoices()
	if err != nil {
		log.Fatalf("unable to fetch invoices: %s", err)
	}
	log.Printf("Found invoices: %s", invoices.String())

	// we should now email invoices that are present in invoices, but not in sentIDs - and then
	// update sendIDs + write to disk
	d := gomail.NewDialer("smtp.gmail.com", 587, os.Getenv("HIPER_EMAIL_FROM"), os.Getenv("HIPER_EMAIL_PASSWORD"))
	for _, invoice := range invoices {
		_, exists := sentIDs[invoice.ID]
		if exists {
			continue
		}

		log.Printf("Will send invoice %d of a total amount of %s", invoice.ID, invoice.Amount)

		m := gomail.NewMessage()
		m.SetHeader("From", os.Getenv("HIPER_EMAIL_FROM"))
		m.SetHeader("To", os.Getenv("HIPER_EMAIL_TO"))
		m.SetHeader("Subject", fmt.Sprintf("Regning %d fra hiper på %s", invoice.ID, invoice.Amount))
		m.SetBody("text/plain", fmt.Sprintf("Hello Mads\n\nDet er tid til at flytte %s til mides konto: %s, vedhæftet er faktura\n\nvh Mide", invoice.Amount, os.Getenv("HIPER_MIDES_ACCOUNTNO")))
		m.Attach(fmt.Sprintf("%d.pdf", invoice.ID), gomail.SetCopyFunc(func(w io.Writer) error {
			// attach the invoice - downloaded directly from hiper
			resp, err := session.Get(invoice.InvoiceURL)
			if err != nil {
				log.Fatalf("unable to get invoice pdf: %s", err)
			}
			defer resp.Body.Close()
			_, err = io.Copy(w, resp.Body)
			if err != nil {
				log.Fatalf("unable to copy invoice pdf to email: %s", err)
			}

			return nil
		}))

		err := d.DialAndSend(m)
		if err != nil {
			log.Fatalf("could not send email: %s", err)
		}
		log.Printf(".... sent")
		// if we reached this far - we should update disk with the new sentIds
		sentIDs[invoice.ID] = true
		StoreIds(sentIDs)
	}

	_ = session.Logout()
}

// StoreIds writes out ids we have sent to disk
func StoreIds(s SentIDs) {
	fd, err := os.Create(os.Getenv("HIPER_ID_DB_PATH"))
	if err != nil {
		log.Fatalf("unable to open database for write: %s", err)
	}

	encoder := json.NewEncoder(fd)
	err = encoder.Encode(s)
	if err != nil {
		log.Fatalf("could not encode json: %s", err)
	}
}
