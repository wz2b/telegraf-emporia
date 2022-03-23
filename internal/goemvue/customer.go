package goemvue

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
)

func (t *EmVueCloudSession) getBasicCustomerInfo()(*CustomerInfo, error){

	u := NewApiUrl("/customers")
	q := u.Query()
	q.Set("email", "cpiggott@gmail.com")
	u.RawQuery = q.Encode()


	response, err := t.apiGet(u)
	if err != nil {
		return nil, err
	}

	bodyBytes, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}

	var info CustomerInfo
	err = json.Unmarshal(bodyBytes, &info)
	if err != nil {
		return nil, err
	}

	return &info, nil
}

func (t *EmVueCloudSession) getFullCustomerInfo(id int) (*CustomerInfoWithDevices, error) {
	var devices CustomerInfoWithDevices

	url := NewApiUrl(fmt.Sprintf("/customers/devices"))

	response, err := t.apiGet(url)
	if err != nil {
		return nil, err
	}

	bodyBytes, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(bodyBytes, &devices)
	if err != nil {
		return nil, err
	}


	return &devices, nil
}
