package main
/*
tq :   2019.7.16
*/
import (
	"net/http"
	"runtime/debug"
	"fmt"
	"github.com/spf13/viper"
	"config"
	"github.com/lexkong/log"
	"pdaCheckPassword"
	"pdaGetPatinfo"
	"pdaGetDoctmark"
	"pdaUpdateExec"
	"pdaGetTwd"
	"pdaUpdateTwd"
	 
)

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
	mux.HandleFunc("/checkpassword", safeHandler(pdaCheckPassword.JsonHandler))
	mux.HandleFunc("/getpatinfo", safeHandler(pdaGetPatinfo.JsonHandler))
	mux.HandleFunc("/getdoctmark", safeHandler(pdaGetDoctmark.JsonHandler))
	mux.HandleFunc("/updateexec", safeHandler(pdaUpdateExec.JsonHandler))
	mux.HandleFunc("/gettwd", safeHandler(pdaGetTwd.JsonHandler))
	mux.HandleFunc("/updatetwd", safeHandler(pdaUpdateTwd.JsonHandler))
	
	fmt.Println("Port:" + viper.GetString("zh.port"))
	err := http.ListenAndServe(":" +  viper.GetString("zh.port") , mux)
	if err != nil {
		log.Infof("无法监听端口: ", err.Error())
	}
	
}
/*
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
*/