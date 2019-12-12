package hiper

import (
	"fmt"
	"io"
	"strconv"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

// Invoices represents a slice of Invoice
type Invoices []Invoice

// FromReader tries to parse invoice data from input that looks something like this
/*
<h1 class="heading heading--big heading--red">Regninger</h1>
<div class="selfcare-page-notifications">
</div>
<table class="table--invoices">
	<tr>
		<th>Regningsnummer</th>
		<th>Type</th>
		<th>Total</th>
		<th>Status</th>
		<th colspan="2">Forfaldsdato</th>
	</tr>
	<tr>
		<td data-label="Regninsnummer">1660653</td>
		<td data-label="Type">Faktura</td>
		<td data-label="Total">308,00&nbsp;kr.</td>
		<td data-label="Status">Betalt</td>
		<td data-label="Forfaldsdato">01/12-2019</td>
		<td class="table--invoices__button">
			<a href="https://www.hiper.dk/mit-hiper/regninger/vis-regning/5098174.pdf" class="icon-link" target="_blank" rel="nofollow">
				<span class="text hide-on-large">Se regningen</span>
				<span title="Åbner en PDF i nyt vindue" class="icon icon-arrow">
					<svg> .. stripped .. </svg>
				</span>
			</a>
		</td>
	</tr>
	<tr>
		<td data-label="Regninsnummer">1644559</td>
		<td data-label="Type">Faktura</td>
		<td data-label="Total">369,60&nbsp;kr.</td>
		<td data-label="Status">Betalt</td>
		<td data-label="Forfaldsdato">16/11-2019</td>
		<td class="table--invoices__button">
			<a href="https://www.hiper.dk/mit-hiper/regninger/vis-regning/5047694.pdf" class="icon-link" target="_blank" rel="nofollow">
				<span class="text hide-on-large">Se regningen</span>
				<span title="Åbner en PDF i nyt vindue" class="icon icon-arrow">
					<svg> .. stripped .. </svg>
				</span>
			</a>
		</td>
	</tr>
</table>
*/
func (i *Invoices) FromReader(r io.Reader) error {

	// Load the HTML document
	doc, err := goquery.NewDocumentFromReader(r)
	if err != nil {
		return fmt.Errorf("unable create goquery document: %w", err)
	}

	skip := true
	doc.Find("tr").Each(func(_ int, s *goquery.Selection) {
		// move past the first element
		if skip {
			skip = false
			return
		}

		invoice := Invoice{}

		ID, err := strconv.Atoi(s.Find("td[data-label='Regninsnummer']").Text())
		if err != nil {
			panic(err)
		}
		invoice.ID = ID

		invoice.Type = s.Find("td[data-label='Type']").Text()
		invoice.Amount = s.Find("td[data-label='Total']").Text()
		invoice.Status = s.Find("td[data-label='Status']").Text()
		invoice.DueDate = s.Find("td[data-label='Forfaldsdato']").Text()

		URL, exists := s.Find("a").Attr("href")
		if !exists {
			panic(fmt.Sprintf("invoice had no pdf download url: %+v", invoice))
		}
		invoice.InvoiceURL = URL

		*i = append(*i, invoice)
	})

	return nil
}

func (i *Invoices) String() string {
	IDs := make([]string, len(*i))
	for index, invoice := range *i {
		IDs[index] = strconv.Itoa(invoice.ID)
	}
	return strings.Join(IDs, ", ")
}
