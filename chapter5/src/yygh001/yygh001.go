package yygh001
import(
	"encoding/xml"
	"net/http"
	"strings"
	"session"
	"github.com/spf13/viper"
	"io/ioutil"
	"github.com/lexkong/log"
)
import (
	"github.com/jinzhu/gorm"
	_ "github.com/denisenkom/go-mssqldb"
	_ "github.com/jinzhu/gorm/dialects/mssql"
  ) 
//4.6.1	上传医院预约挂号号源
type Program struct {
	XMLName xml.Name `xml:"program"`
	Session_id	string `xml:"session_id"`
	Function_id  string `xml:"function_id"`
	Akb020 string `xml:"akb020"`
	Doctor []Doctorsourceinfo `xml:"doctorsourceinfo>row"`  
}

type  Doctorsourceinfo struct {  //医生基本信息
	Aaz307 string  `xml:"aaz307"`  //医疗机构科室编号
	Aaz386 string  `xml:"aaz386"`  //科室名称
	Bka503 string  `xml:"bka503"`  //医生编号  (operid)
	Aac003 string  `xml:"aac003"`  //医生姓名
	Acd231 string  `xml:"acd231"`  //职称
	Bac045 string  `xml:"bac045"`  //专业
	Abk027 string  `xml:"abk027"`  //个人简介
	Source []Sourceinfo  `xml:"sourceinfo>row"` 
}
type Sourceinfo struct {   //医生班次信息
	Aae030 string `xml:"aae030"`   //接诊日期  yyyy.mm.dd
	Bae031 string `xml:"bae031"`   //全天，上午，下午，晚上
	Bae032 string `xml:"bae032"`	//接诊时间段 8:30-9:30
	Bae587 string `xml:"bae587"`  //总号源数
	Bae588 string `xml:"bae588"`	//剩余号源数	
	Bae589  string `xml:"bae589"`	//出诊地点	
	Akc225  string `xml:"akc225"`	//挂号费标准	
	Aae013  string `xml:"aae013"`	//说明	
	Operid  string `xml:"-"`        //只是过滤用，不输出XML
}

type program struct {
	Return_code string   `xml:"return_code"`
    Return_code_message  string   `xml:"return_code_message"`
}
func filterDoctor (x Sourceinfo, thisoperid string) bool{
	if x.Operid == thisoperid {
		//fmt.Println(x.Operid, thisoperid)
		return true
	}
	return false
}

func FilterSlice(s []Sourceinfo, filter func(x Sourceinfo, thisoperid string) bool, operid string) []Sourceinfo {
	// 返回的新切片
	// s[:0] 这种写法是创建了一个 len 为 0，cap 为 len(s) 即和原始切片最大容量一致的切片
	// 因为是过滤，所以新切片的元素总个数一定不大于比原始切片，这样做减少了切片扩容带来的影响
	// 同时，也有一个问题，因为 newS 和 s 共享底层数组，那么过滤后 s 也会被修改！
	newS := s[:0]
	// 遍历，对每个元素执行 filter，符合条件的加入新切片中
	for _, x := range s {
	  if filter(x, operid) {
		//fmt.Println(x)
		newS = append(newS, x)
	  }
	}
	return newS
  }

func returnXML(code string , errTxt string) []byte{
	outv := &program{code, errTxt}
	output, err := xml.MarshalIndent(outv, "  ", "    ")
	if err != nil {
		panic("未知错误" + errTxt)
	}
	 
	return output
	
}

func Process( s_session string, date string) []byte {
  log.Infof("连接数据库", "sqlserver://" + viper.GetString("his.userid") + ":" + viper.GetString("his.password") + "@" + viper.GetString("his.ip") + "?database=" +  viper.GetString("his.database") + ";encrypt=disable;app name=tqtest");
	db, err := gorm.Open("mssql", "sqlserver://" + viper.GetString("his.userid") + ":" + viper.GetString("his.password") + "@" + viper.GetString("his.ip") + "?database=" +  viper.GetString("his.database") + ";encrypt=disable;app name=tqtest")
	if err != nil {
		return returnXML("-1", "院方端无法连接数据库") 
		
	}
	log.Infof("完成连接")
	defer db.Close()
	var result []Doctorsourceinfo   //医生信息
	var bzresult []Sourceinfo //排班信息
 	dbdoc := db.Raw("[sp_getRegSourceInfo_tq]  ?,?,?", date, date, 2).Scan(&result) //见：http://gorm.io/zh_CN/docs/sql_builder.html ,https://github.com/denisenkom/go-mssqldb
	dbbz := db.Raw("[sp_getRegSourceInfo_tq]  ?,?,?", date, date, 3).Scan(&bzresult) //见：http://gorm.io/zh_CN/docs/sql_builder.html ,https://github.com/denisenkom/go-mssqldb
	if dbdoc.Error != nil {
	//	fmt.Println(dbdoc.Error)
		log.Fatal("院方端取数据出库(医生信息)", dbdoc.Error)
		return returnXML("-1", "院方端取数据出库(医生信息)" + dbdoc.Error.Error() )
	}
	if dbbz.Error != nil {
		log.Fatal("院方端取数据出库(排班数据)", dbbz.Error)
		return returnXML("-1", "院方端取数据出库(排班数据)" + dbbz.Error.Error() )
	}
	 
	yh := &Program{Session_id: s_session, Function_id:"yygh001", Akb020:"sc000d"}
	//var doctor *Doctorsourceinfo
	begin1:
	for _, arDoc := range result {
		//过滤当前医生的
		for _, s := range FilterSlice(bzresult, filterDoctor,  arDoc.Bka503){
			arDoc.Source = append(arDoc.Source, s)
		}
		yh.Doctor = append(yh.Doctor, arDoc)
	}
	// var doctor *Doctorsourceinfo
	// doctor = &Doctorsourceinfo{Aaz307:"赵医生", Aaz386:"001"}
	// doctor.Source = append(doctor.Source, Sourceinfo{Aae030: "上午"});
	// doctor.Source = append(doctor.Source, Sourceinfo{Aae030: "下午"});
	// doctor.Source = append(doctor.Source, Sourceinfo{Aae030: "晚上"});
	// yh.Doctor =  append(yh.Doctor, *doctor)

	// doctor = &Doctorsourceinfo{Aaz307:"李医生"}
	// doctor.Source = append(doctor.Source, Sourceinfo{Aae030: "上午"});
	// doctor.Source = append(doctor.Source, Sourceinfo{Aae030: "下午"});
	// doctor.Source = append(doctor.Source, Sourceinfo{Aae030: "晚上"});
	// yh.Doctor =  append(yh.Doctor, *doctor)


	output2, err2 := xml.MarshalIndent(yh, "", "    ")
    if err2 != nil {
        log.Fatal("error: %v\n", err2)
	}
	v := program{} 
	log.Infof(string(output2))
	//下面数据上传
	request, _ := http.NewRequest("POST", viper.GetString("ZHSocialURL"), strings.NewReader(string(output2)))
	//post数据并接收http响应
	resp,err :=http.DefaultClient.Do(request)
	if err!=nil {
		return returnXML("-1", err.Error()) 
	}else {
		respBody,_ := ioutil.ReadAll(resp.Body)
		
		err = xml.Unmarshal(respBody, &v)
		if err != nil {
			return returnXML("-1", err.Error()) 
		}
		if v.Return_code == "-9" {
			//如果返回-9,需要重新取session
			s_session = session.GetSession()
			goto begin1
		}
	}
	return returnXML(v.Return_code,v.Return_code_message)
  }