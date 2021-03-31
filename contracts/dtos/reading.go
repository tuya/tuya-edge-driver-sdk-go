//
// Copyright (C) 2020 IOTech Ltd
//
// SPDX-License-Identifier: Apache-2.0

package dtos

import (
	"bytes"
	"encoding/base64"
	"encoding/binary"
	"encoding/json"
	"errors"
	"fmt"
	"reflect"
	"strconv"
	"strings"
	"time"

	"github.com/google/uuid"

	"github.com/tuya/tuya-edge-driver-sdk-go/contracts"
	"github.com/tuya/tuya-edge-driver-sdk-go/contracts/dtos/common"
	"github.com/tuya/tuya-edge-driver-sdk-go/contracts/models"
)

// BaseReading and its properties are defined in the APIv2 specification:
// https://app.swaggerhub.com/apis-docs/EdgeXFoundry1/core-data/2.x#/BaseReading
type BaseReading struct {
	common.Versionable `json:",inline"`
	Id                 string `json:"id"`
	Created            int64  `json:"created"`
	Origin             int64  `json:"origin" validate:"required"`
	DeviceName         string `json:"deviceName" validate:"required,edgex-dto-rfc3986-unreserved-chars"`
	ResourceName       string `json:"resourceName" validate:"required,edgex-dto-rfc3986-unreserved-chars"`
	ProfileName        string `json:"profileName" validate:"required,edgex-dto-rfc3986-unreserved-chars"`
	ValueType          string `json:"valueType" validate:"required,edgex-dto-value-type"`
	BinaryReading      `json:",inline" validate:"-"`
	SimpleReading      `json:",inline" validate:"-"`
}

// SimpleReading and its properties are defined in the APIv2 specification:
// https://app.swaggerhub.com/apis-docs/EdgeXFoundry1/core-data/2.x#/SimpleReading
type SimpleReading struct {
	Value string `json:"value" validate:"required"`
}

// BinaryReading and its properties are defined in the APIv2 specification:
// https://app.swaggerhub.com/apis-docs/EdgeXFoundry1/core-data/2.x#/BinaryReading
type BinaryReading struct {
	BinaryValue []byte `json:"binaryValue" validate:"gt=0,dive,required"`
	MediaType   string `json:"mediaType" validate:"required"`
}

func newBaseReading(profileName string, deviceName string, resourceName string, valueType string) BaseReading {
	return BaseReading{
		Versionable:  common.NewVersionable(),
		Id:           uuid.NewString(),
		Origin:       time.Now().UnixNano(),
		DeviceName:   deviceName,
		ResourceName: resourceName,
		ProfileName:  profileName,
		ValueType:    valueType,
	}
}

// NewSimpleReading creates and returns a new initialized BaseReading with its SimpleReading initialized
func NewSimpleReading(profileName string, deviceName string, resourceName string, valueType string, value interface{}) (BaseReading, error) {
	stringValue, err := convertInterfaceValue(valueType, value)
	if err != nil {
		return BaseReading{}, err
	}

	reading := newBaseReading(profileName, deviceName, resourceName, valueType)
	reading.SimpleReading = SimpleReading{
		Value: stringValue,
	}
	return reading, nil
}

// NewBinaryReading creates and returns a new initialized BaseReading with its BinaryReading initialized
func NewBinaryReading(profileName string, deviceName string, resourceName string, binaryValue []byte, mediaType string) BaseReading {
	reading := newBaseReading(profileName, deviceName, resourceName, contracts.ValueTypeBinary)
	reading.BinaryReading = BinaryReading{
		BinaryValue: binaryValue,
		MediaType:   mediaType,
	}
	return reading
}

// 将reading中的value string转换成value type类型对应的值
func (br *BaseReading) ConvertValue() (interface{}, error) {
	if br.ValueType == contracts.ValueTypeBinary {
		return br.BinaryValue, nil
	}
	return convertStringValue(br.ValueType, br.Value)
}

func convertStringValue(valueType, value string) (interface{}, error) {
	switch valueType {
	case contracts.ValueTypeBool:
		return strconv.ParseBool(value)

	case contracts.ValueTypeString:
		// string直接返回value
		return value, nil

	case contracts.ValueTypeUint8:
		return parseUint8(value)
	case contracts.ValueTypeUint16:
		return parseUint16(value)
	case contracts.ValueTypeUint32:
		return parseUint32(value)
	case contracts.ValueTypeUint64:
		return strconv.ParseUint(value, 10, 64)

	case contracts.ValueTypeInt8:
		return parseInt8(value)
	case contracts.ValueTypeInt16:
		return parseInt16(value)
	case contracts.ValueTypeInt32:
		return parseInt32(value)
	case contracts.ValueTypeInt64:
		return strconv.ParseInt(value, 10, 64)

	case contracts.ValueTypeFloat32:
		return parseFloat32(value)
	case contracts.ValueTypeFloat64:
		return strconv.ParseFloat(value, 64)

	case contracts.ValueTypeBoolArray:
		var bSli []bool
		err := json.Unmarshal([]byte(value), &bSli)
		return bSli, err
	case contracts.ValueTypeStringArray:
		var sSli []string
		err := json.Unmarshal([]byte(value), &sSli)
		return sSli, err
	case contracts.ValueTypeUint8Array:
		var u8Sli []uint8
		err := json.Unmarshal([]byte(value), &u8Sli)
		return u8Sli, err
	case contracts.ValueTypeUint16Array:
		var u16Sli []uint16
		err := json.Unmarshal([]byte(value), &u16Sli)
		return u16Sli, err
	case contracts.ValueTypeUint32Array:
		var u32Sli []uint32
		err := json.Unmarshal([]byte(value), &u32Sli)
		return u32Sli, err
	case contracts.ValueTypeUint64Array:
		var u64Sli []uint64
		err := json.Unmarshal([]byte(value), &u64Sli)
		return u64Sli, err
	case contracts.ValueTypeInt8Array:
		var i8Sli []int8
		err := json.Unmarshal([]byte(value), &i8Sli)
		return i8Sli, err
	case contracts.ValueTypeInt16Array:
		var i16Sli []int16
		err := json.Unmarshal([]byte(value), &i16Sli)
		return i16Sli, err
	case contracts.ValueTypeInt32Array:
		var i32Sli []int32
		err := json.Unmarshal([]byte(value), &i32Sli)
		return i32Sli, err
	case contracts.ValueTypeInt64Array:
		var i64Sli []int64
		err := json.Unmarshal([]byte(value), &i64Sli)
		return i64Sli, err
		//return parseSimpleArray(valueType, valueType)

	case contracts.ValueTypeFloat32Array:
		var f32Sli []float32
		err := json.Unmarshal([]byte(value), &f32Sli)
		return f32Sli, err
	case contracts.ValueTypeFloat64Array:
		var f64Sli []float64
		err := json.Unmarshal([]byte(value), &f64Sli)
		return f64Sli, err
	default:
		return nil, fmt.Errorf("invalid simple reading type of %s", valueType)
	}
}

func convertInterfaceValue(valueType string, value interface{}) (string, error) {
	switch valueType {
	case contracts.ValueTypeBool:
		return convertSimpleValue(valueType, reflect.Bool, value)
	case contracts.ValueTypeString:
		return convertSimpleValue(valueType, reflect.String, value)

	case contracts.ValueTypeUint8:
		return convertSimpleValue(valueType, reflect.Uint8, value)
	case contracts.ValueTypeUint16:
		return convertSimpleValue(valueType, reflect.Uint16, value)
	case contracts.ValueTypeUint32:
		return convertSimpleValue(valueType, reflect.Uint32, value)
	case contracts.ValueTypeUint64:
		return convertSimpleValue(valueType, reflect.Uint64, value)

	case contracts.ValueTypeInt8:
		return convertSimpleValue(valueType, reflect.Int8, value)
	case contracts.ValueTypeInt16:
		return convertSimpleValue(valueType, reflect.Int16, value)
	case contracts.ValueTypeInt32:
		return convertSimpleValue(valueType, reflect.Int32, value)
	case contracts.ValueTypeInt64:
		return convertSimpleValue(valueType, reflect.Int64, value)

	case contracts.ValueTypeFloat32:
		return convertFloatValue(valueType, reflect.Float32, value)
	case contracts.ValueTypeFloat64:
		return convertFloatValue(valueType, reflect.Float64, value)

	case contracts.ValueTypeBoolArray:
		return convertSimpleArrayValue(valueType, reflect.Bool, value)
	case contracts.ValueTypeStringArray:
		return convertSimpleArrayValue(valueType, reflect.String, value)

	case contracts.ValueTypeUint8Array:
		return convertSimpleArrayValue(valueType, reflect.Uint8, value)
	case contracts.ValueTypeUint16Array:
		return convertSimpleArrayValue(valueType, reflect.Uint16, value)
	case contracts.ValueTypeUint32Array:
		return convertSimpleArrayValue(valueType, reflect.Uint32, value)
	case contracts.ValueTypeUint64Array:
		return convertSimpleArrayValue(valueType, reflect.Uint64, value)

	case contracts.ValueTypeInt8Array:
		return convertSimpleArrayValue(valueType, reflect.Int8, value)
	case contracts.ValueTypeInt16Array:
		return convertSimpleArrayValue(valueType, reflect.Int16, value)
	case contracts.ValueTypeInt32Array:
		return convertSimpleArrayValue(valueType, reflect.Int32, value)
	case contracts.ValueTypeInt64Array:
		return convertSimpleArrayValue(valueType, reflect.Int64, value)

	case contracts.ValueTypeFloat32Array:
		arrayValue, ok := value.([]float32)
		if !ok {
			return "", fmt.Errorf("unable to cast value to []float32 for %s", valueType)
		}

		return convertFloat32ArrayValue(arrayValue)
	case contracts.ValueTypeFloat64Array:
		arrayValue, ok := value.([]float64)
		if !ok {
			return "", fmt.Errorf("unable to cast value to []float64 for %s", valueType)
		}

		return convertFloat64ArrayValue(arrayValue)

	default:
		return "", fmt.Errorf("invalid simple reading type of %s", valueType)
	}
}

func convertSimpleValue(valueType string, kind reflect.Kind, value interface{}) (string, error) {
	if err := validateType(valueType, kind, value); err != nil {
		return "", err
	}

	return fmt.Sprintf("%v", value), nil
}

func convertFloatValue(valueType string, kind reflect.Kind, value interface{}) (string, error) {
	if err := validateType(valueType, kind, value); err != nil {
		return "", err
	}

	buf := new(bytes.Buffer)
	if err := binary.Write(buf, binary.BigEndian, value); err != nil {
		return "", err
	}
	return base64.StdEncoding.EncodeToString(buf.Bytes()), nil
}

func convertSimpleArrayValue(valueType string, kind reflect.Kind, value interface{}) (string, error) {
	if err := validateType(valueType, kind, value); err != nil {
		return "", err
	}

	result := fmt.Sprintf("%v", value)
	result = strings.ReplaceAll(result, " ", ", ")
	return result, nil
}

func parseSimpleArray(valueType, value string) (interface{}, error) {
	var (
		err      error
		newValue []string
	)
	if newValue, err = convertSimpleStringArray(value); err != nil {
		return nil, err
	}
	switch valueType {
	case contracts.ValueTypeBoolArray:
		var bSli = make([]bool, 0, len(newValue))
		for i := range newValue {
			b, err := strconv.ParseBool(newValue[i])
			if err != nil {
				return nil, err
			}
			bSli = append(bSli, b)
		}
		return bSli, nil
	case contracts.ValueTypeStringArray:
		var sSli = make([]string, 0, len(newValue))
		for i := range newValue {
			sSli = append(sSli, newValue[i])
		}
		return sSli, nil

	case contracts.ValueTypeUint8Array:
		var u8Sli = make([]uint8, 0, len(newValue))
		for i := range newValue {
			u8, err := parseUint8(newValue[i])
			if err != nil {
				return nil, err
			}
			u8Sli = append(u8Sli, u8)
		}
		return u8Sli, nil
	case contracts.ValueTypeUint16Array:
		var u16Sli = make([]uint16, 0, len(newValue))
		for i := range newValue {
			u16, err := parseUint16(newValue[i])
			if err != nil {
				return nil, err
			}
			u16Sli = append(u16Sli, u16)
		}
		return u16Sli, nil
	case contracts.ValueTypeUint32Array:
		var u32Sli = make([]uint32, 0, len(newValue))
		for i := range newValue {
			u32, err := parseUint32(newValue[i])
			if err != nil {
				return nil, err
			}
			u32Sli = append(u32Sli, u32)
		}
		return u32Sli, nil
	case contracts.ValueTypeUint64Array:
		var u64Sli = make([]uint64, 0, len(newValue))
		for i := range newValue {
			u64, err := strconv.ParseUint(newValue[i], 10, 64)
			if err != nil {
				return nil, err
			}
			u64Sli = append(u64Sli, u64)
		}
		return u64Sli, nil

	case contracts.ValueTypeInt8Array:
		var i8Sli = make([]int8, 0, len(newValue))
		for i := range newValue {
			i8, err := parseInt8(newValue[i])
			if err != nil {
				return nil, err
			}
			i8Sli = append(i8Sli, i8)
		}
		return i8Sli, nil
	case contracts.ValueTypeInt16Array:
		var i16Sli = make([]int16, 0, len(newValue))
		for i := range newValue {
			i16, err := parseInt16(newValue[i])
			if err != nil {
				return nil, err
			}
			i16Sli = append(i16Sli, i16)
		}
		return i16Sli, nil
	case contracts.ValueTypeInt32Array:
		var i32Sli = make([]int32, 0, len(newValue))
		for i := range newValue {
			i32, err := parseInt32(newValue[i])
			if err != nil {
				return nil, err
			}
			i32Sli = append(i32Sli, i32)
		}
		return i32Sli, nil
	case contracts.ValueTypeInt64Array:
		var i64Sli = make([]int64, 0, len(newValue))
		for i := range newValue {
			i64, err := strconv.ParseInt(newValue[i], 10, 64)
			if err != nil {
				return nil, err
			}
			i64Sli = append(i64Sli, i64)
		}
		return i64Sli, nil
	default:
		return nil, errors.New("wrong simple slice type")
	}
}

func parseUint8(value string) (uint8, error) {
	v, err := strconv.ParseUint(value, 10, 64)
	if err != nil {
		return 0, err
	}
	if v > 0xff {
		return 0, errors.New("value out of range")
	}
	return uint8(v), nil
}

func parseUint16(value string) (uint16, error) {
	v, err := strconv.ParseUint(value, 10, 64)
	if err != nil {
		return 0, err
	}
	if v > 0xffff {
		return 0, errors.New("value out of range")
	}
	return uint16(v), nil
}

func parseUint32(value string) (uint32, error) {
	v, err := strconv.ParseUint(value, 10, 32)
	if err != nil {
		return 0, err
	}
	if v > 0xffffff {
		return 0, errors.New("value out of range")
	}
	return uint32(v), nil
}

func parseFloat32(value string) (float32, error) {
	v, err := strconv.ParseFloat(value, 32)
	if err != nil {
		return 0, err
	}
	return float32(v), nil
}

func parseInt8(value string) (int8, error) {
	v, err := strconv.ParseInt(value, 10, 8)
	return int8(v), err
}

func parseInt16(value string) (int16, error) {
	v, err := strconv.ParseInt(value, 10, 16)
	return int16(v), err
}

func parseInt32(value string) (int32, error) {
	v, err := strconv.ParseInt(value, 10, 32)
	return int32(v), err
}

func convertSimpleStringArray(value string) ([]string, error) {
	var err = errors.New("wrong simple string array")
	if len(value) < 2 {
		return []string{}, err
	}
	if value[0] == 0x5b && value[len(value)-1] == 0x5d { // []
		return strings.Split(value[1:len(value)-1], ", "), nil
	}
	return []string{}, err
}

func convertFloat32ArrayValue(values []float32) (string, error) {
	result := "["
	first := true
	for _, value := range values {
		if first {
			floatValue, err := convertFloatValue(contracts.ValueTypeFloat32, reflect.Float32, value)
			if err != nil {
				return "", err
			}
			result += floatValue
			first = false
			continue
		}

		floatValue, err := convertFloatValue(contracts.ValueTypeFloat32, reflect.Float32, value)
		if err != nil {
			return "", err
		}
		result += ", " + floatValue
	}

	return result, nil
}

func convertFloat64ArrayValue(values []float64) (string, error) {
	result := "["
	first := true
	for _, value := range values {
		if first {
			floatValue, err := convertFloatValue(contracts.ValueTypeFloat64, reflect.Float64, value)
			if err != nil {
				return "", err
			}
			result += floatValue
			first = false
			continue
		}

		floatValue, err := convertFloatValue(contracts.ValueTypeFloat64, reflect.Float64, value)
		if err != nil {
			return "", err
		}
		result += ", " + floatValue
	}

	return result, nil
}

func validateType(valueType string, kind reflect.Kind, value interface{}) error {
	if reflect.TypeOf(value).Kind() == reflect.Slice {
		if kind != reflect.TypeOf(value).Elem().Kind() {
			return fmt.Errorf("slice of type of value `%s` not a match for specified ValueType '%s", kind.String(), valueType)
		}
		return nil
	}

	if kind != reflect.TypeOf(value).Kind() {
		return fmt.Errorf("type of value `%s` not a match for specified ValueType '%s", kind.String(), valueType)
	}

	return nil
}

// Validate satisfies the Validator interface
func (b BaseReading) Validate() error {
	if b.ValueType == contracts.ValueTypeBinary {
		// validate the inner BinaryReading struct
		binaryReading := b.BinaryReading
		if err := contracts.Validate(binaryReading); err != nil {
			return err
		}
	} else {
		// validate the inner SimpleReading struct
		simpleReading := b.SimpleReading
		if err := contracts.Validate(simpleReading); err != nil {
			return err
		}
	}

	return nil
}

// Convert Reading DTO to Reading model
func ToReadingModel(r BaseReading) models.Reading {
	var readingModel models.Reading
	br := models.BaseReading{
		Origin:       r.Origin,
		DeviceName:   r.DeviceName,
		ResourceName: r.ResourceName,
		ProfileName:  r.ProfileName,
		ValueType:    r.ValueType,
	}
	if r.ValueType == contracts.ValueTypeBinary {
		readingModel = models.BinaryReading{
			BaseReading: br,
			BinaryValue: r.BinaryValue,
			MediaType:   r.MediaType,
		}
	} else {
		readingModel = models.SimpleReading{
			BaseReading: br,
			Value:       r.Value,
		}
	}
	return readingModel
}

func FromReadingModelToDTO(reading models.Reading) BaseReading {
	var baseReading BaseReading
	switch r := reading.(type) {
	case models.BinaryReading:
		baseReading = BaseReading{
			Versionable:   common.NewVersionable(),
			Id:            r.Id,
			Created:       r.Created,
			Origin:        r.Origin,
			DeviceName:    r.DeviceName,
			ResourceName:  r.ResourceName,
			ProfileName:   r.ProfileName,
			ValueType:     r.ValueType,
			BinaryReading: BinaryReading{BinaryValue: r.BinaryValue, MediaType: r.MediaType},
		}
	case models.SimpleReading:
		baseReading = BaseReading{
			Versionable:   common.Versionable{ApiVersion: contracts.ApiVersion},
			Id:            r.Id,
			Created:       r.Created,
			Origin:        r.Origin,
			DeviceName:    r.DeviceName,
			ResourceName:  r.ResourceName,
			ProfileName:   r.ProfileName,
			ValueType:     r.ValueType,
			SimpleReading: SimpleReading{Value: r.Value},
		}
	}

	return baseReading
}
