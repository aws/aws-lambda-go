// Copyright 2020 Amazon.com, Inc. or its affiliates. All Rights Reserved.

package events

type ActiveMQEvent struct {
	EventSource    string              `json:"eventSource"`
	EventSourceArn string              `json:"eventSourceArn"`
	Messages       []ActiveMQMessage   `json:"messages"`
}

type ActiveMQMessage struct {
	MessageID       string          `json:"messageID"`
	MessageType     string          `json:"messageType"`
	Timestamp       int64           `json:"timestamp"`
	DeliveryMode    int             `json:"deliveryMode"`
	CorrelationID   string          `json:"correlationID"`
	ReplyTo         string          `json:"replyTo"`
	Destination     Destination     `json:"destination"`
	Redelivered     bool            `json:"redelivered"`
	Type            string          `json:"type"`
	Expiration      int64           `json:"expiration"`
	Priority        int             `json:"priority"`
	Data            string          `json:"data"`
	BrokerInTime    int64           `json:"brokerInTime"`
	BrokerOutTime   int64           `json:"brokerOutTime"`
}

type Destination struct {
	PhysicalName     string    `json:"physicalName"`
}