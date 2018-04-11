// Copyright 2017 Amazon.com, Inc. or its affiliates. All Rights Reserved.

package events

import (
	"encoding/base64"
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestUnmarshalBinary(t *testing.T) {
	input := []byte(`{ "B": "AAEqQQ=="}`)

	var av DynamoDBAttributeValue
	err := json.Unmarshal(input, &av)

	assert.Nil(t, err)
	assert.Equal(t, DataTypeBinary, av.DataType())
	assert.Equal(t, "AAEqQQ==", base64.StdEncoding.EncodeToString(av.Binary()))
}

func TestUnmarshalBoolean(t *testing.T) {
	input := []byte(`{ "BOOL": true}`)

	var av DynamoDBAttributeValue
	err := json.Unmarshal(input, &av)

	assert.Nil(t, err)
	assert.Equal(t, DataTypeBoolean, av.DataType())
	assert.Equal(t, true, av.Boolean())
}

func TestUnmarshalBinarySet(t *testing.T) {
	input := []byte(`{ "BS": ["AAEqQQ==", "AAEqQQ=="] }`)

	var av DynamoDBAttributeValue
	err := json.Unmarshal(input, &av)

	assert.Nil(t, err)
	assert.Equal(t, DataTypeBinarySet, av.DataType())
	assert.Equal(t, 2, len(av.BinarySet()))
	assert.Equal(t, "AAEqQQ==", base64.StdEncoding.EncodeToString(av.BinarySet()[1]))
}

func TestUnmarshalList(t *testing.T) {
	input := []byte(`{ "L": [
            { "S": "Cookies"},
            { "S": "Coffee"},
            { "N": "3.14159" }
        ] }`)

	var av DynamoDBAttributeValue
	err := json.Unmarshal(input, &av)

	assert.Nil(t, err)
	assert.Equal(t, DataTypeList, av.DataType())
	assert.Equal(t, 3, len(av.List()))
	assert.Equal(t, "Cookies", av.List()[0].String())
	assert.Equal(t, "3.14159", av.List()[2].Number())
}

func TestUnmarshalMap(t *testing.T) {
	input := []byte(`
        { "M":
            {
                "Name": { "S": "Joe" },
                "Age":  { "N": "35" }
            }
        }`)

	var av DynamoDBAttributeValue
	err := json.Unmarshal(input, &av)

	assert.Nil(t, err)
	assert.Equal(t, DataTypeMap, av.DataType())
	assert.Equal(t, 2, len(av.Map()))
	assert.Equal(t, "Joe", av.Map()["Name"].String())
	assert.Equal(t, "35", av.Map()["Age"].Number())
}

func TestUnmarshalNumber(t *testing.T) {
	input := []byte(`{ "N": "123.45"}`)

	var av DynamoDBAttributeValue
	err := json.Unmarshal(input, &av)

	assert.Nil(t, err)
	assert.Equal(t, DataTypeNumber, av.DataType())
	assert.Equal(t, "123.45", av.Number())
}

func TestUnmarshalInteger(t *testing.T) {
	input := []byte(`{ "N": "123"}`)

	var av DynamoDBAttributeValue
	err := json.Unmarshal(input, &av)

	assert.Nil(t, err)

	var i int64
	i, err = av.Integer()
	assert.Nil(t, err)
	assert.Equal(t, int64(123), i)
}

func TestUnmarshalFloat(t *testing.T) {
	input := []byte(`{ "N": "123.45"}`)

	var av DynamoDBAttributeValue
	err := json.Unmarshal(input, &av)
	assert.Nil(t, err)

	var f float64
	f, err = av.Float()

	assert.Nil(t, err)
	assert.Equal(t, 123.45, f)
}

func TestUnmarshalIntContainingAFloatString(t *testing.T) {
	input := []byte(`{ "N": "123.45"}`)

	var av DynamoDBAttributeValue
	err := json.Unmarshal(input, &av)

	assert.Nil(t, err)

	var i int64
	i, err = av.Integer()
	assert.Nil(t, err)
	assert.Equal(t, int64(123), i)
}

func TestUnmarshalNumberSet(t *testing.T) {
	input := []byte(`{ "NS": ["1234", "567.8"] }`)

	var av DynamoDBAttributeValue
	err := json.Unmarshal(input, &av)

	assert.Nil(t, err)
	assert.Equal(t, DataTypeNumberSet, av.DataType())
	assert.Equal(t, 2, len(av.NumberSet()))
	assert.Equal(t, "1234", av.NumberSet()[0])
	assert.Equal(t, "567.8", av.NumberSet()[1])
}

func TestUnmarshalNull(t *testing.T) {
	input := []byte(`{ "NULL": true}`)

	var av DynamoDBAttributeValue
	err := json.Unmarshal(input, &av)

	assert.Nil(t, err)
	assert.Equal(t, DataTypeNull, av.DataType())
	assert.True(t, av.IsNull())
}

func TestUnmarshalString(t *testing.T) {
	input := []byte(`{ "S": "Hello"}`)

	var av DynamoDBAttributeValue
	err := json.Unmarshal(input, &av)

	assert.Nil(t, err)
	assert.Equal(t, DataTypeString, av.DataType())
	assert.Equal(t, "Hello", av.String())
}

func TestUnmarshalStringSet(t *testing.T) {
	input := []byte(`{ "SS": [ "Giraffe", "Zebra" ] }`)

	var av DynamoDBAttributeValue
	err := json.Unmarshal(input, &av)

	assert.Nil(t, err)
	assert.Equal(t, 2, len(av.StringSet()))
	assert.Equal(t, DataTypeStringSet, av.DataType())
	assert.Equal(t, "Giraffe", av.StringSet()[0])
	assert.Equal(t, "Zebra", av.StringSet()[1])
}

func TestUnmarshalEmptyStringSet(t *testing.T) {
	input := []byte(`{ "SS": [ ] }`)

	var av DynamoDBAttributeValue
	err := json.Unmarshal(input, &av)

	assert.Nil(t, err)
	assert.Equal(t, DataTypeStringSet, av.DataType())
	assert.Equal(t, 0, len(av.StringSet()))
}

func TestAccessWithWrongTypePanics(t *testing.T) {
	testCases := []struct {
		input         string
		accessor      func(av DynamoDBAttributeValue)
		expectedError IncompatibleDynamoDBTypeError
	}{
		{`{ "B": "AAEqQQ=="}`, func(av DynamoDBAttributeValue) { av.Number() }, IncompatibleDynamoDBTypeError{Requested: DataTypeNumber, Actual: DataTypeBinary}},
		{`{ "BOOL": true}`, func(av DynamoDBAttributeValue) { av.Number() }, IncompatibleDynamoDBTypeError{Requested: DataTypeNumber, Actual: DataTypeBoolean}},
		{`{ "BS": ["AAEqQQ==", "AAEqQQ=="] }`, func(av DynamoDBAttributeValue) { av.Number() }, IncompatibleDynamoDBTypeError{Requested: DataTypeNumber, Actual: DataTypeBinarySet}},
		{`{ "L": [ { "S": "Cookies"} ] }`, func(av DynamoDBAttributeValue) { av.Number() }, IncompatibleDynamoDBTypeError{Requested: DataTypeNumber, Actual: DataTypeList}},
		{`{ "M": { "Name": { "S": "Joe" } } }`, func(av DynamoDBAttributeValue) { av.Number() }, IncompatibleDynamoDBTypeError{Requested: DataTypeNumber, Actual: DataTypeMap}},
		{`{ "N": "123.45"}`, func(av DynamoDBAttributeValue) { av.Boolean() }, IncompatibleDynamoDBTypeError{Requested: DataTypeBoolean, Actual: DataTypeNumber}},
		{`{ "NS": ["1234", "567.8"] }`, func(av DynamoDBAttributeValue) { av.Boolean() }, IncompatibleDynamoDBTypeError{Requested: DataTypeBoolean, Actual: DataTypeNumberSet}},
		{`{ "NULL": true}`, func(av DynamoDBAttributeValue) { av.Number() }, IncompatibleDynamoDBTypeError{Requested: DataTypeNumber, Actual: DataTypeNull}},
		{`{ "S": "Hello"}`, func(av DynamoDBAttributeValue) { av.Number() }, IncompatibleDynamoDBTypeError{Requested: DataTypeNumber, Actual: DataTypeString}},
		{`{ "SS": [ "Giraffe", "Zebra" ] }`, func(av DynamoDBAttributeValue) { av.Number() }, IncompatibleDynamoDBTypeError{Requested: DataTypeNumber, Actual: DataTypeStringSet}},
	}

	for _, testCase := range testCases {
		var av DynamoDBAttributeValue
		err := json.Unmarshal([]byte(testCase.input), &av)
		assert.Nil(t, err)
		// may use PanicsWithValue(expectedError) when it is available
		assertPanicsWithValue(t, testCase.expectedError, func() { testCase.accessor(av) })
	}
}

func assertPanicsWithValue(t *testing.T, expected error, action func()) {
	defer func() {
		r := recover()
		if r == nil {
			t.Errorf("Should have panicked")
		}
		if r != expected {
			t.Errorf("should have panicked with value %v but panicked with value %v", expected, r)
		}
	}()

	action()
}

func TestMarshalAndUnmarshalString(t *testing.T) {
	const inputString = "INPUT STRING"
	inputValue := NewStringAttribute(inputString)
	marshaled, err := json.Marshal(inputValue)
	assert.Nil(t, err)

	var av DynamoDBAttributeValue
	err = json.Unmarshal(marshaled, &av)

	assert.Nil(t, err)
	assert.Equal(t, DataTypeString, av.DataType())
	assert.Equal(t, inputString, av.String())
}

func Test_DynamoDBAttributeValue_NewAttribute(t *testing.T) {
	{
		av := NewBinaryAttribute([]byte{1, 2, 3})
		assert.Equal(t, DataTypeBinary, av.DataType())
		assert.Equal(t, []byte{1, 2, 3}, av.Binary())
	}
	{
		av := NewBooleanAttribute(true)
		assert.Equal(t, DataTypeBoolean, av.DataType())
		assert.Equal(t, true, av.Boolean())
	}
	{
		av := NewBinarySetAttribute([][]byte{[]byte{1, 2, 3}})
		assert.Equal(t, DataTypeBinarySet, av.DataType())
		assert.Equal(t, [][]byte{[]byte{1, 2, 3}}, av.BinarySet())
	}
	{
		av := NewListAttribute([]DynamoDBAttributeValue{
			NewNumberAttribute("1"),
			NewStringAttribute("test"),
		})
		assert.Equal(t, DataTypeList, av.DataType())
		assert.Equal(t, 2, len(av.List()))
	}
	{
		value := map[string]DynamoDBAttributeValue{
			"n": NewNumberAttribute("1"),
			"s": NewStringAttribute("test"),
		}
		av := NewMapAttribute(value)
		assert.Equal(t, DataTypeMap, av.DataType())
		assert.Equal(t, 2, len(av.Map()))
	}
	{
		av := NewNumberAttribute("1")
		assert.Equal(t, DataTypeNumber, av.DataType())
		assert.Equal(t, "1", av.Number())
		v, err := av.Integer()
		assert.Nil(t, err)
		assert.Equal(t, int64(1), v)
	}
	{
		av := NewNumberAttribute("1.1")
		assert.Equal(t, DataTypeNumber, av.DataType())
		assert.Equal(t, "1.1", av.Number())
		v, err := av.Float()
		assert.Nil(t, err)
		assert.Equal(t, float64(1.1), v)
	}
	{
		av := NewNullAttribute()
		assert.Equal(t, DataTypeNull, av.DataType())
		assert.Equal(t, true, av.IsNull())
	}
	{
		av := NewStringAttribute("test")
		assert.Equal(t, DataTypeString, av.DataType())
		assert.Equal(t, "test", av.String())
	}
	{
		av := NewStringSetAttribute([]string{"test", "test"})
		assert.Equal(t, DataTypeStringSet, av.DataType())
		assert.Equal(t, []string{"test", "test"}, av.StringSet())
	}
}
