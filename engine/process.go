package engine

import (
	"encoding/csv"
	"fmt"
	"io"
	"log"
	"mrsydar/tkl/k360/client"
	"mrsydar/tkl/k360/customer"
	"mrsydar/tkl/k360/invoice"
	"mrsydar/tkl/taxpayer"
	"os"
)

func ProcessInvoices(client client.K360Client, csvPath string) error {
	file, err := os.Open(csvPath)
	if err != nil {
		return fmt.Errorf("failed to open file: %v", err)
	}
	defer file.Close()

	reader := csv.NewReader(file)

	prepareNextInvoice := func() (*invoice.Invoice, error) {
		record, err := reader.Read()
		if err != nil {
			return nil, fmt.Errorf("failed to read record: %v", err)
		}

		customerId, err := client.GetCustomerId(customer.Customer{Nip: record[2]})
		if err != nil {
			switch err.(type) {
			case *customer.NotFoundError:
				taxpayer, err := taxpayer.GetTaxpayerData(record[2])
				if err != nil {
					return nil, fmt.Errorf("failed to get taxpayer data with nip %v for invoice %v: %v", record[2], record[0], err)
				}

				newCustomer := customer.Customer{
					Name:        taxpayer.Name,
					Nip:         taxpayer.Nip,
					CountryCode: taxpayer.Address.CountryCode,
					Regon:       taxpayer.Regon,
					Street:      taxpayer.Address.Street,
					PostalCode:  taxpayer.Address.PostalCode,
					City:        taxpayer.Address.City,
					County:      taxpayer.Address.Country,
				}

				customerId, err = client.PostCustomer(newCustomer)
				if err != nil {
					return nil, fmt.Errorf("failed to post customer %v for invoice %v: %v", newCustomer, record[0], err)
				}
			default:
				return nil, fmt.Errorf("failed to get customer id with nip %v for invoice %v: %v", record[2], record[0], err)
			}
		}

		return &invoice.Invoice{
			Customer:        invoice.Customer{Id: customerId},
			DocDate:         record[1],
			DueDate:         record[1],
			TransactionDate: record[1],
			No:              record[0],
			Rows: []invoice.Row{
				{
					TaxId: record[5],
					Item: invoice.Item{
						Code:        record[7],
						Description: record[8],
					},
					Quantity: "1",
					Price:    record[3],
				},
			},
			TaxAmounts: []invoice.TaxAmount{
				{
					TaxId:  record[5],
					Amount: record[4],
				},
			},
			TotalAmount: record[3],
		}, nil
	}

	log.Println("Start processing invoices")
	for invoice, err := prepareNextInvoice(); err != io.EOF; invoice, err = prepareNextInvoice() {
		if err != nil {
			log.Printf("Failed to prepare invoice: %v\n", err)
		}

		err = client.PostInvoice(*invoice)
		if err != nil {
			log.Printf("Failed to post invoice: %v\n", err)
		}
	}
	log.Println("Finished processing invoices")

	return nil
}
