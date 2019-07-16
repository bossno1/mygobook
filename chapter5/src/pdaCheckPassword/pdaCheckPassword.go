package pdaCheckPassword

import (
	"encoding/json"
	"net/http"
	"strconv"
	"fmt"
 	"github.com/spf13/viper"
 	"github.com/lexkong/log"
 	"pub"
)
import (
	"github.com/jinzhu/gorm"
	_ "github.com/denisenkom/go-mssqldb"
	_ "github.com/jinzhu/gorm/dialects/mssql"
  ) 
type  Dictoffice struct {   
	//注意，必须首字母大写表示要输出
	Officeid string   `json:"officeid"`
	Officename string   `json:"officename"`
}
type Dictoper struct {
	//注意，必须首字母大写表示要输出   0正确， 非0错
	Code    int  `json:"code"`
	Message  string `json:"message"`
	Operid	string `json:"operid"`
	Opername  string  `json:"opername"`
	OfficeList []Dictoffice   `json:"officelist"`
}
// Scan
type Result struct {
	Operid int
	Opername  string
	Stage  int
}
func JsonHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("method:", r.Method)
	query := r.URL.Query()
	opercode := query["opercode"]
	//password := query["password"]
	db, err := gorm.Open("mssql", "sqlserver://" + viper.GetString("his.userid") + ":" + viper.GetString("his.password") + "@" + viper.GetString("his.ip") + "?database=" +  viper.GetString("his.database") + ";encrypt=disable;app name=pda_server")
	if err != nil {
		log.Infof("无法连接数据库", err.Error())
		pub.ReturnJSON(-1, "无法连接数据库", w)
		return 
	}
	log.Infof("完成连接")
	defer db.Close()
	//-----取操作员
	ls_sql1 := ` select b.operid,  	b.opername , b.stage 
			from  dictoper b where b.opercode = ? `
	var result Result
	db.Raw(ls_sql1, opercode).Scan(&result)
	if result.Opername == ""{
		pub.ReturnJSON(-1, "用户不存在", w)
		log.Infof("用户不存在")
		return
	}
	if result.Stage != 1 {
		pub.ReturnJSON(-1, "用户状态未审核", w)
		log.Infof("用户状态未审核")
		return
	}
	//-----取操作员相关科室------
	ls_sql :=` 
		select  a.officeName, b.opername ,b.operid,b.officeid  
				from DictOffice a,dictoper b where 
					a.officekind = 1and b.operid = ?
					and a.officeid = b.officeid	and b.stage = 1 
		UNION 
		select a.officeName, b.opername ,c.operid,c.officeid 
		from DictOffice a,dictoper b ,dict_office_oper c
			where a.officekind = 1
			and b.operid = ? and a.officeid = c.officeid 
			and b.operid = c.operid	and b.stage = 1 `
	//var result []Dictoffice   //医生信息
	var officename string
	var opername string
	var officeid int
	var operid int
	rows, err := db.Raw(ls_sql , result.Operid, result.Operid,).Rows() 
	//见：http://gorm.io/zh_CN/docs/sql_builder.html ,https://github.com/denisenkom/go-mssqldb
	if err != nil {
		log.Fatal("取操作员相关科室出错", err) 
		pub.ReturnJSON(-1,  "取操作员相关科室出错" , w)
		return  
	}
	var dictoper Dictoper
	 
	for rows.Next() {
		rows.Scan(&officename, &opername, &operid, &officeid)
		dictoper.Operid = strconv.Itoa(operid)
		dictoper.Opername = opername
		dictoper.OfficeList = append(dictoper.OfficeList,
			 Dictoffice{ Officeid:strconv.Itoa(officeid), Officename: officename})
	}
	rows.Close()
	b, err := json.Marshal(dictoper)
    if err != nil {
        fmt.Println("json err:", err)
    }
    fmt.Println(string(b));
 	//返回结果
	w.Header().Set("Content-Type", "application/json;charset=UTF-8")
	w.Write(b)
	log.Infof("成功，完成操作员的登陆") 
	return
}
 
 