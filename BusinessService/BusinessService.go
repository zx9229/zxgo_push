package businessservice

import (
	"encoding/json"
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
	cache   *CacheData
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
	curData.cache = New_CacheData_Original()
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

const (
	ErrMsgEmpty             = ""
	ErrMsgSUCCESS           = "SUCCESS"
	ErrMsgUserIdNotExist    = "user id not exist"
	ErrMsgIncorrectPassword = "incorrect password"
	ErrMsgInvalidLoginType  = "invalid login type"
	ErrMsgUserHasLoggedIn   = "user has logged in"
	ErrMsgInvalidCategory   = "invalid category"
)

//LoginReq omit
func (thls *BusinessService) LoginReq(conn *wsconnectionmanager.WSConnection, req *txstruct.LoginReq) *txstruct.LoginRsp {
	rsp := new(txstruct.LoginRsp)
	rsp.CALC_TN(true)
	rsp.BaseDataRsp.Code = 0
	rsp.BaseDataRsp.Message = ErrMsgEmpty
	rsp.ReqData = req

	for range "1" {
		if int64(len(thls.cache.AllUser)) < req.UserID {
			rsp.BaseDataRsp.Code = 1
			rsp.BaseDataRsp.Message = ErrMsgUserIdNotExist
			break
		}
		if thls.cache.AllUser[req.UserID].Base.Password != req.Password {
			rsp.BaseDataRsp.Code = 1
			rsp.BaseDataRsp.Message = ErrMsgIncorrectPassword
			break
		}
		if (LoginTypeDEFAULT < req.Way && req.Way < LoginTypeEND) == false {
			rsp.BaseDataRsp.Code = 1
			rsp.BaseDataRsp.Message = ErrMsgInvalidLoginType
			break
		}
		curLoginInfo := thls.cache.AllUser[req.UserID].State[req.Way]
		if curLoginInfo.conn != nil && !req.ForceLogin {
			rsp.BaseDataRsp.Code = 1
			rsp.BaseDataRsp.Message = ErrMsgUserHasLoggedIn
			break
		}
		if curLoginInfo.conn != nil {
			//TODO:关闭之前,发送一个"您被踢下线了"的消息
			curLoginInfo.conn.Close()
		}

		curLoginInfo.conn = conn
		//TODO:给连接附加登录信息的结构体指针

		if 0 <= req.LastMsgID {
			curLoginInfo.LastRecvID = req.LastMsgID
			//TODO:哪些消息尚未推送,把它们推送过去
		}
	}

	if rsp.Message == ErrMsgEmpty {
		rsp.BaseDataRsp.Code = 0
		rsp.BaseDataRsp.Message = ErrMsgSUCCESS
	} else {
		rsp.BaseDataRsp.Code = 1
	}

	return rsp
}

//ReportReq omit
func (thls *BusinessService) ReportReq(conn *wsconnectionmanager.WSConnection, req *txstruct.ReportReq) *txstruct.ReportRsp {
	rsp := new(txstruct.ReportRsp)
	rsp.CALC_TN(true)
	rsp.BaseDataRsp.Code = 0
	rsp.BaseDataRsp.Message = ErrMsgEmpty
	rsp.ReqData = req

	for range "1" {
		//TODO:用户未登录,就校验随附密码,失败就拒绝
		if !thls.cache.SubBase.IsRegistered(req.Category) {
			rsp.BaseDataRsp.Message = ErrMsgInvalidCategory
			break
		}

		reportData := convertToReportData(req)
		reportData.ID = thls.cache.LastPushID + 1
		thls.cache.LastPushID = reportData.ID

		//TODO:待优化,抽离成订阅管理器
		if true {
			jsonByte, _ := json.Marshal(reportData)
			jsonStr := string(jsonByte)

			for _, userData := range thls.cache.AllUser {
				if userData == nil {
					continue
				}
				if !userData.SubInfo.ShouldSend(reportData.UserID, reportData.Category) {
					continue
				}
				for _, stateData := range userData.State {
					if stateData.conn == nil {
						continue
					}
					stateData.conn.Send(jsonStr)
				}
			}
		}
	}

	if rsp.BaseDataRsp.Message == ErrMsgEmpty {
		rsp.BaseDataRsp.Code = 0
		rsp.BaseDataRsp.Message = ErrMsgSUCCESS
	} else {
		rsp.BaseDataRsp.Code = 1
	}

	return rsp
}
