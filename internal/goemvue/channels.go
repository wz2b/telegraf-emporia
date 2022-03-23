package goemvue

func FlattenChannels(customer *CustomerInfoWithDevices) (channels ChannelList, err error) {
	for _, device := range customer.Devices {
		channels = appendDevicesChannels(device, channels)
	}
	return channels, err
}

func appendDevicesChannels(device Device, channels ChannelList) []Channel {
	/*
	 * Add any channels directly attached to this device
	 */
	channels = appendChannels(device.Channels, channels)

	/*
	 * Add any sub-devicess
	 */
	for _, subDevice := range device.Devices {
		channels = appendDevicesChannels(subDevice, channels)
	}
	return channels
}

func appendChannels(channels []Channel, channelList ChannelList) ChannelList {
	for _, channel := range channels {
		channelList = append(channelList, channel)
	}
	return channelList
}


func (channels *ChannelList) Filter(f func(Channel) bool) *ChannelList {

	var r ChannelList

	for _, channel := range(*channels) {
		if f(channel) {
			r = append(r, channel)
		}
	}

	return &r
}