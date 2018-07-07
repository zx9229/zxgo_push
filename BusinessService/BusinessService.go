package businessservice

import (
	"errors"
	"fmt"
	"log"
	"reflect"

	txstruct "github.com/zx9229/zxgo_push/TxStruct"
	wsconnectionmanager "github.com/zx9229/zxgo_push/WSConnectionManager"
)

//BusinessService omit
type BusinessService struct {
	methods map[string]reflect.Value
	parser  *txstruct.TxParser
	manager *wsconnectionmanager.WSConnectionManager
}

//New_BusinessService omit
func New_BusinessService() *BusinessService {
	curData := new(BusinessService)
	//
	curData.methods = make(map[string]reflect.Value)
	curData.parser = txstruct.New_TxParser()
	curData.manager = wsconnectionmanager.New_WSConnectionManager()
	//
	curData.methods = curData.calcMethods()
	curData.manager.CbConnected = curData.handleConnected
	curData.manager.CbDisconnected = curData.handleDisconnected
	curData.manager.CbReceive = curData.handleReceive
	//
	return curData
}

//GetConnectionManager omit
func (thls *BusinessService) GetConnectionManager() *wsconnectionmanager.WSConnectionManager {
	return thls.manager
}

func (thls *BusinessService) calcMethods() map[string]reflect.Value {
	calcMethods := make(map[string]reflect.Value)
	vf := reflect.ValueOf(thls)
	vft := vf.Type()
	mNum := vf.NumMethod()
	for i := 0; i < mNum; i++ {
		mName := vft.Method(i).Name
		calcMethods[mName] = vf.Method(i)
	}
	return calcMethods
}

//HandleConnected omit
func (thls *BusinessService) handleConnected(conn *wsconnectionmanager.WSConnection) {
	log.Println(fmt.Sprintf("[   Connected][%p]", conn))
}

//HandleDisconnected omit
func (thls *BusinessService) handleDisconnected(conn *wsconnectionmanager.WSConnection, err error) {
	log.Println(fmt.Sprintf("[Disconnected][%p]err=%v", conn, err))
}

var (
	ErrParseDataFail       = errors.New("can not parse data")
	ErrFindNotMethod       = errors.New("can not find corresponding method")
	ErrRetValAnomalous     = errors.New("call method and return value anomalous")
	ErrConvertTxStructFail = errors.New("return value convert to tx struct fail")
)

//HandleReceive omit
func (thls *BusinessService) handleReceive(conn *wsconnectionmanager.WSConnection, bytes []byte) {
	var err error
	var responseMessage string

	for range "1" {
		var objData txstruct.TxInterface

		if objData, _, err = thls.parser.ParseByteSlice(bytes); err != nil {
			notice := txstruct.UnknownNotice{Message: ErrParseDataFail.Error(), RawMessage: string(bytes)}
			notice.CALC_TN(true)
			responseMessage = notice.TO_JSON(true)
			break
		}

		var rspData txstruct.TxInterface
		for range "1" {
			var ok bool
			var method reflect.Value
			if method, ok = thls.methods[objData.GET_TN()]; !ok {
				err = ErrFindNotMethod
				break
			}
			sliceIn := []reflect.Value{reflect.ValueOf(objData)}
			sliceOut := method.Call(sliceIn)
			if len(sliceOut) != 1 {
				err = ErrRetValAnomalous
				break
			}
			if rspData, ok = sliceOut[0].Interface().(txstruct.TxInterface); !ok {
				err = ErrConvertTxStructFail
				break
			}
		}
		if rspData == nil {
			notice := txstruct.UnknownNotice{Message: err.Error(), RawMessage: string(bytes)}
			notice.CALC_TN(true)
			responseMessage = notice.TO_JSON(true)
			break
		}

		responseMessage = rspData.TO_JSON(true)
	}

	if err = conn.Send(responseMessage); err != nil {
		conn.Close()
		log.Println(fmt.Sprintf("发送消息失败了"))
	}
}

//LoginReq omit
func (thls *BusinessService) LoginReq(dataReq *txstruct.LoginReq) *txstruct.LoginRsp {
	rspData := new(txstruct.LoginRsp)
	rspData.CALC_TN(true)
	rspData.BaseDataRsp.Code = 0
	rspData.BaseDataRsp.Message = "SUCCESS"
	rspData.ReqData = dataReq
	return rspData
}
