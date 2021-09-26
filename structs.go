package structs

import (
	"errors"
	"fmt"
	"reflect"
	"time"
)

func ToMap(obj interface{}, tagName string) map[string]interface{} {
	return reflectToMap(reflect.ValueOf(obj), tagName)
}

func reflectToMap(iVal reflect.Value, tagName string) map[string]interface{} {
	iVal = reflect.Indirect(iVal)
	iType := iVal.Type()
	//必须是一个结构体
	if iType.Kind() != reflect.Struct {
		panic(errors.New("only support struct or struct pointer"))
	}
	dest := make(map[string]interface{}, iVal.NumField())
	for i := 0; i < iVal.NumField(); i++ {
		var name string
		if iType.Field(i).Anonymous {
			//递归解析嵌套的匿名字段
			for k, v := range reflectToMap(iVal.Field(i), tagName) {
				dest[k] = v
			}
			continue
		}
		if v, ok := iType.Field(i).Tag.Lookup(tagName); ok {
			name = v
		} else {
			name = iType.Field(i).Name
		}
		fmt.Println(iVal.Field(i).Kind())
		switch iVal.Field(i).Kind() {
		case reflect.Struct:
			fv := iVal.Field(i).Interface()
			switch v := fv.(type) {
			case time.Time, *time.Time:
				dest[name] = v
			default:
				dest[name] = reflectToMap(iVal.Field(i), tagName)
			}
		default:
			dest[name] = iVal.Field(i).Interface()
		}
	}
	return dest
}

// MergeMap 合并map
//不处理嵌套情况
func MergeMap(dest map[string]interface{}, ms ...map[string]interface{}) {
	if dest == nil {
		dest = make(map[string]interface{})
	}
	
	for _, m := range ms {
		for k, v := range m {
			dest[k] = v
		}
	}
}

// MergeStruct 合并结构体
//只能合并与目标结构体相同的字段，目标结构体中没有的字段会直接忽略
func MergeStruct(dest interface{}, ms ...interface{}) {
	iVal := reflect.ValueOf(dest)
	//must be a struct pointer
	if iVal.Type().Kind() != reflect.Ptr {
		panic(errors.New("only support struct pointer"))
	}
	iVal = reflect.Indirect(iVal)
	iType := iVal.Type()
	
	is := make([]*Struct, len(ms))
	for i := range ms {
		is[i] = NewStruct(ms[i], "")
	}
	
	for i := 0; i < iVal.NumField(); i++ {
		if !iVal.Field(i).CanSet() {
			continue
		}
		name := iType.Field(i).Name
		for j := len(is) - 1; j >= 0; j-- {
			fv := is[j].FieldByName(name)
			if !fv.IsZero() {
				iVal.Field(i).Set(fv)
				break
			}
		}
	}
}

type Struct struct {
	iVal    reflect.Value
	iType   reflect.Type
	TagName string
}

func NewStruct(i interface{}, tagName string) *Struct {
	iVal := reflect.Indirect(reflect.ValueOf(i))
	return &Struct{
		iVal:    iVal,
		iType:   iVal.Type(),
		TagName: tagName,
	}
}

func (s *Struct) FieldByName(name string) reflect.Value {
	return s.iVal.FieldByName(name)
}
