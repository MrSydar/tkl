package taxpayer

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"regexp"
	"strings"
	"time"

	"mrsydar/tkl/util"
)

type Address struct {
	street     string
	postalCode string
	city       string
}

type Taxpayer struct {
	name    string
	nip     string
	regon   string
	address *Address
	country string
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

	city, err := util.GetFirstSubgroupMatch(splittedAddress[1], cityRegex)
	if err != nil {
		return nil, fmt.Errorf("can't extract city: %v", err)
	}

	return &Address{street, postalCode, city}, nil
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
			return nil, fmt.Errorf("bad response with code %v: %v", response.StatusCode, string(responseBodyByteArray))
		} else {
			return nil, fmt.Errorf("bad response with code %v", response.StatusCode)
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

	var country string
	if isPolishAddress(rawAddress) {
		country = "PL"
	} else {
		return nil, fmt.Errorf("can't confirm customer country: countries other than Poland are not supported by this application")
	}

	return &Taxpayer{name, nip, regon, address, country}, nil
}
