package businessservice

import (
	"fmt"
	"log"
	"reflect"

	txstruct "github.com/zx9229/zxgo_push/TxStruct"
	wscmanager "github.com/zx9229/zxgo_push/WSCManager"
)

//BusinessService omit
type BusinessService struct {
	methods    map[string]reflect.Value        //函数名<=>函数指针
	parser     *txstruct.TxParser              //(通信消息)解析器
	connMngr   *wscmanager.WSConnectionManager //连接管理器
	userMngr   *UserInfoManager                //用户管理器
	catgMngr   *CategoryManager                //类别管理器
	LastPushID int64                           //最后一个推送消息的序号
}

//New_BusinessService omit
func New_BusinessService() *BusinessService {
	curData := new(BusinessService)
	//
	curData.methods = make(map[string]reflect.Value)
	curData.parser = txstruct.New_TxParser()
	curData.connMngr = wscmanager.New_WSConnectionManager()
	//
	curData.methods = curData.calcMethods()
	curData.connMngr.CbConnected = curData.handleConnected
	curData.connMngr.CbDisconnected = curData.handleDisconnected
	curData.connMngr.CbReceive = curData.handleReceive
	//
	curData.userMngr = new_UserInfoManager()
	curData.catgMngr = New_CategoryManager()
	//
	return curData
}

//GetConnectionManager omit
func (thls *BusinessService) GetConnectionManager() *wscmanager.WSConnectionManager {
	return thls.connMngr
}

//calcMethods omit
func (thls *BusinessService) calcMethods() map[string]reflect.Value {
	methodsMap := make(map[string]reflect.Value)
	vf := reflect.ValueOf(thls)
	vft := vf.Type()
	mNum := vf.NumMethod()
	for i := 0; i < mNum; i++ {
		mName := vft.Method(i).Name
		methodsMap[mName] = vf.Method(i)
	}
	return methodsMap
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
		tmpData.sInfo.conn = nil
	}
}

//HandleReceive omit
func (thls *BusinessService) handleReceive(conn *wscmanager.WSConnection, bytes []byte) {
	var err error
	var responseMessage string

	for range "1" {
		var objData txstruct.TxInterface

		if objData, _, err = thls.parser.ParseByteSlice(bytes); err != nil {
			notice := txstruct.UnknownNotice{Message: ErrParseDataFail.Error(), RawMessage: string(bytes)}
			notice.CalcTN(true)
			responseMessage = notice.ToJSON(true)
			break
		}

		var rspData txstruct.TxInterface
		for range "1" {
			var ok bool
			var method reflect.Value
			if method, ok = thls.methods[objData.GetTN()]; !ok {
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
			notice.CalcTN(true)
			responseMessage = notice.ToJSON(true)
			break
		}

		responseMessage = rspData.ToJSON(true)
	}

	if err = conn.Send(responseMessage); err != nil {
		conn.Close()
		log.Println(fmt.Sprintf("发送消息失败了"))
	}
}

//LoginReq omit
func (thls *BusinessService) LoginReq(conn *wscmanager.WSConnection, req *txstruct.LoginReq) *txstruct.LoginRsp {
	rsp := new(txstruct.LoginRsp)
	rsp.CalcTN(true)
	rsp.BaseDataRsp.Code = 0
	rsp.BaseDataRsp.Message = ErrMsgEmpty
	rsp.ReqData = req

	rsp.BaseDataRsp.Message = thls.userMngr.LoginUser(conn, req)

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
	rsp.CalcTN(true)
	rsp.BaseDataRsp.Code = 0
	rsp.BaseDataRsp.Message = emptyString
	rsp.ReqData = req

	for range "1" {
		if conn.ExtraData == nil {
			//用户未登录,就校验随附密码
			if !thls.userMngr.PasswordOk(req.OnceUID, req.OncePwd) {
				rsp.BaseDataRsp.Message = ErrMsgNotLoginAndOncePwdErr
				break
			}
		}

		if !thls.catgMngr.IsRegistered(req.Category) {
			rsp.BaseDataRsp.Message = ErrMsgInvalidCategory
			break
		}

		reportData := convertToReportData(req)
		reportData.ID = thls.LastPushID + 1
		thls.LastPushID = reportData.ID

		thls.userMngr.PushData(reportData)
	}

	if rsp.BaseDataRsp.Message == emptyString {
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
	rsp.CalcTN(true)
	rsp.BaseDataRsp.Code = 0
	rsp.BaseDataRsp.Message = ErrMsgEmpty
	rsp.ReqData = req

	InitialPassword := "pwd"
	if rsp.NewUserID = thls.userMngr.CreateUser(InitialPassword); 0 < rsp.NewUserID {
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
	rsp.CalcTN(true)
	rsp.BaseDataRsp.Code = 0
	rsp.BaseDataRsp.Message = ErrMsgEmpty
	rsp.ReqData = req

	for range "1" {
		var subInfo *UserSubscriptionInfo
		if conn.ExtraData == nil {
			if !thls.userMngr.PasswordOk(req.OnceUID, req.OncePwd) {
				rsp.BaseDataRsp.Message = ErrMsgNotLoginAndOncePwdErr
				break
			}
			subInfo = thls.userMngr.allUser[req.OnceUID-1].SubInfo
		} else {
			subInfo = conn.ExtraData.(*UserTempData).tInfo.SubInfo
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

//ActionCategoryReq omit
func (thls *BusinessService) ActionCategoryReq(conn *wscmanager.WSConnection, req *txstruct.ActionCategoryReq) *txstruct.ActionCategoryRsp {
	rsp := new(txstruct.ActionCategoryRsp)
	rsp.CalcTN(true)
	rsp.BaseDataRsp.Code = 0
	rsp.BaseDataRsp.Message = emptyString
	rsp.ReqData = req

	for range "1" {
		if conn.ExtraData == nil {
			if !thls.userMngr.PasswordOk(req.OnceUID, req.OncePwd) {
				rsp.BaseDataRsp.Message = ErrMsgNotLoginAndOncePwdErr
				break
			}
		}
		if 0 < req.Action {
			thls.catgMngr.AddCategory(req.Category)
		} else if 0 < req.Action {
			thls.catgMngr.DelCategory(req.Category)
		} else {
			rsp.QryResult = thls.catgMngr.QryCategory(req.Category)
		}
	}

	if rsp.BaseDataRsp.Message == emptyString {
		rsp.BaseDataRsp.Code = 0
		rsp.BaseDataRsp.Message = ErrMsgSUCCESS
	} else {
		rsp.BaseDataRsp.Code = 1
	}
	return rsp
}
