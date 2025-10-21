## SQL To Struct 数据库到Golang Struct转换工具

## 使用方式
1. 下载包然后执行 `go install github.com/aveyuan/sql2struct@latest`
2. 在项目文件下，执行 `sql2struct`

## 配置文件
>如果没有config.yml，需要执行 `sql2struct init` 创建，放置到项目下

`sql2struct -t init`

`config.yml`文件内容如下

```yml
config:
  struct_tpl: #自定义模板
  out_dir: ./models #输目录
  package: models #输出包名称
  format: crlf #格式化字符，默认lf(linux)，windows默认是crlf
  notname: false #用于模板判断是否需要生成tablename
  camelcase: false #是否生成驼峰命名json
db:
  db_name: api #db名称
  db_tables: [] #表选择
  dsn: root:123456@tcp(127.0.0.1:3306)/api?charset=utf8mb4&parseTime=true  #数据库连接

```

如果需要自定义模板，需要在`config.yml`中配置`struct_tpl`

`sql2struct -t template`

`struct.tpl`自定义模板内容如下
>如需自定义结构文件，可以自行创建文件定义

```go
package {{.Package}}

{{if .ImportTime}}import "time"{{end}}

type {{.StructName}} struct { {{range .MessageDetail}}
    {{.StructName}} {{.StructType}} `gorm:"column:{{.SQlCOLUMNNAME}};type:{{.SQLCOLUMNTYPE}}{{if .SQLCOLUMNDEFAULT}};default:{{.SQLCOLUMNDEFAULT}}{{end}}{{if eq .SQLCOLUMNKEY "PRI"}};primaryKey{{end}}{{if eq .SQLISNULLABLE "NO"}};not null{{end}}{{if ne .SQLCOLUMNCOMMENT ""}};comment:{{.SQLCOLUMNCOMMENT}}{{end}}" json:"{{.SQlCOLUMNNAME}}" form:"{{.SQlCOLUMNNAME}}"` {{end}} 
}

{{if eq .NotName false}}
func ({{.StructName}}) TableName() string {
	return "{{.TableName}}"
}
{{end}}

```

## 配置讲解

### 数据源配置
1. 首先是DB连接，配置一下DSN `db.dsn`
    > 单独的配置`db_name` 是为了更好的获取到数据库schema信息
2. 配置了需要导出那些表，如果`db_tables`为空，默认是导出所有

### 数据输出配置
- `config.struct_tpl` 指定了需要生成的文件模板。
    > 这使得我们可以灵活的定制文件的内容,针对不同的项目，可以使用不同的模板
- `config.out_dir` 指定了输出文件的路径
- `config.package` 指定model的包名称
- `config.format` 指定格式化文件字符
- `config.notname` 布尔值,用于模板判断是否需要生成tablename