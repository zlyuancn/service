package mysql_binlog

import (
	"encoding"
	"encoding/base64"
	"fmt"
	"reflect"
	"strings"
	"time"

	jsoniter "github.com/json-iterator/go"
	"github.com/zlyuancn/zstr"
)

// 目标字段
type targetField struct {
	value  *reflect.Value
	key    string
	string bool // 将类型视为字符串, 列入CHAR, VARCHAR, TINYBLOB, TINYTEXT, BLOB, TEXT, MEDIUMBLOB, MEDIUMTEXT, LONGBLOB, LONGTEXT, 原始数据必须是 string或nil
	json   bool // 将类型视为JSON, 原始数据必须是 string或nil
	point  bool // 将类型视为POINT, 原始数据必须是 长度为2的[]float64或nil
	binary bool // 将类型视为BINARY, 原始数据必须是 base64string或nil
}

func newScanField(fv *reflect.Value, t *reflect.StructField) (*targetField, bool) {
	k, ok := t.Tag.Lookup("scan")
	if !ok {
		k, ok = t.Tag.Lookup("json")
	}
	if !ok {
		k = t.Name
	}
	if k == "" || k == "-" {
		return nil, false
	}

	ks := strings.Split(k, ",")
	sf := &targetField{
		value: fv,
		key:   ks[0],
	}
	for _, opt := range ks[1:] {
		switch opt {
		case "string":
			sf.string = true
		case "json":
			sf.json = true
		case "point":
			sf.point = true
		case "binary":
			sf.binary = true
		}
	}
	return sf, true
}

// 扫描时间类型
func scanTimeType(a interface{}, t *time.Time) (err error) {
	s := zstr.AnyToStr(a)
	switch len(s) {
	case 8: // TIME
		*t, err = time.ParseInLocation("15:04:05", s, time.Local)
	case 10: // DATE
		*t, err = time.ParseInLocation("2006-01-02", s, time.Local)
	case 19: // DATETIME TIMESTAMP
		*t, err = time.ParseInLocation("2006-01-02 15:04:05", s, time.Local)
	default:
		return fmt.Errorf("数据(%v)不能转换为*time.Time", a)
	}
	if err != nil {
		return fmt.Errorf("数据(%v)不能转换为*time.Time", a)
	}
	return
}

// 将json数据扫描到outPtr
func scanJsonType(a interface{}, outPtr interface{}) error {
	switch s := a.(type) {
	case string:
		return jsoniter.UnmarshalFromString(s, outPtr)
	}
	return fmt.Errorf("数据(%v)不能转换为%T", a, outPtr)
}

// 扫描切片(原始数据是切片)
func scanSlice(a []interface{}, outValue *reflect.Value, outPtrRealType reflect.Type) error {
	itemFieldType := outPtrRealType.Elem()
	itemFieldTypeIsPtr := itemFieldType.Kind() == reflect.Ptr
	if itemFieldTypeIsPtr {
		itemFieldType = itemFieldType.Elem()
	}

	items := make([]reflect.Value, len(a))
	for i, temp := range a {
		item := reflect.New(itemFieldType)
		err := zstr.ScanAny(temp, item.Interface())
		if err != nil {
			return err
		}
		if !itemFieldTypeIsPtr {
			item = item.Elem()
		}
		items[i] = item
	}
	out := outValue.Elem()
	values := reflect.Append(out, items...)
	out.Set(values)
	return nil
}

// 将字段数据扫描到outPtr中
func scanFieldData(a interface{}, outPtr interface{}, outValue *reflect.Value, outPtrRealType reflect.Type, tf *targetField) error {
	// 指定类型扫描
	if tf.string || tf.json {
		s, ok := a.(string)
		if !ok {
			return fmt.Errorf("json标签的原始数据必须是string, 但是得到%T", a)
		}
		switch tp := outPtr.(type) {
		case encoding.BinaryUnmarshaler:
			return tp.UnmarshalBinary(zstr.StringToBytes(&s))
		}

		switch tp := outPtr.(type) {
		case zstr.AnyUnmarshaler:
			return tp.UnmarshalAny(s)
		}

		if tf.string {
			return zstr.Scan(s, outPtr)
		}
		return scanJsonType(a, outPtr)
	}
	if tf.point {
		var f64s []float64

		temp, ok := a.([]interface{})
		if ok {
			if len(temp) != 2 {
				return fmt.Errorf("point标签的原始数据必须是长度为2的[]float64或[]interface{}, 但是收到长度%d", len(temp))
			}

			f64s = make([]float64, 2)
			for i, tt := range temp {
				if v, ok := tt.(float64); ok { // float64
					f64s[i] = v
					continue
				}
				return fmt.Errorf("point标签的原始数据必须是[]float64或[]interface{}, 但是切片中其中一个类型为%T", tt)
			}
		} else {
			temp, ok := a.([]float64)
			if !ok {
				return fmt.Errorf("point标签的原始数据必须是[]float64或[]interface{}, 但是得到%T", a)
			}
			if len(temp) != 2 {
				return fmt.Errorf("point标签的原始数据必须是长度为2的[]float64或[]interface{}, 但是收到长度%d", len(temp))
			}
			f64s = append(make([]float64, 0, 2), temp...)
		}

		switch t := outPtr.(type) {
		case zstr.AnyUnmarshaler:
			return t.UnmarshalAny(f64s)
		}
		return scanSlice(temp, outValue, outPtrRealType)
	}
	if tf.binary {
		s, ok := a.(string)
		if !ok {
			return fmt.Errorf("binary标签的原始数据必须是string, 但是得到%T", a)
		}
		bs, err := base64.StdEncoding.DecodeString(s)
		if err != nil {
			return fmt.Errorf("binary标签的原始数据使用base64解码失败: %v", err)
		}

		switch tp := outPtr.(type) {
		case encoding.BinaryUnmarshaler:
			return tp.UnmarshalBinary(zstr.StringToBytes(&s))
		}

		switch tp := outPtr.(type) {
		case zstr.AnyUnmarshaler:
			return tp.UnmarshalAny(s)
		}

		switch tp := outPtr.(type) {
		case *[]byte:
			*tp = bs
		case *string:
			*tp = *zstr.BytesToString(bs)
		default:
			return fmt.Errorf("binary 数据无法解码到(%T), 考虑为它实现encoding.BinaryUnmarshaler接口或zstr.AnyUnmarshaler", outPtr)
		}
		return nil
	}

	// 这些类型要优先于自定义扫描
	switch t := outPtr.(type) {
	case *time.Time:
		return scanTimeType(a, t)
	}

	// 自定义扫描
	switch ta := a.(type) {
	case []byte:
		switch tp := outPtr.(type) {
		case encoding.BinaryUnmarshaler:
			return tp.UnmarshalBinary(ta)
		}
	case string:
		switch tp := outPtr.(type) {
		case encoding.BinaryUnmarshaler:
			return tp.UnmarshalBinary(zstr.StringToBytes(&ta))
		}
	}
	switch tp := outPtr.(type) {
	case zstr.AnyUnmarshaler:
		return tp.UnmarshalAny(a)
	}

	switch ta := a.(type) {
	case []interface{}:
		return scanSlice(ta, outValue, outPtrRealType)
	case string: // 考虑json
		switch outPtrRealType.Kind() {
		case reflect.Struct, reflect.Map, reflect.Slice, reflect.Interface: // 只有json数据才行
			return scanJsonType(a, outPtr)
		}
	}

	return zstr.ScanAny(a, outPtr)
}

// 将mysql数据扫描到结构体中, 字段标签优先取 scan, 然后取 json, 都没有则取字段名, 如果标签为空或短横线-则忽略
func ScanMysqlData(m map[string]interface{}, outPtr interface{}) (err error) {
	aType := reflect.TypeOf(outPtr)
	if aType.Kind() != reflect.Ptr {
		return fmt.Errorf("outPtr必须是带指针的结构体")
	}
	aType = aType.Elem()
	if aType.Kind() != reflect.Struct {
		return fmt.Errorf("outPtr必须是带指针的结构体")
	}
	aValue := reflect.ValueOf(outPtr).Elem()

	// 扫描结构体
	fields := make(map[string]*targetField, aType.NumField())
	for i := 0; i < aType.NumField(); i++ {
		field := aType.Field(i)
		if field.PkgPath != "" {
			continue
		}
		v := aValue.Field(i)
		tf, ok := newScanField(&v, &field)
		if !ok {
			continue
		}
		fields[tf.key] = tf
	}

	// 循环写入数据
	for k, v := range m {
		sf, ok := fields[k]
		if !ok {
			continue
		}

		fieldType := sf.value.Type()
		fieldIsPtr := fieldType.Kind() == reflect.Ptr
		if fieldIsPtr {
			fieldType = fieldType.Elem()
		}

		temp := reflect.New(fieldType) // 根据实际类型创建临时值

		// 空值
		if v == nil {
			unmarshaler, ok := temp.Interface().(zstr.AnyUnmarshaler)
			if !ok { // 直接忽略空值
				continue
			}
			err = unmarshaler.UnmarshalAny(nil)
		} else {
			err = scanFieldData(v, temp.Interface(), &temp, fieldType, sf) // 将数据扫描到临时值
		}

		if err != nil {
			return err
		}
		if !fieldIsPtr {
			temp = temp.Elem() // 如果field不是ptr则只需要实际内容
		}
		sf.value.Set(temp)
	}
	return nil
}
