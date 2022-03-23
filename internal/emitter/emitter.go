package emitter

import (
	lineProtocol "github.com/influxdata/line-protocol"
	"math"
	"strings"
	"time"
	"telegraf-vue2/internal/goemvue"
)

type VueEmitter struct {
	MeasurementName string
}

func Create() *VueEmitter {
	return &VueEmitter{
		MeasurementName: "power",
	}
}

func generateFieldName(channel *goemvue.Channel) string {
	if channel.Name != "" {
		fieldName := strings.ReplaceAll(channel.Name, " ", "")
		return strings.TrimSpace(fieldName)
	} else if channel.ChannelNum == "1,2,3" {
		return "total"
	} else {
		return "ch" + channel.ChannelNum
	}
}

func (v *VueEmitter) HistToMetricList(channel *goemvue.Channel, hist *goemvue.History, slot time.Duration) ([]lineProtocol.MutableMetric, error) {
	var metrics = make([]lineProtocol.MutableMetric, 0)

	fieldName := generateFieldName(channel)

	/*
	 * The number that comes in is units of kWh, so turning that into watts requires
	 * multiplying the time interval by one hour.  Compute a scaling factor once.
	 */
	timeSlotScaleFactor := (time.Hour.Seconds() / slot.Seconds())

	for index, value := range hist.Data {
		if value != nil {

			// Convert from kilowatts to watts
			fValue := *value * timeSlotScaleFactor
			fValue = fValue * 1000.0
			fValue =  math.Round(fValue * 1000.0)/1000.0
			recordTime := hist.Start.Add(time.Duration(index) * time.Second)
			tags := make(map[string]string)
			fields := make(map[string]interface{})
			fields[fieldName] = fValue
			metric, err := lineProtocol.New(v.MeasurementName, tags, fields, recordTime)
			if err != nil {
				return metrics, nil
			} else {
				metrics = append(metrics, metric)
			}
		}
	}
	return metrics, nil

}
