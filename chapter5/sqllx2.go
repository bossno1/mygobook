package main

import (
    "database/sql"
    "fmt"
	"log"
	"net/url"
	/*_ "github.com/lib/pq"*/
	_ "github.com/denisenkom/go-mssqldb"
    /* 这个star 比较少_ "github.com/mattn/go-adodb"*/
    "github.com/jmoiron/sqlx"
)
var schema = `
--可控病种列表（启用标志+有相关诊断）
--select distinct 可控病种ID,可控病种自编码,受控病种 from v_cp_kkbz where ICD10_CODE = 'xxxxx'
create table doctmark_cp_diff3
(
    ID        numeric               identity,
    autonumb  int                   not null,
    ddate     datetime              not null,
    sdate     varchar(8)            not null,
    operid    int                   not null,
    hospid    int                   not null,
    mediid    int                   not null,
    mediname  varchar(255)          not null,
    Reasonid  int                   not null,
    Reason    varchar(255)          not null,
    constraint PK_DOCTMARK_CP_DIFF primary key (ID)
)
;

/* ============================================================ */
/*   Index: doctmark_cp_diff_i1                                 */
/* ============================================================ */
create index doctmark_cp_diff_i3 on doctmark_cp_diff3 (autonumb, sdate)
;
CREATE TABLE person (
    first_name varchar(50),
    last_name  varchar(50),
    email  varchar(50)
);

CREATE TABLE place (
    country varchar(50),
    city varchar(50) NULL,
    telcode integer
)`

type Person struct {
    FirstName string `db:"first_name"`
    LastName  string `db:"last_name"`
    Email     string
}

type Place struct {
    Country string
    City    sql.NullString
    TelCode int
}

func main() {
	query := url.Values{}
  	query.Add("app name", "MyAppName")

  	u := &url.URL{
      Scheme:   "sqlserver",
      User:     url.UserPassword("sa", "146-164-156-"),
      Host:     fmt.Sprintf("%s:%d", "127.0.0.1", 51798),
      // Path:  instance, // if connecting to an instance instead of a port
      RawQuery: query.Encode(),
  	}
    // this Pings the database trying to connect, panics on error
    // use sqlx.Open() for sql.Open() semantics
	//db, err := sqlx.Connect("postgres", "user=foo dbname=bar sslmode=disable")
	//db, err := sqlx.Connect("adodb", "Provider=SQLOLEDB;Data Source=192.168.31.144,51798;Initial Catalog=his_yb;user id=sa;password=146-164-156-;")
	fmt.Println( u.String())
    //db, err := sqlx.Connect("sqlserver", u.String()) //"sqlserver://sa:146-164-156-@localhost:51798?database=his_yb&connection+timeout=30")
    //注意连接SQL SERVER 2008R2出错（可能需要打SP3补丁），而SQL SERVER 2017  port:52813 没有问题  
    db, err := sqlx.Connect("sqlserver", "sqlserver://sa:146-164-156-@127.0.0.1:51798?database=master;encrypt=disable;app name=tqtest")
	//db, err := sqlx.Connect("sqlserver", "server=192.168.31.144;port=51798;user id=sa;password=146-164-156-;database=his_yb") // ://sa:146-164-156-@localhost:51798?database=his_yb&connection+timeout=30"") //"")
	//server=localhost\\SQLExpress;user id=sa;database=master;app name=MyAppName
    if err != nil {
		log.Fatalln(err)
		return;
    }

    // exec the schema or fail; multi-statement Exec behavior varies between
    // database drivers;  pq will exec them all, sqlite3 won't, ymmv
    db.MustExec(schema)
    /*
    以下如果用dbMustExec 则没有事务
    
    

    db.MustExec("INSERT INTO person (first_name, last_name, email) VALUES (@p1, @p2, @p3)", "Jason", "Moiron", "jmoiron@jmoiron.net")
    db.MustExec("INSERT INTO person (first_name, last_name, email) VALUES (@p1, @p2, @p3)", "John", "Doe", "johndoeDNE@gmail.net")
    */

    /*
    以下用tx 则是有事务
    */
    tx := db.MustBegin()
    tx.MustExec("INSERT INTO person (first_name, last_name, email) VALUES (@p1, @p2, @p3)", "Jason", "Moiron", "jmoiron@jmoiron.net")
    tx.MustExec("INSERT INTO person (first_name, last_name, email) VALUES (@p1, @p2, @p3)", "John", "Doe", "johndoeDNE@gmail.net")
    tx.MustExec("INSERT INTO place (country, city, telcode) VALUES (@p1, @p2, @p3)", "United States", "New York", 1)
    tx.MustExec("INSERT INTO place (country, telcode) VALUES (@p1, @p2)", "Hong Kong", 852)
    tx.MustExec("INSERT INTO place (country, telcode) VALUES (@p1, @p2)", "Singapore", 65)
    //Named queries can use structs, so if you have an existing struct (i.e. person := &Person{}) that you have populated, you can pass it in as &person
    //tx.NamedExec("INSERT INTO person (first_name, last_name, email) VALUES (:first_name, :last_name, :email)", &Person{"Jane", "Citizen", "jane.citzen@example.com"})
    // var FirstName, SecondNames string
	// fmt.Printf("Please enter your full name: ")
	// fmt.Scanln(&FirstName, &SecondNames) 
	// fmt.Printf("Hi %s %s!\n", FirstName, SecondNames)
    tx.Commit()

    // Query the database, storing results in a []Person (wrapped in []interface{})
    people := []Person{}
    db.Select(&people, "SELECT * FROM person ORDER BY first_name ASC")
    fmt.Printf("count: %d ", len(people))
    jason, john := people[0], people[1]

    fmt.Printf("%#v\n%#v", jason, john)
    // Person{FirstName:"Jason", LastName:"Moiron", Email:"jmoiron@jmoiron.net"}
    // Person{FirstName:"John", LastName:"Doe", Email:"johndoeDNE@gmail.net"}

    // You can also get a single result, a la QueryRow
    jason = Person{}
    err = db.Get(&jason, "SELECT * FROM person WHERE first_name=$1", "Jason")
    fmt.Printf("%#v\n", jason)
    // Person{FirstName:"Jason", LastName:"Moiron", Email:"jmoiron@jmoiron.net"}

    // if you have null fields and use SELECT *, you must use sql.Null* in your struct
    places := []Place{}
    err = db.Select(&places, "SELECT * FROM place ORDER BY telcode ASC")
    if err != nil {
        fmt.Println(err)
        return
    }
    usa, singsing, honkers := places[0], places[1], places[2]
    
    fmt.Printf("%#v\n%#v\n%#v\n", usa, singsing, honkers)
    // Place{Country:"United States", City:sql.NullString{String:"New York", Valid:true}, TelCode:1}
    // Place{Country:"Singapore", City:sql.NullString{String:"", Valid:false}, TelCode:65}
    // Place{Country:"Hong Kong", City:sql.NullString{String:"", Valid:false}, TelCode:852}

    // Loop through rows using only one struct
    place := Place{}
    rows, err := db.Queryx("SELECT * FROM place")
    for rows.Next() {
        err := rows.StructScan(&place)
        if err != nil {
            log.Fatalln(err)
        } 
        fmt.Printf("%#v\n", place)
    }
    // Place{Country:"United States", City:sql.NullString{String:"New York", Valid:true}, TelCode:1}
    // Place{Country:"Hong Kong", City:sql.NullString{String:"", Valid:false}, TelCode:852}
    // Place{Country:"Singapore", City:sql.NullString{String:"", Valid:false}, TelCode:65}

    // Named queries, using `:name` as the bindvar.  Automatic bindvar support
    // which takes into account the dbtype based on the driverName on sqlx.Open/Connect
    _, err = db.NamedExec(`INSERT INTO person (first_name,last_name,email) VALUES (:first,:last,:email)`, 
        map[string]interface{}{
            "first": "Bin",
            "last": "Smuth",
            "email": "bensmith@allblacks.nz",
    })

    // Selects Mr. Smith from the database
    rows, err = db.NamedQuery(`SELECT * FROM person WHERE first_name=:fn`, map[string]interface{}{"fn": "Bin"})

    // Named queries can also use structs.  Their bind names follow the same rules
    // as the name -> db mapping, so struct fields are lowercased and the `db` tag
    // is taken into consideration.
    rows, err = db.NamedQuery(`SELECT * FROM person WHERE first_name=:first_name`, jason)
}