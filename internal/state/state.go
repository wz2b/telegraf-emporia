package state

import "time"

type ChannelLastTimeMap map[string](*time.Time)

type AgentState struct {
	deviceLastTimeMap map[int](*ChannelLastTimeMap)
}

func CreateAgentState() *AgentState {
	return &AgentState{
		deviceLastTimeMap: make(map[int](*ChannelLastTimeMap), 0),
	}
}


func (a *AgentState) GetLastTime(deviceGid int, channel string) *time.Time {
	if a.deviceLastTimeMap[deviceGid] == nil {
		newMap := make(ChannelLastTimeMap)
		a.deviceLastTimeMap[deviceGid] = &newMap
	}
	return (*a.deviceLastTimeMap[deviceGid])[channel]
}

func (a *AgentState) SetLastTime(deviceGid int, channel string, newTime *time.Time) {
	if a.deviceLastTimeMap[deviceGid] == nil {
		newMap := make(ChannelLastTimeMap)
		a.deviceLastTimeMap[deviceGid] = &newMap
	}
	(*a.deviceLastTimeMap[deviceGid])[channel] = newTime
}