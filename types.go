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
type ZInt64Arr []int64
type ZStrArr []string
type StrArrArr [][]string // [["a","b"],["c","d"]]

type IntArrArr [][]int
type Int3Arr [][][]int   //[[[],[]],[[],[]]]
type Int64ArrArr [][]int64
type FloatArrArr [][]float64
type SimpleObj map[string]interface{}      // {"a":1,"b":true,"c":"string","d":3.44}
type ArrSimpleObj []map[string]interface{} // [{"a":1,"b":true,"c":"string","d":3.44}, {"f":123}]
type ZArrObj []map[string]interface{}
type ZIntObj map[int]int
type IntMap map[int]interface{}
type ZTime time.Time
type ZItemList []ZItem
type ZItem struct {
	Id  int `json:"i"`
	Num int `json:"n"`
}

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

func (t IntMap) Value() (driver.Value, error) {
	tstring, err := json.Marshal(t)
	return tstring, err
}

func (t *IntMap) Scan(v interface{}) error {
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

func (t *ZIntArr) Join(sep string) string {
	str := ""
	for idx, loc := range *t {
		str += strconv.Itoa(loc)
		if idx != len(*t)-1 {
			str += sep
		}
	}
	return str
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

func (t ZInt64Arr) Value() (driver.Value, error) {
	strarr := []string{}
	for _, loc := range t {
		strarr = append(strarr, strconv.FormatInt(loc, 10))
	}
	result := strings.Join(strarr, ",")
	return result, nil
}

func (t *ZInt64Arr) Join(sep string) string {
	str := ""
	for idx, loc := range *t {
		str += strconv.FormatInt(loc, 10)
		if idx != len(*t)-1 {
			str += sep
		}
	}
	return str
}

func (t *ZInt64Arr) Scan(v interface{}) error {
	vtype := reflect.TypeOf(v)
	switch vtype.String() {
	case "[]uint8":
		str := string(v.([]byte))
		if str == "" {
			*t = []int64{}
			return nil
		}
		strarr := strings.Split(str, ",")
		for _, i2 := range strarr {
			i, _ := strconv.ParseInt(i2, 10, 64)
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

// 二维数组
func (t IntArrArr) Value() (driver.Value, error) {
	strarr := []string{}
	for _, arr := range t {
		arr2str := []string{}
		for _, ii := range arr {
			arr2str = append(arr2str, strconv.Itoa(ii))
		}
		strarr = append(strarr, strings.Join(arr2str, ","))
	}
	result := strings.Join(strarr, ";")
	return result, nil
}

func (t *IntArrArr) Scan(v interface{}) error {
	vtype := reflect.TypeOf(v)
	switch vtype.String() {
	case "[]uint8":
		str := string(v.([]byte))
		if str == "" {
			*t = [][]int{}
			return nil
		}
		strarr := strings.Split(str, ";")
		var arrArr [][]int
		for _, str := range strarr {
			arr := strings.Split(str, ",")
			intArr := []int{}
			for _, strInArr := range arr {
				ss, _ := strconv.Atoi(strInArr)
				intArr = append(intArr, ss)
			}
			arrArr = append(arrArr, intArr)
		}
		*t = arrArr
	}
	return nil
}

// 三维数组
func (t Int3Arr) Value() (driver.Value, error) {
	strarr := []string{}
	for _, arr := range t {
		arr2str := []string{}
		for _, ii := range arr {
			arr3str := []string{}
			for _, iii := range ii {
				arr3str = append(arr3str, strconv.Itoa(iii))
			}
			arr2str = append(arr2str, strings.Join(arr3str, ","))
		}
		strarr = append(strarr, strings.Join(arr2str, ";"))
	}
	result := strings.Join(strarr, "|")
	return result, nil
}

func (t *Int3Arr) Scan(v interface{}) error {
	vtype := reflect.TypeOf(v)
	switch vtype.String() {
	case "[]uint8":
		str := string(v.([]byte))
		if str == "" {
			*t = [][][]int{}
			return nil
		}
		strarr := strings.Split(str, "|")
		var arrArr [][][]int
		for _, str := range strarr { // 1,1;2,2
			arr := strings.Split(str, ";")
			intArr := [][]int{}
			for _, strInArr := range arr {
				arr3 := strings.Split(strInArr, ",")
				int3Arr := []int{}
				for _, sub3 := range arr3 {
					ss, _ := strconv.Atoi(sub3)
					int3Arr = append(int3Arr, ss)
				}
				intArr = append(intArr, int3Arr)
			}
			arrArr = append(arrArr, intArr)
		}
		*t = arrArr
	}
	return nil
}

func (t Int64ArrArr) Value() (driver.Value, error) {
	strarr := []string{}
	for _, arr := range t {
		arr2str := []string{}
		for _, ii := range arr {
			arr2str = append(arr2str, strconv.FormatInt(int64(ii), 10))
		}
		strarr = append(strarr, strings.Join(arr2str, ","))
	}
	result := strings.Join(strarr, ";")
	return result, nil
}

func (t *Int64ArrArr) Scan(v interface{}) error {
	vtype := reflect.TypeOf(v)
	switch vtype.String() {
	case "[]uint8":
		str := string(v.([]byte))
		if str == "" {
			*t = [][]int64{}
			return nil
		}
		strarr := strings.Split(str, ";")
		var arrArr [][]int64
		for _, str := range strarr {
			arr := strings.Split(str, ",")
			intArr := []int64{}
			for _, strInArr := range arr {
				ss, _ := strconv.ParseInt(strInArr, 10, 64)
				intArr = append(intArr, ss)
			}
			arrArr = append(arrArr, intArr)
		}
		*t = arrArr
	}
	return nil
}

// 二维浮点数组
func (t FloatArrArr) Value() (driver.Value, error) {
	strarr := []string{}
	for _, arr := range t {
		arr2str := []string{}
		for _, ii := range arr {
			arr2str = append(arr2str, fmt.Sprintf("%.4f", ii))
		}
		strarr = append(strarr, strings.Join(arr2str, ","))
	}
	result := strings.Join(strarr, ";")
	return result, nil
}

func (t *FloatArrArr) Scan(v interface{}) error {
	vtype := reflect.TypeOf(v)
	switch vtype.String() {
	case "[]float64":
		str := string(v.([]byte))
		if str == "" {
			*t = [][]float64{}
			return nil
		}
		strarr := strings.Split(str, ";")
		var arrArr [][]float64
		for _, str := range strarr {
			arr := strings.Split(str, ",")
			intArr := []float64{}
			for _, strInArr := range arr {
				ss, _ := strconv.ParseFloat(strInArr, 64)
				intArr = append(intArr, ss)
			}
			arrArr = append(arrArr, intArr)
		}
		*t = arrArr
	}
	return nil
}

// 二维字符串数组
func (t StrArrArr) Value() (driver.Value, error) {
	strarr := []string{}
	for _, arr := range t {
		strarr = append(strarr, strings.Join(arr, ","))
	}
	result := strings.Join(strarr, ";")
	return result, nil
}

func (t *StrArrArr) Scan(v interface{}) error {
	vtype := reflect.TypeOf(v)
	switch vtype.String() {
	case "[]uint8":
		str := string(v.([]byte))
		if str == "" {
			*t = [][]string{}
			return nil
		}
		strarr := strings.Split(str, ";")
		var arrArr [][]string
		for _, str := range strarr {
			arr := strings.Split(str, ",")
			arrArr = append(arrArr, arr)
		}
		*t = arrArr
	}
	return nil
}

// 简写json(object)
func (t SimpleObj) Value() (driver.Value, error) {
	var result string
	for k1, v1 := range t {
		vtype := reflect.TypeOf(v1)
		switch vtype.String() {
		case "string":
			result += k1 + ":" + v1.(string) + ","
			break
		case "float64":
			result += k1 + ":" + strconv.FormatFloat(v1.(float64), 'f', -1, 64) + ","
			break
		}
	}
	result = strings.Replace(result, ",", "", -2)
	return result, nil
}

func (t *SimpleObj) Scan(v interface{}) error {
	vtype := reflect.TypeOf(v)
	switch vtype.String() {
	case "[]uint8":
		str := string(v.([]byte))
		if str == "" {
			*t = map[string]interface{}{}
			return nil
		}
		strarr := strings.Split(str, ",")
		var obj map[string]interface{}
		for _, str := range strarr {
			arr := strings.Split(str, ":")
			var strValue interface{}
			if arr[1] == "null" || arr[1] == "" {
				strValue = nil
			} else if arr[1] == "true" {
				strValue = true
			} else if arr[1] == "false" {
				strValue = false
			} else {
				intValue, err := strconv.ParseInt(arr[1], 10, 64)
				if err != nil {
					floatValue, err := strconv.ParseFloat(arr[1], 64)
					if err != nil {
						strValue = arr[1]
					} else {
						strValue = floatValue
					}
				} else {
					strValue = intValue
				}
			}
			obj[arr[0]] = strValue
		}
		*t = obj
	}
	return nil
}

// obj数组
func (t ArrSimpleObj) Value() (driver.Value, error) {
	var result string
	for _, m := range t {
		for k1, v1 := range m {
			vtype := reflect.TypeOf(v1)
			switch vtype.String() {
			case "string":
				result += k1 + ":" + v1.(string) + ","
				break
			case "float64":
				result += k1 + ":" + strconv.FormatFloat(v1.(float64), 'f', -1, 64) + ","
				break
			}
		}
	}
	result = strings.Replace(result, ",", "", -2)
	return result, nil
}

func (t *ArrSimpleObj) Scan(v interface{}) error {
	vtype := reflect.TypeOf(v)
	switch vtype.String() {
	case "[]uint8":
		str := string(v.([]byte))
		if str == "" {
			*t = []map[string]interface{}{}
			return nil
		}
		strarr := strings.Split(str, ",")
		var obj []map[string]interface{}
		for _, str := range strarr {
			arr := strings.Split(str, ":")
			obj = append(obj, map[string]interface{}{arr[0]: arr[1]})
		}
		*t = obj
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

func (t ZIntObj) Value() (driver.Value, error) {
	tstring, err := json.Marshal(t)
	return tstring, err
}

func (t *ZIntObj) Scan(v interface{}) error {
	json.Unmarshal(v.([]byte), t)
	return nil
}

func (t ZItemList) Value() (driver.Value, error) {
	tstring, err := json.Marshal(t)
	return tstring, err
}

func (t *ZItemList) Scan(v interface{}) error {
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

// 集合 set 
// 定义一个泛型类型的集合
type Set[T comparable] map[T]struct{}

// 创建一个新的集合
func NewSet[T comparable]() *Set[T] {
	set := make(Set[T]) // 新建一个空集合
	return &set
}

func (t Set[T]) Value() (driver.Value, error) {
	list := t.List()
	tstring, err := json.Marshal(list)
	return tstring, err
}

func (t *Set[T]) Scan(v interface{}) error {
	list := []T{}
	json.Unmarshal(v.([]byte), &list)
	if (*t) == nil {
		(*t) = Set[T]{}
	}
	(*t).AddList(list)
	return nil
}

// 添加元素
func (set *Set[T]) Add(key T) {
	(*set)[key] = struct{}{}
}

// 移除元素
func (set *Set[T]) Remove(key T) {
	delete(*set, key)
}

// 检查元素是否存在
func (set *Set[T]) Contains(key T) bool {
	_, exists := (*set)[key]
	return exists
}

// 获取集合的长度
func (set *Set[T]) Len() int {
	return len(*set)
}

// 批量添加元素
func (set *Set[T]) AddList(list []T) {
	for _, key := range list {
		set.Add(key)
	}
}

// 列出所有元素
func (set *Set[T]) List() []T {
	list := make([]T, 0, len(*set))
	for key := range *set {
		list = append(list, key)
	}
	return list
}

// 交集
func (set *Set[T]) Intersect(another Set[T]) (crossSet *Set[T]) {
	intersection := Set[T]{}
	for elem := range another {
		if set.Contains(elem) {
			intersection.Add(elem)
		}
	}
	return &intersection
}
