package pdaGetDoctmark

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
type  Doctmark struct {   
	//注意，必须首字母大写表示要输出
	Marktype string   `json:"marktype"`
	Mediname string   `json:"mediname"`
	Begidate string   `json:"begidate"`
	Doctname string   `json:"doctname"`
	Opername string   `json:"opername"`
	Curnum   int   `json:"curnum"`
	Thegroup int   `json:"thegroup"`
	Theoquan string   `json:"theoquan"`  //100ml
	Times string   `json:"times"`  //TID, BID
	Usenum string   `json:"usenum"`
	Usagetype string   `json:"usagetype"`  //大类中文：口服，输液，注射，其他
	Usagename string   `json:"usagename"`  //具体用法PO, IVD
	Dictententname string   `json:"dictententname"`  //嘱托
	Autonumb string   `json:"autonumb"`
	Detailid string   `json:"detailid"`
	Execcount int   `json:"execcount"`  //当天执行次数
	Execid string   `json:"execid"`
	Exectime string   `json:"exectime"`
	Exectime2 string   `json:"exectime2"`
	Checkoperid string   `json:"checkoperid"`
	Checkopername string   `json:"checkopername"`
	Isexec string   `json:"isexec"` 
}
type Patinfo struct {
	//注意，必须首字母大写表示要输出   0正确， 非0错
	Code    int  `json:"code"`
	Message  string `json:"message"`
	Hospid	string `json:"hospid"`
	Hospcode  string  `json:"hospcode"`
	Name	string `json:"name"`
	Sex  string  `json:"sex"`
	Bed   string `json:"bed"`
	Ddate  string  `json:"ddate"`  //入院日期 hospdate
	Doctmarks []Doctmark   `json:"data"`
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
	sdate := query["sdate"][0]  //yyyy-mm-dd 格式
	
	db, err := gorm.Open("mssql", "sqlserver://" + viper.GetString("his.userid") + ":" + viper.GetString("his.password") + "@" + viper.GetString("his.ip") + ":" + viper.GetString("his.sqlport") + "?database=" +  viper.GetString("his.database") + ";encrypt=disable;app name=pda_server")
	if err != nil {
		log.Infof("无法连接数据库", err.Error())
		pub.ReturnJSON(-1, "无法连接数据库", w)
		return 
	}
	log.Infof("完成连接")
	defer db.Close()
	//-----取患者基本资料
	ls_sql1 := ` select  patinfo.hospid , patinfo.hospcode, patinfo.name, 
				CASE patinfo.sex WHEN 1 THEN '男' when 2 THEN '女' ELSE '未登记' end sex,
				bed = DictMedi.MediCode ,
				convert(varchar(20),patinfo.hospdate ,120) ddate
				
			from  patinfo, doctmark, dictmedi
			where 
			Patinfo.HospID = DoctMark.HospID  and  
			DoctMark.MediID = DictMedi.MediId  and  
			doctmark.isMain = 1  AND  	
			doctmark.enddate is null  AND  
			doctmark.marktype = 4  AND  
			dictmedi.item = 2  	and 		
			patinfo.hospid = ? `
	var result Result
	db.Raw(ls_sql1, hospid).Scan(&result)
	if result.Name == ""{
		pub.ReturnJSON(-1, "患者不存在", w)
		log.Infof("患者不存在")
		return
	}
	fmt.Println(result.Hospcode, result.Sex);
		
	//-----取医嘱------
	/*
	ls_sql :=`   SELECT 
				marktype = CASE doctmark.marktype WHEN 1 THEN '长嘱' when 2 THEN '临嘱' when 5 THEN '出院带药'ELSE '其他' end ,
				Mediname = 
					CASE
								when iscancel = 1   then '(取消)' +  mediname  + ' ' + spec  + '(' + isnull(doctmark.remark, '') + ')' 
								else
									case 
									WHEN isnull(doctmark.remark,'') = '' THEN mediname  + ' ' + spec
									
									ELSE 
													case 
														when len(mediname) < 2 then doctmark.remark 
														else
															mediname + ' ' + spec + '(' + doctmark.remark+ ')' 
													end
								end 
					END ,
				BegiDate = convert(varchar(20), DoctMark.BegiDate, 120),
				Doctname = isnull(x1.opername, ''),
				Opername = isnull(x2.opername, ''),	
				Thegroup = DoctMark.TheGroup,
				Theoquan = DoctMark.TheoQuan + isnull(DoctMark.unit, '')  , --剂量+剂量单位
				Times = isnull(dicttimes.hz, ''),
				usenum = convert(varchar(20), DoctMark.thenum) + isnull(dictmedi.ambunit1, '')    , --领用数量
				usagetype = CASE dictusage.class 
							WHEN 1 THEN '口服'
							WHEN 2 THEN '输液'
							WHEN 3 THEN '注射'
							ELSE  '其他' end,
				Usagename = dictusage.usagename,
				Dictententname = isnull(dictent.entname,''),
				autonumb = convert(varchar(20), doctmark.autonumb) ,
				detailid = convert(varchar(20), hosp_detail_exec.detailid),
				execid = convert(varchar(20), execid),
				Exectime = convert(varchar(16), Exectime, 120),
				Exectime2 = isnull(convert(varchar(20), Exectime2, 120), ''),
				Checkoperid = isnull(convert(varchar(20), Checkoperid), ''),
				Checkopername =  isnull(convert(varchar(20), x5.opername), ''),
				isExec = CASE WHEN  Checkoperid > 0 THEN '1' ELSE '0' END
			FROM 	DoctMark   
					left join DictOper x2 on doctmark.operid = x2.operid 
					left join DictOper x3 on doctmark.doctid2 = x3.operid    
					left join DictOper x4 on doctmark.operid2 = x4.operid      
					left join DictOper x1 on doctmark.doctid = x3.operid     
					left join dictent   on doctmark.entid = dictent.id   
					left join dictusage on doctmark.usageid = dictusage.usageid    
					left join dicttimes on doctmark.entid = dicttimes.timeid   ,
					DictMedi,   
					dictusage_class,
					
					hosp_detail_exec 
					left join DictOper x5 on hosp_detail_exec.Checkoperid = x5.operid ,
					patinfo  
			WHERE ( DoctMark.MediID = DictMedi.MediId )    
				and dictusage.class = dictusage_class.class
				and doctmark.marktype in (1,2,5)
				and doctmark.iscancel = 0  --不显示取消的 
				and doctmark.hospid = patinfo.hospid 
				and patinfo.hospid = ?
				and doctmark.autonumb = hosp_detail_exec.autonumb  
				and (doctmark.ischeck = 1)  
				and hosp_detail_exec.exectime between  ? and ? 
			order by doctmark.marktype, doctmark.curnum, doctmark.thegroup, doctmark.serial, hosp_detail_exec.exectime   `
*/			
	ls_sql :=`   SELECT 
			marktype = CASE doctmark.marktype WHEN 1 THEN '长嘱' when 2 THEN '临嘱' when 5 THEN '出院带药'ELSE '其他' end ,
			Mediname = 
				CASE
							when iscancel = 1   then '(取消)' +  mediname  + ' ' + spec  + '(' + isnull(doctmark.remark, '') + ')' 
							else
								case 
								WHEN isnull(doctmark.remark,'') = '' THEN mediname  + ' ' + spec
								
								ELSE 
												case 
													when len(mediname) < 2 then doctmark.remark 
													else
														mediname + ' ' + spec + '(' + doctmark.remark+ ')' 
												end
							end 
				END ,
			BegiDate = convert(varchar(20), DoctMark.BegiDate, 120),
			Doctname = isnull(x1.opername, ''),
			Opername = isnull(x2.opername, ''),	
			Curnum = Doctmark.curnum,
			Thegroup = DoctMark.TheGroup,
			Theoquan = CASE WHEN len(DoctMark.TheoQuan) > 0 THEN DoctMark.TheoQuan + isnull(DoctMark.unit, '') else '' end , --剂量+剂量单位
			Times = isnull(dicttimes.hz, ''),
			usenum = convert(varchar(20), DoctMark.thenum) + isnull(dictmedi.ambunit1, '')    , --领用数量
			usagetype = CASE dictusage.class 
						WHEN 1 THEN '口服'
						WHEN 2 THEN '输液'
						WHEN 3 THEN '注射'
						ELSE  '其他' end,
			Usagename = isnull(dictusage.usagename, ''),
			Dictententname = isnull(dictent.entname,''),
			autonumb = convert(varchar(20), doctmark.autonumb) ,
			execcount = isnull(doctmark_exec.execcount, 0),
			--当天需要执行次数，要首日，平时，末日来辨断
			maxcount = CASE WHEN dicttimes.Periods =  1  THEN 
									case when convert(char(8),doctmark.BegiDate,112) = CONVERT(char(8),getdate(), 112) THEN doctmark.firstnum
										 when convert(char(8),doctmark.EndDate,112) = CONVERT(char(8),getdate(), 112) THEN doctmark.lastnum
										 ELSE DictTimes.Times END
								ELSE  0  --如果是2天，3天或者7天一次的，当天执行次数辨断较为复杂，时间关系，先暂用放0不做辨断
							END
		FROM 	DoctMark   
			left join DictOper x1 on doctmark.doctid = x1.operid     
			left join DictOper x2 on doctmark.operid = x2.operid 
			left join DictOper x3 on doctmark.doctid2 = x3.operid    
			left join DictOper x4 on doctmark.operid2 = x4.operid      
			left join dictent   on doctmark.entid = dictent.id   
			left join 
			(select dictusage.usageid, dictusage.usagename,dictusage.class  from dictusage ,dictusage_class
					where dictusage.class = dictusage_class.class ) dictusage on doctmark.usageid = dictusage.usageid    
			left join dicttimes on doctmark.timeid = dicttimes.timeid  
			left join Doctmark_exec on doctmark.autonumb = Doctmark_exec.autonumb   ,
			DictMedi,   
			--dictusage_class,
			patinfo  
		WHERE ( DoctMark.MediID = DictMedi.MediId )    
			and doctmark.marktype in (1,2,5)
			and doctmark.iscancel = 0  --不显示取消的 
			and doctmark.hospid = patinfo.hospid 
			and patinfo.hospid = ?
			and (doctmark.ischeck = 1)  
			and doctmark.begidate < ? 
			and isnull(doctmark.enddate, getdate()) < ?
		order by doctmark.marktype, doctmark.curnum, doctmark.thegroup, doctmark.serial   `			
			fmt.Println(sdate); // +
	rows, err := db.Raw(ls_sql , result.Hospid, sdate + " 23:59:59",  sdate + " 23:59:59").Rows() 
	
	//见：http://gorm.io/zh_CN/docs/sql_builder.html ,https://github.com/denisenkom/go-mssqldb
	if err != nil {
		log.Fatal("取医嘱出错", err) 
		pub.ReturnJSON(-1,  "取医嘱出错" , w)
		return  
	}
	
	var patinfo Patinfo
	patinfo.Hospid = strconv.Itoa(result.Hospid);
	patinfo.Hospcode = result.Hospcode
	patinfo.Name = result.Name
	patinfo.Sex = result.Sex
	patinfo.Ddate = result.Ddate
	patinfo.Bed = result.Bed
	for rows.Next() {
		var doctmark Doctmark
		///db.ScanRows(rows, &doctmark) 
		error := rows.Scan(&doctmark.Marktype, &doctmark.Mediname, &doctmark.Begidate, &doctmark.Doctname, &doctmark.Opername,
				&doctmark.Curnum,&doctmark.Thegroup, &doctmark.Theoquan, &doctmark.Times, &doctmark.Usenum, &doctmark.Usagetype, &doctmark.Usagename,
				&doctmark.Dictententname, &doctmark.Autonumb, &doctmark.Detailid, &doctmark.Execcount)
				//&doctmark.Execid, &doctmark.Exectime, &doctmark.Exectime2,
				//&doctmark.Checkoperid, &doctmark.Checkopername, &doctmark.Isexec) 
	 
		if error != nil {
			fmt.Println(error);
		} 
		patinfo.Doctmarks = append(patinfo.Doctmarks,
			 doctmark)
	}
	rows.Close()
	b, err := json.Marshal(patinfo)
    if err != nil {
        fmt.Println("json err:", err)
    }
   // fmt.Println(string(b));
 	//返回结果
	w.Header().Set("Content-Type", "application/json;charset=UTF-8")
	w.Write(b)
	log.Infof("成功，完成医嘱获取") 
	return
}
 
 