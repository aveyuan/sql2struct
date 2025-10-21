package {{.Package}}

{{if .ImportTime}}import "time"{{end}}

type {{.StructName}} struct { {{range .MessageDetail}}
    {{.StructName}} {{.StructType}} `gorm:"column:{{.SQlCOLUMNNAME}};type:{{.SQLCOLUMNTYPE}}{{if .SQLCOLUMNDEFAULT}};default:{{.SQLCOLUMNDEFAULT}}{{end}}{{if eq .SQLCOLUMNKEY "PRI"}};primaryKey{{end}}{{if eq .SQLISNULLABLE "NO"}};not null{{end}}{{if ne .SQLCOLUMNCOMMENT ""}};comment:{{.SQLCOLUMNCOMMENT}}{{end}}" json:"{{.SQlCOLUMNNAMEFMT}}" form:"{{.SQlCOLUMNNAMEFMT}}"` {{end}} 
}

{{if eq .NotName false}}
func ({{.StructName}}) TableName() string {
	return "{{.TableName}}"
}
{{end}}