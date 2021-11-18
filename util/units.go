package util

import (
	"errors"
	"fmt"
	"math"
)

// Multiplier value
const (
	BaseDecimal = 1000.0
	BaseBinary  = 1024.0
)

// Units value
const (
	BYTE = iota
	KILO
	MEGA
	GIGA
	TERA
	PETA
	EXA
	ZETTA
	YOTTA
)

var units = map[float64]string{
	BYTE: "B",

	BaseDecimal * KILO:  "KB",
	BaseDecimal * MEGA:  "MB",
	BaseDecimal * GIGA:  "GB",
	BaseDecimal * TERA:  "TB",
	BaseDecimal * PETA:  "PB",
	BaseDecimal * EXA:   "EB",
	BaseDecimal * ZETTA: "ZB",
	BaseDecimal * YOTTA: "YB",

	BaseBinary * KILO:  "KiB",
	BaseBinary * MEGA:  "MiB",
	BaseBinary * GIGA:  "GiB",
	BaseBinary * TERA:  "TiB",
	BaseBinary * PETA:  "PiB",
	BaseBinary * EXA:   "EiB",
	BaseBinary * ZETTA: "ZiB",
	BaseBinary * YOTTA: "YiB",
}

// CastToFloat64 converts the given value to float64
func CastToFloat64(value interface{}) (float64, error) {
	var v float64
	switch value.(type) {
	case float64:
		v = value.(float64)
	case float32:
		v = float64(value.(float32))
	case uint:
		v = float64(value.(uint))
	case uint64:
		v = float64(value.(uint64))
	case uint32:
		v = float64(value.(uint32))
	case uint16:
		v = float64(value.(uint16))
	case byte:
		v = float64(value.(byte))
	case int:
		v = float64(value.(int))
	case int64:
		v = float64(value.(int64))
	case int32:
		v = float64(value.(int32))
	case int16:
		v = float64(value.(int16))
	case int8:
		v = float64(value.(int8))
	default:
		return 0.0, fmt.Errorf("unsupported type")
	}
	return v, nil
}

// ConvUnitDec returns the string representation of value in 1000 base with given unit.
func ConvUnitDec(value interface{}, target float64) (string, error) {
	v, err := CastToFloat64(value)
	if err != nil {
		return "", err
	}
	return ConvUnit(v, BaseDecimal, target, false)
}

// ConvUnitBin returns the string representation of value in 1024 base with given unit.
func ConvUnitBin(value interface{}, target float64) (string, error) {
	v, err := CastToFloat64(value)
	if err != nil {
		return "", err
	}
	return ConvUnit(v, BaseBinary, target, false)
}

// ConvUnitDecFit returns the string representation of value in 1000 base with given unit.
// The unit will be optimized automatically.
func ConvUnitDecFit(value interface{}, target float64) (string, error) {
	v, err := CastToFloat64(value)
	if err != nil {
		return "", err
	}
	return ConvUnit(v, BaseDecimal, target, true)
}

// ConvUnitBinFit returns the string representation of value in 1024 base with given unit.
// The unit will be optimized automatically.
func ConvUnitBinFit(value interface{}, target float64) (string, error) {
	v, err := CastToFloat64(value)
	if err != nil {
		return "", err
	}
	return ConvUnit(v, BaseBinary, target, true)
}

// ConvUnit returns the string representation of value in given base with given unit.
// The unit will be optimized automatically if fit is true.
func ConvUnit(value float64, mlt float64, target float64, fit bool) (string, error) {
	if mlt != BaseBinary && mlt != BaseDecimal {
		return "", errors.New("invalid multipler")
	}

	if value == 0 || target == BYTE {
		return fmt.Sprintf("%.0fB", value), nil
	}

	if value < 0 {
		return "", errors.New("invalid value")
	}

	if _, ok := units[mlt*target]; !ok {
		return "", errors.New("unsupported range")
	}

	// denom never be 0 because mlt must be greater than 0
	res := math.Floor(value / math.Pow(mlt, target))
	if res < 1.0 {
		return "0B", nil
	}

	if mlt >= res || !fit {
		return fmt.Sprintf("%.1f%s", res, units[mlt*target]), nil
	}

	tgt := target
	for {
		if _, ok := units[mlt*tgt]; !ok {
			break
		}
		res = res / mlt
		tgt++
		if mlt >= res {
			break
		}
		res = math.Floor(res)
	}
	return fmt.Sprintf("%.1f%s", res, units[mlt*tgt]), nil
}
