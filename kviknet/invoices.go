package kviknet

import (
	"fmt"
	"io"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

// Invoices represents a slice of Invoice
type Invoices []Invoice

// FromReader tries to parse invoice data from input that looks something like this
/*
<div class="pure-u-1">
<table class="data_table">
	<thead>
		<tr>
			<th style="text-align: left;">Type</th>
			<th style="text-align: left;">#</th>
			<th style="text-align: left;" class="account_invoices_hide_row_sm">Dato</th>
			<th style="text-align: left;">Betalingsfrist</th>
			<th class="account_invoices_hide_row_sm">Status</th>
			<th>Beløb</th>
			<th style="width: 15px"></th>
		</tr>
	</thead>
	<tbody>

									<tr >
		<td style="text-align: left;">Indbetaling</td>
		<td style="text-align: left;">1329340</td>
		<td style="text-align: left;" class="account_invoices_hide_row_sm">28-02-2022</td>
		<td style="text-align: left;">


			<span class="due_date_large"></span>



		</td>
		<td class="account_invoices_hide_row_sm"></td>
		<td>329,00 kr.</td>
		<td>
			<a data-overlay-width="400" title="Download" href="/konto/fakturaoversigt/indbetaling/download/1013202-c1c073bb8ba1eb7387f1ab589bdbb4ee"><i class="fa fa-download fa-fw"></i></a>
		</td>
	</tr>

										<tr >
		<td style="text-align: left;">Faktura</td>
		<td style="text-align: left;">KN1021267</td>
		<td style="text-align: left;" class="account_invoices_hide_row_sm">21-02-2022</td>
		<td style="text-align: left;">


			<span class="due_date_large">28-02-2022</span>
			<span class="due_date_small">28-02-22<span>


		</td>
		<td class="account_invoices_hide_row_sm">Betalt</td>
		<td>-329,00 kr.</td>
		<td>
			<a data-overlay-width="400" title="Download" href="/konto/fakturaoversigt/faktura/download/1021267-b9f06c899f51d575713f03a702f975e1"><i class="fa fa-download fa-fw"></i></a>
		</td>
	</tr>

*/
func (i *Invoices) FromReader(r io.Reader) error {

	// Load the HTML document
	doc, err := goquery.NewDocumentFromReader(r)
	if err != nil {
		return fmt.Errorf("unable create goquery document: %w", err)
	}

	doc.Find(".data_table tbody tr").Each(func(_ int, s *goquery.Selection) {
		invoice := Invoice{}

		invoice.ID = strings.Trim(s.Find("td:nth-child(2)").Text(), " \t\n")
		invoice.Type = strings.Trim(s.Find("td:nth-child(1)").Text(), " \t\n")
		// ultra hack: also remove leading minus sign from amount
		invoice.Amount = strings.Trim(s.Find("td:nth-child(6)").Text(), "- \t\n")
		invoice.Status = strings.Trim(s.Find("td:nth-child(5)").Text(), " \t\n")
		invoice.DueDate = strings.Trim(s.Find("td:nth-child(4) .due_date_large").Text(), " \t\n")

		URL, exists := s.Find("a").Attr("href")
		if !exists {
			panic(fmt.Sprintf("invoice had no pdf download url: %+v", invoice))
		}
		invoice.InvoiceURL = fmt.Sprintf("https://www.kviknet.dk%s", URL)

		*i = append(*i, invoice)
	})

	return nil
}

func (i *Invoices) String() string {
	IDs := make([]string, len(*i))
	for index, invoice := range *i {
		IDs[index] = invoice.ID
	}
	return strings.Join(IDs, ", ")
}
