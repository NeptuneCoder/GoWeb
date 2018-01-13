/**
根据配置文件，校验接口
1. 根据配置文件`Config.json`配置文件，设置基本信息
	baseUrl
	请求头参数名
	执行次数
	接口配置文件的路径
	接口的执行的顺序定义
2. 执行接口
	接口说明
	解析json文件，读取接口名，参数列表
	返回执行的结果
	将执行结果写入本地
	一个接口花费的时间
	所有接口执行总共花费的时长
	统计花时间最长的接口
3. 根据上述结果，生成`时间.json`文件
使用过程中的问题：
1. post请求，地址后面和body都需要跟参数
2. 接口中需要token和account时，全部在Config.json文件中处理
3. 异常时，记录请求参数

*/
package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"
)

const (
	BaseConfigName = "/Config.json"
)

/**
基本的配置信息
*/
type BaseConfig struct {
	BaseUrl         string   //基本的base url 地址
	ResultFilePath  string   //结果文件的路径
	ResultFileName  string   //结果文件的路径
	Size            int      //TODO 用于执行多少次，暂时不处理
	Path            string   `interface file path` //读取接口参数的路径
	IsCreateLogFile bool     //是否生成新的文件
	Order           []string //执行接口的顺序
	IsLogin         bool     //是否登录
	LoginItf        string   //刷新token的接口
	RefreshToken    string   //刷新token的值
	TokenName       string   //token name
	Headers         map[string]string
}

/**
	接口测试的次数,当写入的数小于等于0时，赋值为1次。既接口执行一次
*/
func (baseConfig *BaseConfig) getRunCount() int {
	if baseConfig.Size <= 0 {
		return 1
	} else {
		return baseConfig.Size
	}
}

/**
接口请求的参数
*/
type ItfParams struct {
	Method string                 //POST,还是GET
	Handle string                 //短的接口名
	Param  map[string]interface{} //参数列表
}

/**
请求接口成功后返回的信息
*/
type BaseResponse struct {
	Code int
	Msg  string
	Data map[string]interface{}
}

var f *os.File

func main() {
	currentPath := GetCurrentDirectory()
	fmt.Println("path = ", currentPath)
	bc, err := ReadBaseConfig(currentPath + BaseConfigName)
	if err != nil {
		log.Fatal("sorroy ,need config,Please contact yh ")
		return
	}
	//创建日志文件
	if bc.IsCreateLogFile {
		f, _ = CreateFile(bc.ResultFilePath, bc.ResultFileName)
		if err != nil {
			log.Fatal("create  log file failed ,will can't record log")
		}
	}
	if err != nil {
		log.Fatal("read base config failed")
		return
	}
	fmt.Println("total size interface = ", len(bc.Order))
	var record []byte
	for _, path := range bc.Order {
		fmt.Println("pa = ", path)
		//获取接口的参数
		itf, err := ReadItfParam(bc.Path + path)

		if err != nil {
			log.Fatal("read ReadItfParam config failed = ", err)
		} else {
			start := time.Now().Nanosecond() //执行前的时间
			re := httpRequest(bc, itf)
			if bc.RefreshToken == path { //刷新token的接口和读取的参数文件名一致时，刷新token
				//当时登录接口时，将token值赋给bc的tokenValue属性上。
				str := fmt.Sprintf("%v", re.Data["token"])
				bc.Headers[bc.TokenName] = string(str)
				fmt.Println("token", string(str))
			}
			end := time.Now().Nanosecond() //执行结束的时间
			record, err = AssembleJson(itf.Handle, ((end - start) / 1e6), re.Code)
		}
	}
	defer func() { //使用匿名函数处理最后必须执行的功能
		//当创建日志文件时，写入日志文件否者输出到终端
		if bc.IsCreateLogFile {
			WriteResultInfo(f, record)
		} else {
			fmt.Println("log record:", string(record))
		}
	}()
}

func Write2Log(bc *BaseConfig, f *os.File, record []byte) {

}

/**
获取当前的相对路径
*/
func GetCurrentDirectory() string {
	dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		log.Fatal(err)
	}
	return strings.Replace(dir, "\\", "/", -1)
}

/**
读取请求的配置信息
*/
func ReadBaseConfig(path string) (*BaseConfig, error) {
	fmt.Println("path = ", path)
	res, err := ioutil.ReadFile(path)
	if err != nil {
		log.Fatal("read base config failed")
		return nil, err
	}
	//定义baseconfig文件
	var bc BaseConfig
	//读取base文件，并且将其转换为对象
	err = json.Unmarshal(res, &bc)
	if err != nil {
		fmt.Println("parse to struct failed!", err)
		return nil, err
	}

	return &bc, nil
}

/**
读取接口请求的参数信息
*/
func ReadItfParam(path string) (*ItfParams, error) {
	fmt.Println("path = = ", path)
	res, err := ioutil.ReadFile(path)
	if err != nil {
		fmt.Println("err = = ", err)
		return nil, err
	}
	var itf ItfParams
	//读取base文件，并且将其转换为对象
	err = json.Unmarshal(res, &itf)
	if err != nil {
		fmt.Println("err  == = = ", err)
		return nil, err
	}
	return &itf, nil
}

/**
请求网络接口
增加对支持POST 和 GET请求的支持
*/
func httpRequest(bc *BaseConfig, param *ItfParams) *BaseResponse {
	client := &http.Client{}
	var req *http.Request
	var err error
	if param.Method == "GET" {
		var pm string
		for k, v := range param.Param {
			pm += (k + "=" + fmt.Sprintf("%v", v)) + "&"
		}
		pm = string([]byte(pm)[0 : len(pm)-1]) //处理拼接在最后的"&"符号，利用切片
		fmt.Println(pm)
		url := bc.BaseUrl + param.Handle + "?" + pm
		req, _ = http.NewRequest(param.Method, url, nil)
	} else if param.Method == "POST" {
		param1, _ := json.Marshal(param.Param)
		req, err = http.NewRequest(param.Method, bc.BaseUrl+param.Handle, bytes.NewBuffer(param1))
	}else{
		req, err = http.NewRequest(param.Method, bc.BaseUrl+param.Handle,nil)
	}

	if err != nil {
		return nil
	}
	//添加头信息
	for k, v := range bc.Headers {
		req.Header.Set(k, v)
	}
	resp, err := client.Do(req)
	if err != nil {
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		// handle error
		return nil
	}
	var br BaseResponse
	json.Unmarshal(body, &br)
	fmt.Println(string(body))
	//返回请求的结果
	return &br

}

/**
组装每个请求的结果，整合为json数据
*/
var cache = make(map[string]interface{})

func AssembleJson(handle string, time, code int) ([]byte, error) {

	v := make(map[string]interface{})
	v["time"] = time
	v["code"] = code
	cache[handle] = v
	re, err := json.Marshal(cache)
	if err != nil {
		return nil, err
	}
	return re, nil

}
func WriteResultInfo(file *os.File, result []byte) {
	file.Write([]byte(result))
}
func CreateFile(logFilePath, fileName string) (*os.File, error) {
	os.MkdirAll(logFilePath,7777)//创建文件夹

	tm := string(time.Now().Format("2006-01-02_15:04:05"))
	fmt.Println("logFilePath = ", logFilePath+tm+"_"+fileName)
	file, err := os.Create(logFilePath + tm + "_" + fileName)
	if err != nil {
		log.Fatal("create log file failed", err)
		return nil, err
	}
	return file, nil

}
