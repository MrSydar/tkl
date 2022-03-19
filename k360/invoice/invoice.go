package invoice

type Customer struct {
	Id string `json:"id"`
}

type Item struct {
	Code        string `json:"Code"`
	Description string `json:"Description"`
}

type Row struct {
	TaxId    string `json:"TaxId"`
	Item     Item   `json:"Item"`
	Quantity string `json:"Quantity"`
	Price    string `json:"Price"`
}

type TaxAmount struct {
	TaxId  string `json:"TaxId"`
	Amount string `json:"Amount"`
}

type Invoice struct {
	Customer        Customer    `json:"Customer"`
	DocDate         string      `json:"DocDate"`
	DueDate         string      `json:"DueDate"`
	TransactionDate string      `json:"TransactionDate"`
	No              string      `json:"InvoiceNo"`
	Rows            []Row       `json:"InvoiceRow"`
	TaxAmounts      []TaxAmount `json:"TaxAmount"`
	TotalAmount     string      `json:"TotalAmount"`
}
