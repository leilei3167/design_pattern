package db

import (
	"reflect"
	"strings"
)

//record代表表中的数据,利用反射模拟了将结构体字段的储存(类似表的分列)

type record struct {
	primaryKey any
	fields     map[string]int
	values     []any
}

func (r record) convertByValue(result any) (e error) {
	defer func() { //捕获反射时的panic
		if err := recover(); err != nil {
			e = ErrRecordTypeInvalid
		}
	}()

	rType := reflect.TypeOf(result)
	rVal := reflect.ValueOf(result)

	if rType.Kind() == reflect.Pointer { //如果是指针类型,则获取其真正的值
		rType = rType.Elem()
		rVal = rVal.Elem()
	}

	for i := 0; i < rType.NumField(); i++ { //将record字段值还原到字段中
		field := rVal.Field(i)
		field.Set(reflect.ValueOf(r.values[i]))
	}
	return nil
}

func (r record) convertByType(rType reflect.Type) (result any, e error) {
	defer func() { //捕获反射时的panic
		if err := recover(); err != nil {
			e = ErrRecordTypeInvalid
		}
	}()

	if rType.Kind() == reflect.Pointer {
		rType = rType.Elem()
	}
	rVal := reflect.New(rType)
	return rVal, nil
}

//创建一条记录
func recordFrom(key, value any) (r record, e error) {
	defer func() { //捕获反射时的panic
		if err := recover(); err != nil {
			e = ErrRecordTypeInvalid
		}
	}()

	//获取value的类型和值
	vType := reflect.TypeOf(value)
	vVal := reflect.ValueOf(value)
	if vVal.Type().Kind() == reflect.Pointer {
		vType = vType.Elem()
		vVal = vVal.Elem()
	}

	record := record{
		primaryKey: key,
		fields:     make(map[string]int, vVal.NumField()),
		values:     make([]any, vVal.NumField()), //相当于列,单独存每个字段的值
	}
	for i := 0; i < vVal.NumField(); i++ {
		fieldType := vType.Field(i)
		fieldVal := vVal.Field(i)
		name := strings.ToLower(fieldType.Name)
		record.fields[name] = i
		record.values[i] = fieldVal.Interface()
	}
	return record, nil
}
