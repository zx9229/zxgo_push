package businessservice

//go get -u -v golang.org/x/net/websocket
//响应消息,要带过来请求消息的结构体.
//校验用户名密码=>在内存数据库查count(id=?&&pwd=?)
//用[当前接收到哪个推送序号了]作为业务心跳消息.
//把所有用户数据都加载到内存中

import (
	"errors"
	"time"

	txstruct "github.com/zx9229/zxgo_push/TxStruct"
)

func convertToReportData(src *txstruct.ReportReq) *txstruct.ReportData {
	dst := new(txstruct.ReportData)
	dst.ID = 0
	dst.Time = time.Now()
	dst.UserID = src.UserID
	dst.UserTagID = src.UserTagID
	dst.UserTagTime = src.UserTagTime
	dst.AttachedI = src.AttachedI
	dst.Message = src.Message
	dst.Category = src.Category
	return dst
}

const emptyString = ""

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

var (
	ErrConnectionIsLoggedOn = "connection is logged on"
	ErrLogicError           = "logic error"
	ErrUserIdNotExist       = "user id not exist"
	ErrIncorrectPassword    = "incorrect password"
	ErrInvalidLoginType     = "invalid login type"
	ErrUserHasLoggedIn      = "user has logged in"
	ErrInvalidCategory      = "invalid category"
)

var (
	ErrParseDataFail       = errors.New("can not parse data")
	ErrFindNotMethod       = errors.New("can not find corresponding method")
	ErrRetValAnomalous     = errors.New("call method and return value anomalous")
	ErrConvertTxStructFail = errors.New("return value convert to tx struct fail")
)
