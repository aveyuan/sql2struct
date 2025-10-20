package main

type Message struct {
	StructName    string
	TableName     string
	Package       string
	ImportTime    bool
	NotName       bool
	MessageDetail []TableFied
}
type Field struct {
	TypeName string
	AttrName string
	Num      int
}

type TableFied struct {
	SQlCaseCOLUMNNAME string
	SQlCOLUMNNAME     string
	SQLCOLUMNTYPE     string
	SQLDATATYPE       string
	SQLCOLUMNKEY      string
	SQLISNULLABLE     string
	SQLCOLUMNCOMMENT  string
	SQLCOLUMNDEFAULT  string
	Num               int
	StructName        string
	StructType        string
}

type DDLFied struct {
	Table      string
	CrateTable string
}

var typeArrStruct = map[string]string{
	"int":       "int",
	"tinyint":   "int8",
	"smallint":  "int16",
	"mediumint": "int32",
	"enum":      "int",
	"bigint":    "int64",
	"char":      "string",
	"varchar":   "string",
	"text":      "string",
	"longtext":  "string",
	"timestamp": "time.Time",
	"date":      "time.Time",
	"datetime":  "time.Time",
	"double":    "float64",
	"decimal":   "float64",
	"float":     "float64",
	// 对于无符号的支持
	"uint":       "uint",
	"utinyint":   "uint8",
	"usmallint":  "uint16",
	"umediumint": "uint32",
	"uenum":      "uint32",
	"ubigint":    "uint64",
}
