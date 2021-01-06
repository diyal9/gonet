package db

import (
	"encoding/json"
	"fmt"
	"github.com/golang/protobuf/proto"
	"gonet/base"
	"math"
	"reflect"
	"strconv"
	"strings"
	"unicode"
)

func insertsql(sqlData *SqlData, p *Properties, val string){
	sqlData.SqlValue += fmt.Sprintf("'%s',", val)
	sqlData.SqlName += fmt.Sprintf("`%s`,", p.Name)
}

func insertsqlblob(sqlData *SqlData, p *Properties, val []byte){
	sqlData.SqlValue += fmt.Sprintf("'0x%s',", val)
	sqlData.SqlName += fmt.Sprintf("`%s`,", p.Name)
}

func insertsqlarray(sqlData *SqlData, p *Properties, val string, i int){
	if sqlData.bitMap != nil && !sqlData.bitMap.Test(i){
		return
	}
	sqlData.SqlValue += fmt.Sprintf("'%s',", val)
	sqlData.SqlName += fmt.Sprintf("`%s%d`,", p.Name, i)
}

func getInsertSql(classField reflect.StructField, classVal reflect.Value, sqlData *SqlData) (bool) {
	p := getProperties(classField)
	sType := getTypeString(classField, classVal)
	if p.IsJson(){
		data, _ := json.Marshal(classVal.Interface())
		insertsql(sqlData, p, string(data))
		return true
	}else if p.IsBlob(){
		for classVal.Kind() == reflect.Ptr {
			classVal = classVal.Elem()
		}
		data, _ := proto.Marshal(classVal.Addr().Interface().(proto.Message))
		insertsqlblob(sqlData, p, data)
		return true
	}else if p.IsIgnore(){
		return true
	}

	switch sType {
	case "*bool":
		value := bool(false)
		if !classVal.IsNil() {
			value = *classVal.Interface().(*bool)
		}
		insertsql(sqlData, p, strconv.FormatBool(value))
	case "*string":
		value := string("")
		if !classVal.IsNil() {
			value = *classVal.Interface().(*string)
		}
		insertsql(sqlData, p, value)
	case "*float32":
		value := float32(0)
		if !classVal.IsNil() {
			value = *classVal.Interface().(*float32)
		}
		insertsql(sqlData, p, strconv.FormatFloat(float64(value), 'f', -1, 32))
	case "*float64":
		value := float64(0)
		if !classVal.IsNil() {
			value = *classVal.Interface().(*float64)
		}
		insertsql(sqlData, p, strconv.FormatFloat(value, 'f', -1, 64))
	case "*int":
		value := int(0)
		if !classVal.IsNil() {
			value = *classVal.Interface().(*int)
		}
		insertsql(sqlData, p, strconv.FormatInt(int64(value),10))
	case "*int8":
		value := int8(0)
		if !classVal.IsNil() {
			value = *classVal.Interface().(*int8)
		}
		insertsql(sqlData, p, strconv.FormatInt(int64(value),10))
	case "*int16":
		value := int16(0)
		if !classVal.IsNil() {
			value = *classVal.Interface().(*int16)
		}
		insertsql(sqlData, p, strconv.FormatInt(int64(value),10))
	case "*int32":
		value := int32(0)
		if !classVal.IsNil() {
			value = *classVal.Interface().(*int32)
		}
		insertsql(sqlData, p, strconv.FormatInt(int64(value),10))
	case "*int64":
		value := int64(0)
		if !classVal.IsNil() {
			value = *classVal.Interface().(*int64)
		}
		if !p.IsDatetime(){
			insertsql(sqlData, p, strconv.FormatInt(int64(value),10))
		}else{
			insertsql(sqlData, p, GetDBTimeString(int64(value)))
		}
	case "*uint":
		value := uint(0)
		if !classVal.IsNil() {
			value = *classVal.Interface().(*uint)
		}
		insertsql(sqlData, p, strconv.FormatUint(uint64(value),10))
	case "*uint8":
		value := uint8(0)
		if !classVal.IsNil() {
			value = *classVal.Interface().(*uint8)
		}
		insertsql(sqlData, p, strconv.FormatUint(uint64(value),10))
	case "*uint16":
		value := uint16(0)
		if !classVal.IsNil() {
			value = *classVal.Interface().(*uint16)
		}
		insertsql(sqlData, p, strconv.FormatUint(uint64(value),10))
	case "*uint32":
		value := uint32(0)
		if !classVal.IsNil() {
			value = *classVal.Interface().(*uint32)
		}
		insertsql(sqlData, p, strconv.FormatUint(uint64(value),10))
	case "*uint64":
		value := uint64(0)
		if !classVal.IsNil() {
			value = *classVal.Interface().(*uint64)
		}
		insertsql(sqlData, p, strconv.FormatUint(uint64(value),10))
	case "*struct":
		if !classVal.IsNil() {
			value := classVal.Elem().Interface()
			parseInserSql(value, sqlData)
		}



	case "bool":
		insertsql(sqlData, p, strconv.FormatBool(classVal.Bool()))
	case "string":
		insertsql(sqlData, p, classVal.String())
	case "float32":
		insertsql(sqlData, p, strconv.FormatFloat(classVal.Float(), 'f', -1, 32))
	case "float64":
		insertsql(sqlData, p, strconv.FormatFloat(classVal.Float(), 'f', -1, 64))
	case "int":
		insertsql(sqlData, p, strconv.FormatInt(classVal.Int(),10))
	case "int8":
		insertsql(sqlData, p, strconv.FormatInt(classVal.Int(),10))
	case "int16":
		insertsql(sqlData, p, strconv.FormatInt(classVal.Int(),10))
	case "int32":
		insertsql(sqlData, p, strconv.FormatInt(classVal.Int(),10))
	case "int64":
		if !p.IsDatetime(){
			insertsql(sqlData, p, strconv.FormatInt(classVal.Int(),10))
		}else{
			insertsql(sqlData, p, GetDBTimeString(classVal.Int()))
		}
	case "uint":
		insertsql(sqlData, p, strconv.FormatUint(classVal.Uint(),10))
	case "uint8":
		insertsql(sqlData, p, strconv.FormatUint(classVal.Uint(),10))
	case "uint16":
		insertsql(sqlData, p, strconv.FormatUint(classVal.Uint(),10))
	case "uint32":
		insertsql(sqlData, p, strconv.FormatUint(classVal.Uint(), 10))
	case "uint64":
		insertsql(sqlData, p, strconv.FormatUint(classVal.Uint(),10))
	case "struct":
		parseInserSql(classVal.Interface(), sqlData)



	case "[]bool":
		value := []bool{}
		if !classVal.IsNil() {
			value = classVal.Interface().([]bool)
		}
		for i,v := range value{
			insertsqlarray(sqlData, p, strconv.FormatBool(v), i)
		}
	case "[]string":
		value := []string{}
		if !classVal.IsNil() {
			value = classVal.Interface().([]string)
		}
		for i,v := range value{
			insertsqlarray(sqlData, p, v, i)
		}
	case "[]float32":
		value := []float32{}
		if !classVal.IsNil() {
			value = classVal.Interface().([]float32)
		}
		for i,v := range value{
			insertsqlarray(sqlData, p, strconv.FormatFloat(float64(v), 'f', -1, 32), i)
		}
	case "[]float64":
		value := []float64{}
		if !classVal.IsNil() {
			value = classVal.Interface().([]float64)
		}
		for i,v := range value{
			insertsqlarray(sqlData, p, strconv.FormatFloat(v, 'f', -1, 64), i)
		}
	case "[]int":
		value := []int{}
		if !classVal.IsNil() {
			value = classVal.Interface().([]int)
		}
		for i,v := range value{
			insertsqlarray(sqlData, p, strconv.FormatInt(int64(v), 10), i)
		}
	case "[]int8":
		value := []int8{}
		if !classVal.IsNil() {
			value = classVal.Interface().([]int8)
		}
		for i,v := range value{
			insertsqlarray(sqlData, p, strconv.FormatInt(int64(v), 10), i)
		}
	case "[]int16":
		value := []int16{}
		if !classVal.IsNil() {
			value = classVal.Interface().([]int16)
		}
		for i,v := range value{
			insertsqlarray(sqlData, p, strconv.FormatInt(int64(v), 10), i)
		}
	case "[]int32":
		value := []int32{}
		if !classVal.IsNil() {
			value = classVal.Interface().([]int32)
		}
		for i,v := range value{
			insertsqlarray(sqlData, p, strconv.FormatInt(int64(v), 10), i)
		}
	case "[]int64":
		value := []int64{}
		if !classVal.IsNil() {
			value = classVal.Interface().([]int64)
		}
		for i,v := range value{
			if !p.IsDatetime(){
				insertsqlarray(sqlData, p, strconv.FormatInt(int64(v), 10), i)
			}else{
				insertsqlarray(sqlData, p, GetDBTimeString(v), i)
			}
		}
	case "[]uint":
		value := []uint{}
		if !classVal.IsNil() {
			value = classVal.Interface().([]uint)
		}
		for i,v := range value{
			insertsqlarray(sqlData, p, strconv.FormatUint(uint64(v), 10), i)
		}
	case "[]uint8":
		value := []uint8{}
		if !classVal.IsNil() {
			value = classVal.Interface().([]uint8)
		}
		for i,v := range value{
			insertsqlarray(sqlData, p, strconv.FormatUint(uint64(v), 10), i)
		}
	case "[]uint16":
		value := []uint16{}
		if !classVal.IsNil() {
			value = classVal.Interface().([]uint16)
		}
		for i,v := range value{
			insertsqlarray(sqlData, p, strconv.FormatUint(uint64(v), 10), i)
		}
	case "[]uint32":
		value := []uint32{}
		if !classVal.IsNil() {
			value = classVal.Interface().([]uint32)
		}
		for i,v := range value{
			insertsqlarray(sqlData, p, strconv.FormatUint(uint64(v), 10), i)
		}
	case "[]uint64":
		value := []uint64{}
		if !classVal.IsNil() {
			value = classVal.Interface().([]uint64)
		}
		for i,v := range value{
			insertsqlarray(sqlData, p, strconv.FormatUint(uint64(v), 10), i)
		}
	case "[]struct":
		for i := 0;  i < classVal.Len(); i++{
			parseInserSql(classVal.Index(i).Interface(), sqlData)
		}



	case "[*]bool":
		for i := 0; i < classVal.Len(); i++{
			insertsqlarray(sqlData, p, strconv.FormatBool(classVal.Index(i).Bool()), i)
		}
	case "[*]string":
		for i := 0; i < classVal.Len(); i++{
			insertsqlarray(sqlData, p, classVal.Index(i).String(), i)
		}
	case "[*]float32":
		for i := 0; i < classVal.Len(); i++{
			insertsqlarray(sqlData, p, strconv.FormatFloat(classVal.Index(i).Float(), 'f', -1, 64), i)
		}
	case "[*]float64":
		for i := 0; i < classVal.Len(); i++{
			insertsqlarray(sqlData, p, strconv.FormatFloat(classVal.Index(i).Float(), 'f', -1, 64), i)
		}
	case "[*]int":
		for i := 0; i < classVal.Len(); i++{
			insertsqlarray(sqlData, p, strconv.FormatInt(classVal.Index(i).Int(), 10), i)
		}
	case "[*]int8":
		for i := 0; i < classVal.Len(); i++{
			insertsqlarray(sqlData, p, strconv.FormatInt(classVal.Index(i).Int(), 10), i)
		}
	case "[*]int16":
		for i := 0; i < classVal.Len(); i++{
			insertsqlarray(sqlData, p, strconv.FormatInt(classVal.Index(i).Int(), 10), i)
		}
	case "[*]int32":
		for i := 0; i < classVal.Len(); i++{
			insertsqlarray(sqlData, p, strconv.FormatInt(classVal.Index(i).Int(), 10), i)
		}
	case "[*]int64":
		for i := 0; i < classVal.Len(); i++{
			insertsqlarray(sqlData, p, strconv.FormatInt(classVal.Index(i).Int(), 10), i)
		}
	case "[*]uint":
		for i := 0; i < classVal.Len(); i++{
			insertsqlarray(sqlData, p, strconv.FormatUint(classVal.Index(i).Uint(), 10), i)
		}
	case "[*]uint8":
		for i := 0; i < classVal.Len(); i++{
			insertsqlarray(sqlData, p, strconv.FormatUint(classVal.Index(i).Uint(), 10), i)
		}
	case "[*]uint16":
		for i := 0; i < classVal.Len(); i++{
			insertsqlarray(sqlData, p, strconv.FormatUint(classVal.Index(i).Uint(), 10), i)
		}
	case "[*]uint32":
		for i := 0; i < classVal.Len(); i++{
			insertsqlarray(sqlData, p, strconv.FormatUint(classVal.Index(i).Uint(), 10), i)
		}
	case "[*]uint64":
		for i := 0; i < classVal.Len(); i++{
			insertsqlarray(sqlData, p, strconv.FormatUint(classVal.Index(i).Uint(), 10), i)
		}
	case "[*]struct":
		for i := 0;  i < classVal.Len(); i++{
			parseInserSql(classVal.Index(i).Interface(), sqlData)
		}
	default:
		fmt.Println("getInsertSql type not supported", sType,  classField.Type)
		panic("getInsertSql type not supported")
		return false
		//}
	}
	return true
}

func parseInserSql(obj interface{}, sqlData *SqlData){
	classVal := reflect.ValueOf(obj)
	for classVal.Kind() == reflect.Ptr {
		classVal = classVal.Elem()
	}
	classType := classVal.Type()

	for i := 0; i < classType.NumField(); i++{
		if !classVal.Field(i).CanInterface(){
			continue
		}

		bRight:= getInsertSql(classType.Field(i), classVal.Field(i), sqlData)
		if !bRight{
			errorStr := fmt.Sprintf("parseInserSql type not supported %s", classType.Name())
			panic(errorStr)
			return //丢弃这个包
		}
	}
}

func insertSqlStr(sqltable string, sqlData *SqlData) string{
	sqlname := sqlData.SqlName
	sqlvalue := sqlData.SqlValue
	index := strings.LastIndex(sqlname, ",")
	if index!= -1{
		sqlname = sqlname[:index]
	}

	index = strings.LastIndex(sqlvalue, ",")
	if index!= -1{
		sqlvalue = sqlvalue[:index]
	}
	return "insert into "+ sqltable + " (" + sqlname+") VALUES (" + sqlvalue + ")"
}

//--- struct to sql
func InsertSql(obj interface{}, sqltable string,)string{
	defer func() {
		if err := recover(); err != nil {
			base.TraceCode(err)
		}
	}()

	sqlData := &SqlData{}
	parseInserSql(obj, sqlData)
	return  insertSqlStr(sqltable, sqlData)
}

func InsertSqlEx(obj interface{}, sqltable string, params ...string) string {
	defer func() {
		if err := recover(); err != nil {
			base.TraceCode(err)
		}
	}()

	classVal := reflect.ValueOf(obj)
	for classVal.Kind() == reflect.Ptr {
		classVal = classVal.Elem()
	}
	classType := classVal.Type()

	sqlData := &SqlData{}
	nameMap := make(map[string] *base.BitMap)//name index[for array]
	for _,v := range params{
		nIndex, i := 0, 0
		v1 := strings.ToLower(v)
		v2 := strings.TrimRightFunc(v, func(r rune) bool {
			if unicode.IsNumber(r){
				nIndex = int(r - '0') * int(math.Pow(10, float64(i))) + nIndex
				i++
				return true
			}
			return false
		})
		if v1 != v2{
			bitMap, bOk := nameMap[v2]
			if !bOk{
				bitMap = base.NewBitMap(MAX_ARRAY_LENGTH)
				nameMap[v2] = bitMap
			}
			bitMap.Set(nIndex)
		}else{
			nameMap[v1] = nil
		}
	}
	for i := 0; i < classType.NumField(); i++ {
		if !classVal.Field(i).CanInterface() {//private成员不能读取
			continue
		}

		sf := classType.Field(i)
		p := getProperties(sf)
		bitMap, exist := nameMap[p.Name]
		if exist{
			sqlData.bitMap = bitMap
			bRight := getInsertSql(sf, classVal.Field(i), sqlData)
			if !bRight{
				errorStr := fmt.Sprintf("InsertSqlEx error %s", reflect.TypeOf(obj).Name())
				panic(errorStr)
				return ""//丢弃这个包
			}
		}
	}
	return insertSqlStr(sqltable, sqlData)
}



