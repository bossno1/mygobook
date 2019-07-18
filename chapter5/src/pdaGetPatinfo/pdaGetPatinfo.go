package pdaGetPatinfo
 
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

type Patinfo struct {
	//注意，必须首字母大写表示要输出   0正确， 非0错
	
	Hospid	string `json:"hospid"`
	Hospcode  string  `json:"hospcode"`
	Bed	string `json:"bed"`
	Name	string `json:"name"`
	Sex  string  `json:"sex"`
	Age	string `json:"age"`  //xx岁或xx天
	Officeid  string  `json:"officeid"`
	Officename  string  `json:"officename"`
	Doctorname  string  `json:"doctorname"`
	Hospdate  string  `json:"hospdate"`  //入院日期 hospdate
	Nurselevel string  `json:"nurselevel"` //护理级别
	Icd string  `json:"icd"` //入院诊断
}
// Scan
type Result struct {
	Code    int  `json:"code"`
	Message  string `json:"message"`
	Data []Patinfo   `json:"data"`
}
func JsonHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("method:", r.Method)
	query := r.URL.Query()
	officeid := query["officeid"][0]
	db, err := gorm.Open("mssql", "sqlserver://" + viper.GetString("his.userid") + ":" + viper.GetString("his.password") + "@" + viper.GetString("his.ip") + ":" + viper.GetString("his.sqlport") + "?database=" +  viper.GetString("his.database") + ";encrypt=disable;app name=pda_server")
	if err != nil {
		log.Infof("无法连接数据库", err.Error())
		pub.ReturnJSON(-1, "无法连接数据库", w)
		return 
	}
	log.Infof("完成连接")
	defer db.Close()
	//-----科室患者清单
	ls_sql := `  
		SELECT  hospid = convert(varchar(20), patinfo.hospid) , 
				hospcode = patinfo.hospcode,
				bed = DictMedi.MediCode,   
				name = Patinfo.NAME,   
				sex = CASE  Patinfo.SEX WHEN 1 THEN '男' WHEN 2 THEN '女'  ELSE '未登记' End, 
				age = dbo.f_calcAge( patinfo.age, patinfo.agemonth, patinfo.ageday,    patinfo.agehh,  patinfo.agemm) ,
				officeid =doctmark.requoffiid,
				officename = dictoffice.officename,		
				doctorname = x1.opername,
				hospdate = convert(varchar(20), Patinfo.HospDate, 120),
				nurselevel = CASE  patinfo.nurselevel WHEN 1 THEN '特级' 
								WHEN 2 THEN '一级'  
								WHEN 3 THEN '二级'  
								WHEN 4 THEN '三级'  
								WHEN 5 THEN '常规'  
								WHEN 6 THEN '特殊疾病护理'  
								WHEN 7 THEN '新生儿护理'  
								ELSE '未知' End,
				icd = nIncomple.Name
		FROM Patinfo with (nolock) 
			left join nclintype on Patinfo.patitype = nclintype.clintypeid
			left join patientinfo on Patinfo.caseid = patientinfo.inpatino
			left join nIncomple on Patinfo.HospResu = nIncomple.incompleId
			,
			DictMedi,   
			DoctMark   
			left join dictoffice on DoctMark.requoffiid = dictoffice.officeid 
			left join dictoper x1 on doctmark.doctid = x1.operid 
			left join dictoper x2 on doctmark.operid = x2.operid  
		
		WHERE  Patinfo.HospID = DoctMark.HospID  and  
			DoctMark.MediID = DictMedi.MediId  and  
		--	  Patinfo.orderid = 0  AND  
			Patinfo.curstate = 3  AND  
			patinfo.HospPerf = 0  AND  
			doctmark.isMain = 1  AND  
			doctmark.enddate is null  AND  
			doctmark.marktype = 4  AND  
			dictmedi.item = 2 and   
			doctmark.requoffiid = ?
			order  by DictMedi.MediCode
`
	 
	rows, err := db.Raw(ls_sql , officeid).Rows() 
	
	//见：http://gorm.io/zh_CN/docs/sql_builder.html ,https://github.com/denisenkom/go-mssqldb
	if err != nil {
		log.Fatal("取医嘱出错", err) 
		pub.ReturnJSON(-1,  "取医嘱出错" , w)
		return  
	}
	var result Result
	for rows.Next() {
		var patinfo Patinfo
		///db.ScanRows(rows, &doctmark) 
		error := rows.Scan(&patinfo.Hospid, &patinfo.Hospcode, &patinfo.Bed, &patinfo.Name, &patinfo.Sex,
				&patinfo.Age, &patinfo.Officeid, &patinfo.Officename, &patinfo.Doctorname, &patinfo.Hospdate, &patinfo.Nurselevel ,
				&patinfo.Icd) 
	 
		if error != nil {
			fmt.Println(error);
		} 
		result.Data = append(result.Data,
			patinfo)
	}
	rows.Close()
	b, err := json.Marshal(result)
    if err != nil {
        fmt.Println("json err:", err)
    }
    fmt.Println(string(b));
 	//返回结果
	w.Header().Set("Content-Type", "application/json;charset=UTF-8")
	w.Write(b)
	log.Infof("成功，完成在院患者获取") 
	return
}
 
 