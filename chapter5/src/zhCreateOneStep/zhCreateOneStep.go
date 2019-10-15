package zhCreateOneStep

import (
	"context"
	"database/sql"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
	"zh_func1"

	"github.com/bitly/go-simplejson"
	mssql "github.com/denisenkom/go-mssqldb"
	"github.com/lexkong/log"
	"github.com/spf13/viper"
)

func JsonHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("method: CreateOneStep ", r.Method)
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		fmt.Printf("read body err, %v\n", err)
		zh_func1.ReturnJSON("0", err.Error(), w)
		return
	}
	js, err := simplejson.NewJson(body)
	log.Infof(string(body))
	//fmt.Println(string(body))
	//解析参数

	ClubId := js.Get("clubId").MustString()
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
	var rs_fldclientid = "" //返回

	//下面这个方法也可以
	_, err1 := txn.ExecContext(ctx, "sp_get_client_no",
		sql.Named("as_linkcode", ClubId),
		sql.Named("al_type", 2),
		sql.Named("rs_fldclientid", sql.Out{Dest: &rs_fldclientid}),
	)
	if err1 != nil {
		log.Infof(err1.Error())
		zh_func1.ReturnJSON("0", err1.Error(), w)
		return
	}
	//fmt.Println("rs_autonumb is %s", rs_autonumb)

	var rs_errormsg = ""
	var rs_errorcode = 0

	//entityManager.getTransaction().begin();Double.parseDouble
	customerCode := js.Get("customerCode").MustString() //CRM会员ID
	customerName := js.Get("customerName").MustString() //CRM会员名称
	mobile := js.Get("mobile").MustString()             //手机号
	sex := js.Get("sex").MustString()
	birthday := js.Get("birthday").MustString()         //
	identityType := js.Get("identityType").MustString() //证件类型
	identityNo := js.Get("identityNo").MustString()     //证件号码
	email := js.Get("email").MustString()               //
	source := js.Get("source").MustString()             //信息来源
	userCode := js.Get("userCode").MustString()         //归属员工工号
	remark := js.Get("remark").MustString()             //备注
	club := js.Get("club").MustString()                 //归属门店
	packageId := js.Get("packageId").MustString()       //CRM会员名称
	price := js.Get("price").MustString()               //
	discount := js.Get("discount").MustString()         //
	cash := js.Get("cash").MustString()                 //
	bank := js.Get("bank").MustString()                 //
	transfer := js.Get("transfer").MustString()         //
	weixin := js.Get("weixin").MustString()             //
	zfb := js.Get("zfb").MustString()                   //
	agency := js.Get("agency").MustString()             //
	securitycard := js.Get("securitycard").MustString() //
	dj := js.Get("dj").MustString()                     //
	card := js.Get("card").MustString()                 //
	free := js.Get("free").MustString()                 //
	integral := js.Get("integral").MustString()         //
	coupon := js.Get("coupon").MustString()             //
	userId1 := js.Get("userId1").MustString()           //
	userperform1 := js.Get("userperform1").MustString() //
	userId2 := js.Get("userId2").MustString()           //
	userperform2 := js.Get("userperform2").MustString() //

	_, err1 = txn.ExecContext(ctx, "sp_g017_tran",
		sql.Named("fldclientid", rs_fldclientid),
		sql.Named("customerCode", customerCode),
		sql.Named("customerName", customerName),
		sql.Named("mobile", mobile),
		sql.Named("sex", sex),
		sql.Named("birthday", birthday),
		sql.Named("identityType", identityType),
		sql.Named("identityNo", identityNo),
		sql.Named("email", email),
		sql.Named("source", source),
		sql.Named("userCode", userCode),
		sql.Named("remark", remark),

		sql.Named("club", club),
		sql.Named("packageId", packageId),
		sql.Named("price", price),
		sql.Named("discount", discount),
		sql.Named("cash", cash),
		sql.Named("bank", bank),
		sql.Named("transfer", transfer),
		sql.Named("weixin", weixin),
		sql.Named("zfb", zfb),
		sql.Named("agency", agency),
		sql.Named("securitycard", securitycard),
		sql.Named("dj", dj),
		sql.Named("card", card),
		sql.Named("free", free),
		sql.Named("integral", integral),
		sql.Named("coupon", coupon),
		sql.Named("userId1", userId1),
		sql.Named("userperform1", userperform1),
		sql.Named("userId2", userId2),
		sql.Named("userperform2", userperform2),

		sql.Named("errormsg", sql.Out{Dest: &rs_errormsg}),
		sql.Named("errorcode", sql.Out{Dest: &rs_errorcode}),
	)
	if err1 != nil {
		log.Infof(err1.Error())
		zh_func1.ReturnJSON("0", err1.Error(), w)
		return
	}
	if rs_errorcode > 0 {
		err = txn.Commit()
		if err != nil {
			log.Infof("commit失败:" + err.Error())
			zh_func1.ReturnJSON("0", err.Error(), w)
			return
		}
		log.Infof("成功，会员id:" + rs_fldclientid)
		zh_func1.ReturnJSON("1", rs_fldclientid, w)
	} else {
		log.Infof("失败:" + rs_errormsg)
		zh_func1.ReturnJSON(strconv.Itoa(rs_errorcode), rs_errormsg, w)
	}
	return
}
