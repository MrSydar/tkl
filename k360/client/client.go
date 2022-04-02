package client

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"time"

	"mrsydar/tkl/k360/customer"
	"mrsydar/tkl/k360/invoice"
)

type K360Client struct {
	apiId  string
	apiKey string
}

func New(apiId, apiKey string) *K360Client {
	return &K360Client{apiId, apiKey}
}

func (client *K360Client) GetCustomerId(data customer.Customer) (string, error) {
	url := url.URL{
		Scheme:   "https",
		Host:     "program.360ksiegowosc.pl",
		Path:     "api/v1/getcustomers",
		RawQuery: fmt.Sprintf("ApiId=%s", client.apiId),
	}

	response, err := client.post(url, data)
	if err != nil {
		return "", err
	}

	foundCustomers := []struct {
		Id string `json:"CustomerId"`
	}{}

	err = unmarshalBody(*response, &foundCustomers)
	if err != nil {
		return "", err
	}

	if len(foundCustomers) != 1 {
		if len(foundCustomers) == 0 {
			return "", customer.ErrNotFound
		} else {
			return "", errors.New("too many customers found")
		}
	}

	return foundCustomers[0].Id, nil
}

func (client *K360Client) PostCustomer(data customer.Customer) (string, error) {
	url := url.URL{
		Scheme:   "https",
		Host:     "program.360ksiegowosc.pl",
		Path:     "api/v2/sendcustomer",
		RawQuery: fmt.Sprintf("ApiId=%s", client.apiId),
	}

	response, err := client.post(url, data)
	if err != nil {
		return "", err
	}

	addedCustomer := struct {
		Id string `json:"Id"`
	}{}

	err = unmarshalBody(*response, &addedCustomer)
	if err != nil {
		return "", err
	}

	return addedCustomer.Id, nil
}

func (client *K360Client) PostInvoice(invoiceData invoice.Invoice) error {
	url := url.URL{
		Scheme:   "https",
		Host:     "program.360ksiegowosc.pl",
		Path:     "api/v1/sendinvoice",
		RawQuery: fmt.Sprintf("ApiId=%s", client.apiId),
	}

	_, err := client.post(url, invoiceData)
	if err != nil {
		return err
	}

	return nil
}

func (client *K360Client) post(url url.URL, data interface{}) (*http.Response, error) {
	jsonBody, err := json.Marshal(data)
	if err != nil {
		return nil, err
	}

	timestampf := time.Now().Format("20060102150405")

	h := hmac.New(sha256.New, []byte(client.apiKey))
	h.Write([]byte(client.apiId + timestampf + string(jsonBody)))
	signature := hex.EncodeToString(h.Sum(nil))

	if len(url.RawQuery) != 0 {
		url.RawQuery += "&"
	}

	url.RawQuery += fmt.Sprintf("timestamp=%s&signature=%s", timestampf, signature)

	response, err := http.Post(url.String(), "application/json", bytes.NewBuffer(jsonBody))
	if err != nil {
		return nil, err
	}

	if response.StatusCode != 200 {
		body, err := ioutil.ReadAll(response.Body)
		if err != nil {
			return nil, fmt.Errorf("bad response: code: %v", response.StatusCode)
		} else {
			return nil, fmt.Errorf("bad response: code: %v, body: %q", response.StatusCode, string(body))
		}
	}

	return response, nil
}
