package yygh001
//4.6.1	上传医院预约挂号号源
type Program struct {
	Session_id	string `xml:"session_id"`
	Function_id  string `xml:"function_id"`
	Akb020 string `xml:"akb020"`
	Doctor []Doctorsourceinfo `xml:"doctorsourceinfo>row"`  
}

type  Doctorsourceinfo struct {  //医生基本信息
	Aaz307 string  `xml:"aaz307"`
	Aaz386 string  `xml:"aaz386"`
	Bka503 string  `xml:"bka503"`
	Aac003 string  `xml:"aac003"`
	Acd231 string  `xml:"acd231"`
	Bac045 string  `xml:"bac045"`
	Abk027 string  `xml:"abk027"`
	Source []Sourceinfo  `xml:"sourceinfo>row"` 
}
type Sourceinfo struct {   //医生班次信息
	Aae030 string `xml:"aae030"`
	Bae031 string `xml:"bae031"`   //全天，上午，下午，晚上
}