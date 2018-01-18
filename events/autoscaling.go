package events

import (
	"time"
)

// AutoScalingEvent struct is used to parse the json for auto scaling event types //
type AutoScalingEvent struct {
	Version    string       `json:"version"`     // The version of event data
	ID         string       `json:"id"`          // The unique ID of the event
	DetailType string       `json:"detail-type"` //Details about event type
	Source     string       `json:"source"`      //Source of the event
	AccountID  string       `json:"account"`     //AccountId
	Time       time.Time    `json:"time"`        //Event timestamp
	Region     string       `json:"region"`      //Region of event
	Resources  []string     `json:"resources"`   //Information about resources impacted by event
	Detail     EventDetails `json:"detail"`
}

// EventDetails struct is used to parse the nested event Details available within the Autoscaling event
type EventDetails struct {
	StatusCode           string      `json:"StatusCode,omitempty"`           //Status code for Event
	AutoscalingGroupName string      `json:"AutoScalingGroupName,omitempty"` //AutoScalingGroup name
	ActivityId           string      `json:"ActivityId,omitempty"`           //ActivityId
	Details              DetailsInfo `json:"Details,omitempty"`              //Additional details about event such as AZ and SubnetId of instance
	RequestId            string      `json:"RequestId,omitempty"`            //RequestId of Request
	StatusMessage        string      `json:"StatusMessage,omitempty"`        //StatusMEssage
	EndTime              time.Time   `json:"EndTime,omitempty"`              //EndTime of Event
	EC2InstanceId        string      `json:"EC2InstanceId,omitempty"`        //InstanceID of ec2 instance impacted by event
	StartTime            time.Time   `json:"StartTime,omitempty"`            //StartTime of the event
	Cause                string      `json:"Cause,omitempty"`                //Cause of event
	LifecycleActionToken string      `json:"LifecycleActionToken,omitempty"` //LifecycleActionToken details
	LifecycleHookName    string      `json:"LifecycleHookName,omitempty"`    //LifecycleHookName
	LifecycleTransition  string      `json:"LifecycleTransition,omitempty"`  //LifecycleTransition
	NotificationMetadata string      `json:"NotificationMetadata,omitempty"` //Notification metadata
}

// DetailsInfo contains AZ and SubnetId for instance impacted by the autoscaling event //
type DetailsInfo struct {
	AvailabilityZone string `json:"Availability Zone",omitempty` //AvailabilityZone for impacted instance
	SubnetID         string `json:"Subnet ID",omitempty`         //SubnetID for impacted instance
}
