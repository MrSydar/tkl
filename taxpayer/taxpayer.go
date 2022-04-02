package taxpayer

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"regexp"
	"strings"
	"time"
)

type Address struct {
	Street      string
	PostalCode  string
	City        string
	Country     string
	CountryCode string
}

type Taxpayer struct {
	Name    string
	Nip     string
	Regon   string
	Address *Address
}

type BufferedTaxpayerDataLoader struct {
	RetrievedTaxpayers map[string]*Taxpayer

	nipBuffer []string
}

func NewBufferedTaxpayerDataLoader() *BufferedTaxpayerDataLoader {
	return &BufferedTaxpayerDataLoader{
		RetrievedTaxpayers: make(map[string]*Taxpayer),
		nipBuffer:          make([]string, 0, 30),
	}
}

func (loader *BufferedTaxpayerDataLoader) LoadTaxpayerData(nip string) error {
	loader.nipBuffer = append(loader.nipBuffer, nip)
	if len(loader.nipBuffer) == cap(loader.nipBuffer) {
		err := loader.Flush()
		if err != nil {
			return fmt.Errorf("flush error: %v", err)
		}
	}
	return nil
}

func (loader *BufferedTaxpayerDataLoader) Flush() error {
	loc, _ := time.LoadLocation("Europe/Warsaw")
	plTime := time.Now().In(loc)

	if len(loader.nipBuffer) == 0 {
		return nil
	}

	url := fmt.Sprintf("https://wl-api.mf.gov.pl/api/search/nips/%s?date=%d-%02d-%02d", strings.Join(loader.nipBuffer, ","), plTime.Year(), plTime.Month(), plTime.Day())
	loader.nipBuffer = loader.nipBuffer[:0]

	response, err := http.Get(url)
	if err != nil {
		return err
	} else if response.StatusCode < 200 || response.StatusCode >= 300 {
		if responseBodyByteArray, err := ioutil.ReadAll(response.Body); err != nil {
			return fmt.Errorf("bad response with code %v", response.StatusCode)
		} else {
			return fmt.Errorf("bad response with code %v: %v", response.StatusCode, string(responseBodyByteArray))
		}
	}

	taxpayersJson := make(map[string]interface{})
	if err = json.NewDecoder(response.Body).Decode(&taxpayersJson); err != nil {
		return fmt.Errorf("can't decode response body: %v", err)
	}

	resultJson, ok := taxpayersJson["result"].(map[string]interface{})
	if !ok {
		return fmt.Errorf("can't find result in taxpayer data: %v", taxpayersJson)
	}

	entriesJson, ok := resultJson["entries"].([]interface{})
	if !ok {
		return fmt.Errorf("can't find entries in result: %v", resultJson)
	}

	for entryJson := range entriesJson {
		subjectsJson, ok := entriesJson[entryJson].([]interface{})
		if !ok || len(subjectsJson) != 1 {
			log.Println("ignoring taxpayer: can't find subjects in entry:", entriesJson[entryJson])
			continue
		}

		subjectJson, ok := subjectsJson[0].(map[string]interface{})
		if !ok {
			log.Println("ignoring taxpayer: can't find first subject in subjects:", subjectsJson)
			continue
		}

		name, ok := subjectJson["name"].(string)
		if !ok || name == "" {
			log.Println("ignoring taxpayer: can't find name in subject:", subjectJson)
			continue
		}

		regon, ok := subjectJson["regon"].(string)
		if !ok || regon == "" {
			log.Println("ignoring taxpayer: can't find regon in subject:", subjectJson)
			continue
		}

		nip, ok := subjectJson["nip"].(string)
		if !ok || nip == "" {
			log.Println("ignoring taxpayer: can't find nip in subject:", subjectJson)
			continue
		}

		rawAddress, ok := subjectJson["workingAddress"].(string)
		if !ok || rawAddress == "" {
			log.Println("ignoring taxpayer: can't find workingAddress in subject:", subjectJson)
			continue
		}

		address, err := parseAddress(rawAddress)
		if err != nil {
			log.Println("ignoring taxpayer: can't parse address:", err)
			continue
		}

		if isPolishAddress(rawAddress) {
			address.CountryCode = "PL"
			address.Country = "POLSKA"
		} else {
			log.Println("ignoring taxpayer: address is not polish:", rawAddress)
			continue
		}

		loader.RetrievedTaxpayers[nip] = &Taxpayer{name, nip, regon, address}
	}

	return nil
}

func isPolishAddress(address string) bool {
	lowerAddress := strings.ToLower(address)
	for _, city := range plCities {
		if strings.Contains(lowerAddress, city) {
			return true
		}
	}
	return false
}

var postalCodeRegex = regexp.MustCompile("[0-9]{2}-[0-9]{3}")
var cityRegex = regexp.MustCompile("[0-9]{2}-[0-9]{3} (.*)")

func parseAddress(address string) (*Address, error) {
	splittedAddress := strings.Split(address, ",")
	if len(splittedAddress) != 2 {
		return nil, fmt.Errorf("too many splitted address parts: %v", len(splittedAddress))
	}
	street := splittedAddress[0]

	postalCode := postalCodeRegex.FindString(splittedAddress[1])
	if len(postalCode) != 6 {
		return nil, fmt.Errorf("can't extract postal code")
	}

	city, err := getFirstSubgroupMatch(splittedAddress[1], cityRegex)
	if err != nil {
		return nil, fmt.Errorf("can't extract city: %v", err)
	}

	return &Address{Street: street, PostalCode: postalCode, City: city}, nil
}

func GetTaxpayerData(nip string) (*Taxpayer, error) {
	loc, _ := time.LoadLocation("Europe/Warsaw")
	plTime := time.Now().In(loc)

	url := fmt.Sprintf("https://wl-api.mf.gov.pl/api/search/nip/%s?date=%d-%02d-%02d", nip, plTime.Year(), plTime.Month(), plTime.Day())

	response, err := http.Get(url)
	if err != nil {
		return nil, err
	} else if response.StatusCode < 200 || response.StatusCode >= 300 {
		if responseBodyByteArray, err := ioutil.ReadAll(response.Body); err != nil {
			return nil, fmt.Errorf("bad response with code %v", response.StatusCode)
		} else {
			return nil, fmt.Errorf("bad response with code %v: %v", response.StatusCode, string(responseBodyByteArray))
		}
	}

	taxpayerJson := make(map[string]interface{})
	if err = json.NewDecoder(response.Body).Decode(&taxpayerJson); err != nil {
		return nil, err
	}

	resultJson, ok := taxpayerJson["result"].(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("can't find result")
	}

	subjectJson, ok := resultJson["subject"].(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("can't find subject")
	}

	name, ok := subjectJson["name"].(string)
	if !ok || name == "" {
		return nil, fmt.Errorf("can't find name")
	}

	regon, ok := subjectJson["regon"].(string)
	if !ok || regon == "" {
		return nil, fmt.Errorf("can't find name")
	}

	rawAddress, ok := subjectJson["workingAddress"].(string)
	if !ok || rawAddress == "" {
		return nil, fmt.Errorf("can't find workingAddress")
	}

	address, err := parseAddress(rawAddress)
	if err != nil {
		return nil, fmt.Errorf("can't parse address")
	}

	if isPolishAddress(rawAddress) {
		address.CountryCode = "PL"
		address.Country = "POLSKA"
	} else {
		return nil, fmt.Errorf("can't confirm customer country: countries other than Poland are not supported by this application")
	}

	return &Taxpayer{name, nip, regon, address}, nil
}
