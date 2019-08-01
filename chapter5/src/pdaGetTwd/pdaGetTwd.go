package pdaGetTwd

import (
	"encoding/json"
	"net/http"
//	"strconv"
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
type  TwdDetailData struct {   
	//注意，必须首字母大写表示要输出
	Hospid string   `json:"hospid"`
	Twddate string   `json:"twddate"` //三测单日期(yyyy-mm-dd格式)
	Timepoint string   `json:"timepoint"`
	Twstatus string   `json:"twstatus"`
	T1 string   `json:"t1"`  //口表t1
	T2   string   `json:"t2"`   //腋表t2
	T3 string   `json:"t3"`   //肛表t3
	Mb string   `json:"mb"`  //脉搏mb
	Xl string   `json:"xl"`  // ??
	Hx string   `json:"hx"`   //呼吸hx
	M1 string   `json:"m1"`  //主要事件m1
	M2 string   `json:"m2"`  //辅助事件m2
	Xy string   `json:"xy"`  //血压
	Xy2 string   `json:"xy2"`  //血氧
	Operid string   `json:"operid"`  // 
	Opername string   `json:"opername"`
	 
}
type TwdData struct {
	Code    int  `json:"code"`
	Message  string `json:"message"`
	//注意，必须首字母大写表示要输出   0正确， 非0错
	Hospid    string  `json:"hospid"`
	Hospcode    string  `json:"hospcode"`
	Name	string `json:"name"`
	Sex  string  `json:"sex"`
	Twddate  string `json:"twddate"`
	//Age	string `json:"age"`
	Xy1  string  `json:"xy1"` //血压(mmHg)xy1
	Xy2   string `json:"xy2"`//血压(mmHg)xy2
	Tz  string  `json:"tz"`  //体重(Kg)tz
	Ps   string `json:"ps"`   //皮试 ps
	Sryl   string `json:"sryl"`//入量(ml)sryl
	Db   string `json:"db"`   //出量(ml)db
	Nl   string `json:"nl"`  //大便(次/日)nl
	Col1   string `json:"col1"`
	Col1value   string `json:"col1value"`
	Col2   string `json:"col2"`
	Col2value   string `json:"col2value"`
	Col3   string `json:"col3"`
	Col3value   string `json:"col3value"`
	Col4   string `json:"col4"`
	Col4value   string `json:"col4value"`
	Col5   string `json:"col5"`
	Col5value   string `json:"col5value"`
	 

	TwdDetailDatas []TwdDetailData   `json:"data"`
}
// Scan
type Result struct {
	Hospid int
	Hospcode  string
	Name  string
	Sex   string
	Bed   string
	Ddate  string
}
func JsonHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("method:", r.Method)
	query := r.URL.Query()
	hospid := query["hospid"][0]
	twddate := query["twddate"][0]  //yyyy-mm-dd 格式
	 
	db, err := gorm.Open("mssql", "sqlserver://" + viper.GetString("his.userid") + ":" + viper.GetString("his.password") + "@" + viper.GetString("his.ip") + ":" + viper.GetString("his.sqlport") + "?database=" +  viper.GetString("his.database") + ";encrypt=disable;app name=pda_server")
	if err != nil {
		log.Infof("无法连接数据库", err.Error())
		pub.ReturnJSON(-1, "无法连接数据库", w)
		return 
	}
	log.Infof("完成连接")
	defer db.Close()
	//-----取患者基本资料+三测单master表
	ls_sql1 := ` select  patinfo.hospid , patinfo.hospcode, patinfo.name, 
				CASE patinfo.sex WHEN 1 THEN '男' when 2 THEN '女' ELSE '未登记' end sex,
				twddate = convert(varchar(10),pat_twd_master.recording_date, 120),   
				pat_twd_master.xy1,   
				pat_twd_master.xy2,   
				pat_twd_master.tz,   
				pat_twd_master.ps,   
				pat_twd_master.sryl,   
				pat_twd_master.db,   
				pat_twd_master.nl,   
				pat_twd_master.col1,   
				pat_twd_master.col1_value,   
				pat_twd_master.col2,   
				pat_twd_master.col2_value,   
				pat_twd_master.col3,   
				pat_twd_master.col3_value,   
				pat_twd_master.col4,   
				pat_twd_master.col4_value,   
				pat_twd_master.col5,   
				pat_twd_master.col5_value 
			from  patinfo			
				 left join pat_twd_master   on  patinfo.hospid = pat_twd_master.zyh     
				     and pat_twd_master.recording_date = ?
		    where  patinfo.hospid = ? `
			
	var twdData TwdData
    db.Raw(ls_sql1, twddate,  hospid).Scan(&twdData)
	if twdData.Name == ""{
		pub.ReturnJSON(-1, "患者不存在", w)
		log.Infof("患者不存在")
		return
	}
	//见：http://gorm.io/zh_CN/docs/sql_builder.html ,https://github.com/denisenkom/go-mssqldb

	//-----三测单明细内容------

	ls_sql :=` select   pat_twd_item.zyh,   
	Twddate = convert(varchar(10), pat_twd_item.recording_date, 120),   
	Timepoint = pat_twd_item.time_point,   
	Twstatus  = pat_twd_item.tw_status  ,
	isnull(pat_twd_item.t1, ''),   
	isnull(pat_twd_item.t2, ''),   
	isnull(pat_twd_item.t3, ''),   
	isnull(pat_twd_item.mb, ''),   
	isnull(pat_twd_item.xl, ''),   
	isnull(pat_twd_item.hx, ''),   
	m1 = isnull(pat_twd_item.m1, ''),   
	m2 = isnull(pat_twd_item.m2, ''),
	isnull(pat_twd_item.xy, ''),
	isnull(pat_twd_item.xy2, ''),

	isnull(pat_twd_item.operid, '0'),
	isnull(pat_twd_item.opername, '')
   		
	   FROM pat_twd_item  
	   WHERE pat_twd_item.zyh = ?   and pat_twd_item.recording_date = ? 
	   ORDER BY  recording_date, time_point, tw_status  `			
		 
	rows, err := db.Raw(ls_sql , hospid, twddate).Rows() 
	
	//见：http://gorm.io/zh_CN/docs/sql_builder.html ,https://github.com/denisenkom/go-mssqldb
	if err != nil {
		log.Fatal("取医嘱出错", err) 
		pub.ReturnJSON(-1,  "取医嘱出错" , w)
		return  
	}
	
	
	for rows.Next() {
		var twdDetailData TwdDetailData
		///db.ScanRows(rows, &doctmark) 
		error := rows.Scan(&twdDetailData.Hospid, &twdDetailData.Twddate, &twdDetailData.Timepoint, 
			&twdDetailData.Twstatus, &twdDetailData.T1,
			&twdDetailData.T2,&twdDetailData.T3, &twdDetailData.Mb, &twdDetailData.Xl, &twdDetailData.Hx, 
			&twdDetailData.M1, &twdDetailData.M2,
			&twdDetailData.Xy, &twdDetailData.Xy2,
			&twdDetailData.Operid, &twdDetailData.Opername ) 
	 
		if error != nil {
			fmt.Println(error);
		} 
		twdData.TwdDetailDatas = append(twdData.TwdDetailDatas,
			twdDetailData)
	}
	rows.Close()
	b, err := json.Marshal(twdData)
    if err != nil {
        fmt.Println("json err:", err)
    }
   // fmt.Println(string(b));
 	//返回结果
	w.Header().Set("Content-Type", "application/json;charset=UTF-8")
	w.Write(b)
	log.Infof("成功，完成三测单获取") 
	
	return
}
 
 