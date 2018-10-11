// Copyright 2017 Amazon.com, Inc. or its affiliates. All Rights Reserved.

package events

type IotButtonEvent struct {
	SrialNumber    string `json:"serialNumber"`
	ClickType      string `json:"clickType"`
	BatteryVoltage string `json:"batteryVoltage"`
}
