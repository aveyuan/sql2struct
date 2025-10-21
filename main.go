package main

import (
	"bufio"
	"bytes"
	"database/sql"
	_ "embed"
	"flag"
	"fmt"
	"html/template"
	"log"
	"os"

	"os/exec"
	"strings"

	_ "github.com/go-sql-driver/mysql"
	"github.com/spf13/viper"
)

//go:embed "struct.tpl"
var TplStruct string

var DB *sql.DB
var Dsn string
var DBName string
var OutDir string
var InTplStruct string
var Tables []string
var ExTables []string
var Package string
var Format string
var NotName bool
var SQLDir string
var RunFlag string
var CamelCase bool

func init() {
	flag.StringVar(&RunFlag, "t", "sql2struct", "sql2struct 生成模型 init 输出配置文件 template 输出模板配置文件")
	flag.Parse()
}

func main() {

	if RunFlag == "sql2struct" {
		Run()
		log.Println("generate struct success")
	}

	if RunFlag == "init" {
		InitConfig()
		log.Println("init config success")
	}

	if RunFlag == "template" {
		GenarateTemplate()
		log.Println("generate template success")
	}

}

func InitConfig() {
	os.WriteFile("config.yml", []byte(`config:
  struct_tpl: #自定义模板
  out_dir: ./models #输目录
  package: models #输出包名称
  format: crlf #格式化字符，默认lf(linux)，windows默认是crlf
  notname: false #用于模板判断是否需要生成tablename
  camelcase: false #是否生成驼峰命名json
db:
  db_name: api #db名称
  db_tables: [] #表选择
  dsn: root:123456@tcp(127.0.0.1:3306)/?charset=utf8mb4&parseTime=true  #数据库连接`), 0644)
}

func GenarateTemplate() {
	os.WriteFile("struct.tpl", []byte(TplStruct), 0644)
}

func Run() {

	viper.SetConfigFile("config.yml")
	err := viper.ReadInConfig()
	if err != nil {
		log.Fatal(err)
	}
	Dsn = viper.GetString("db.dsn")
	DBName = viper.GetString("db.db_name")
	OutDir = viper.GetString("config.out_dir")
	if OutDir == "" {
		OutDir = "message"
	}

	InTplStruct = viper.GetString("config.struct_tpl")

	Package = viper.GetString("config.package")
	Format = viper.GetString("config.format")

	Tables = viper.GetStringSlice("db.db_tables")
	ExTables = viper.GetStringSlice("db.db_extables")
	NotName = viper.GetBool("config.notname")
	CamelCase = viper.GetBool("config.camelcase")
	SQLDir = viper.GetString("sql.out_dir")

	db, err := Connect("mysql", Dsn)
	if err != nil {
		log.Fatal("db连接失败")
	}
	DB = db

	all := GetTables(DBName)

	if len(Tables) > 0 {
		var b []string
		for _, v := range Tables {
			for _, v2 := range all {
				if v2 == v {
					b = append(b, v2)
					break
				}
			}
		}
		all = b
	}

	if len(all) == 0 {
		return
	}

	var messages []Message
	for _, v := range all {
		var one Message
		one.StructName = UderscoreToUpperCamelCase(v) //首字母大写
		one.TableName = v
		one.NotName = NotName
		all := GetStruct(v, DBName)
		for k, v2 := range all {
			all[k].StructType = TypeMToStruct(func() string {
				// 对无符号进行判别
				if strings.Contains(v2.SQLCOLUMNTYPE, "unsigned") {
					return "u" + v2.SQLDATATYPE
				}
				return v2.SQLDATATYPE
			}())
			all[k].StructName = UderscoreToUpperCamelCase(v2.SQlCOLUMNNAME)
			if CamelCase {
				all[k].SQlCOLUMNNAMEFMT = ToLowUpperCamelCase(v2.SQlCOLUMNNAME)
			} else {
				all[k].SQlCOLUMNNAMEFMT = all[k].SQlCOLUMNNAME
			}

			// 判断是否有time
			if strings.Contains(all[k].StructType, "time.Time") {
				one.ImportTime = true
			}
		}
		one.MessageDetail = all
		messages = append(messages, one)
	}
	GenarateStruct(OutDir, messages)

}

func GenarateStruct(dir string, all []Message) {
	var tmpl *template.Template
	var err error

	if InTplStruct != "" {
		tmpl, err = template.ParseFiles(InTplStruct)
		if err != nil {
			log.Fatal(err)
		}
	} else {
		tmpl, err = template.New("index").Parse(TplStruct)
		if err != nil {
			log.Fatal(err)
		}
	}

	_, err = os.Stat(dir)
	if err != nil {
		// 创建文件夹
		os.MkdirAll(dir, 0755)
	}
	// 循环，创建具体的文件
	for _, v := range all {
		filepath := fmt.Sprintf("%v/%v.go", dir, v.TableName)

		file, err := os.OpenFile(filepath, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0755)
		if err != nil {
			log.Fatal(err)
		}
		defer file.Close()
		v.Package = Package
		err = tmpl.Execute(file, v)
		if err != nil {
			log.Fatal(err)
		}

		// 执行代码格式化
		cmd := exec.Command("gofmt", "-w", filepath)
		if err := cmd.Run(); err != nil {
			log.Print(err)
		}

		// gofmt会对文件格式crlf转换为lf，现在重新转换
		if Format == "CRLF" {
			if err := Rewite(filepath); err != nil {
				log.Fatal(err)
			}
		}

	}
}

func Rewite(filepath string) error {
	writer := bytes.NewBuffer(nil)
	fx := func() error {
		f, err := os.Open(filepath)
		if err != nil {
			return err
		}
		defer f.Close()

		scanner := bufio.NewScanner(f)
		for scanner.Scan() {
			line := scanner.Text() + "\r\n"
			_, err := writer.WriteString(line)
			if err != nil {
				return err
			}
		}

		// 检查是否有错误发生
		if err := scanner.Err(); err != nil {
			return err
		}
		return nil

	}
	if err := fx(); err != nil {
		return err
	}
	if err := os.WriteFile(filepath, writer.Bytes(), 0755); err != nil {
		return err
	}

	return nil
}

func GetTables(dbName string) []string {
	rows, err := DB.Query(fmt.Sprintf(`SELECT TABLE_NAME FROM information_schema.TABLES t WHERE  t.TABLE_SCHEMA = "%v"`, dbName))
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()
	var all []string
	for rows.Next() {
		var name string
		rows.Scan(&name)
		all = append(all, name)
	}

	if err = rows.Err(); err != nil {
		log.Fatal(err)
	}
	return all
}

func GetStruct(dbTable, dbName string) []TableFied {
	rows, err := DB.Query("SELECT  c.COLUMN_NAME,c.COLUMN_TYPE,c.DATA_TYPE,c.COLUMN_KEY ,c.IS_NULLABLE ,c.COLUMN_COMMENT,c.COLUMN_DEFAULT  FROM INFORMATION_SCHEMA.Columns c WHERE c.`TABLE_SCHEMA`='" + dbName + "' AND c.TABLE_NAME = '" + dbTable + "' ORDER BY c.ORDINAL_POSITION")
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()
	var all []TableFied
	var num = 1
	for rows.Next() {
		var one TableFied
		rows.Scan(&one.SQlCOLUMNNAME, &one.SQLCOLUMNTYPE, &one.SQLDATATYPE, &one.SQLCOLUMNKEY, &one.SQLISNULLABLE, &one.SQLCOLUMNCOMMENT, &one.SQLCOLUMNDEFAULT)
		one.Num = num
		all = append(all, one)
		num++
	}
	return all

}

func Connect(driverName, dsn string) (*sql.DB, error) {
	db, err := sql.Open(driverName, dsn)

	if err != nil {
		log.Fatalln(err)
	}
	db.SetMaxIdleConns(0)
	db.SetMaxOpenConns(2)
	if err := db.Ping(); err != nil {
		log.Fatalln(err)
	}
	return db, err
}

func TypeMToStruct(m string) string {
	if _, ok := typeArrStruct[m]; ok {
		return typeArrStruct[m]
	}
	return "string"
}

// 下划线单词转为大写驼峰单词
func UderscoreToUpperCamelCase(s string) string {
	s = strings.Replace(s, "_", " ", -1)
	s = strings.Title(s)
	return strings.Replace(s, " ", "", -1)
}

// 头部小写，其他大写
func ToLowUpperCamelCase(s string) string {
	s = strings.Replace(s, "_", " ", -1)
	s = strings.Title(s)
	return strings.ToLower(s[:1]) + strings.Replace(s, " ", "", -1)[1:]
}
