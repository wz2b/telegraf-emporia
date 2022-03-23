package main

import (
	"bytes"
	"flag"
	"fmt"
	lineProtocol "github.com/influxdata/line-protocol"
	"log"
	"os"
	"telegraf-vue2/internal/emitter"
	"telegraf-vue2/internal/goemvue"
	"telegraf-vue2/internal/state"
	"telegraf-vue2/internal/tclogger"
	"time"
)

var logWriter *tclogger.TelegrafCompatibleLogger = tclogger.Create().Start()
var agentState *state.AgentState

const MEASUREMENT_NAME = "power"

func main() {
	logWriter.Writer = os.Stderr

	user := flag.String("user", "", "emporia cloud login")
	password := flag.String("password", "", "emporia cloud password")

	flag.Parse()

	state := state.CreateAgentState()
	emp := goemvue.NewEmVueCloud(*user, *password)
	emp.DebugLog = log.New(logWriter, "", 0)

	metricEmitter := emitter.Create()

	err := emp.Initialize()
	if err != nil {
		log.Fatal(err)
	}

	//
	// Build a list of channels we care about
	//
	channels, err := goemvue.FlattenChannels(emp.CustomerInfo)

	fmt.Println(channels)

	namedChannels := channels.Filter(func(c goemvue.Channel) bool {
		return c.ChannelNum == "1,2,3" || c.Name != ""
	})

	/*
	 * Run
	 */

	for {

		for _, channel := range *namedChannels {

			var start = state.GetLastTime(channel.DeviceGid, channel.ChannelNum)
			if start == nil {
				newStart := time.Now().Add(-1 * time.Hour)
				start = &newStart
			}

			stop := start.Add(1 * time.Hour)

			hist, err := emp.GetChannelHistory(&channel, start, &stop)

			if err != nil {
				log.Print(err)
			} else {
				metrics, err := metricEmitter.HistToMetricList(&channel, hist, 1*time.Second)
				if err != nil {
					log.Print(err)
				} else {
					lastTime := metrics[len(metrics)-1].Time().Add(1 * time.Second)

					log.Printf("start=%s stop=%s next=%s\n", start.UTC(), stop.UTC(), lastTime.UTC())
					state.SetLastTime(channel.DeviceGid, channel.ChannelNum, &lastTime)

					dump(metrics)
					fmt.Printf("# Received %d metrics for %d/%s (%s)\n",
						len(metrics),
						channel.DeviceGid,
						channel.ChannelNum,
						channel.Name)
				}
			}
		}

		time.Sleep(5 * time.Minute)
	}

}

func dump(metrics []lineProtocol.MutableMetric) {
	buf := &bytes.Buffer{}
	serializer := lineProtocol.NewEncoder(buf)
	serializer.SetMaxLineBytes(-1)
	serializer.SetFieldTypeSupport(lineProtocol.UintSupport)

	for _, metric := range metrics {
		serializer.Encode(metric)
	}

	fmt.Print(buf.String())
}
