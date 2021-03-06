package util

import (
	"errors"
	"github.com/XH-JMC/cta/common/logs"
	"reflect"
	"strconv"
	"time"
)

func Protect(f func()) {
	defer func() {
		if err := recover(); err != nil {
			logs.Warnf("recover panic: %s", err)
		}
	}()
	f()
}

// try once first, and if it fails, retry ${retryTimes} times
func Retry(retryTimes int, retryInterval time.Duration, f func() bool) bool {
	if f() {
		return true
	}
	for i := 0; i < retryTimes; i++ {
		if f() {
			return true
		}
		time.Sleep(retryInterval)
	}
	return false
}

func Interface2Int(item interface{}) (int, error) {
	i64, err := Interface2Int64(item)
	return int(i64), err
}

func Interface2Int64(item interface{}) (int64, error) {
	switch val := item.(type) {
	case byte:
		return int64(val), nil
	case bool:
		if val {
			return 1, nil
		} else {
			return 0, nil
		}
	case int:
		return int64(val), nil
	case int8:
		return int64(val), nil
	case int16:
		return int64(val), nil
	case int32:
		return int64(val), nil
	case int64:
		return val, nil
	case float32:
		return int64(val), nil
	case float64:
		return int64(val), nil
	case string:
		return strconv.ParseInt(val, 10, 64)
	case []byte:
		return strconv.ParseInt(string(val), 10, 64)
	default:
		return 0, errors.New("unsupported type")
	}
}

func Interface2Float64(item interface{}) (float64, error) {
	switch val := item.(type) {
	case byte:
		return float64(val), nil
	case bool:
		if val {
			return 1, nil
		} else {
			return 0, nil
		}
	case int:
		return float64(val), nil
	case int8:
		return float64(val), nil
	case int16:
		return float64(val), nil
	case int32:
		return float64(val), nil
	case int64:
		return float64(val), nil
	case float32:
		return float64(val), nil
	case float64:
		return val, nil
	case string:
		return strconv.ParseFloat(val, 64)
	case []byte:
		return strconv.ParseFloat(string(val), 64)
	default:
		return 0, errors.New("unsupported type")
	}
}

func Interface2Bool(item interface{}) (bool, error) {
	switch val := item.(type) {
	case byte:
		return val != 0, nil
	case bool:
		return val, nil
	case int:
		return val != 0, nil
	case int8:
		return val != 0, nil
	case int16:
		return val != 0, nil
	case int32:
		return val != 0, nil
	case int64:
		return val != 0, nil
	case float32:
		return val != 0, nil
	case float64:
		return val != 0, nil
	case string:
		return strconv.ParseBool(val)
	case []byte:
		return strconv.ParseBool(string(val))
	default:
		return false, errors.New("unsupported type")
	}
}

func Interface2String(item interface{}) (string, error) {
	switch val := item.(type) {
	case byte:
		return strconv.FormatInt(int64(val), 10), nil
	case bool:
		if val {
			return "true", nil
		} else {
			return "false", nil
		}
	case int:
		return strconv.FormatInt(int64(val), 10), nil
	case int8:
		return strconv.FormatInt(int64(val), 10), nil
	case int16:
		return strconv.FormatInt(int64(val), 10), nil
	case int32:
		return strconv.FormatInt(int64(val), 10), nil
	case int64:
		return strconv.FormatInt(val, 10), nil
	case float32:
		return strconv.FormatFloat(float64(val), 'f', 20, 64), nil
	case float64:
		return strconv.FormatFloat(val, 'f', 20, 64), nil
	case string:
		return val, nil
	case []byte:
		return string(val), nil
	default:
		return "", errors.New("unsupported type")
	}
}

func interface2String(item interface{}) (string, bool) {
	switch val := item.(type) {
	case string:
		return val, true
	case []byte:
		return string(val), true
	default:
		return "", false
	}
}

func InterfaceEqual(x, y interface{}) bool {
	xstr, xok := interface2String(x)
	ystr, yok := interface2String(y)
	if xok && yok {
		return xstr == ystr
	}
	if xok || yok {
		return false
	}

	xfloat, xerr := Interface2Float64(x)
	yfloat, yerr := Interface2Float64(y)
	if xerr == nil && yerr == nil {
		return xfloat == yfloat
	}
	if xerr == nil || yerr == nil {
		return false
	}

	return reflect.DeepEqual(x, y)
}

// 将所有类型为map[interface{}]interface{}的子孙元素转换为map[string]interface{}
func ConvertMapInterfaceInterface2MapStringInterface(container interface{}) interface{} {
	switch c := container.(type) {
	case []interface{}:
		for i, item := range c {
			c[i] = ConvertMapInterfaceInterface2MapStringInterface(item)
		}
		return c
	case map[string]interface{}:
		for key, value := range c {
			c[key] = ConvertMapInterfaceInterface2MapStringInterface(value)
		}
		return c
	case map[interface{}]interface{}:
		ret := make(map[string]interface{})
		for key, value := range c {
			keyStr, _ := interface2String(key)
			ret[keyStr] = ConvertMapInterfaceInterface2MapStringInterface(value)
		}
		return ret
	default:
		return c
	}
}
