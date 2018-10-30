package jgin

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"reflect"
	"strconv"
	"time"
	_ "github.com/go-sql-driver/mysql"

	"html/template"
	"errors"
	
)

var restFuncMap template.FuncMap = make(template.FuncMap)

func init(){
	restFuncMap["ctxpath"]=ctxpath
	restFuncMap["pageurl"]=pageurl
	restFuncMap["apiurl"]=apiurl
	restFuncMap["version"]=version
	restFuncMap["hello"]=hello
	restFuncMap["asset"]=asset
}

func asset() string{
	cfg := GetCfg()

	return  cfg.App["asset"]


}

func GetFuncMap()(template.FuncMap){
	return restFuncMap
}


func hello(d string) string{
	return "hello "+d
}


func ctxpath() string{
	cfg := GetCfg()
	url := cfg.App["protocal"]+"://" + cfg.App["domain"]
	if (cfg.App["port"]!="80"){
		url += (":"+cfg.App["port"])
	}
	return url
}

func pageurl(uri string) string{
	cfg := GetCfg()
	url := cfg.App["protocal"]+"://" + cfg.App["domain"]
	if (cfg.App["port"]!="80"){
		url += (":"+cfg.App["port"])
	}
	return url+"/"+uri +".shtml"
}

func apiurl(uri string) string{
	return uri
}


func version() string{
	cfg := GetCfg()
	if len(cfg.App["version"])==0{
		return  strconv.FormatInt(time.Now().Unix(),10)
	}else{
		return  cfg.App["version"]
	}

}


func getMd5(s string) string {
	signByte := []byte(s)
	hash := md5.New()
	hash.Write(signByte)
	return hex.EncodeToString(hash.Sum(nil))
}


//用map填充结构
func FillStruct(data map[string]interface{}, obj interface{}) error {
	for k, v := range data {
		err := SetField(obj, k, v)
		if err != nil {
			return err
		}
	}
	return nil
}

//用map的值替换结构的值
func SetField(obj interface{}, name string, value interface{}) error {
	structValue := reflect.ValueOf(obj).Elem()        //结构体属性值
	structFieldValue := structValue.FieldByName(name) //结构体单个属性值

	if !structFieldValue.IsValid() {
		return fmt.Errorf("No such field: %s in obj", name)
	}

	if !structFieldValue.CanSet() {
		return fmt.Errorf("Cannot set %s field value", name)
	}

	structFieldType := structFieldValue.Type() //结构体的类型
	val := reflect.ValueOf(value)              //map值的反射值

	var err error
	if structFieldType != val.Type() {
		val, err = TypeConversion(fmt.Sprintf("%v", value), structFieldValue.Type().Name()) //类型转换
		if err != nil {
			return err
		}
	}

	structFieldValue.Set(val)
	return nil
}

//类型转换
func TypeConversion(value string, ntype string) (reflect.Value, error) {
	if ntype == "string" {
		return reflect.ValueOf(value), nil
	} else if ntype == "time.Time" {
		t, err := time.ParseInLocation("2006-01-02 15:04:05", value, time.Local)
		return reflect.ValueOf(t), err
	} else if ntype == "Time" {
		t, err := time.ParseInLocation("2006-01-02 15:04:05", value, time.Local)
		return reflect.ValueOf(t), err
	} else if ntype == "int" {
		i, err := strconv.Atoi(value)
		return reflect.ValueOf(i), err
	} else if ntype == "int8" {
		i, err := strconv.ParseInt(value, 10, 64)
		return reflect.ValueOf(int8(i)), err
	} else if ntype == "int32" {
		i, err := strconv.ParseInt(value, 10, 64)
		return reflect.ValueOf(int64(i)), err
	} else if ntype == "int64" {
		i, err := strconv.ParseInt(value, 10, 64)
		return reflect.ValueOf(i), err
	} else if ntype == "float32" {
		i, err := strconv.ParseFloat(value, 64)
		return reflect.ValueOf(float32(i)), err
	} else if ntype == "float64" {
		i, err := strconv.ParseFloat(value, 64)
		return reflect.ValueOf(i), err
	}

	//else if .......增加其他一些类型的转换

	return reflect.ValueOf(value), errors.New("未知的类型：" + ntype)
}

//结构体转为map
func Struct2Map(obj interface{}) map[string]interface{} {
	t := reflect.TypeOf(obj)
	v := reflect.ValueOf(obj)

	var data= make(map[string]interface{})
	for i := 0; i < t.NumField(); i++ {
		data[t.Field(i).Name] = v.Field(i).Interface()
	}
	return data
}