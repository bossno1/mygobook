package main

import (
   
	"log"
	"time"
	_ "github.com/denisenkom/go-mssqldb"
	"github.com/go-xorm/xorm"
	
    /* 这个star 比较少_ "github.com/mattn/go-adodb"*/
   
)  
var engine *xorm.Engine
func main() {
	 
	var err error
	engine, err = xorm.NewEngine("sqlserver", "sqlserver://sa:146-164-156-@127.0.0.2:52813?database=master")

  //  db, err := sqlx.Connect("sqlserver", "sqlserver://sa:146-164-156-@127.0.0.2:52813?database=master")
   if err != nil {
		log.Fatalln(err)
		return;
    }
	type User struct {
		Id int64
		Name string
		Salt string
		Age int
		Passwd string `xorm:"varchar(200)"`
		Created time.Time `xorm:"created"`
		Updated time.Time `xorm:"updated"`
	}
	
	err1 := engine.Sync2(new(User))
	if err1 != nil {
		log.Fatalln(err1)
		return;
    }
     
}