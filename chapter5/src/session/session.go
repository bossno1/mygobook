package session
//taoqing 2019.3.26  gorm  ，如何调用存储过程呢？
import (
 
	"net/http"
	"strings"
	"fmt"
	"io/ioutil"
	"config"
	"github.com/spf13/viper"
	"encoding/xml"
)

type program struct {
	XMLName  xml.Name `xml:"program"`
	Return_code  string  `xml:"return_code"`
	Session_id   string  `xml:"session_id"`
	Notify_msg   string  `xml:"notify_msg"`
	Return_code_message string  `xml:"return_code_message"`
}

func GetSession() string{
	//初始化配置文件
	if err := config.Init(""); err != nil{
		panic(err)
	}
	var parm string
	parm = `<program>
			<function_id>` + viper.GetString("function_id") + `</function_id>
			<userid>` + viper.GetString("userid") + `</userid>
			<password>` + viper.GetString("password") + `</password>
			</program>
			`
	fmt.Printf("获取session:%v\n",viper.GetString("ZHSocialURL"))
	request, _ := http.NewRequest("POST", viper.GetString("ZHSocialURL"), strings.NewReader(parm))
    //post数据并接收http响应
    resp,err :=http.DefaultClient.Do(request)
    if err!=nil{
		fmt.Printf("错误，无法获取Session:%v\n",err)
		return ""
    }else {
      
		respBody,_ := ioutil.ReadAll(resp.Body)
	
		v := program{}
		err = xml.Unmarshal(respBody, &v)
		if err != nil {
			fmt.Printf("无法解析Session返回的内容: %v", err)
			return ""
		}
		//fmt.Printf("return:%v\n", string(respBody))
		if len(v.Notify_msg) > 0 {
			fmt.Printf("通知信息:%v\n",v.Notify_msg)
		}
		fmt.Printf("Session:%v\n",v.Session_id)
		return v.Session_id
	}
	

	 
}