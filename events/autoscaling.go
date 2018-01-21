package events

import (
	"time"
)

// AutoScalingEvent struct is used to parse the json for auto scaling event types //
type AutoScalingEvent struct {
	Version    string                  `json:"version"`     // The version of event data
	ID         string                  `json:"id"`          // The unique ID of the event
	DetailType string                  `json:"detail-type"` //Details about event type
	Source     string                  `json:"source"`      //Source of the event
	AccountID  string                  `json:"account"`     //AccountId
	Time       time.Time               `json:"time"`        //Event timestamp
	Region     string                  `json:"region"`      //Region of event
	Resources  []string                `json:"resources"`   //Information about resources impacted by event
	Detail     AutoScalingEventDetails `json:"detail"`
}

// AutoScalingEventDetails struct is used to parse the nested event Details available within the Autoscaling event
type AutoScalingEventDetails struct {
	StatusCode           string                       `json:"StatusCode,omitempty"`           //Status code for Event
	AutoScalingGroupName string                       `json:"AutoScalingGroupName,omitempty"` //AutoScalingGroup name
	ActivityID           string                       `json:"ActivityId,omitempty"`           //ActivityId
	Details              *AutoScalingEventDetailsInfo `json:"Details,omitempty"`              //Additional details about event such as AZ and SubnetId of instance
	RequestID            string                       `json:"RequestId,omitempty"`            //RequestId of Request
	StatusMessage        *string                      `json:"StatusMessage,omitempty"`        //StatusMessage
	EndTime              *time.Time                   `json:"EndTime,omitempty"`              //EndTime of Event
	EC2InstanceID        string                       `json:"EC2InstanceId,omitempty"`        //InstanceID of ec2 instance impacted by event
	StartTime            *time.Time                   `json:"StartTime,omitempty"`            //StartTime of the event
	Cause                string                       `json:"Cause,omitempty"`                //Cause of event
	LifecycleActionToken string                       `json:"LifecycleActionToken,omitempty"` //LifecycleActionToken details
	LifecycleHookName    string                       `json:"LifecycleHookName,omitempty"`    //LifecycleHookName
	LifecycleTransition  string                       `json:"LifecycleTransition,omitempty"`  //LifecycleTransition
	NotificationMetadata string                       `json:"NotificationMetadata,omitempty"` //Notification metadata
	Description          string                       `json:"Description,omitempty"`          //Description in detail block
}

// AutoScalingEventDetailsInfo contains AZ and SubnetId for instance impacted by the autoscaling event //
type AutoScalingEventDetailsInfo struct {
	AvailabilityZone string `json:"Availability Zone",omitempty` //AvailabilityZone for impacted instance
	SubnetID         string `json:"Subnet ID",omitempty`         //SubnetID for impacted instance
}
