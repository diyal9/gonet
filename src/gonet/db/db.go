package db

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"gonet/base"
	"reflect"
	"strconv"
	"strings"
	"time"
)

const (
	MAX_ARRAY_LENGTH = 64
)

type(
	Datetime int64

	Row struct {
		m_Resut map[string] string
	}

	IRow interface {
		init()
		Set(key, val string)
		Get(key string) string
		String(key string) string
		Int(key string) int
		Int64(key string) int64
		Float32(key string) float32
		Float64(key string) float64
		Bool(key string) bool
		Time(key string) int64
		Byte(key string) []byte
		Obj(obj interface{}) bool
	}

	Rows struct {
		m_Rows []*Row
		m_posRow int
	}

	IRows interface {
		Next() bool
		Row() *Row
		Obj(obj interface{}) bool
	}

	Properties struct {
		Name string
		Primary bool
		DateTime bool
		Blob bool
		Json bool
		Ignore bool
		tag string
	}

	SqlData struct {
		SqlName string
		SqlValue string
		bitMap *base.BitMap
	}
)

//主键 `sql:"primary"`
func (this *Properties) IsPrimary ()bool{
	return this.Primary
}

//日期 `sql:"datetime"`
func (this *Properties) IsDatetime ()bool{
	return this.DateTime
}

//二进制 `sql:"blob"`
func (this *Properties) IsBlob ()bool{
	return this.Blob
}

//json `sql:"json"`
func (this *Properties) IsJson ()bool{
	return this.Json
}

//ignore `sql:"-"`
func (this *Properties) IsIgnore ()bool{
	return this.Ignore
}

//---获取datetime时间
func  GetDBTimeString(t int64)string{
	tm := time.Unix(t, 0)
	return  tm.Format("2006-01-02 15:04:05")
}

func OpenDB(svr string, usr string, pwd string, db string) *sql.DB {
	sqlstr := fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=utf8mb4", usr, pwd, svr, db)
	mydb, err := sql.Open("mysql", sqlstr)
	base.ChechErr(err)
	return mydb
}

func getProperties(sf reflect.StructField) *Properties{
	p := &Properties{}
	p.tag = sf.Tag.Get("sql")
	fields := strings.Split(p.tag, ";")
	for _, v := range fields{
		switch v {
		case "primary" :
			p.Primary = true
		case "datetime":
			p.DateTime = true
		case "blob":
			p.Blob = true
		case "json":
			p.Json = true
		case "-":
			p.Ignore = true
		default:
			if strings.Contains(v, "name:"){
				p.Name = v[5:]
			}
		}
	}
	return p
}

func (this *Row) init() {
	this.m_Resut = make(map[string] string)
}

func (this *Row) Set(key, val string){
	this.m_Resut[key] = val
}

func (this *Row) Get(key string) string{
	//key = strings.ToLower(key)
	v, exist := this.m_Resut[key]
	if exist{
		return v
	}

	return ""
}

func (this *Row) String(key string) string{
	return this.Get(key)
}

func (this *Row) Int(key string) int{
	n, _ := strconv.Atoi(this.Get(key))
	return n
}

func (this *Row) Int64(key string) int64{
	n, _ := strconv.ParseInt(this.Get(key), 0, 64)
	return n
}

func (this *Row) Float32(key string) float32{
	n, _ := strconv.ParseFloat(this.Get(key), 32)
	return float32(n)
}

func (this *Row) Float64(key string) float64{
	n, _ := strconv.ParseFloat(this.Get(key), 64)
	return n
}

func (this *Row) Bool(key string) bool{
	n, _ := strconv.ParseBool(this.Get(key))
	return n
}

func (this *Row) Time(key string) int64{
	return base.GetDBTime(this.Get(key)).Unix()
}

func (this *Row) Byte(key string) []byte{
	return []byte(this.Get(key))
}

func (this *Row) Obj(obj interface{}) bool{
	return LoadObjSql(obj, this)
}

func (this *Rows) init(){
	this.m_posRow = 0
}

func (this *Rows) Next() bool{
	if this.m_posRow < len(this.m_Rows){
		this.m_posRow++
		return true
	}
	return false
}

func (this *Rows) Row() *Row{
	nPos := this.m_posRow-1
	if nPos >= 0 && nPos < len(this.m_Rows){
		return this.m_Rows[nPos]
	}

	return NewRow()
}

func (this *Rows) Obj(obj interface{}) bool{
	defer func() {
		if err := recover(); err != nil {
			base.TraceCode(err)
		}
	}()

	r := reflect.Indirect(reflect.ValueOf(obj))
	isPtr := false
	if kind := r.Kind(); kind == reflect.Slice {
		rType := r.Type().Elem()
		if rType.Kind() == reflect.Ptr {
			isPtr = true
			rType = rType.Elem()
		}
		for this.Next(){
			elem := reflect.New(rType).Elem()
			LoadObjSql(elem.Addr().Interface(), this.Row())
			if isPtr{
				r.Set(reflect.Append(r, elem.Addr()))
			}else{
				r.Set(reflect.Append(r, elem))
			}
		}
	}
	return true
}

func NewRow() *Row{
	row := &Row{}
	row.init()
	return row
}

func Query(rows *sql.Rows, err error) *Rows{
	rs := &Rows{}
	rs.init()
	if rows != nil && err == nil{
		cloumns, err := rows.Columns()
		cloumnsLen := len(cloumns)
		if err == nil && cloumnsLen > 0{
			for rows.Next(){
				r := NewRow()
				value := make([]*string, cloumnsLen)
				value1 := make([]interface{}, cloumnsLen)
				for i, _ := range value{
					value[i] = new(string)
					value1[i] = value[i]
				}
				rows.Scan(value1...)
				for i, v := range value{
					r.m_Resut[cloumns[i]] = *v
				}
				rs.m_Rows = append(rs.m_Rows, r)
			}
		}
		rows.Close()
	}
	return rs
}
//--------------------note存储过程----------------------//
//mysql存储过程多变更集的时候要用 NextResultSet()
/*rows, err := this.m_db.Query(fmt.Sprintf("call `sp_checkcreatePlayer`(%d)", this.AccountId))
if err == nil && rows != nil{
	if rows.NextResultSet(){//import
	rs := db.Query(rows, err)
		if rs.Next(){
			err := rs.Row().Int("@err")
		}
	}
}*/