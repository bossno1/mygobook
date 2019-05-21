package main
//taoqing 2019.3.26  gorm  ，如何调用存储过程呢？
/*
/SelfService
*/
import (
	"context"
	"net/http"
	"os"
	"runtime/debug"
	"io/ioutil"
	"fmt"
	//"encoding/xml"
	//"yygh001"
	//"strings"
	"net/url"
	"github.com/spf13/viper"
	"encoding/json"
	"config"
	"github.com/bitly/go-simplejson"
	"github.com/lexkong/log"
	"database/sql"
	mssql "github.com/denisenkom/go-mssqldb"
)
//---------------
/*
//	"github.com/jinzhu/gorm"
SelfService 入参：
{  
  
   "clubId":"1",
   "Remark1":"remark",
   "consumeTime":"2018/12/30 02:05:00",
   "infoList": [
	   {
	   "customerId":"0011191847",
	   "packageId":"001001025745",
	   "itemId":"12",
	   "itemBzPrice":"100.00",
	   "itemPrice":"20.00",
	   "servicePersonal":"1",
	   "type":"1",
	   "relaId":"1",
	   "count":"1"
	   }
   ]
}
*/
type SelfService struct {
    ClubId string
	Remark1   string
	ConsumeTime  string
}
type InfoList struct {
    CustomerId string
	Remark1   string
	ItemId  string
	ItemBzPrice string
	ItemPrice string 
	ServicePersonal string
	Type string
	RelaId string
	Count string
}

type InfoListlice struct {
    InfoLists []InfoList
}
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
 
func jsonHandler(w http.ResponseWriter, r *http.Request) {
	var s InfoListlice
	fmt.Println("method:", r.Method)
	body, err := ioutil.ReadAll(r.Body)
    if err != nil {
        fmt.Printf("read body err, %v\n", err)
        return
	}
	js, err := simplejson.NewJson(body);

	//fmt.Println(string(body))
	//解析参数
	err = json.Unmarshal([]byte(body), &s)
    if err != nil {
        fmt.Printf("error: %v", err)
        return
	}
	 
	ClubId := js.Get("clubId").MustString()
	remark1 := js.Get("Remark1").MustString()
	ConsumeTime := js.Get("consumeTime").MustString()
	//fmt.Println(ClubId, remark1, ConsumeTime)
	arr, _ := js.Get("infoList").Array()

	var info interface{}
	for index:=0; index < len(arr) ; index++ {
		info = arr[index]
		//type interface {} does not support indexing， 要用以下方式获取
		fmt.Println(info.(map[string]interface{})["customerId"])

    }
	//fmt.Println(arr[0](map[string]interface{})["customerId"])//["customerId"])
	//fmt.Println(arr)
 
	query := url.Values{}
	query.Add("app name", "MyAppName")
	query.Add("encrypt", "disable")
	query.Add("database", viper.GetString("zh.database"))
	
	u := &url.URL{
		Scheme:   "sqlserver" ,
		User:     url.UserPassword(viper.GetString("zh.userid"), viper.GetString("zh.password")),
		Host:     fmt.Sprintf("%s:%s", viper.GetString("zh.ip"), viper.GetString("zh.sqlport")),
		//Path:  ".\\sql2008" , //instance, // if connecting to an instance instead of a port
		RawQuery: query.Encode(),
	}
	connector, err := mssql.NewConnector(u.String())
	if err != nil {
		log.Infof(u.String(), err.Error())
		return
	}
	db := sql.OpenDB(connector)
	defer db.Close() 
	//err = db.Ping()
	if err != nil {
		log.Infof("无法连接数据库:", err.Error())
		return
	}
	log.Infof("完成连接")
	txn, err := db.Begin() //txn
	if err != nil {
		log.Infof(u.String(), err.Error());
		return
	}
	defer txn.Rollback()
	//----开始处理-----
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	//connector.SessionInitSQL = "set implicit_transactions off"
	var rs_autonumb = "" //返回流水号
	// rows, err1 := txn.QueryContext(ctx, "sp_get_invoinfo",
	// sql.Named("al_item", 1),
	// sql.Named("as_linkcode", ClubId),
	// sql.Named("as_date", ConsumeTime),
	// sql.Named("rs_autonumb", sql.Out{Dest: &rs_autonumb}),
	// )
	//  var strrow string
	// for rows.Next() {
	// 	err = rows.Scan(&strrow)
	// }
	//下面这个方法也可以
	_, err1 := txn.ExecContext(ctx, "sp_get_invoinfo",
		sql.Named("al_item", 1),
		sql.Named("as_linkcode", ClubId),
		sql.Named("as_date", ConsumeTime),
		sql.Named("rs_autonumb", sql.Out{Dest: &rs_autonumb}),
	)	
	if err1 != nil {
		log.Infof(err1.Error());
		return
	}
 	fmt.Println("rs_autonumb is %s", rs_autonumb)
	err = txn.Commit()
	if err != nil {
		log.Infof(err.Error());
		return
	}

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
				log.Infof("WARN: panic fired in %v.panic - %v", fn, e)
				log.Infof(string(debug.Stack()))
			}
		}()
		fn(w, r)
	}
}
 
 
func main() {
	//初始化配置文件
	if err := config.Init(""); err != nil{
		panic(err)
	}
	//定义替换规则
	/*
	rep := strings.NewReplacer("&","&amp;", 
							   "<", "&lt;", 
							   ">", "&gt;", 
							   "'", "&apos;",  
							   "\"", "&quot;") 
	abc := rep.Replace("\"<12345y>");
	fmt.Println(abc)
	*/
	mux := http.NewServeMux()
	mux.HandleFunc("/SelfService", safeHandler(jsonHandler))
	fmt.Println("Port:" + viper.GetString("zh.port"))
	err := http.ListenAndServe(":" +  viper.GetString("zh.port") , mux)
	if err != nil {
		log.Infof("无法监听端口: ", err.Error())
	}
	
}
