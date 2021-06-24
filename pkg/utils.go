package pkg

import "reflect"

func Tags2Map(obj interface{}) map[string]string {
	t := reflect.TypeOf(obj)
	v := reflect.ValueOf(obj)

	var data = make(map[string]string)
	for i := 0; i < t.NumField(); i++ {
		data[t.Field(i).Name] = v.Field(i).String()
	}
	return data
}

func Fields2Map(obj interface{}) map[string]interface{} {
	t := reflect.TypeOf(obj)
	v := reflect.ValueOf(obj)

	var data = make(map[string]interface{})
	for i := 0; i < t.NumField(); i++ {
		data[t.Field(i).Name] = v.Field(i).Interface()
	}
	return data
}
