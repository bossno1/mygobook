package pdaUpdateExec

import (
	"context"
	"net/http"
	"io/ioutil"
	"fmt"
	"net/url"
	"github.com/spf13/viper"
	"github.com/bitly/go-simplejson"
	"github.com/lexkong/log"
	"database/sql"
	mssql "github.com/denisenkom/go-mssqldb"
	//"strconv"
	"pub"
)

func JsonHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("method:", r.Method)
	body, err := ioutil.ReadAll(r.Body)
    if err != nil {
		fmt.Printf("read body err, %v\n", err)
		pub.ReturnJSON(-1, err.Error(), w)
		
        return
	}
	js, err := simplejson.NewJson(body);
	log.Infof(string(body))
	//fmt.Println(string(body))
	//解析参数
	 
	// ClubId := js.Get("clubId").MustString()
	// remark1 := js.Get("Remark1").MustString()
	// ConsumeTime := js.Get("consumeTime").MustString()
	//fmt.Println(ClubId, remark1, ConsumeTime)
	arr, _ := js.Array() //arr

	//设置数据库链接参数
	query := url.Values{}
	query.Add("app name", "pdaserver")
	query.Add("encrypt", "disable")
	query.Add("database", viper.GetString("his.database"))
	u := &url.URL{
		Scheme:   "sqlserver" ,
		User:     url.UserPassword(viper.GetString("his.userid"), viper.GetString("his.password")),
		Host:     fmt.Sprintf("%s:%s", viper.GetString("his.ip"), viper.GetString("his.sqlport")),
		//Path:  ".\\sql2008" , //instance, // if connecting to an instance instead of a port
		RawQuery: query.Encode(),
	}
	connector, err := mssql.NewConnector(u.String())
	if err != nil {
		log.Infof(u.String(), err.Error())
		pub.ReturnJSON(-1, err.Error(), w)
		return
	}
	db := sql.OpenDB(connector)
	defer db.Close() 
	//启动事务
	txn, err := db.Begin() 
	if err != nil {
		log.Infof(u.String(), err.Error());
		pub.ReturnJSON(-1, err.Error(), w)
		return
	}
	defer txn.Rollback()
	//----开始处理-----
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	//connector.SessionInitSQL = "set implicit_transactions off"
	 
	
	var info interface{}
	var rs_errormsg = ""
	 
	for index:=0; index < len(arr) ; index++ {
		 info = arr[index]
		 //type interface {} does not support indexing， 要用以下方式获取
		//  count, _  = strconv.ParseFloat(info.(map[string]interface{})["count"].(string), 64) 
		//  itemprice, _ = strconv.ParseFloat(info.(map[string]interface{})["itemPrice"].(string), 64)
		//  itemprice = zh_func1.Decimal(itemprice ,4) 
		//  price = zh_func1.Decimal(count * itemprice,2) 
		//  totalprice = totalprice + price;
		//  _, err1 := txn.ExecContext(ctx, "update hosp_detail_exec set exectime2 = getdate(), Checkoperid = @p1 where execid = @p2",
		// 	 info.(map[string]interface{})["checkoperid"],  info.(map[string]interface{})["execid"],
		// 	 )
		_, err1 := txn.ExecContext(ctx,  `insert into  Doctmark_exec (autonumb, execdate,execoperid, execopername, exectime, remark, execcount) 
				values (@p1, convert(char(10), getdate(), 120),  @p2,  @p3, getdate(), '', 1 )`	,
				info.(map[string]interface{})["autonumb"],  info.(map[string]interface{})["checkoperid"],
				info.(map[string]interface{})["checkopername"])
		if err1 != nil {
			log.Infof(err1.Error());
			pub.ReturnJSON(-1, err1.Error(), w)
			return
		}
		_, err2 := txn.ExecContext(ctx, ` update  Doctmark_exec set execcount = x2.sumcount from 
					(select autonumb, execdate , count(*) as sumcount from Doctmark_exec	where 
							autonumb = @p1 and execdate = convert(char(10), getdate(), 120) 
							group by autonumb, execdate ) 
					x2 where Doctmark_exec.autonumb = x2.autonumb  and  Doctmark_exec.execdate = x2.execdate`,
				info.(map[string]interface{})["autonumb"])
		if err2 != nil {
			log.Infof(err2.Error());
			pub.ReturnJSON(-1, err2.Error(), w)
			return
		}
		
		if len(rs_errormsg) > 0 {
			log.Infof(rs_errormsg);
			pub.ReturnJSON(-1, rs_errormsg, w)
			return
		}
		  
	}
	
	err = txn.Commit()
	if err != nil {
		log.Infof(err.Error());
		pub.ReturnJSON(-1, "保存数据出错", w)
	 
		return
	}
	log.Infof("成功更新执行状态") 
	pub.ReturnJSON(0, "", w)
	return


}
 
 