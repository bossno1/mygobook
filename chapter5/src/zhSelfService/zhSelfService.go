package zhSelfService

import (
	"context"
	"net/http"

	"database/sql"
	"fmt"
	"io/ioutil"
	"net/url"
	"strconv"
	"zh_func1"

	"github.com/bitly/go-simplejson"
	mssql "github.com/denisenkom/go-mssqldb"
	"github.com/lexkong/log"
	"github.com/spf13/viper"
)

func JsonHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("method:", r.Method)
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		fmt.Printf("read body err, %v\n", err)
		zh_func1.ReturnJSON("0", err.Error(), w)
		return
	}
	log.Infof("----收到请求------")
	js, err := simplejson.NewJson(body)
	//log.Infof(string(body))
	//fmt.Println(string(body))
	//解析参数

	ClubId := js.Get("clubId").MustString()
	remark1 := js.Get("Remark1").MustString()
	ConsumeTime := js.Get("consumeTime").MustString()
	//fmt.Println(ClubId, remark1, ConsumeTime)
	arr, _ := js.Get("infoList").Array()

	//设置数据库链接参数
	query := url.Values{}
	query.Add("app name", "MyAppName")
	query.Add("encrypt", "disable")
	query.Add("database", viper.GetString("zh.database"))
	u := &url.URL{
		Scheme: "sqlserver",
		User:   url.UserPassword(viper.GetString("zh.userid"), viper.GetString("zh.password")),
		Host:   fmt.Sprintf("%s:%s", viper.GetString("zh.ip"), viper.GetString("zh.sqlport")),
		//Path:  ".\\sql2008" , //instance, // if connecting to an instance instead of a port
		RawQuery: query.Encode(),
	}
	connector, err := mssql.NewConnector(u.String())
	if err != nil {
		log.Infof(u.String(), err.Error())
		zh_func1.ReturnJSON("0", err.Error(), w)
		return
	}
	db := sql.OpenDB(connector)
	defer db.Close()
	//启动事务
	txn, err := db.Begin()
	if err != nil {
		log.Infof(u.String(), err.Error())
		zh_func1.ReturnJSON("0", err.Error(), w)
		return
	}
	defer txn.Rollback()
	//----开始处理-----
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	//connector.SessionInitSQL = "set implicit_transactions off"
	var rs_autonumb = "" //返回流水号

	//下面这个方法也可以 (db.ExecContext 不在事务加)
	_, err1 := txn.ExecContext(ctx, "sp_get_invoinfo",
		sql.Named("al_item", 1),
		sql.Named("as_linkcode", ClubId),
		sql.Named("as_date", ConsumeTime),
		sql.Named("rs_autonumb", sql.Out{Dest: &rs_autonumb}),
	)
	log.Infof("取得流水号：" + rs_autonumb)
	if err1 != nil {
		log.Infof(err1.Error())
		zh_func1.ReturnJSON("0", err1.Error(), w)
		return
	}
	//fmt.Println("rs_autonumb is %s", rs_autonumb)
	var info interface{}
	var rs_errormsg = ""
	var price, count, itemprice, totalprice float64

	for index := 0; index < len(arr); index++ {
		info = arr[index]
		//type interface {} does not support indexing， 要用以下方式获取
		count, _ = strconv.ParseFloat(info.(map[string]interface{})["count"].(string), 64)
		itemprice, _ = strconv.ParseFloat(info.(map[string]interface{})["itemPrice"].(string), 64)
		itemprice = zh_func1.Decimal(itemprice, 4)
		price = zh_func1.Decimal(count*itemprice, 2)
		totalprice = totalprice + price

		_, err1 := txn.ExecContext(ctx, "sp_g012",
			sql.Named("newAutonumb", rs_autonumb),
			sql.Named("serial", index),
			sql.Named("fldclientid", info.(map[string]interface{})["customerId"]),
			sql.Named("carid", info.(map[string]interface{})["packageId"]),
			sql.Named("clubId", ClubId),
			sql.Named("itemId", info.(map[string]interface{})["itemId"]),
			sql.Named("itemBzPrice", info.(map[string]interface{})["itemBzPrice"]),
			sql.Named("itemPrice", price),
			sql.Named("relaId", info.(map[string]interface{})["relaId"]),
			sql.Named("servicePersonal", info.(map[string]interface{})["servicePersonal"]),
			sql.Named("consumeTime", ConsumeTime),
			sql.Named("type", info.(map[string]interface{})["type"]),
			sql.Named("QUANTITY", info.(map[string]interface{})["count"]),
			sql.Named("totalPrice", 0),
			sql.Named("remark1", remark1),
			sql.Named("saveType", 2), //2写表明细
			sql.Named("errormsg", sql.Out{Dest: &rs_errormsg}),
		)
		if err1 != nil {
			log.Infof(err1.Error())
			zh_func1.ReturnJSON("0", err1.Error(), w)
			return
		}
		if len(rs_errormsg) > 0 {
			log.Infof(rs_errormsg)
			zh_func1.ReturnJSON("0", rs_errormsg, w)
			return
		}

	}
	log.Infof("完成明细写入" + rs_autonumb)

	//写表头
	_, err1 = txn.ExecContext(ctx, "sp_g012",
		sql.Named("newAutonumb", rs_autonumb),
		sql.Named("serial", 0),
		sql.Named("fldclientid", info.(map[string]interface{})["customerId"]),
		sql.Named("carid", info.(map[string]interface{})["packageId"]),
		sql.Named("clubId", ClubId),
		sql.Named("itemId", info.(map[string]interface{})["itemId"]),
		sql.Named("itemBzPrice", info.(map[string]interface{})["itemBzPrice"]),
		sql.Named("itemPrice", 0), //写表头不用这个
		sql.Named("relaId", info.(map[string]interface{})["relaId"]),
		sql.Named("servicePersonal", info.(map[string]interface{})["servicePersonal"]),
		sql.Named("consumeTime", ConsumeTime),
		sql.Named("type", info.(map[string]interface{})["type"]),
		sql.Named("QUANTITY", info.(map[string]interface{})["count"]),
		sql.Named("totalPrice", totalprice),
		sql.Named("remark1", remark1),
		sql.Named("saveType", 1), //1写表头
		sql.Named("errormsg", sql.Out{Dest: &rs_errormsg}),
	)
	if err1 != nil {
		log.Infof(err1.Error())
		zh_func1.ReturnJSON("-2", err1.Error(), w)
		return
	}
	if len(rs_errormsg) > 0 {
		log.Infof(rs_errormsg)
		zh_func1.ReturnJSON("0", rs_errormsg, w)
		return
	}
	log.Infof("完成表头写入" + rs_autonumb)

	err = txn.Commit()
	if err != nil {
		log.Infof(err.Error())
		zh_func1.ReturnJSON("-1", err.Error(), w)
		return
	}

	log.Infof("----成功完成----" + rs_autonumb)

	zh_func1.ReturnJSON("1", rs_autonumb, w)

	return

}
