package main
//taoqing 2019.3.26  gorm  ，如何调用存储过程呢？
import (
	"log"
	"net/http"
	"os"
	"runtime/debug"
	"io/ioutil"
	"fmt"
	"encoding/xml"
	"yygh001"
	"session"
	"strings"
)

//---------------
const (
	ListDir      = 0x0001
	UPLOAD_DIR   = "./uploads"
	TEMPLATE_DIR = "./views"
)
 
func check(err error) {
	if err != nil {
		panic(err)
	}
}
 
func isExists(path string) bool {
	_, err := os.Stat(path)
	if err == nil {
		return true
	}
	return os.IsExist(err)
}
type parmXml struct {
	Program    xml.Name `xml:"program"`
    Function_id  string `xml:"function_id"`
	Akb020     string   `xml:"akb020"`
	Bka601   string    `xml:"bka601"`
	Bka017   string    `xml:"bka017"`
  
}
type program struct {
	Return_code string   `xml:"return_code"`
    Return_code_message  string   `xml:"return_code_message"`
}
func xmlHandler(w http.ResponseWriter, r *http.Request) {
	v := parmXml{}
	fmt.Println("method:", r.Method)
	body, err := ioutil.ReadAll(r.Body)
    if err != nil {
        fmt.Printf("read body err, %v\n", err)
        return
	}
	//解析参数
    err = xml.Unmarshal(body, &v)
    if err != nil {
        fmt.Printf("error: %v", err)
        return
	}
	/*bka601指令类型
	1	上传号源
	2	上传就诊和费用信息
	3	上传诊疗结果
	4	上传取药信息
	5	取消费用信息
	*/
	if v.Bka601 == "1"  {
		fmt.Println("收到入参: %v", v)
		fmt.Println("s_session: %v", s_session)

		date := v.Bka017
		date = string([]byte(date)[0:4]) + "-" +  string([]byte(date)[4:6]) + "-" + string([]byte(date)[6:8])
		fmt.Println("日期条件: %v" , date)
		output2 := yygh001.Process(s_session, date)  //bka017 = 日期yyyymmdd
		//返回结果
		w.Header().Set("Content-Type", "application/xml;charset=UTF-8")
		w.Write([]byte(output2))
	}
	

	/*
	outv := &program{"1", "中拓香洲测试返回，收到function_id=" + v.Function_id}
    output, err := xml.MarshalIndent(outv, "  ", "    ")
    if err != nil {
        fmt.Printf("error: %v\n", err)
    }
	fmt.Println(v)
	*/

	//fmt.Println(body)
	//http.ServeFile(w, r, imagePath)
}
 
 
func safeHandler(fn http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if e, ok := recover().(error); ok {
				http.Error(w, e.Error(), http.StatusInternalServerError)

				// 或者输出自定义的 50x 错误页面
				// w.WriteHeader(http.StatusInternalServerError)
				// renderHtml(w, "error", e.Error())

				// logging
				log.Println("WARN: panic fired in %v.panic - %v", fn, e)
				log.Println(string(debug.Stack()))
			}
		}()
		fn(w, r)
	}
}
 
// Scan
type Result struct {
	Opername string
	Operid  int
}
var s_session string
func main() {
	//定义替换规则
	rep := strings.NewReplacer("&","&amp;", 
							   "<", "&lt;", 
							   ">", "&gt;", 
							   "'", "&apos;",  /*单引号*/
							   "\"", "&quot;")  //定义替换规则
	abc := rep.Replace("\"<12345y>");
	fmt.Println(abc)
	
	s_session = session.GetSession()
	
	mux := http.NewServeMux()
	mux.HandleFunc("/ZTWebService.asmx/FacadeService", safeHandler(xmlHandler))
	
	err := http.ListenAndServe(":8018", mux)
	if err != nil {
		log.Fatal("ListenAndServe: ", err.Error())
	}
	
}
