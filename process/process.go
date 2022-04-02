package process

import (
	"encoding/csv"
	"errors"
	"fmt"
	"io"
	"log"
	"mrsydar/tkl/k360/client"
	"mrsydar/tkl/k360/customer"
	"mrsydar/tkl/k360/invoice"
	"mrsydar/tkl/taxpayer"
	"os"
)

func getInvoiceFromRecord(record []string, customerId string) invoice.Invoice {
	return invoice.Invoice{
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
	}
}

func ProcessInvoices(client client.K360Client, csvPath string) error {
	file, err := os.Create("skipped_invoices.csv")
	if err != nil {
		return fmt.Errorf("failed to create file for skipped invoices: %v", err)
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	file, err = os.Open(csvPath)
	if err != nil {
		return fmt.Errorf("failed to open file: %v", err)
	}
	defer file.Close()

	reader := csv.NewReader(file)

	header, err := reader.Read()
	if err != nil {
		return fmt.Errorf("failed to skip header: %v", err)
	}
	writer.Write(header)

	taxpayerLoader := taxpayer.NewBufferedTaxpayerDataLoader()
	csvRecordsUnknownNipInvoices := make([][]string, 0)

	log.Println("start processing invoices without nip")

	for {
		record, err := reader.Read()
		if err != nil {
			if errors.Is(err, io.EOF) {
				break
			}
			return fmt.Errorf("failed to read record: %v", err)
		}

		nip := record[2]

		var customerId string
		if nip != "" {
			customerId, err = client.GetCustomerId(customer.Customer{Nip: nip})
			if err != nil {
				if errors.Is(err, customer.ErrNotFound) {
					err = taxpayerLoader.LoadTaxpayerData(nip)
					if err != nil {
						log.Printf("failed to load taxpayer data with nip %v: %v\n", nip, err)
					}

					csvRecordsUnknownNipInvoices = append(csvRecordsUnknownNipInvoices, record)
					continue
				} else {
					log.Printf("failed to get customer id with nip %v for invoice %v: %v\n", nip, record[0], err)
					writer.Write(record)
					continue
				}
			}
		} else {
			customerId = record[6]
		}

		invoice := getInvoiceFromRecord(record, customerId)

		err = client.PostInvoice(invoice)
		if err != nil {
			log.Printf("failed to post invoice %v: %v\n", invoice, err)
			writer.Write(record)
		}
	}

	log.Println("end processing invoices without nip")

	err = taxpayerLoader.Flush()
	if err != nil {
		log.Println("failed to flush taxpayer loader:", err)
	}

	log.Println("start processing invoices with nip")

	for _, record := range csvRecordsUnknownNipInvoices {
		if taxpayerLoader.RetrievedTaxpayers[record[2]] == nil {
			log.Printf("failed to get taxpayer info with nip %v for invoice %v\n", record[2], record[0])
		} else {
			taxpayer := taxpayerLoader.RetrievedTaxpayers[record[2]]

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

			customerId, err := client.PostCustomer(newCustomer)
			if err != nil {
				log.Printf("failed to post customer %v for invoice %v: %v", newCustomer, record[0], err)
				continue
			}

			invoice := getInvoiceFromRecord(record, customerId)

			err = client.PostInvoice(invoice)
			if err != nil {
				log.Printf("failed to post invoice %v: %v\n", invoice, err)
			}
		}
	}

	log.Println("end processing invoices with nip")

	return nil
}
