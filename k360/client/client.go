package client

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"time"
)

type K360Client struct {
	apiId  string
	apiKey string
}

func New(apiId, apiKey string) *K360Client {
	return &K360Client{apiId, apiKey}
}

func (client *K360Client) GetCustomerId(name, nip string) (string, error) {
	customerData := struct {
		Name string `json:"Name"`
		Nip  string `json:"vatRegNo"`
	}{name, nip}

	customerJson, err := json.Marshal(customerData)
	if err != nil {
		return "", err
	}

	url := url.URL{
		Scheme:   "https",
		Host:     "program.360ksiegowosc.pl",
		Path:     "api/v1/getcustomers",
		RawQuery: fmt.Sprintf("ApiId=%s", client.apiId),
	}

	response, err := client.post(url, string(customerJson))
	if err != nil {
		return "", err
	}

	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return "", err
	}

	foundCustomers := []struct {
		Id string `json:"CustomerId"`
	}{}

	err = json.Unmarshal(body, &foundCustomers)
	if err != nil {
		return "", err
	}

	if len(foundCustomers) != 1 {
		if len(foundCustomers) == 0 {
			return "", &CustomerNotFoundError{fmt.Sprintf("no customers found with name %q and nip %q", name, nip)}
		} else {
			return "", fmt.Errorf("too many customers with name %q and nip %q", name, nip)
		}
	}

	return foundCustomers[0].Id, nil
}

func (client *K360Client) PostCustomer(name, nip, countryCode, regon, street, postalCode, city, county string) (string, error) {
	customerData := struct {
		Name        string `json:"Name"`
		Nip         string `json:"VatRegNo"`
		CountryCode string `json:"CountryCode"`
		Regon       string `json:"RegNo,omitempty"`
		Street      string `json:"Address,omitempty"`
		PostalCode  string `json:"PostalCode,omitempty"`
		City        string `json:"City,omitempty"`
		County      string `json:"County,omitempty"`
	}{name, nip, countryCode, regon, street, postalCode, city, county}

	customerJson, err := json.Marshal(customerData)
	if err != nil {
		return "", err
	}

	url := url.URL{
		Scheme:   "https",
		Host:     "program.360ksiegowosc.pl",
		Path:     "api/v2/sendcustomer",
		RawQuery: fmt.Sprintf("ApiId=%s", client.apiId),
	}

	response, err := client.post(url, string(customerJson))
	if err != nil {
		return "", err
	}

	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return "", err
	}

	addedCustomer := struct {
		Id string `json:"Id"`
	}{}

	err = json.Unmarshal(body, &addedCustomer)
	if err != nil {
		return "", err
	}

	return addedCustomer.Id, nil
}

func (client *K360Client) post(url url.URL, jsonBody string) (*http.Response, error) {
	timestampf := time.Now().Format("20060102150405")

	h := hmac.New(sha256.New, []byte(client.apiKey))
	h.Write([]byte(client.apiId + timestampf + jsonBody))
	signature := hex.EncodeToString(h.Sum(nil))

	if len(url.RawQuery) != 0 {
		url.RawQuery += "&"
	}

	url.RawQuery += fmt.Sprintf("timestamp=%s&signature=%s", timestampf, signature)

	response, err := http.Post(url.String(), "application/json", bytes.NewBuffer([]byte(jsonBody)))
	if err != nil {
		return nil, err
	}

	if response.StatusCode != 200 {
		body, err := ioutil.ReadAll(response.Body)
		if err != nil {
			return nil, fmt.Errorf("bad response: code: %v", response.StatusCode)
		} else {
			return nil, fmt.Errorf("bad response: code: %v, body: %q", response.StatusCode, body)
		}
	}

	return response, nil
}
