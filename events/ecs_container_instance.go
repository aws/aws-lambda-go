package events

import (
	"encoding/json"
	"time"
)

type ECSContainerInstanceEvent struct {
	Version    string                              `json:"version"`
	ID         string                              `json:"id"`
	DetailType string                              `json:"detail-type"`
	Source     string                              `json:"source"`
	Account    string                              `json:"account"`
	Time       time.Time                           `json:"time"`
	Region     string                              `json:"region"`
	Resources  []string                            `json:"resources"`
	Detail     ECSContainerInstanceEventDetailType `json:"detail"`
}

type ECSContainerInstanceEventDetailType struct {
	AgentConnected       bool        `json:"agentConnected"`
	Attributes           []Attribute `json:"attributes"`
	ClusterARN           string      `json:"clusterArn"`
	ContainerInstanceARN string      `json:"containerInstanceArn"`
	EC2InstanceID        string      `json:"ec2InstanceId"`
	RegisteredResources  []Resource  `json:"registeredResources"`
	RemainingResources   []Resource  `json:"remainingResources"`
	Status               string      `json:"status"`
	Version              int         `json:"version"`
	VersionInfo          VersionInfo `json:"versionInfo"`
	UpdatedAt            time.Time   `json:"updatedAt"`
}

type Attribute struct {
	Name string `json:"name"`
}

type Resource struct {
	Name           string    `json:"name"`
	Type           string    `json:"type"`
	IntegerValue   int       `json:"integerValue,omitempty"`
	StringSetValue []*string `json:"stringSetValue,omitempty"`
}

type VersionInfo struct {
	AgentHash     string `json:"agentHash"`
	AgentVersion  string `json:"agentVersion"`
	DockerVersion string `json:"dockerVersion"`
}

// MarshalJSON implements cuustom marshaling to marshal the struct into JSON format while preserving an empty string slice in `StringSetValue` field.
func (r Resource) MarshalJSON() ([]byte, error) {
	type Alias Resource
	aux := struct {
		StringSetValue json.RawMessage `json:"stringSetValue,omitempty"`
		Alias
	}{
		Alias: (Alias)(r),
	}

	if r.StringSetValue != nil {
		b, err := json.Marshal(r.StringSetValue)
		if err != nil {
			return nil, err
		}
		aux.StringSetValue = b
	}

	return json.Marshal(&aux)
}
