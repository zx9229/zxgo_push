package businessservice

import (
	"time"

	wsconnectionmanager "github.com/zx9229/zxgo_push/WSConnectionManager"
)

const (
	LoginTypeDEFAULT = iota //(登录方式)0,默认值,无效值
	LoginTypeWeb            //(登录方式)网页
	LoginTypeMobile         //(登录方式)手机
	LoginTypePC             //(登录方式)电脑
	LoginTypeEND            //(登录方式)结束值,最大值的下一个
)

//UserBaseInfo 用户的基础信息
type UserBaseInfo struct {
	UserID     int64
	Password   string
	Memo       string    //备注.
	CreateTime time.Time //创建时刻.
	UpdateTime time.Time //更新时刻.
}

//UserStateInfo 用户的状态信息
type UserStateInfo struct {
	conn       *wsconnectionmanager.WSConnection //有值,表示在线.(这个字段不往外导出)
	LoginType  int                               //电脑登录,网页登录,APP登录,等.
	MaxPushID  int64                             //推送给它的最大的推送序号
	LastRecvID int64                             //它上报的自己接收到的最后一个序号
}

//UserSummary 用户的汇总信息
type UserSummary struct {
	Base    UserBaseInfo
	State   []*UserStateInfo
	SubInfo *SubscribeUserInfo
}

func New_UserSummary() *UserSummary {
	curData := new(UserSummary)
	//
	curData.State = make([]*UserStateInfo, 0)
	for i := LoginTypeDEFAULT + 1; i < LoginTypeEND; i++ {
		stateInfo := new(UserStateInfo)
		stateInfo.LoginType = i
		curData.State = append(curData.State, stateInfo)
	}
	curData.SubInfo = New_SubscribeUserInfo()
	//
	return curData
}
