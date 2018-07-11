package businessservice

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"reflect"

	txstruct "github.com/zx9229/zxgo_push/TxStruct"
	wscmanager "github.com/zx9229/zxgo_push/WSCManager"
)

//BusinessService omit
type BusinessService struct {
	methods    map[string]reflect.Value
	parser     *txstruct.TxParser
	manager    *wscmanager.WSConnectionManager
	cache      *TotalUserManager
	subBase    *SubscribeBaseInfo
	LastPushID int64 //最后一个推送消息的序号
}

//New_BusinessService omit
func New_BusinessService() *BusinessService {
	curData := new(BusinessService)
	//
	curData.methods = make(map[string]reflect.Value)
	curData.parser = txstruct.New_TxParser()
	curData.manager = wscmanager.New_WSConnectionManager()
	//
	curData.methods = curData.calcMethods()
	curData.manager.CbConnected = curData.handleConnected
	curData.manager.CbDisconnected = curData.handleDisconnected
	curData.manager.CbReceive = curData.handleReceive
	//
	curData.cache = New_TotalUserManager()
	curData.subBase = New_SubscribeBaseInfo()
	curData.subBase.AddCategory("cat") //TODO:临时调试代码
	//
	return curData
}

//GetConnectionManager omit
func (thls *BusinessService) GetConnectionManager() *wscmanager.WSConnectionManager {
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
func (thls *BusinessService) handleConnected(conn *wscmanager.WSConnection) {
	log.Println(fmt.Sprintf("[   Connected][%p]", conn))
}

//HandleDisconnected omit
func (thls *BusinessService) handleDisconnected(conn *wscmanager.WSConnection, err error) {
	log.Println(fmt.Sprintf("[Disconnected][%p]err=%v", conn, err))
	if conn.ExtraData != nil {
		tmpData := conn.ExtraData.(*UserTempData)
		tmpData.state.conn = nil
	}
}

var (
	ErrParseDataFail       = errors.New("can not parse data")
	ErrFindNotMethod       = errors.New("can not find corresponding method")
	ErrRetValAnomalous     = errors.New("call method and return value anomalous")
	ErrConvertTxStructFail = errors.New("return value convert to tx struct fail")
)

//HandleReceive omit
func (thls *BusinessService) handleReceive(conn *wscmanager.WSConnection, bytes []byte) {
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
			sliceIn := []reflect.Value{reflect.ValueOf(conn), reflect.ValueOf(objData)}
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
	ErrMsgEmpty                 = ""
	ErrMsgSUCCESS               = "SUCCESS"
	ErrMsgUserIdNotExist        = "user id not exist"
	ErrMsgIncorrectPassword     = "incorrect password"
	ErrMsgInvalidLoginType      = "invalid login type"
	ErrMsgUserHasLoggedIn       = "user has logged in"
	ErrMsgInvalidCategory       = "invalid category"
	ErrMsgNotLoginAndOncePwdErr = "not login and once password error"
)

//LoginReq omit
func (thls *BusinessService) LoginReq(conn *wscmanager.WSConnection, req *txstruct.LoginReq) *txstruct.LoginRsp {
	rsp := new(txstruct.LoginRsp)
	rsp.CALC_TN(true)
	rsp.BaseDataRsp.Code = 0
	rsp.BaseDataRsp.Message = ErrMsgEmpty
	rsp.ReqData = req

	if err := thls.cache.LoginUser(conn, req); err != nil {
		rsp.BaseDataRsp.Message = err.Error()
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
func (thls *BusinessService) ReportReq(conn *wscmanager.WSConnection, req *txstruct.ReportReq) *txstruct.ReportRsp {
	rsp := new(txstruct.ReportRsp)
	rsp.CALC_TN(true)
	rsp.BaseDataRsp.Code = 0
	rsp.BaseDataRsp.Message = ErrMsgEmpty
	rsp.ReqData = req

	for range "1" {
		//TODO:用户未登录,就校验随附密码,失败就拒绝
		if !thls.subBase.IsRegistered(req.Category) {
			rsp.BaseDataRsp.Message = ErrMsgInvalidCategory
			break
		}

		reportData := convertToReportData(req)
		reportData.ID = thls.LastPushID + 1
		thls.LastPushID = reportData.ID

		//TODO:待优化,抽离成订阅管理器
		if true {
			jsonByte, _ := json.Marshal(reportData)
			jsonStr := string(jsonByte)

			for _, userData := range thls.cache.allUser {
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

//AddUserReq omit
func (thls *BusinessService) AddUserReq(conn *wscmanager.WSConnection, req *txstruct.AddUserReq) *txstruct.AddUserRsp {
	rsp := new(txstruct.AddUserRsp)
	rsp.CALC_TN(true)
	rsp.BaseDataRsp.Code = 0
	rsp.BaseDataRsp.Message = ErrMsgEmpty
	rsp.ReqData = req

	InitialPassword := "pwd"
	if rsp.NewUserID = thls.cache.CreateUser(InitialPassword); 0 < rsp.NewUserID {
		rsp.BaseDataRsp.Code = 0
		rsp.BaseDataRsp.Message = ErrMsgSUCCESS
		rsp.NewPassword = InitialPassword
	} else {
		rsp.BaseDataRsp.Code = 1
		rsp.BaseDataRsp.Message = "FAIL"
	}

	return rsp
}

//SubscribeReq omit
func (thls *BusinessService) SubscribeReq(conn *wscmanager.WSConnection, req *txstruct.SubscribeReq) *txstruct.SubscribeRsp {
	rsp := new(txstruct.SubscribeRsp)
	rsp.CALC_TN(true)
	rsp.BaseDataRsp.Code = 0
	rsp.BaseDataRsp.Message = ErrMsgEmpty
	rsp.ReqData = req

	for range "1" {
		var subInfo *SubscribeUserInfo
		if conn.ExtraData == nil {
			if !thls.cache.UserAndPasswordIsOk(req.OnceUID, req.OncePwd) {
				rsp.BaseDataRsp.Message = ErrMsgNotLoginAndOncePwdErr
				break
			}
			subInfo = thls.cache.allUser[req.OnceUID-1].SubInfo
		} else {
			subInfo = conn.ExtraData.(UserTempData).summary.SubInfo
		}
		if 0 < req.SubUID {
			subInfo.SubUser(req.SubUID)
		}
		if 0 < len(req.SubData) {
			subInfo.SubCategory(req.SubData)
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
