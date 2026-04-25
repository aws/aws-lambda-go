// Copyright 2024 Amazon.com, Inc. or its affiliates. All Rights Reserved.

package events

import (
	"encoding/json"
	"fmt"
)

// UnmarshalStreamImage unmarshals a stream image (a map of DynamoDBAttributeValue
// keyed by attribute name, as found in DynamoDBStreamRecord.Keys, NewImage and
// OldImage) into the provided destination value.
//
// The destination is decoded using the standard encoding/json package, so
// `json` struct tags work as expected. This sidesteps the long-standing
// limitation that DynamoDBAttributeValue is not directly compatible with the
// dynamodbattribute helpers shipped by the AWS SDK (see issue #58).
//
// Typical usage:
//
//	type MyItem struct {
//	    ID   string `json:"id"`
//	    Name string `json:"name"`
//	}
//
//	var item MyItem
//	if err := events.UnmarshalStreamImage(record.Change.NewImage, &item); err != nil {
//	    return err
//	}
//
// Note that this helper does NOT honor `dynamodbav` tags used by the AWS SDK's
// dynamodbattribute package. For SDK-tag-aware decoding, use ToDynamoDBJSON to
// emit canonical DynamoDB JSON and feed it to your SDK of choice.
func UnmarshalStreamImage(image map[string]DynamoDBAttributeValue, out interface{}) error {
	flat := make(map[string]interface{}, len(image))
	for k, v := range image {
		raw, err := flattenAttributeValue(v)
		if err != nil {
			return fmt.Errorf("UnmarshalStreamImage: %q: %w", k, err)
		}
		flat[k] = raw
	}

	encoded, err := json.Marshal(flat)
	if err != nil {
		return fmt.Errorf("UnmarshalStreamImage: encode: %w", err)
	}
	if err := json.Unmarshal(encoded, out); err != nil {
		return fmt.Errorf("UnmarshalStreamImage: decode: %w", err)
	}
	return nil
}

// ToDynamoDBJSON returns the canonical DynamoDB JSON wire form of an attribute
// value (for example {"S":"hello"} or {"N":"123"}). The returned bytes are
// directly compatible with json.Unmarshal-ing into the AttributeValue type
// defined by either aws-sdk-go (service/dynamodb.AttributeValue) or
// aws-sdk-go-v2 (service/dynamodb/types.AttributeValueMemberX), via the
// standard encoding/json package.
//
// This avoids forcing aws-lambda-go to take a hard dependency on either SDK
// version while still giving callers a stable bridge into SDK types.
func (av DynamoDBAttributeValue) ToDynamoDBJSON() ([]byte, error) {
	return av.MarshalJSON()
}

// ToDynamoDBJSONMap returns the canonical DynamoDB JSON wire form of an
// attribute-value map, suitable for json.Unmarshal-ing into
// map[string]*dynamodb.AttributeValue (SDK v1) or
// map[string]types.AttributeValue (SDK v2).
func ToDynamoDBJSONMap(image map[string]DynamoDBAttributeValue) ([]byte, error) {
	return json.Marshal(image)
}

// ToDynamoDBJSON returns the canonical DynamoDB JSON wire form of this stream
// record's Keys, NewImage and OldImage, structured as a top-level object with
// those three keys (any of which may be omitted when empty). This is a
// convenience for callers who want a single round-trip into SDK types:
//
//	raw, _ := record.Change.ToDynamoDBJSON()
//	// json.Unmarshal(raw, &someSDKShape)
func (r DynamoDBStreamRecord) ToDynamoDBJSON() ([]byte, error) {
	envelope := struct {
		Keys     map[string]DynamoDBAttributeValue `json:"Keys,omitempty"`
		NewImage map[string]DynamoDBAttributeValue `json:"NewImage,omitempty"`
		OldImage map[string]DynamoDBAttributeValue `json:"OldImage,omitempty"`
	}{
		Keys:     r.Keys,
		NewImage: r.NewImage,
		OldImage: r.OldImage,
	}
	return json.Marshal(envelope)
}

// flattenAttributeValue converts a DynamoDBAttributeValue into a plain Go
// value (string, float64, bool, []byte, []interface{}, map[string]interface{}
// or nil) so it can be re-encoded as ordinary JSON for downstream
// json.Unmarshal calls. Numbers are decoded with json.Number to preserve
// precision when the destination uses json.Number / json.Decoder.UseNumber.
func flattenAttributeValue(av DynamoDBAttributeValue) (interface{}, error) {
	switch av.DataType() {
	case DataTypeNull:
		return nil, nil
	case DataTypeString:
		return av.String(), nil
	case DataTypeNumber:
		return json.Number(av.Number()), nil
	case DataTypeBoolean:
		return av.Boolean(), nil
	case DataTypeBinary:
		// Mirror DynamoDB's wire shape: binaries are base64 strings on the wire.
		return av.Binary(), nil
	case DataTypeStringSet:
		return av.StringSet(), nil
	case DataTypeNumberSet:
		ns := av.NumberSet()
		out := make([]json.Number, len(ns))
		for i, n := range ns {
			out[i] = json.Number(n)
		}
		return out, nil
	case DataTypeBinarySet:
		return av.BinarySet(), nil
	case DataTypeList:
		list := av.List()
		out := make([]interface{}, len(list))
		for i, item := range list {
			v, err := flattenAttributeValue(item)
			if err != nil {
				return nil, err
			}
			out[i] = v
		}
		return out, nil
	case DataTypeMap:
		m := av.Map()
		out := make(map[string]interface{}, len(m))
		for k, item := range m {
			v, err := flattenAttributeValue(item)
			if err != nil {
				return nil, err
			}
			out[k] = v
		}
		return out, nil
	default:
		return nil, fmt.Errorf("unsupported DynamoDB data type %v", av.DataType())
	}
}
