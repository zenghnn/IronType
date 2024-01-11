package IronType

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"
	"reflect"
	"strconv"
	"strings"
	"time"
)

type ZJson map[string]interface{}
type ZIntArr []int
type ZStrArr []string
type ZArrObj []map[string]interface{}
type ZTime time.Time

const (
	TimeFormat1 = "2006-01-02 15:04:05.999"
	TimeFormat2 = "20060102"
	TimeFormat3 = "2006-01-02 15:04:05"
)

func (t ZJson) Value() (driver.Value, error) {
	tstring, err := json.Marshal(t)
	return tstring, err
}

func (t *ZJson) Scan(v interface{}) error {
	json.Unmarshal(v.([]byte), t)
	return nil
}

func (t ZIntArr) Value() (driver.Value, error) {
	strarr := []string{}
	for _, loc := range t {
		strarr = append(strarr, strconv.Itoa(loc))
	}
	result := strings.Join(strarr, ",")
	return result, nil
}

func (t *ZIntArr) Scan(v interface{}) error {
	vtype := reflect.TypeOf(v)
	switch vtype.String() {
	case "[]uint8":
		str := string(v.([]byte))
		if str == "" {
			*t = []int{}
			return nil
		}
		strarr := strings.Split(str, ",")
		for _, i2 := range strarr {
			i, _ := strconv.Atoi(i2)
			*t = append(*t, i)
		}
	}
	return nil
}

func (t ZStrArr) Value() (driver.Value, error) {
	strarr := []string{}
	for _, loc := range t {
		strarr = append(strarr, loc)
	}
	result := strings.Join(strarr, ",")
	return result, nil
}

func (t *ZStrArr) Scan(v interface{}) error {
	vtype := reflect.TypeOf(v)
	switch vtype.String() {
	case "[]uint8":
		str := string(v.([]byte))
		if str == "" {
			*t = []string{}
			return nil
		}
		strarr := strings.Split(str, ",")
		*t = strarr
	}
	return nil
}

func (t ZArrObj) Value() (driver.Value, error) {
	tstring, err := json.Marshal(t)
	return tstring, err
}

func (t *ZArrObj) Scan(v interface{}) error {
	json.Unmarshal(v.([]byte), t)
	return nil
}

func (t ZTime) Value() (driver.Value, error) {
	tstring := t.String()
	return tstring, nil
}

// 出
func (t *ZTime) Scan(v interface{}) error {
	if v == nil {
		t = nil
		return nil
	} else {
		temp := ZTime(v.(time.Time))
		*t = temp
		return nil
	}
}

func (t *ZTime) String() string {
	//如果是空值就返回null
	if time.Time(*t).IsZero() {
		return ""
	}
	return time.Time(*t).Format(TimeFormat1)
}
func (t *ZTime) SqlString() string {
	if time.Time(*t).IsZero() {
		return "null"
	}
	return "'" + time.Time(*t).Format(TimeFormat1) + "'"
}
func (t *ZTime) IsEmpty() bool {
	return (t) == nil
}
func (t *ZTime) ToTime() time.Time {
	return time.Time(*t)
}
func (t ZTime) MarshalJSON() ([]byte, error) {
	var stamp = fmt.Sprintf("\"%s\"", t.String())
	return []byte(stamp), nil
}
func (t *ZTime) UnmarshalJSON(b []byte) error {
	//b = bytes.Trim(b, "\"")
	ext, err := time.Parse(TimeFormat1, string(b))
	if err != nil {
		return err
	}
	*t = ZTime(ext)
	return nil
}
