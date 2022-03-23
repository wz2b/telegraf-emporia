package goemvue

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"strconv"
	"time"
)


func (t *EmVueCloudSession) GetChannelHistory(channel *Channel, start *time.Time, end *time.Time) (*History, error) {
	u := NewApiUrl("/AppAPI")
	q := u.Query()
	q.Set("apiMethod", "getChartUsage")
	q.Set("deviceGid", strconv.Itoa(channel.DeviceGid))
	q.Set("channel", channel.ChannelNum)
	q.Set("start", start.UTC().Format("2006-01-02T15:04:05Z"))
	q.Set("end", end.UTC().Format("2006-01-02T15:04:05Z"))
	q.Set("scale", "1S")
	q.Set("energyUnit", "KilowattHours")
	u.RawQuery = q.Encode()

	response, err := t.apiGet(u)

	if err != nil {
		return nil, err
	}

	responseBytes, _ := ioutil.ReadAll(response.Body)
	responseString := string(responseBytes)

	if response.StatusCode != 200 {
		return nil, errors.New(fmt.Sprintf("HTTP %d %s %s\n", response.StatusCode, response.Status, responseString))
	}

	var history History
	err = json.Unmarshal(responseBytes, &history)
	if err != nil {
		return nil, err
	}

	return &history, nil
}
