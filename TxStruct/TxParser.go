package txstruct

import (
	"encoding/json"
	"errors"
	"reflect"
)

//TxParser 通信结构体的解析器
type TxParser struct {
	mapStr2Type map[string]reflect.Type
}

//New_TxParser omit
func New_TxParser() *TxParser {
	newData := new(TxParser)
	newData.mapStr2Type = CalcMapStr2Type()
	return newData
}

//CalcMapStr2Type omit
func CalcMapStr2Type() map[string]reflect.Type {
	_, sliceData := inner_check_by_compile()

	cacheData := map[string]reflect.Type{}
	for _, element := range sliceData {
		curType := reflect.ValueOf(element).Type()
		cacheData[curType.Name()] = curType
	}

	return cacheData
}

//ParseString omit
func (thls *TxParser) ParseString(jsonStr string) (objData interface{}, objType reflect.Type, err error) {
	return thls.ParseByteSlice([]byte(jsonStr))
}

var (
	ErrFindNotTypeByName         = errors.New("can not find type by type name")
	ErrConvertToTransmitDataFail = errors.New("convert to transmit data fail")
)

//ParseByteSlice 返回值含义  objData:反序列化jsonByte后,得到的对象; objType:对象的类型; err:错误的详细情况.
func (thls *TxParser) ParseByteSlice(jsonByte []byte) (objData TxInterface, objType reflect.Type, err error) {
	objData = nil
	objType = nil
	err = nil

	baseData := &BaseDataTx{}
	if err = json.Unmarshal(jsonByte, baseData); err != nil {
		return
	}

	var ok bool
	if objType, ok = thls.mapStr2Type[baseData.GetTN()]; !ok {
		err = ErrFindNotTypeByName
		return
	}

	if objData, ok = reflect.New(objType).Interface().(TxInterface); !ok {
		err = ErrConvertToTransmitDataFail
		return
	}

	if err = json.Unmarshal(jsonByte, objData); err != nil {
		return
	}

	return
}
