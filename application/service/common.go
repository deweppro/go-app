package service

import "reflect"

var srvType = reflect.TypeOf(new(IService)).Elem()

//AsService service interface check
func AsService(v reflect.Value) (IService, bool) {
	if v.Type().AssignableTo(srvType) {
		return v.Interface().(IService), true
	}
	return nil, false
}

//IsService service interface check
func IsService(v interface{}) bool {
	if _, ok := v.(IService); ok {
		return true
	}
	if _, ok := v.(IServiceCtx); ok {
		return true
	}
	return false
}

var srvTypeCtx = reflect.TypeOf(new(IServiceCtx)).Elem()

//AsServiceCtx service interface check
func AsServiceCtx(v reflect.Value) (IServiceCtx, bool) {
	if v.Type().AssignableTo(srvTypeCtx) {
		return v.Interface().(IServiceCtx), true
	}
	return nil, false
}
