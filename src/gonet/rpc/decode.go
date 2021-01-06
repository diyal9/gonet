package rpc

import (
	"context"
	"github.com/golang/protobuf/proto"
	"github.com/json-iterator/go"
	"gonet/base"
	"gonet/message"
	"reflect"
	"strings"
)

//rpc UnmarshalHead
func UnmarshalHead(buff []byte) (*message.RpcPacket, RpcHead){
	rpcPacket := &message.RpcPacket{}
	proto.Unmarshal(buff, rpcPacket)
	if rpcPacket.RpcHead == nil{
		rpcPacket.RpcHead = &message.RpcHead{}
	}
	rpcPacket.FuncName = strings.ToLower(rpcPacket.FuncName)
	return rpcPacket, *(*RpcHead)(rpcPacket.RpcHead)
}

//rpc Unmarshal
//pFuncType for RegisterCall func
func UnmarshalBody(rpcPacket *message.RpcPacket, pFuncType reflect.Type) []interface{}{
	bitstream := base.NewBitStream(rpcPacket.RpcBody, len(rpcPacket.RpcBody))
	nCurLen := int(rpcPacket.ArgLen)
	params := make([]interface{}, nCurLen+1)
	params[0] = context.WithValue(context.Background(), "rpcHead", *(*RpcHead)(rpcPacket.RpcHead))
	for i := 1; i <= nCurLen; i++  {
		switch bitstream.ReadInt(8) {
		case RPC_BOOL:
			params[i] = bitstream.ReadFlag()
		case RPC_STRING:
			params[i] = bitstream.ReadString()
		case RPC_FLOAT32:
			params[i] = bitstream.ReadFloat()
		case RPC_FLOAT64:
			params[i] = bitstream.ReadFloat64()
		case RPC_INT8:
			params[i] = int8(bitstream.ReadInt(8))
		case RPC_INT:
			params[i] = bitstream.ReadInt(32)
		case RPC_INT16:
			params[i] = int16(bitstream.ReadInt(16))
		case RPC_INT32:
			params[i] = int32(bitstream.ReadInt(32))
		case RPC_INT64:
			params[i] = int64(bitstream.ReadInt64(64))
		case RPC_UINT:
			params[i] = uint(bitstream.ReadInt(32))
		case RPC_UINT8:
			params[i] = uint8(bitstream.ReadInt(8))
		case RPC_UINT16:
			params[i] = uint16(bitstream.ReadInt(16))
		case RPC_UINT32:
			params[i] = uint32(bitstream.ReadInt(32))
		case RPC_UINT64:
			params[i] = uint64(bitstream.ReadInt64(64))


		case RPC_BOOL_SLICE:
			nLen := bitstream.ReadInt(16)
			val := make([]bool, nLen)
			for i := 0; i < nLen; i++ {
				val[i] = bitstream.ReadFlag()
			}
			params[i] = val
		case RPC_STRING_SLICE:
			nLen := bitstream.ReadInt(16)
			val := make([]string, nLen)
			for i := 0; i < nLen; i++ {
				val[i] = bitstream.ReadString()
			}
			params[i] = val
		case RPC_FLOAT32_SLICE:
			nLen := bitstream.ReadInt(16)
			val := make([]float32, nLen)
			for i := 0; i < nLen; i++ {
				val[i] = bitstream.ReadFloat()
			}
			params[i] = val
		case RPC_FLOAT64_SLICE:
			nLen := bitstream.ReadInt(16)
			val := make([]float64, nLen)
			for i := 0; i < nLen; i++ {
				val[i] = bitstream.ReadFloat64()
			}
			params[i] = val
		case RPC_INT_SLICE:
			nLen := bitstream.ReadInt(16)
			val := make([]int, nLen)
			for i := 0; i < nLen; i++ {
				val[i] = bitstream.ReadInt(32)
			}
			params[i] = val
		case RPC_INT8_SLICE:
			nLen := bitstream.ReadInt(16)
			val := make([]int8, nLen)
			for i := 0; i < nLen; i++ {
				val[i] = int8(bitstream.ReadInt(8))
			}
			params[i] = val
		case RPC_INT16_SLICE:
			nLen := bitstream.ReadInt(16)
			val := make([]int16, nLen)
			for i := 0; i < nLen; i++ {
				val[i] = int16(bitstream.ReadInt(16))
			}
			params[i] = val
		case RPC_INT32_SLICE:
			nLen := bitstream.ReadInt(16)
			val := make([]int32, nLen)
			for i := 0; i < nLen; i++ {
				val[i] = int32(bitstream.ReadInt(32))
			}
			params[i] = val
		case RPC_INT64_SLICE:
			nLen := bitstream.ReadInt(16)
			val := make([]int64, nLen)
			for i := 0; i < nLen; i++ {
				val[i] = int64(bitstream.ReadInt64(64))
			}
			params[i] = val
		case RPC_UINT_SLICE:
			nLen := bitstream.ReadInt(16)
			val := make([]uint, nLen)
			for i := 0; i < nLen; i++ {
				val[i] = uint(bitstream.ReadInt(32))
			}
			params[i] = val
		case RPC_UINT8_SLICE:
			nLen := bitstream.ReadInt(16)
			val := make([]uint8, nLen)
			for i := 0; i < nLen; i++ {
				val[i] = uint8(bitstream.ReadInt(8))
			}
			params[i] = val
		case RPC_UINT16_SLICE:
			nLen := bitstream.ReadInt(16)
			val := make([]uint16, nLen)
			for i := 0; i < nLen; i++ {
				val[i] = uint16(bitstream.ReadInt(16))
			}
			params[i] = val
		case RPC_UINT32_SLICE:
			nLen := bitstream.ReadInt(16)
			val := make([]uint32, nLen)
			for i := 0; i < nLen; i++ {
				val[i] = uint32(bitstream.ReadInt(32))
			}
			params[i] = val
		case RPC_UINT64_SLICE:
			nLen := bitstream.ReadInt(16)
			val := make([]uint64, nLen)
			for i := 0; i < nLen; i++ {
				val[i] = uint64(bitstream.ReadInt64(64))
			}
			params[i] = val



		case RPC_BOOL_ARRAY:
			nLen := bitstream.ReadInt(16)
			aa := bool(false)
			tVal := reflect.ArrayOf(nLen, reflect.TypeOf(aa))
			val := reflect.New(tVal).Elem()
			for i := 0; i < nLen; i++ {
				val.Index(i).SetBool(bitstream.ReadFlag())
			}
			params[i] = val.Interface()
		case RPC_STRING_ARRAY:
			nLen := bitstream.ReadInt(16)
			aa := string("")
			tVal := reflect.ArrayOf(nLen, reflect.TypeOf(aa))
			val := reflect.New(tVal).Elem()
			for i := 0; i < nLen; i++ {
				val.Index(i).SetString(bitstream.ReadString())
			}
			params[i] = val.Interface()
		case RPC_FLOAT32_ARRAY:
			nLen := bitstream.ReadInt(16)
			aa := float32(0)
			tVal := reflect.ArrayOf(nLen, reflect.TypeOf(aa))
			val := reflect.New(tVal).Elem()
			for i := 0; i < nLen; i++ {
				val.Index(i).SetFloat(float64(bitstream.ReadFloat()))
			}
			params[i] = val.Interface()
		case RPC_FLOAT64_ARRAY:
			nLen := bitstream.ReadInt(16)
			aa := float64(0)
			tVal := reflect.ArrayOf(nLen, reflect.TypeOf(aa))
			val := reflect.New(tVal).Elem()
			for i := 0; i < nLen; i++ {
				val.Index(i).SetFloat(bitstream.ReadFloat64())
			}
			params[i] = val.Interface()
		case RPC_INT_ARRAY:
			nLen := bitstream.ReadInt(16)
			aa := int(0)
			tVal := reflect.ArrayOf(nLen, reflect.TypeOf(aa))
			val := reflect.New(tVal).Elem()
			for i := 0; i < nLen; i++ {
				val.Index(i).SetInt(int64(bitstream.ReadInt(32)))
			}
			params[i] = val.Interface()
		case RPC_INT8_ARRAY:
			nLen := bitstream.ReadInt(16)
			aa := int8(0)
			tVal := reflect.ArrayOf(nLen, reflect.TypeOf(aa))
			val := reflect.New(tVal).Elem()
			for i := 0; i < nLen; i++ {
				val.Index(i).SetInt(int64(bitstream.ReadInt(8)))
			}
			params[i] = val.Interface()
		case RPC_INT16_ARRAY:
			nLen := bitstream.ReadInt(16)
			aa := int16(0)
			tVal := reflect.ArrayOf(nLen, reflect.TypeOf(aa))
			val := reflect.New(tVal).Elem()
			for i := 0; i < nLen; i++ {
				val.Index(i).SetInt(int64(bitstream.ReadInt(16)))
			}
			params[i] = val.Interface()
		case RPC_INT32_ARRAY:
			nLen := bitstream.ReadInt(16)
			aa := int32(0)
			tVal := reflect.ArrayOf(nLen, reflect.TypeOf(aa))
			val := reflect.New(tVal).Elem()
			for i := 0; i < nLen; i++ {
				val.Index(i).SetInt(int64(bitstream.ReadInt(32)))
			}
			params[i] = val.Interface()
		case RPC_INT64_ARRAY:
			nLen := bitstream.ReadInt(16)
			aa := int64(0)
			tVal := reflect.ArrayOf(nLen, reflect.TypeOf(aa))
			val := reflect.New(tVal).Elem()
			for i := 0; i < nLen; i++ {
				val.Index(i).SetInt(int64(bitstream.ReadInt64(64)))
			}
			params[i] = val.Interface()
		case RPC_UINT_ARRAY:
			nLen := bitstream.ReadInt(16)
			aa := uint(0)
			tVal := reflect.ArrayOf(nLen, reflect.TypeOf(aa))
			val := reflect.New(tVal).Elem()
			for i := 0; i < nLen; i++ {
				val.Index(i).SetInt(int64(bitstream.ReadInt(32)))
			}
			params[i] = val.Interface()
		case RPC_UINT8_ARRAY:
			nLen := bitstream.ReadInt(16)
			aa := uint8(0)
			tVal := reflect.ArrayOf(nLen, reflect.TypeOf(aa))
			val := reflect.New(tVal).Elem()
			for i := 0; i < nLen; i++ {
				val.Index(i).SetInt(int64(bitstream.ReadInt(8)))
			}
			params[i] = val.Interface()
		case RPC_UINT16_ARRAY:
			nLen := bitstream.ReadInt(16)
			aa := uint16(0)
			tVal := reflect.ArrayOf(nLen, reflect.TypeOf(aa))
			val := reflect.New(tVal).Elem()
			for i := 0; i < nLen; i++ {
				val.Index(i).SetInt(int64(bitstream.ReadInt(16)))
			}
			params[i] = val.Interface()
		case RPC_UINT32_ARRAY:
			nLen := bitstream.ReadInt(16)
			aa := uint32(0)
			tVal := reflect.ArrayOf(nLen, reflect.TypeOf(aa))
			val := reflect.New(tVal).Elem()
			for i := 0; i < nLen; i++ {
				val.Index(i).SetInt(int64(bitstream.ReadInt(32)))
			}
			params[i] = val.Interface()
		case RPC_UINT64_ARRAY:
			nLen := bitstream.ReadInt(16)
			aa := uint64(0)
			tVal := reflect.ArrayOf(nLen, reflect.TypeOf(aa))
			val := reflect.New(tVal).Elem()
			for i := 0; i < nLen; i++ {
				val.Index(i).SetInt(int64(bitstream.ReadInt64(64)))
			}
			params[i] = val.Interface()



		case RPC_BOOL_PTR:
			val := new(bool)
			*val = bitstream.ReadFlag()
			params[i] = val
		case RPC_STRING_PTR:
			val := new(string)
			*val = bitstream.ReadString()
			params[i] = val
		case RPC_FLOAT32_PTR:
			val := new(float32)
			*val = bitstream.ReadFloat()
			params[i] = val
		case RPC_FLOAT64_PTR:
			val := new(float64)
			*val = bitstream.ReadFloat64()
			params[i] = val
		case RPC_INT_PTR:
			val := new(int)
			*val = bitstream.ReadInt(32)
			params[i] = val
		case RPC_INT8_PTR:
			val := new(int8)
			*val = int8(bitstream.ReadInt(8))
			params[i] = val
		case RPC_INT16_PTR:
			val := new(int16)
			*val = int16(bitstream.ReadInt(16))
			params[i] = val
		case RPC_INT32_PTR:
			val := new(int32)
			*val = int32(bitstream.ReadInt(32))
			params[i] = val
		case RPC_INT64_PTR:
			val := new(int64)
			*val = int64(bitstream.ReadInt64(64))
			params[i] = val
		case RPC_UINT_PTR:
			val := new(uint)
			*val = uint(bitstream.ReadInt(32))
			params[i] = val
		case RPC_UINT8_PTR:
			val := new(uint8)
			*val = uint8(bitstream.ReadInt(8))
			params[i] = val
		case RPC_UINT16_PTR:
			val := new(uint16)
			*val = uint16(bitstream.ReadInt(16))
			params[i] = val
		case RPC_UINT32_PTR:
			val := new(uint32)
			*val = uint32(bitstream.ReadInt(32))
			params[i] = val
		case RPC_UINT64_PTR:
			val := new(uint64)
			*val = uint64(bitstream.ReadInt64(64))
			params[i] = val



		case RPC_BOOL_PTR_SLICE:
			nLen := bitstream.ReadInt(16)
			val := make([]*bool, nLen)
			for i := 0; i < nLen; i++ {
				val[i] = new(bool)
				*val[i] = bitstream.ReadFlag()
			}
			params[i] = val
		case RPC_STRING_PTR_SLICE:
			nLen := bitstream.ReadInt(16)
			val := make([]*string, nLen)
			for i := 0; i < nLen; i++ {
				val[i] = new(string)
				*val[i] = bitstream.ReadString()
			}
			params[i] = val
		case RPC_FLOAT32_PTR_SLICE:
			nLen := bitstream.ReadInt(16)
			val := make([]*float32, nLen)
			for i := 0; i < nLen; i++ {
				val[i] = new(float32)
				*val[i] = bitstream.ReadFloat()
			}
			params[i] = val
		case RPC_FLOAT64_PTR_SLICE:
			nLen := bitstream.ReadInt(16)
			val := make([]*float64, nLen)
			for i := 0; i < nLen; i++ {
				val[i] = new(float64)
				*val[i] = bitstream.ReadFloat64()
			}
			params[i] = val
		case RPC_INT_PTR_SLICE:
			nLen := bitstream.ReadInt(16)
			val := make([]*int, nLen)
			for i := 0; i < nLen; i++ {
				val[i] = new(int)
				*val[i] = bitstream.ReadInt(32)
			}
			params[i] = val
		case RPC_INT8_PTR_SLICE:
			nLen := bitstream.ReadInt(16)
			val := make([]*int8, nLen)
			for i := 0; i < nLen; i++ {
				val[i] = new(int8)
				*val[i] = int8(bitstream.ReadInt(8))
			}
			params[i] = val
		case RPC_INT16_PTR_SLICE:
			nLen := bitstream.ReadInt(16)
			val := make([]*int16, nLen)
			for i := 0; i < nLen; i++ {
				val[i] = new(int16)
				*val[i] = int16(bitstream.ReadInt(16))
			}
			params[i] = val
		case RPC_INT32_PTR_SLICE:
			nLen := bitstream.ReadInt(16)
			val := make([]*int32, nLen)
			for i := 0; i < nLen; i++ {
				val[i] = new(int32)
				*val[i] = int32(bitstream.ReadInt(32))
			}
			params[i] = val
		case RPC_INT64_PTR_SLICE:
			nLen := bitstream.ReadInt(16)
			val := make([]*int64, nLen)
			for i := 0; i < nLen; i++ {
				val[i] = new(int64)
				*val[i] = int64(bitstream.ReadInt64(64))
			}
			params[i] = val
		case RPC_UINT_PTR_SLICE:
			nLen := bitstream.ReadInt(16)
			val := make([]*uint, nLen)
			for i := 0; i < nLen; i++ {
				val[i] = new(uint)
				*val[i] = uint(bitstream.ReadInt(32))
			}
			params[i] = val
		case RPC_UINT8_PTR_SLICE:
			nLen := bitstream.ReadInt(16)
			val := make([]*uint8, nLen)
			for i := 0; i < nLen; i++ {
				val[i] = new(uint8)
				*val[i] = uint8(bitstream.ReadInt(8))
			}
			params[i] = val
		case RPC_UINT16_PTR_SLICE:
			nLen := bitstream.ReadInt(16)
			val := make([]*uint16, nLen)
			for i := 0; i < nLen; i++ {
				val[i] = new(uint16)
				*val[i] = uint16(bitstream.ReadInt(16))
			}
			params[i] = val
		case RPC_UINT32_PTR_SLICE:
			nLen := bitstream.ReadInt(16)
			val := make([]*uint32, nLen)
			for i := 0; i < nLen; i++ {
				val[i] = new(uint32)
				*val[i] = uint32(bitstream.ReadInt(32))
			}
			params[i] = val
		case RPC_UINT64_PTR_SLICE:
			nLen := bitstream.ReadInt(16)
			val := make([]*uint64, nLen)
			for i := 0; i < nLen; i++ {
				val[i] = new(uint64)
				*val[i] = uint64(bitstream.ReadInt64(64))
			}
			params[i] = val



		case RPC_BOOL_PTR_ARRAY:
			nLen := bitstream.ReadInt(16)
			aa := bool(false)
			tVal := reflect.ArrayOf(nLen, reflect.TypeOf(&aa))
			val := reflect.New(tVal).Elem()
			for i := 0; i < nLen; i++ {
				val1 := bitstream.ReadFlag()
				val.Index(i).Set(reflect.ValueOf(&val1))
			}
			params[i] = val.Interface()
		case RPC_STRING_PTR_ARRAY:
			nLen := bitstream.ReadInt(16)
			aa := string("")
			tVal := reflect.ArrayOf(nLen, reflect.TypeOf(&aa))
			val := reflect.New(tVal).Elem()
			for i := 0; i < nLen; i++ {
				val1 := bitstream.ReadString()
				val.Index(i).Set(reflect.ValueOf(&val1))
			}
			params[i] = val.Interface()
		case RPC_FLOAT32_PTR_ARRAY:
			nLen := bitstream.ReadInt(16)
			aa := float32(0)
			tVal := reflect.ArrayOf(nLen, reflect.TypeOf(&aa))
			val := reflect.New(tVal).Elem()
			for i := 0; i < nLen; i++ {
				val1 := bitstream.ReadFloat()
				val.Index(i).Set(reflect.ValueOf(&val1))
			}
			params[i] = val.Interface()
		case RPC_FLOAT64_PTR_ARRAY:
			nLen := bitstream.ReadInt(16)
			aa := float64(0)
			tVal := reflect.ArrayOf(nLen, reflect.TypeOf(&aa))
			val := reflect.New(tVal).Elem()
			for i := 0; i < nLen; i++ {
				val1 := bitstream.ReadFloat64()
				val.Index(i).Set(reflect.ValueOf(&val1))
			}
			params[i] = val.Interface()
		case RPC_INT_PTR_ARRAY:
			nLen := bitstream.ReadInt(16)
			aa := int(0)
			tVal := reflect.ArrayOf(nLen, reflect.TypeOf(&aa))
			val := reflect.New(tVal).Elem()
			for i := 0; i < nLen; i++ {
				val1 := bitstream.ReadInt(32)
				val.Index(i).Set(reflect.ValueOf(&val1))
			}
			params[i] = val.Interface()
		case RPC_INT8_PTR_ARRAY:
			nLen := bitstream.ReadInt(16)
			aa := int8(0)
			tVal := reflect.ArrayOf(nLen, reflect.TypeOf(&aa))
			val := reflect.New(tVal).Elem()
			for i := 0; i < nLen; i++ {
				val1 := int8(bitstream.ReadInt(8))
				val.Index(i).Set(reflect.ValueOf(&val1))
			}
			params[i] = val.Interface()
		case RPC_INT16_PTR_ARRAY:
			nLen := bitstream.ReadInt(16)
			aa := int16(0)
			tVal := reflect.ArrayOf(nLen, reflect.TypeOf(&aa))
			val := reflect.New(tVal).Elem()
			for i := 0; i < nLen; i++ {
				val1 := int16(bitstream.ReadInt(16))
				val.Index(i).Set(reflect.ValueOf(&val1))
			}
			params[i] = val.Interface()
		case RPC_INT32_PTR_ARRAY:
			nLen := bitstream.ReadInt(16)
			aa := int32(0)
			tVal := reflect.ArrayOf(nLen, reflect.TypeOf(&aa))
			val := reflect.New(tVal).Elem()
			for i := 0; i < nLen; i++ {
				val1 := int32(bitstream.ReadInt(32))
				val.Index(i).Set(reflect.ValueOf(&val1))
			}
			params[i] = val.Interface()
		case RPC_INT64_PTR_ARRAY:
			nLen := bitstream.ReadInt(16)
			aa := int64(0)
			tVal := reflect.ArrayOf(nLen, reflect.TypeOf(&aa))
			val := reflect.New(tVal).Elem()
			for i := 0; i < nLen; i++ {
				val1 := int64(bitstream.ReadInt64(64))
				val.Index(i).Set(reflect.ValueOf(&val1))
			}
			params[i] = val.Interface()
		case RPC_UINT_PTR_ARRAY:
			nLen := bitstream.ReadInt(16)
			aa := uint(0)
			tVal := reflect.ArrayOf(nLen, reflect.TypeOf(&aa))
			val := reflect.New(tVal).Elem()
			for i := 0; i < nLen; i++ {
				val1 := uint(bitstream.ReadInt(32))
				val.Index(i).Set(reflect.ValueOf(&val1))
			}
			params[i] = val.Interface()
		case RPC_UINT8_PTR_ARRAY:
			nLen := bitstream.ReadInt(16)
			aa := uint8(0)
			tVal := reflect.ArrayOf(nLen, reflect.TypeOf(&aa))
			val := reflect.New(tVal).Elem()
			for i := 0; i < nLen; i++ {
				val1 := uint8(bitstream.ReadInt(8))
				val.Index(i).Set(reflect.ValueOf(&val1))
			}
			params[i] = val.Interface()
		case RPC_UINT16_PTR_ARRAY:
			nLen := bitstream.ReadInt(16)
			aa := uint16(0)
			tVal := reflect.ArrayOf(nLen, reflect.TypeOf(&aa))
			val := reflect.New(tVal).Elem()
			for i := 0; i < nLen; i++ {
				val1 := uint16(bitstream.ReadInt(16))
				val.Index(i).Set(reflect.ValueOf(&val1))
			}
			params[i] = val.Interface()
		case RPC_UINT32_PTR_ARRAY:
			nLen := bitstream.ReadInt(16)
			aa := uint32(0)
			tVal := reflect.ArrayOf(nLen, reflect.TypeOf(&aa))
			val := reflect.New(tVal).Elem()
			for i := 0; i < nLen; i++ {
				val1 := uint32(bitstream.ReadInt(32))
				val.Index(i).Set(reflect.ValueOf(&val1))
			}
			params[i] = val.Interface()
		case RPC_UINT64_PTR_ARRAY:
			nLen := bitstream.ReadInt(16)
			aa := uint64(0)
			tVal := reflect.ArrayOf(nLen, reflect.TypeOf(&aa))
			val := reflect.New(tVal).Elem()
			for i := 0; i < nLen; i++ {
				val1 := uint64(bitstream.ReadInt64(64))
				val.Index(i).Set(reflect.ValueOf(&val1))
			}
			params[i] = val.Interface()


		case RPC_BOOL_SLICE_PTR:
			nLen := bitstream.ReadInt(16)
			val := make([]bool, nLen)
			for i := 0; i < nLen; i++ {
				val[i] = bitstream.ReadFlag()
			}
			params[i] = &val
		case RPC_STRING_SLICE_PTR:
			nLen := bitstream.ReadInt(16)
			val := make([]string, nLen)
			for i := 0; i < nLen; i++ {
				val[i] = bitstream.ReadString()
			}
			params[i] = &val
		case RPC_FLOAT32_SLICE_PTR:
			nLen := bitstream.ReadInt(16)
			val := make([]float32, nLen)
			for i := 0; i < nLen; i++ {
				val[i] = bitstream.ReadFloat()
			}
			params[i] = &val
		case RPC_FLOAT64_SLICE_PTR:
			nLen := bitstream.ReadInt(16)
			val := make([]float64, nLen)
			for i := 0; i < nLen; i++ {
				val[i] = bitstream.ReadFloat64()
			}
			params[i] = &val
		case RPC_INT_SLICE_PTR:
			nLen := bitstream.ReadInt(16)
			val := make([]int, nLen)
			for i := 0; i < nLen; i++ {
				val[i] = bitstream.ReadInt(32)
			}
			params[i] = &val
		case RPC_INT8_SLICE_PTR:
			nLen := bitstream.ReadInt(16)
			val := make([]int8, nLen)
			for i := 0; i < nLen; i++ {
				val[i] = int8(bitstream.ReadInt(8))
			}
			params[i] = &val
		case RPC_INT16_SLICE_PTR:
			nLen := bitstream.ReadInt(16)
			val := make([]int16, nLen)
			for i := 0; i < nLen; i++ {
				val[i] = int16(bitstream.ReadInt(16))
			}
			params[i] = &val
		case RPC_INT32_SLICE_PTR:
			nLen := bitstream.ReadInt(16)
			val := make([]int32, nLen)
			for i := 0; i < nLen; i++ {
				val[i] = int32(bitstream.ReadInt(32))
			}
			params[i] = &val
		case RPC_INT64_SLICE_PTR:
			nLen := bitstream.ReadInt(16)
			val := make([]int64, nLen)
			for i := 0; i < nLen; i++ {
				val[i] = int64(bitstream.ReadInt64(64))
			}
			params[i] = &val
		case RPC_UINT_SLICE_PTR:
			nLen := bitstream.ReadInt(16)
			val := make([]uint, nLen)
			for i := 0; i < nLen; i++ {
				val[i] = uint(bitstream.ReadInt(32))
			}
			params[i] = &val
		case RPC_UINT8_SLICE_PTR:
			nLen := bitstream.ReadInt(16)
			val := make([]uint8, nLen)
			for i := 0; i < nLen; i++ {
				val[i] = uint8(bitstream.ReadInt(8))
			}
			params[i] = &val
		case RPC_UINT16_SLICE_PTR:
			nLen := bitstream.ReadInt(16)
			val := make([]uint16, nLen)
			for i := 0; i < nLen; i++ {
				val[i] = uint16(bitstream.ReadInt(16))
			}
			params[i] = &val
		case RPC_UINT32_SLICE_PTR:
			nLen := bitstream.ReadInt(16)
			val := make([]uint32, nLen)
			for i := 0; i < nLen; i++ {
				val[i] = uint32(bitstream.ReadInt(32))
			}
			params[i] = &val
		case RPC_UINT64_SLICE_PTR:
			nLen := bitstream.ReadInt(16)
			val := make([]uint64, nLen)
			for i := 0; i < nLen; i++ {
				val[i] = uint64(bitstream.ReadInt64(64))
			}
			params[i] = &val



		case RPC_BOOL_PTR_SLICE_PTR:
			nLen := bitstream.ReadInt(16)
			val := make([]*bool, nLen)
			for i := 0; i < nLen; i++ {
				val[i] = new(bool)
				*val[i] = bitstream.ReadFlag()
			}
			params[i] = &val
		case RPC_STRING_PTR_SLICE_PTR:
			nLen := bitstream.ReadInt(16)
			val := make([]*string, nLen)
			for i := 0; i < nLen; i++ {
				val[i] = new(string)
				*val[i] = bitstream.ReadString()
			}
			params[i] = &val
		case RPC_FLOAT32_PTR_SLICE_PTR:
			nLen := bitstream.ReadInt(16)
			val := make([]*float32, nLen)
			for i := 0; i < nLen; i++ {
				val[i] = new(float32)
				*val[i] = bitstream.ReadFloat()
			}
			params[i] = &val
		case RPC_FLOAT64_PTR_SLICE_PTR:
			nLen := bitstream.ReadInt(16)
			val := make([]*float64, nLen)
			for i := 0; i < nLen; i++ {
				val[i] = new(float64)
				*val[i] = bitstream.ReadFloat64()
			}
			params[i] = &val
		case RPC_INT_PTR_SLICE_PTR:
			nLen := bitstream.ReadInt(16)
			val := make([]*int, nLen)
			for i := 0; i < nLen; i++ {
				val[i] = new(int)
				*val[i] = bitstream.ReadInt(32)
			}
			params[i] = &val
		case RPC_INT8_PTR_SLICE_PTR:
			nLen := bitstream.ReadInt(16)
			val := make([]*int8, nLen)
			for i := 0; i < nLen; i++ {
				val[i] = new(int8)
				*val[i] = int8(bitstream.ReadInt(8))
			}
			params[i] = &val
		case RPC_INT16_PTR_SLICE_PTR:
			nLen := bitstream.ReadInt(16)
			val := make([]*int16, nLen)
			for i := 0; i < nLen; i++ {
				val[i] = new(int16)
				*val[i] = int16(bitstream.ReadInt(16))
			}
			params[i] = &val
		case RPC_INT32_PTR_SLICE_PTR:
			nLen := bitstream.ReadInt(16)
			val := make([]*int32, nLen)
			for i := 0; i < nLen; i++ {
				val[i] = new(int32)
				*val[i] = int32(bitstream.ReadInt(32))
			}
			params[i] = &val
		case RPC_INT64_PTR_SLICE_PTR:
			nLen := bitstream.ReadInt(16)
			val := make([]*int64, nLen)
			for i := 0; i < nLen; i++ {
				val[i] = new(int64)
				*val[i] = int64(bitstream.ReadInt64(64))
			}
			params[i] = &val
		case RPC_UINT_PTR_SLICE_PTR:
			nLen := bitstream.ReadInt(16)
			val := make([]*uint, nLen)
			for i := 0; i < nLen; i++ {
				val[i] = new(uint)
				*val[i] = uint(bitstream.ReadInt(32))
			}
			params[i] = &val
		case RPC_UINT8_PTR_SLICE_PTR:
			nLen := bitstream.ReadInt(16)
			val := make([]*uint8, nLen)
			for i := 0; i < nLen; i++ {
				val[i] = new(uint8)
				*val[i] = uint8(bitstream.ReadInt(8))
			}
			params[i] = &val
		case RPC_UINT16_PTR_SLICE_PTR:
			nLen := bitstream.ReadInt(16)
			val := make([]*uint16, nLen)
			for i := 0; i < nLen; i++ {
				val[i] = new(uint16)
				*val[i] = uint16(bitstream.ReadInt(16))
			}
			params[i] = &val
		case RPC_UINT32_PTR_SLICE_PTR:
			nLen := bitstream.ReadInt(16)
			val := make([]*uint32, nLen)
			for i := 0; i < nLen; i++ {
				val[i] = new(uint32)
				*val[i] = uint32(bitstream.ReadInt(32))
			}
			params[i] = &val
		case RPC_UINT64_PTR_SLICE_PTR:
			nLen := bitstream.ReadInt(16)
			val := make([]*uint64, nLen)
			for i := 0; i < nLen; i++ {
				val[i] = new(uint64)
				*val[i] = uint64(bitstream.ReadInt64(64))
			}
			params[i] = &val



		case RPC_BOOL_ARRAY_PTR:
			nLen := bitstream.ReadInt(16)
			aa := bool(false)
			tVal := reflect.ArrayOf(nLen, reflect.TypeOf(aa))
			val := reflect.New(tVal)
			for i := 0; i < nLen; i++ {
				val.Elem().Index(i).SetBool(bitstream.ReadFlag())
			}
			params[i] = val.Interface()
		case RPC_STRING_ARRAY_PTR:
			nLen := bitstream.ReadInt(16)
			aa := string("")
			tVal := reflect.ArrayOf(nLen, reflect.TypeOf(aa))
			val := reflect.New(tVal)
			for i := 0; i < nLen; i++ {
				val.Elem().Index(i).SetString(bitstream.ReadString())
			}
			params[i] = val.Interface()
		case RPC_FLOAT32_ARRAY_PTR:
			nLen := bitstream.ReadInt(16)
			aa := float32(0)
			tVal := reflect.ArrayOf(nLen, reflect.TypeOf(aa))
			val := reflect.New(tVal)
			for i := 0; i < nLen; i++ {
				val.Elem().Index(i).SetFloat(float64(bitstream.ReadFloat()))
			}
			params[i] = val.Interface()
		case RPC_FLOAT64_ARRAY_PTR:
			nLen := bitstream.ReadInt(16)
			aa := float64(0)
			tVal := reflect.ArrayOf(nLen, reflect.TypeOf(aa))
			val := reflect.New(tVal)
			for i := 0; i < nLen; i++ {
				val.Elem().Index(i).SetFloat(bitstream.ReadFloat64())
			}
			params[i] = val.Interface()
		case RPC_INT_ARRAY_PTR:
			nLen := bitstream.ReadInt(16)
			aa := int(0)
			tVal := reflect.ArrayOf(nLen, reflect.TypeOf(aa))
			val := reflect.New(tVal)
			for i := 0; i < nLen; i++ {
				val.Elem().Index(i).SetInt(int64(bitstream.ReadInt(32)))
			}
			params[i] = val.Interface()
		case RPC_INT8_ARRAY_PTR:
			nLen := bitstream.ReadInt(16)
			aa := int8(0)
			tVal := reflect.ArrayOf(nLen, reflect.TypeOf(aa))
			val := reflect.New(tVal)
			for i := 0; i < nLen; i++ {
				val.Elem().Index(i).SetInt(int64(bitstream.ReadInt(8)))
			}
			params[i] = val.Interface()
		case RPC_INT16_ARRAY_PTR:
			nLen := bitstream.ReadInt(16)
			aa := int16(0)
			tVal := reflect.ArrayOf(nLen, reflect.TypeOf(aa))
			val := reflect.New(tVal)
			for i := 0; i < nLen; i++ {
				val.Elem().Index(i).SetInt(int64(bitstream.ReadInt(16)))
			}
			params[i] = val.Interface()
		case RPC_INT32_ARRAY_PTR:
			nLen := bitstream.ReadInt(16)
			aa := int32(0)
			tVal := reflect.ArrayOf(nLen, reflect.TypeOf(aa))
			val := reflect.New(tVal)
			for i := 0; i < nLen; i++ {
				val.Elem().Index(i).SetInt(int64(bitstream.ReadInt(32)))
			}
			params[i] = val.Interface()
		case RPC_INT64_ARRAY_PTR:
			nLen := bitstream.ReadInt(16)
			aa := int64(0)
			tVal := reflect.ArrayOf(nLen, reflect.TypeOf(aa))
			val := reflect.New(tVal)
			for i := 0; i < nLen; i++ {
				val.Elem().Index(i).SetInt(int64(bitstream.ReadInt64(64)))
			}
			params[i] = val.Interface()
		case RPC_UINT_ARRAY_PTR:
			nLen := bitstream.ReadInt(16)
			aa := uint(0)
			tVal := reflect.ArrayOf(nLen, reflect.TypeOf(aa))
			val := reflect.New(tVal)
			for i := 0; i < nLen; i++ {
				val.Elem().Index(i).SetInt(int64(bitstream.ReadInt(32)))
			}
			params[i] = val.Interface()
		case RPC_UINT8_ARRAY_PTR:
			nLen := bitstream.ReadInt(16)
			aa := uint8(0)
			tVal := reflect.ArrayOf(nLen, reflect.TypeOf(aa))
			val := reflect.New(tVal)
			for i := 0; i < nLen; i++ {
				val.Elem().Index(i).SetInt(int64(bitstream.ReadInt(8)))
			}
			params[i] = val.Interface()
		case RPC_UINT16_ARRAY_PTR:
			nLen := bitstream.ReadInt(16)
			aa := uint16(0)
			tVal := reflect.ArrayOf(nLen, reflect.TypeOf(aa))
			val := reflect.New(tVal)
			for i := 0; i < nLen; i++ {
				val.Elem().Index(i).SetInt(int64(bitstream.ReadInt(16)))
			}
			params[i] = val.Interface()
		case RPC_UINT32_ARRAY_PTR:
			nLen := bitstream.ReadInt(16)
			aa := uint32(0)
			tVal := reflect.ArrayOf(nLen, reflect.TypeOf(aa))
			val := reflect.New(tVal)
			for i := 0; i < nLen; i++ {
				val.Elem().Index(i).SetInt(int64(bitstream.ReadInt(32)))
			}
			params[i] = val.Interface()
		case RPC_UINT64_ARRAY_PTR:
			nLen := bitstream.ReadInt(16)
			aa := uint64(0)
			tVal := reflect.ArrayOf(nLen, reflect.TypeOf(aa))
			val := reflect.New(tVal)
			for i := 0; i < nLen; i++ {
				val.Elem().Index(i).SetInt(int64(bitstream.ReadInt64(64)))
			}
			params[i] = val.Interface()



		case RPC_BOOL_PTR_ARRAY_PTR:
			nLen := bitstream.ReadInt(16)
			aa := bool(false)
			tVal := reflect.ArrayOf(nLen, reflect.TypeOf(&aa))
			val := reflect.New(tVal).Elem()
			for i := 0; i < nLen; i++ {
				val1 := bitstream.ReadFlag()
				val.Index(i).Set(reflect.ValueOf(&val1))
			}
			params[i] = val.Addr().Interface()
		case RPC_STRING_PTR_ARRAY_PTR:
			nLen := bitstream.ReadInt(16)
			aa := string("")
			tVal := reflect.ArrayOf(nLen, reflect.TypeOf(&aa))
			val := reflect.New(tVal).Elem()
			for i := 0; i < nLen; i++ {
				val1 := bitstream.ReadString()
				val.Index(i).Set(reflect.ValueOf(&val1))
			}
			params[i] = val.Addr().Interface()
		case RPC_FLOAT32_PTR_ARRAY_PTR:
			nLen := bitstream.ReadInt(16)
			aa := float32(0)
			tVal := reflect.ArrayOf(nLen, reflect.TypeOf(&aa))
			val := reflect.New(tVal).Elem()
			for i := 0; i < nLen; i++ {
				val1 := bitstream.ReadFloat()
				val.Index(i).Set(reflect.ValueOf(&val1))
			}
			params[i] = val.Addr().Interface()
		case RPC_FLOAT64_PTR_ARRAY_PTR:
			nLen := bitstream.ReadInt(16)
			aa := float64(0)
			tVal := reflect.ArrayOf(nLen, reflect.TypeOf(&aa))
			val := reflect.New(tVal).Elem()
			for i := 0; i < nLen; i++ {
				val1 := bitstream.ReadFloat64()
				val.Index(i).Set(reflect.ValueOf(&val1))
			}
			params[i] = val.Addr().Interface()
		case RPC_INT_PTR_ARRAY_PTR:
			nLen := bitstream.ReadInt(16)
			aa := int(0)
			tVal := reflect.ArrayOf(nLen, reflect.TypeOf(&aa))
			val := reflect.New(tVal).Elem()
			for i := 0; i < nLen; i++ {
				val1 := bitstream.ReadInt(32)
				val.Index(i).Set(reflect.ValueOf(&val1))
			}
			params[i] = val.Addr().Interface()
		case RPC_INT8_PTR_ARRAY_PTR:
			nLen := bitstream.ReadInt(16)
			aa := int8(0)
			tVal := reflect.ArrayOf(nLen, reflect.TypeOf(&aa))
			val := reflect.New(tVal).Elem()
			for i := 0; i < nLen; i++ {
				val1 := int8(bitstream.ReadInt(8))
				val.Index(i).Set(reflect.ValueOf(&val1))
			}
			params[i] = val.Addr().Interface()
		case RPC_INT16_PTR_ARRAY_PTR:
			nLen := bitstream.ReadInt(16)
			aa := int16(0)
			tVal := reflect.ArrayOf(nLen, reflect.TypeOf(&aa))
			val := reflect.New(tVal).Elem()
			for i := 0; i < nLen; i++ {
				val1 := int16(bitstream.ReadInt(16))
				val.Index(i).Set(reflect.ValueOf(&val1))
			}
			params[i] = val.Addr().Interface()
		case RPC_INT32_PTR_ARRAY_PTR:
			nLen := bitstream.ReadInt(16)
			aa := int32(0)
			tVal := reflect.ArrayOf(nLen, reflect.TypeOf(&aa))
			val := reflect.New(tVal).Elem()
			for i := 0; i < nLen; i++ {
				val1 := int32(bitstream.ReadInt(32))
				val.Index(i).Set(reflect.ValueOf(&val1))
			}
			params[i] = val.Addr().Interface()
		case RPC_INT64_PTR_ARRAY_PTR:
			nLen := bitstream.ReadInt(16)
			aa := int64(0)
			tVal := reflect.ArrayOf(nLen, reflect.TypeOf(&aa))
			val := reflect.New(tVal).Elem()
			for i := 0; i < nLen; i++ {
				val1 := int64(bitstream.ReadInt64(64))
				val.Index(i).Set(reflect.ValueOf(&val1))
			}
			params[i] = val.Addr().Interface()
		case RPC_UINT_PTR_ARRAY_PTR:
			nLen := bitstream.ReadInt(16)
			aa := uint(0)
			tVal := reflect.ArrayOf(nLen, reflect.TypeOf(&aa))
			val := reflect.New(tVal).Elem()
			for i := 0; i < nLen; i++ {
				val1 := uint(bitstream.ReadInt(32))
				val.Index(i).Set(reflect.ValueOf(&val1))
			}
			params[i] = val.Addr().Interface()
		case RPC_UINT8_PTR_ARRAY_PTR:
			nLen := bitstream.ReadInt(16)
			aa := uint8(0)
			tVal := reflect.ArrayOf(nLen, reflect.TypeOf(&aa))
			val := reflect.New(tVal).Elem()
			for i := 0; i < nLen; i++ {
				val1 := uint8(bitstream.ReadInt(8))
				val.Index(i).Set(reflect.ValueOf(&val1))
			}
			params[i] = val.Addr().Interface()
		case RPC_UINT16_PTR_ARRAY_PTR:
			nLen := bitstream.ReadInt(16)
			aa := uint16(0)
			tVal := reflect.ArrayOf(nLen, reflect.TypeOf(&aa))
			val := reflect.New(tVal).Elem()
			for i := 0; i < nLen; i++ {
				val1 := uint16(bitstream.ReadInt(16))
				val.Index(i).Set(reflect.ValueOf(&val1))
			}
			params[i] = val.Addr().Interface()
		case RPC_UINT32_PTR_ARRAY_PTR:
			nLen := bitstream.ReadInt(16)
			aa := uint32(0)
			tVal := reflect.ArrayOf(nLen, reflect.TypeOf(&aa))
			val := reflect.New(tVal).Elem()
			for i := 0; i < nLen; i++ {
				val1 := uint32(bitstream.ReadInt(32))
				val.Index(i).Set(reflect.ValueOf(&val1))
			}
			params[i] = val.Addr().Interface()
		case RPC_UINT64_PTR_ARRAY_PTR:
			nLen := bitstream.ReadInt(16)
			aa := uint64(0)
			tVal := reflect.ArrayOf(nLen, reflect.TypeOf(&aa))
			val := reflect.New(tVal).Elem()
			for i := 0; i < nLen; i++ {
				val1 := uint64(bitstream.ReadInt64(64))
				val.Index(i).Set(reflect.ValueOf(&val1))
			}
			params[i] = val.Addr().Interface()



		case RPC_MESSAGE://protobuf
			packet, err := unmarshalPB(bitstream)
			if err == nil{
				params[i] = packet
			}



		case RPC_GOB://gob
			nLen := bitstream.ReadInt(32)
			packetBuf := bitstream.ReadBits(nLen << 3)

			if pFuncType != nil{
				if i < pFuncType.NumIn() {
					val := reflect.New(pFuncType.In(i))
					json := jsoniter.ConfigCompatibleWithStandardLibrary
					err := json.Unmarshal(packetBuf, val.Interface())
					/*buf := bytes.NewBuffer(packetBuf)
					enc := gob.NewDecoder(buf)
					err := enc.DecodeValue(val)*/
					if err == nil{
						params[i] = val.Elem().Interface()
					}
				}
			}

		default:
			panic("func [%s] params type not supported")
		}
	}
	return params
}

//rpc  UnmarshalPB
func unmarshalPB(bitstream *base.BitStream) (proto.Message, error) {
	packetName := bitstream.ReadString()
	nLen := bitstream.ReadInt(32)
	packetBuf := bitstream.ReadBits(nLen << 3)
	packet := reflect.New(proto.MessageType(packetName).Elem()).Interface().(proto.Message)
	err := proto.Unmarshal(packetBuf, packet)
	return  packet, err
}