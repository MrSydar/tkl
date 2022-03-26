package client

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
)

func unmarshalBody(response http.Response, body interface{}) error {
	data, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return err
	}

	err = json.Unmarshal(data, body)
	if err != nil {
		return err
	}

	return nil
}
