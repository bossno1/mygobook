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
)
import (
	"github.com/jinzhu/gorm"
	_ "github.com/denisenkom/go-mssqldb"
	_ "github.com/jinzhu/gorm/dialects/mssql"
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
	
	db, err := gorm.Open("mssql", "sqlserver://sa:146-164-156-@127.0.0.1:51798?database=his_yb;encrypt=disable;app name=tqtest")
	if err != nil {
	panic("failed to connect database")
	}
	defer db.Close()
	 
	/*
	alter Proc [dbo].[sp_tqtest] 
		@al_operid int
	as
	begin 
		select opername , operid from dictoper --where operid= @al_operid
	end 
	*/
	var result []Result
	db.Raw("[sp_tqtest]  ?", 1).Scan(&result)  //见：http://gorm.io/zh_CN/docs/sql_builder.html ,https://github.com/denisenkom/go-mssqldb
	fmt.Println(result)  //返回数组
	
	

	//返回结果

	w.Header().Set("Content-Type", "application/xml;charset=UTF-8")
	/*
	outv := &program{"1", "中拓香洲测试返回，收到function_id=" + v.Function_id}
	 
    output, err := xml.MarshalIndent(outv, "  ", "    ")
    if err != nil {
        fmt.Printf("error: %v\n", err)
    }
	fmt.Println(v)
	*/
	yh := &yygh001.Program{Session_id:"234234", Function_id:"234234", Akb020:"akb020"}
	//yh.Doctor =  append(yh.Doctor, yygh001.Doctorsourceinfo{Aaz307:"Aaz307"})
	var doctor *yygh001.Doctorsourceinfo
	doctor = &yygh001.Doctorsourceinfo{Aaz307:"赵医生"}
	//sourceinfo只出现一次
	doctor.Source = append(doctor.Source, yygh001.Sourceinfo{Aae030: "上午"});
	doctor.Source = append(doctor.Source, yygh001.Sourceinfo{Aae030: "下午"});
	doctor.Source = append(doctor.Source, yygh001.Sourceinfo{Aae030: "晚上"});
	yh.Doctor =  append(yh.Doctor, *doctor)

	doctor = &yygh001.Doctorsourceinfo{Aaz307:"李医生"}
	//sourceinfo只出现一次
	doctor.Source = append(doctor.Source, yygh001.Sourceinfo{Aae030: "上午"});
	doctor.Source = append(doctor.Source, yygh001.Sourceinfo{Aae030: "下午"});
	doctor.Source = append(doctor.Source, yygh001.Sourceinfo{Aae030: "晚上"});
	yh.Doctor =  append(yh.Doctor, *doctor)


	output2, err2 := xml.MarshalIndent(yh, "  ", "    ")
    if err2 != nil {
        fmt.Printf("error: %v\n", err2)
    }
	w.Write([]byte(output2))
	
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

var schema = `
--可控病种列表（启用标志+有相关诊断）
--select distinct 可控病种ID,可控病种自编码,受控病种 from v_cp_kkbz where ICD10_CODE = 'xxxxx'
create table doctmark_cp_diff3
(
    ID        numeric               identity,
    autonumb  int                   not null,
    ddate     datetime              not null,
    sdate     varchar(8)            not null,
    operid    int                   not null,
    hospid    int                   not null,
    mediid    int                   not null,
    mediname  varchar(255)          not null,
    Reasonid  int                   not null,
    Reason    varchar(255)          not null,
    constraint PK_DOCTMARK_CP_DIFF primary key (ID)
)
;

/* ============================================================ */
/*   Index: doctmark_cp_diff_i1                                 */
/* ============================================================ */
create index doctmark_cp_diff_i3 on doctmark_cp_diff3 (autonumb, sdate)
;
CREATE TABLE person (
    first_name varchar(50),
    last_name  varchar(50),
    email  varchar(50)
);

CREATE TABLE place (
    country varchar(50),
    city varchar(50) NULL,
    telcode integer
)`
type Product struct {
	gorm.Model
	Code string `gorm:"size:20"`
	Price uint
  }
// Scan
type Result struct {
	Opername string
	Operid  int
}

func main() {
	//db, err := sqlx.Connect("sqlserver", "sqlserver://sa:146-164-156-@127.0.0.1:52813?database=master;encrypt=disable;app name=tqtest")
	//db, err := gorm.Open("sqlserver", "sqlserver://sa:146-164-156-@127.0.0.1:51798?database=master;encrypt=disable;app name=tqtest")
	
	mux := http.NewServeMux()
	mux.HandleFunc("/ZTWebService.asmx/FacadeService", safeHandler(xmlHandler))
	
	err := http.ListenAndServe(":8018", mux)
	if err != nil {
		log.Fatal("ListenAndServe: ", err.Error())
	}
}
