// Copyright 2024 Amazon.com, Inc. or its affiliates. All Rights Reserved.

package events

import (
	"encoding/json"
	"io/ioutil"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestUnmarshalStreamImage_ScalarFields(t *testing.T) {
	image := map[string]DynamoDBAttributeValue{
		"id":     NewStringAttribute("abc-123"),
		"count":  NewNumberAttribute("42"),
		"score":  NewNumberAttribute("3.14"),
		"active": NewBooleanAttribute(true),
		"tags":   NewStringSetAttribute([]string{"alpha", "beta"}),
	}

	type item struct {
		ID     string   `json:"id"`
		Count  int      `json:"count"`
		Score  float64  `json:"score"`
		Active bool     `json:"active"`
		Tags   []string `json:"tags"`
	}

	var got item
	require.NoError(t, UnmarshalStreamImage(image, &got))
	assert.Equal(t, "abc-123", got.ID)
	assert.Equal(t, 42, got.Count)
	assert.InDelta(t, 3.14, got.Score, 1e-9)
	assert.True(t, got.Active)
	assert.ElementsMatch(t, []string{"alpha", "beta"}, got.Tags)
}

func TestUnmarshalStreamImage_NestedMapAndList(t *testing.T) {
	image := map[string]DynamoDBAttributeValue{
		"user": NewMapAttribute(map[string]DynamoDBAttributeValue{
			"name": NewStringAttribute("Joe"),
			"age":  NewNumberAttribute("35"),
		}),
		"items": NewListAttribute([]DynamoDBAttributeValue{
			NewStringAttribute("Cookies"),
			NewStringAttribute("Coffee"),
		}),
	}

	type user struct {
		Name string `json:"name"`
		Age  int    `json:"age"`
	}
	type item struct {
		User  user     `json:"user"`
		Items []string `json:"items"`
	}

	var got item
	require.NoError(t, UnmarshalStreamImage(image, &got))
	assert.Equal(t, "Joe", got.User.Name)
	assert.Equal(t, 35, got.User.Age)
	assert.Equal(t, []string{"Cookies", "Coffee"}, got.Items)
}

func TestUnmarshalStreamImage_NullAttribute(t *testing.T) {
	image := map[string]DynamoDBAttributeValue{
		"deleted_at": NewNullAttribute(),
		"name":       NewStringAttribute("ada"),
	}

	type item struct {
		Name      string  `json:"name"`
		DeletedAt *string `json:"deleted_at"`
	}

	var got item
	require.NoError(t, UnmarshalStreamImage(image, &got))
	assert.Equal(t, "ada", got.Name)
	assert.Nil(t, got.DeletedAt)
}

func TestUnmarshalStreamImage_FromTestdata(t *testing.T) {
	raw, err := ioutil.ReadFile("./testdata/dynamodb-event.json")
	require.NoError(t, err)

	var evt DynamoDBEvent
	require.NoError(t, json.Unmarshal(raw, &evt))
	require.NotEmpty(t, evt.Records)

	first := evt.Records[0].Change.NewImage
	type partial struct {
		Val string `json:"val"`
		Key string `json:"key"`
	}
	var got partial
	require.NoError(t, UnmarshalStreamImage(first, &got))
	assert.Equal(t, "data", got.Val)
	assert.Equal(t, "binary", got.Key)
}

func TestToDynamoDBJSON_AttributeValueRoundTrip(t *testing.T) {
	av := NewStringAttribute("hello")
	raw, err := av.ToDynamoDBJSON()
	require.NoError(t, err)

	var decoded DynamoDBAttributeValue
	require.NoError(t, json.Unmarshal(raw, &decoded))
	assert.Equal(t, "hello", decoded.String())
}

func TestToDynamoDBJSONMap_PreservesShape(t *testing.T) {
	image := map[string]DynamoDBAttributeValue{
		"id":  NewStringAttribute("k1"),
		"qty": NewNumberAttribute("7"),
	}

	raw, err := ToDynamoDBJSONMap(image)
	require.NoError(t, err)

	var decoded map[string]DynamoDBAttributeValue
	require.NoError(t, json.Unmarshal(raw, &decoded))
	assert.Equal(t, "k1", decoded["id"].String())
	assert.Equal(t, "7", decoded["qty"].Number())
}

func TestStreamRecord_ToDynamoDBJSON_OmitsEmpty(t *testing.T) {
	rec := DynamoDBStreamRecord{
		Keys: map[string]DynamoDBAttributeValue{
			"pk": NewStringAttribute("k1"),
		},
	}

	raw, err := rec.ToDynamoDBJSON()
	require.NoError(t, err)

	var envelope map[string]json.RawMessage
	require.NoError(t, json.Unmarshal(raw, &envelope))
	_, hasKeys := envelope["Keys"]
	_, hasNew := envelope["NewImage"]
	_, hasOld := envelope["OldImage"]
	assert.True(t, hasKeys)
	assert.False(t, hasNew)
	assert.False(t, hasOld)
}

func TestStreamRecord_ToDynamoDBJSON_AllSections(t *testing.T) {
	rec := DynamoDBStreamRecord{
		Keys: map[string]DynamoDBAttributeValue{
			"pk": NewStringAttribute("k1"),
		},
		NewImage: map[string]DynamoDBAttributeValue{
			"pk":  NewStringAttribute("k1"),
			"qty": NewNumberAttribute("7"),
		},
		OldImage: map[string]DynamoDBAttributeValue{
			"pk":  NewStringAttribute("k1"),
			"qty": NewNumberAttribute("3"),
		},
	}

	raw, err := rec.ToDynamoDBJSON()
	require.NoError(t, err)

	var envelope struct {
		Keys     map[string]DynamoDBAttributeValue `json:"Keys"`
		NewImage map[string]DynamoDBAttributeValue `json:"NewImage"`
		OldImage map[string]DynamoDBAttributeValue `json:"OldImage"`
	}
	require.NoError(t, json.Unmarshal(raw, &envelope))
	assert.Equal(t, "k1", envelope.Keys["pk"].String())
	assert.Equal(t, "7", envelope.NewImage["qty"].Number())
	assert.Equal(t, "3", envelope.OldImage["qty"].Number())
}
