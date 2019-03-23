package main
//taoqing 2019.3.2 d views/list.html 增加了 visjs， 增加对于PB端的互动
/*
运行photoweb.exe 
visjs演示：
http://localhost:8080/view/
*/
import (
	"html/template"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path"
	"runtime/debug"
	"database/sql"
	"strings"
	"fmt"
)
//go get github.com/mattn/go-adodb
import (
	_ "github.com/mattn/go-adodb"
)
type Mssql struct {
	*sql.DB
	dataSource string
	database string
	windows bool
	sa   SA
}
type SA struct {
	user string
	passwd string
	port int 

}
func (m *Mssql) Open() (err error){
	var conf []string
	conf = append(conf, "Provider=SQLOLEDB")
	conf = append(conf, "Data Source=" +m.dataSource)
	conf = append(conf, "Initial Catalog=" +m.database)

	//以当前windows系统用户身份登录SQL， 如果服务不支持这个方式登录，就会出错
	if m.windows {
		conf = append(conf, "intergrated security=SSPI")
	} else {
		conf = append(conf, "hostname=ffff")
		conf = append(conf, "user id=" + m.sa.user)
		conf = append(conf, "password=" + m.sa.passwd)
		//conf = append(conf, "port="+ fmt.Sprint(m.sa.port))
	}
	m.DB, err = sql.Open("adodb", strings.Join(conf, ";"))
	if err != nil {
		return err
	}
	return nil
}

func main() {
	db := Mssql {
		//DESKTOP-LN3T12M\SQL2008
		//127.0.0.1
		//127.0.0.1\\SQL2008
		dataSource: "192.168.31.144,51798",
		database: "his_yb",
		//windows: true为windows身份验证, false必须设置sa帐号和密码
		windows: false,
		sa: SA{
			user:"sa",
			passwd: "146-164-156-",
			port: 51798, //这个参数发现没用， 而是在端口写在datastore用 ,号分隔即可
		},
	}
	err := db.Open() //发现这里并不会马上打开，而是等到Query发起才会去连接数据库
	if err != nil{
		fmt.Println("sql open:", err)
		return
	}
	defer db.Close()  //最好在close之前加个等待命令，看是否真的没有断开链接
	rows, err := db.Query("select opername from dictoper")
	if err != nil {
		fmt.Println("query:", err)
		return
	}
	for rows.Next(){
		var name string
		rows.Scan(&name)
		fmt.Printf("Name: %s\n", name)
	}
}

//---------------
const (
	ListDir      = 0x0001
	UPLOAD_DIR   = "./uploads"
	TEMPLATE_DIR = "./views"
)

var templates = make(map[string]*template.Template)

func init() {
	fileInfoArr, err := ioutil.ReadDir(TEMPLATE_DIR)
	check(err)
	var templateName, templatePath string
	for _, fileInfo := range fileInfoArr {
		templateName = fileInfo.Name()
		if ext := path.Ext(templateName); ext != ".html" {
			continue
		}
		templatePath = TEMPLATE_DIR + "/" + templateName
		log.Println("Loading template:", templatePath)
		t := template.Must(template.ParseFiles(templatePath))
		templates[templateName] = t
	}
}

func check(err error) {
	if err != nil {
		panic(err)
	}
}

func renderHtml(w http.ResponseWriter, tmpl string, locals map[string]interface{}) {
	tmpl += ".html"
	err := templates[tmpl].Execute(w, locals)
	check(err)
}

func isExists(path string) bool {
	_, err := os.Stat(path)
	if err == nil {
		return true
	}
	return os.IsExist(err)
}

func uploadHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		renderHtml(w, "upload", nil)
	}
	if r.Method == "POST" {
		f, h, err := r.FormFile("image")
		check(err)
		filename := h.Filename
		defer f.Close()
		t, err := os.Create(UPLOAD_DIR + "/" + filename)
		check(err)
		defer t.Close()
		_, err = io.Copy(t, f)
		check(err)
		http.Redirect(w, r, "/view?id="+filename, http.StatusFound)
	}
}

func viewHandler(w http.ResponseWriter, r *http.Request) {
	imageId := r.FormValue("id")
	imagePath := UPLOAD_DIR + "/" + imageId
	if ok := isExists(imagePath); !ok {
		http.NotFound(w, r)
		return
	}

	w.Header().Set("Content-Type", "image")
	http.ServeFile(w, r, imagePath)
}

func listHandler(w http.ResponseWriter, r *http.Request) {
	fileInfoArr, err := ioutil.ReadDir("./uploads")
	check(err)
	locals := make(map[string]interface{})
	images := []string{}
	for _, fileInfo := range fileInfoArr {
		images = append(images, fileInfo.Name())
	}
	locals["images"] = images
	renderHtml(w, "list", locals)
}

func safeHandler(fn http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if e, ok := recover().(error); ok {
				http.Error(w, e.Error(), http.StatusInternalServerError)

				// 或者输出自定义的 50x 错误页面
				// w.WriteHeader(http.StatusInternalServerError)
				// renderHtml(w, "error", e.Error())

				// logging
				log.Println("WARN: panic fired in %v.panic - %v", fn, e)
				log.Println(string(debug.Stack()))
			}
		}()
		fn(w, r)
	}
}

func staticDirHandler(mux *http.ServeMux, prefix string, staticDir string, flags int) {
	mux.HandleFunc(prefix, func(w http.ResponseWriter, r *http.Request) {
		file := staticDir + r.URL.Path[len(prefix)-1:]
		if (flags & ListDir) == 0 {
			fi, err := os.Stat(file)
			if err != nil || fi.IsDir() {
				http.NotFound(w, r)
				return
			}
		}
		http.ServeFile(w, r, file)
	})
}

func main1() {
	mux := http.NewServeMux()
	staticDirHandler(mux, "/assets/", "./public", 0)
	mux.HandleFunc("/", safeHandler(listHandler))
	mux.HandleFunc("/view", safeHandler(viewHandler))
	mux.HandleFunc("/upload", safeHandler(uploadHandler))
	err := http.ListenAndServe(":8080", mux)
	if err != nil {
		log.Fatal("ListenAndServe: ", err.Error())
	}
}
