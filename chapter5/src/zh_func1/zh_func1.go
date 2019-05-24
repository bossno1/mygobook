package zh_func1
/*
tq : SelfService 2019.5.22
*/
import (
	"net/http"
	"fmt"
	"encoding/json"
	"strconv"
)
type returnJson struct {
	// ID 不会导出到JSON中
	//ID int `json:"-"`
	//这样表示会进行二次JSON编码  	Message string `json:"message,string"`
	Result  string `json:"status"`
	Message string `json:"message"`
	// 如果 ServerIP 为空，则不输出到JSON串中
	//ServerIP   string `json:"serverIP,omitempty"`
}

func ReturnJSON(code string , errTxt string, w http.ResponseWriter) []byte{
	s := returnJson {
		Result:  code,  //0失败 1成功 
		Message: errTxt, 
	}
	b, _ := json.Marshal(s)
	//返回结果
	w.Header().Set("Content-Type", "application/json;charset=UTF-8")
	w.Write(b)
	return b
}
func check(err error) {
	if err != nil {
		panic(err)
	}
}

func Decimal(value float64, dec int) float64 {
	value, _ = strconv.ParseFloat(fmt.Sprintf("%." +  strconv.Itoa(dec) + "f", value), 64)
	return value
}
