package funcmap

import (
	"errors"
	"reflect"
)

var (
	ErrParamsNotAdapted = errors.New("The number of params does not match.")
)

type FuncMap map[string]reflect.Value

func New() FuncMap {
	return make(FuncMap)
}

func (f FuncMap) Bind(name string, fn interface{}) (err error) {
	defer func() {
		if e := recover(); e != nil {
			err = errors.New(name + " is not callable func.")
		}
	}()

	v := reflect.ValueOf(fn)
	if v.Type().Kind() != reflect.Func {
		err = errors.New(name + " is not callable func.")
		return
	}

	//v.Type().NumIn() //panic if it's not a func
	f[name] = v
	return
}

func (f FuncMap) Call(name string, params ...interface{}) (result []reflect.Value, err error) {
	if _, ok := f[name]; !ok {
		err = errors.New(name + " does not bind.")
		return
	}
	if len(params) != f[name].Type().NumIn() {
		err = ErrParamsNotAdapted
		return
	}
	in := make([]reflect.Value, len(params))
	for k, param := range params {
		in[k] = reflect.ValueOf(param)
	}
	result = f[name].Call(in)
	return
}
