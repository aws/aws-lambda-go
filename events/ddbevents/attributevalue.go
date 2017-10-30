// Copyright 2017 Amazon.com, Inc. or its affiliates. All Rights Reserved.

package ddbevents

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"errors"
	"strconv"
)

// AttributeValue provides convenient access for a value stored in DynamoDB.
// For more information,  please see http://docs.aws.amazon.com/amazondynamodb/latest/APIReference/API_AttributeValue.html
type AttributeValue struct {
	value    anyValue
	dataType int
}

// Binary() provides access to an attribute of type Binary.
// Method panics if the attribute is not of type Binary.
func (av AttributeValue) Binary() []byte {
	av.ensureType(DataTypeBinary)
	return av.value.([]byte)
}

// Boolean() provides access to an attribute of type Boolean.
// Method panics if the attribute is not of type Boolean.
func (av AttributeValue) Boolean() bool {
	av.ensureType(DataTypeBoolean)
	return av.value.(bool)
}

// BinarySet() provides access to an attribute of type Binary Set.
// Method panics if the attribute is not of type BinarySet.
func (av AttributeValue) BinarySet() [][]byte {
	av.ensureType(DataTypeBinarySet)
	return av.value.([][]byte)
}

// List() provides access to an attribute of type List. Each element
// of the list is an AttributeValue itself.
// Method panics if the attribute is not of type List.
func (av AttributeValue) List() []AttributeValue {
	av.ensureType(DataTypeList)
	return av.value.([]AttributeValue)
}

// Map() provides access to an attribute of type Map. They Keys are strings
// and the values are AttributeValue instances.
// Method panics if the attribute is not of type Map.
func (av AttributeValue) Map() map[string]AttributeValue {
	av.ensureType(DataTypeMap)
	return av.value.(map[string]AttributeValue)
}

// Number() provides access to an attribute of type Number.
// DynamoDB sends the values as strings. For convenience please see also
// the methods Integer() and Float().
// Method panics if the attribute is not of type Number.
func (av AttributeValue) Number() string {
	av.ensureType(DataTypeNumber)
	return av.value.(string)
}

// Integer() provides access to an attribute of type Number.
// DynamoDB sends the values as strings. For convenience this method
// provides conversion to int. If the value cannot be represented by
// a signed integer, err.Err = ErrRange and the returned value is the maximum magnitude integer
// of an int64 of the appropriate sign.
// Method panics if the attribute is not of type Number.
func (av AttributeValue) Integer() (int64, error) {
	s, err := strconv.ParseFloat(av.Number(), 64)
	return int64(s), err
}

// Float() provides access to an attribute of type Number.
// DynamoDB sends the values as strings. For convenience this method
// provides conversion to float64.
// The returned value is the nearest floating point number rounded using IEEE754 unbiased rounding.
// If the number is more than 1/2 ULP away from the largest floating point number of the given size,
// the value returned is ±Inf, err.Err = ErrRange.
// Method panics if the attribute is not of type Number.
func (av AttributeValue) Float() (float64, error) {
	s, err := strconv.ParseFloat(av.Number(), 64)
	return s, err
}

// NumberSet() provides access to an attribute of type Number Set.
// DynamoDB sends the numbers as strings.
// Method panics if the attribute is not of type Number.
func (av AttributeValue) NumberSet() []string {
	av.ensureType(DataTypeNumberSet)
	return av.value.([]string)
}

// String() provides access to an attribute of type String.
// Method panics if the attribute is not of type String.
func (av AttributeValue) String() string {
	av.ensureType(DataTypeString)
	return av.value.(string)
}

// StringSet() provides access to an attribute of type String Set.
// Method panics if the attribute is not of type String Set.
func (av AttributeValue) StringSet() []string {
	av.ensureType(DataTypeStringSet)
	return av.value.([]string)
}

// IsNull() returns true if the attribute is of type Null.
func (av AttributeValue) IsNull() bool {
	return av.value == nil
}

// Provides access to the DynamoDB type of the attribute
func (av AttributeValue) DataType() int {
	return av.dataType
}

// Creates an AttributeValue containing a String
func NewStringAttribute(value string) AttributeValue {
	var av AttributeValue
	av.value = value
	av.dataType = DataTypeString
	return av
}

// Data types for attributes supported natively by DynamoDB
const (
	DataTypeBinary = iota
	DataTypeBoolean
	DataTypeBinarySet
	DataTypeList
	DataTypeMap
	DataTypeNumber
	DataTypeNumberSet
	DataTypeNull
	DataTypeString
	DataTypeStringSet
)

type anyValue interface{}

var ErrUnsupportedType = errors.New("attributevalue: unsupported type")
var ErrInvalidType = errors.New("attributevalue: accessor called for incompatible type")

func (av *AttributeValue) ensureType(expectedType int) {
	if av.dataType != expectedType {
		panic(ErrInvalidType)
	}
}

// MarshalJSON implements custom marshaling to be used by the standard json/encoding package
func (av AttributeValue) MarshalJSON() ([]byte, error) {

	var buff bytes.Buffer
	var err error
	var b []byte

	switch av.dataType {
	case DataTypeBinary:
		buff.WriteString(`{ "B":`)
		b, err = json.Marshal(av.value.([]byte))
		buff.Write(b)

	case DataTypeBoolean:
		buff.WriteString(`{ "BOOL":`)
		b, err = json.Marshal(av.value.(bool))
		buff.Write(b)

	case DataTypeBinarySet:
		buff.WriteString(`{ "BS":`)
		b, err = json.Marshal(av.value.([][]byte))
		buff.Write(b)

	case DataTypeList:
		buff.WriteString(`{ "L":`)
		b, err = json.Marshal(av.value.([]AttributeValue))
		buff.Write(b)

	case DataTypeMap:
		buff.WriteString(`{ "M":`)
		b, err = json.Marshal(av.value.(map[string]AttributeValue))
		buff.Write(b)

	case DataTypeNumber:
		buff.WriteString(`{ "N":`)
		b, err = json.Marshal(av.value.(string))
		buff.Write(b)

	case DataTypeNumberSet:
		buff.WriteString(`{ "NS":`)
		b, err = json.Marshal(av.value.([]string))
		buff.Write(b)

	case DataTypeNull:
		buff.WriteString(`{ "NULL": true `)

	case DataTypeString:
		buff.WriteString(`{ "S":`)
		b, err = json.Marshal(av.value.(string))
		buff.Write(b)

	case DataTypeStringSet:
		buff.WriteString(`{ "SS":`)
		b, err = json.Marshal(av.value.([]string))
		buff.Write(b)

	default:
		err = ErrUnsupportedType
	}

	buff.WriteString(`}`)
	return buff.Bytes(), err
}

func unmarshalNull(target *AttributeValue) error {
	target.value = nil
	target.dataType = DataTypeNull
	return nil
}

func unmarshalString(target *AttributeValue, value interface{}) error {
	var ok bool
	target.value, ok = value.(string)
	target.dataType = DataTypeString
	if !ok {
		return errors.New("attributevalue: S type should contain a string")
	}
	return nil
}

func unmarshalBinary(target *AttributeValue, value interface{}) error {
	stringValue, ok := value.(string)
	if !ok {
		return errors.New("attributevalue: B type should contain a base64 string")
	}

	binaryValue, err := base64.StdEncoding.DecodeString(stringValue)
	if err != nil {
		return err
	}

	target.value = binaryValue
	target.dataType = DataTypeBinary
	return nil
}

func unmarshalBoolean(target *AttributeValue, value interface{}) error {
	booleanValue, ok := value.(bool)
	if !ok {
		return errors.New("attributevalue: BOOL type should contain a boolean")
	}

	target.value = booleanValue
	target.dataType = DataTypeBoolean
	return nil
}

func unmarshalBinarySet(target *AttributeValue, value interface{}) error {
	list, ok := value.([]interface{})
	if !ok {
		return errors.New("attributevalue: BS type should contain a list of base64 strings")
	}

	binarySet := make([][]byte, len(list), len(list))

	for index, element := range list {
		var err error
		elementString := element.(string)
		binarySet[index], err = base64.StdEncoding.DecodeString(elementString)
		if err != nil {
			return err
		}
	}

	target.value = binarySet
	target.dataType = DataTypeBinarySet
	return nil
}

func unmarshalList(target *AttributeValue, value interface{}) error {
	list, ok := value.([]interface{})
	if !ok {
		return errors.New("attributevalue: L type should contain a list")
	}

	attributeValues := make([]AttributeValue, len(list), len(list))
	for index, element := range list {

		elementMap, ok := element.(map[string]interface{})
		if !ok {
			return errors.New("attributeValue: element of a list is not an AttributeValue")
		}

		var elementAttributeValue AttributeValue
		err := unmarshalAttributeValueMap(&elementAttributeValue, elementMap)
		if err != nil {
			return errors.New("attributeValue: unmarshal of child AttributeValue failed")
		}
		attributeValues[index] = elementAttributeValue
	}
	target.value = attributeValues
	target.dataType = DataTypeList
	return nil
}

func unmarshalMap(target *AttributeValue, value interface{}) error {
	m, ok := value.(map[string]interface{})
	if !ok {
		return errors.New("attributevalue: M type should contain a map")
	}

	attributeValues := make(map[string]AttributeValue)
	for k, v := range m {

		elementMap, ok := v.(map[string]interface{})
		if !ok {
			return errors.New("attributeValue: element of a map is not an AttributeValue")
		}

		var elementAttributeValue AttributeValue
		err := unmarshalAttributeValueMap(&elementAttributeValue, elementMap)
		if err != nil {
			return errors.New("attributeValue: unmarshal of child AttributeValue failed")
		}
		attributeValues[k] = elementAttributeValue
	}
	target.value = attributeValues
	target.dataType = DataTypeMap
	return nil
}

func unmarshalNumber(target *AttributeValue, value interface{}) error {
	var ok bool
	target.value, ok = value.(string)
	target.dataType = DataTypeNumber
	if !ok {
		return errors.New("attributevalue: N type should contain a string")
	}
	return nil
}

func unmarshalNumberSet(target *AttributeValue, value interface{}) error {
	list, ok := value.([]interface{})
	if !ok {
		return errors.New("attributevalue: NS type should contain a list of strings")
	}

	numberSet := make([]string, len(list), len(list))

	for index, element := range list {
		numberSet[index], ok = element.(string)
		if !ok {
			return errors.New("attributevalue: NS type should contain a list of strings")
		}
	}

	target.value = numberSet
	target.dataType = DataTypeNumberSet
	return nil
}

func unmarshalStringSet(target *AttributeValue, value interface{}) error {
	list, ok := value.([]interface{})
	if !ok {
		return errors.New("attributevalue: SS type should contain a list of strings")
	}

	stringSet := make([]string, len(list), len(list))

	for index, element := range list {
		stringSet[index], ok = element.(string)
		if !ok {
			return errors.New("attributevalue: SS type should contain a list of strings")
		}
	}

	target.value = stringSet
	target.dataType = DataTypeStringSet
	return nil
}

func unmarshalAttributeValue(target *AttributeValue, typeLabel string, jsonValue interface{}) error {

	switch typeLabel {
	case "NULL":
		return unmarshalNull(target)
	case "B":
		return unmarshalBinary(target, jsonValue)
	case "BOOL":
		return unmarshalBoolean(target, jsonValue)
	case "BS":
		return unmarshalBinarySet(target, jsonValue)
	case "L":
		return unmarshalList(target, jsonValue)
	case "M":
		return unmarshalMap(target, jsonValue)
	case "N":
		return unmarshalNumber(target, jsonValue)
	case "NS":
		return unmarshalNumberSet(target, jsonValue)
	case "S":
		return unmarshalString(target, jsonValue)
	case "SS":
		return unmarshalStringSet(target, jsonValue)
	default:
		target.value = nil
		target.dataType = DataTypeNull
		return ErrUnsupportedType
	}
}

func (target *AttributeValue) UnmarshalJSON(b []byte) error {
	var m map[string]interface{}

	err := json.Unmarshal(b, &m)
	if err != nil {
		return err
	}

	return unmarshalAttributeValueMap(target, m)
}

func unmarshalAttributeValueMap(target *AttributeValue, m map[string]interface{}) error {
	if m == nil {
		return errors.New("attributevalue: does not contain a map")
	}

	if len(m) != 1 {
		return errors.New("attributevalue: map must contain a single type")
	}

	for k, v := range m {
		return unmarshalAttributeValue(target, k, v)
	}

	return nil
}
