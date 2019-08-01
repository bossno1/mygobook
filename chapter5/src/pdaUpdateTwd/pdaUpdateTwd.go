package pdaUpdateTwd

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
// Scan
type Result struct {
	count int 
}
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
	fmt.Println(string(body))
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
	
	
	for index:=0; index < len(arr) ; index++ {
		 info = arr[index]
		 hospid := info.(map[string]interface{})["hospid"]
		 twddate := info.(map[string]interface{})["twddate"]
		 time_point := info.(map[string]interface{})["time_point"]
		 tw_status := info.(map[string]interface{})["tw_status"]
		 
		 //辨断记录是否有，决定insert还是update 
		var count int
	 
		ls_sql1 := `select count(*) from  pat_twd_master where zyh = @p1 and  recording_date = @p2`
		rows, err := db.QueryContext(ctx, ls_sql1, hospid, twddate)
		if err != nil {
			log.Infof(err.Error())
			pub.ReturnJSON(-1, err.Error(), w)
			return
		}
		rows.Next()
		err = rows.Scan(&count)
		if err != nil {
			log.Infof(err.Error())
			pub.ReturnJSON(-1, err.Error(), w)
			return
		}
		var err1 error
		 if count == 0 {
			ls_sql1 := `insert into pat_twd_master (zyh,recording_date, xy1, xy2, tz, ps, sryl,
				db, nl, col1, col1_value, col2, col2_value, col3, col3_value, col4, col4_value, col5, col5_value	) 
				values (@p1, @p2 , @p3 , @p4 , @p5, @p6, @p7, 
					@p8,@p9,@p10,@p11,@p12,@p13,@p14,@p15,@p16,@p17,@p18,@p19) `

				 fmt.Println(ls_sql1);
			_, err1   = txn.ExecContext(ctx,  ls_sql1,
				info.(map[string]interface{})["hospid"],  info.(map[string]interface{})["twddate"],
				info.(map[string]interface{})["mxy1"], info.(map[string]interface{})["mxy2"], 
				info.(map[string]interface{})["mtz"], info.(map[string]interface{})["mps"], 
				info.(map[string]interface{})["msryl"], info.(map[string]interface{})["mdb"], 
				info.(map[string]interface{})["mnl"],  
				info.(map[string]interface{})["mcol1"], info.(map[string]interface{})["mcol1value"], 
				info.(map[string]interface{})["mcol2"], info.(map[string]interface{})["mcol2value"], 
				info.(map[string]interface{})["mcol3"], info.(map[string]interface{})["mcol3value"], 
				info.(map[string]interface{})["mcol4"], info.(map[string]interface{})["mcol4value"], 
				info.(map[string]interface{})["mcol5"], info.(map[string]interface{})["mcol5value"])				 
			 
		 }else{
			ls_sql1 := `update  pat_twd_master  set  xy1 = @p1, xy2 = @p2, tz = @p3, ps = @p4, sryl = @p5,
							db = @p6, nl = @p7, col1 = @p8, col1_value = @p9, col2 = @p10, col2_value = @p11, 
							col3 = @p12, col3_value = @p13, col4 = @p14, col4_value = @p15, col5 = @p16, col5_value = @p17 
						where 	zyh = @p18 and 	recording_date = @p19 `
						fmt.Println(ls_sql1);
			_, err1  = txn.ExecContext(ctx,  ls_sql1,
				info.(map[string]interface{})["mxy1"], info.(map[string]interface{})["mxy2"], 
				info.(map[string]interface{})["mtz"], info.(map[string]interface{})["mps"], 
				info.(map[string]interface{})["msryl"], info.(map[string]interface{})["mdb"], 
				info.(map[string]interface{})["mnl"],  
				info.(map[string]interface{})["mcol1"], info.(map[string]interface{})["mcol1value"], 
				info.(map[string]interface{})["mcol2"], info.(map[string]interface{})["mcol2value"], 
				info.(map[string]interface{})["mcol3"], info.(map[string]interface{})["mcol3value"], 
				info.(map[string]interface{})["mcol4"], info.(map[string]interface{})["mcol4value"], 
				info.(map[string]interface{})["mcol5"], info.(map[string]interface{})["mcol5value"],
				info.(map[string]interface{})["hospid"],  info.(map[string]interface{})["twddate"])				 
		 }

		if err1 != nil {
			log.Infof(err1.Error());
			pub.ReturnJSON(-1, err1.Error(), w)
			return
		}
		
		 //type interface {} does not support indexing， 要用以下方式获取
		//  count, _  = strconv.ParseFloat(info.(map[string]interface{})["count"].(string), 64) 
		//  itemprice, _ = strconv.ParseFloat(info.(map[string]interface{})["itemPrice"].(string), 64)
		//  itemprice = zh_func1.Decimal(itemprice ,4) 
		//  price = zh_func1.Decimal(count * itemprice,2) 
		//  totalprice = totalprice + price;
		//  _, err1 := txn.ExecContext(ctx, "update hosp_detail_exec set exectime2 = getdate(), Checkoperid = @p1 where execid = @p2",
		// 	 info.(map[string]interface{})["checkoperid"],  info.(map[string]interface{})["execid"],
		// 	 )
	
		//-------------现在写三测单明细表---------------------------------
		ls_sql1 = `select count(*) from  pat_twd_item where zyh = @p1 and  recording_date = @p2 and time_point = @p3 
					and tw_status = @p4`
		rows, err = db.QueryContext(ctx, ls_sql1, hospid, twddate, time_point, tw_status)
		if err != nil {
			log.Infof(err.Error())
			pub.ReturnJSON(-1, err.Error(), w)
			return
		}
		rows.Next()
		err = rows.Scan(&count)
		if err != nil {
			log.Infof(err.Error())
			pub.ReturnJSON(-1, err.Error(), w)
			return
		}
		
		 if count == 0 {
			ls_sql1 := `insert into pat_twd_item (zyh,recording_date, time_point, tw_status, t1, t2, t3, mb,
				xl, hx, m1, m2, xy, xy2, operid, opername 	) 
				values (@p1, @p2 , @p3 , @p4 , @p5, @p6, @p7, 
					@p8,@p9,@p10,@p11,@p12,@p13,@p14,@p15,@p16) `

		 	fmt.Println(ls_sql1);
			_, err1   = txn.ExecContext(ctx,  ls_sql1,
				info.(map[string]interface{})["hospid"],  info.(map[string]interface{})["twddate"],
				info.(map[string]interface{})["time_point"], info.(map[string]interface{})["tw_status"], 
				info.(map[string]interface{})["t1"], info.(map[string]interface{})["t2"], 
				info.(map[string]interface{})["t3"], info.(map[string]interface{})["mb"], 
				info.(map[string]interface{})["xl"],  
				info.(map[string]interface{})["hx"], info.(map[string]interface{})["m1"], 
				info.(map[string]interface{})["m2"], info.(map[string]interface{})["xy"], 
				info.(map[string]interface{})["xy2"], info.(map[string]interface{})["operid"], 
				info.(map[string]interface{})["opername"] )				 
			 
		 }else{
			ls_sql1 := `update  pat_twd_item  set  t1 = @p1, t2 = @p2, t3 = @p3, mb = @p4, xl = @p5,
						hx = @p6, m1 = @p7, m2 = @p8, xy = @p9, xy2 = @p10, operid = @p11, opername = @p12 
						where 	zyh = @p13 and 	recording_date = @p14 and time_point = @p15 and tw_status = @p16`
			fmt.Println(ls_sql1);
			_, err1  = txn.ExecContext(ctx,  ls_sql1,
				info.(map[string]interface{})["t1"], info.(map[string]interface{})["t2"], 
				info.(map[string]interface{})["t3"], info.(map[string]interface{})["mb"], 
				info.(map[string]interface{})["xl"],  
				info.(map[string]interface{})["hx"], info.(map[string]interface{})["m1"], 
				info.(map[string]interface{})["m2"], info.(map[string]interface{})["xy"], 
				info.(map[string]interface{})["xy2"], info.(map[string]interface{})["operid"], 
				info.(map[string]interface{})["opername"] ,
				info.(map[string]interface{})["hospid"],  info.(map[string]interface{})["twddate"],
				info.(map[string]interface{})["time_point"], info.(map[string]interface{})["tw_status"])				 
		 }
	
		if err1 != nil {
			log.Infof(err1.Error());
			pub.ReturnJSON(-1, err1.Error(), w)
			return
		}
	}
	
	err = txn.Commit()
	if err != nil {
		log.Infof(err.Error());
		pub.ReturnJSON(-1, "保存数据出错", w)
	 
		return
	}
	log.Infof("成功更新三测单数据") 
	pub.ReturnJSON(0, "", w)
	return


}
 
 