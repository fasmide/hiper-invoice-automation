package kviknet

import (
	"fmt"
)

// Invoice represents an invoice that can be downloaded
type Invoice struct {
	ID         string
	Type       string
	Amount     string
	Status     string
	DueDate    string
	InvoiceURL string
}

// String representation of this invoice
// "1660653" "Faktura" "308,00 kr." "Betalt" "01/12-2019"
func (i *Invoice) String() string {
	return fmt.Sprintf("%s %s %s %s %s %s",
		i.ID,
		i.Type,
		i.Amount,
		i.Status,
		i.DueDate,
		i.InvoiceURL,
	)
}
