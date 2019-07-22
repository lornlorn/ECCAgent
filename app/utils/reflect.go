package utils

import (
    "errors"
    "fmt"
    "reflect"

    "github.com/cihub/seelog"
)

// FunctionMap 定义函数映射类型
type FunctionMap map[string]reflect.Value

// type FunctionMap map[string]interface{}

// FuncMap 声明函数映射
var FuncMap FunctionMap

/*
InitFunctionMap func()
初始化函数映射表
*/
func InitFunctionMap(objs ...interface{}) {

    // var ajaxapi udfuncs.API
    FuncMap = make(FunctionMap, 0)

    for idx, obj := range objs {

        seelog.Debugf("Reflect Type %v : %v", idx, reflect.TypeOf(obj))

        //创建反射变量，注意这里需要传入ruTest变量的地址；
        //不传入地址就只能反射Routers静态定义的方法
        apiObj := reflect.ValueOf(obj)
        apiObjType := apiObj.Type()
        //读取方法数量
        methodNum := apiObj.NumMethod()
        // fmt.Println("NumMethod:", methodNum)

        //遍历路由器的方法，并将其存入控制器映射变量中
        for i := 0; i < methodNum; i++ {
            methodName := apiObjType.Method(i).Name
            FuncMap[methodName] = apiObj.Method(i)
            seelog.Debugf("Reflect Method Index : %v, Name : %v", i, methodName)
        }
    }
    // seelog.Debugf("Reflect Function Map : %v", FuncMap)

}

/*
FuncCall func(mName string, param ...interface{})
*/
func FuncCall(mName string, params ...interface{}) ([]reflect.Value, error) {

    seelog.Debugf("Call Reflect Method Name : %v, Params : %v", mName, params)

    // 判断调用方法是否存在
    if _, ok := FuncMap[mName]; ok {

    } else {
        return nil, fmt.Errorf("Function Name [%v] Not Exists In FuncMap", mName)
    }

    // f := reflect.ValueOf(FuncMap[mName])

    // if len(params) != f.Type().NumIn() {
    //  log.Println("The number of input params not match")
    //  return nil
    // }

    if len(params) != FuncMap[mName].Type().NumIn() {
        return nil, errors.New("The number of input params not match")
    }

    in := make([]reflect.Value, len(params))

    for k, v := range params {
        in[k] = reflect.ValueOf(v)
    }

    //使用方法名字符串调用指定方法
    // return f.Call(in)
    return FuncMap[mName].Call(in), nil

}

// ReflectCall func(m map[string]interface{}, name string, params ...interface{}) ([]reflect.Value, error)
func ReflectCall(m map[string]interface{}, name string, params ...interface{}) ([]reflect.Value, error) {

    f := reflect.ValueOf(m[name])
    if len(params) != f.Type().NumIn() {
        return nil, errors.New("the number of input params not match")
    }

    in := make([]reflect.Value, len(params))
    for k, v := range params {
        in[k] = reflect.ValueOf(v)
    }

    return f.Call(in), nil

}
