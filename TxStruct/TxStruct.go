package txstruct

import (
	"reflect"
	"time"
)

//TxInterface 通信结构体的接口
type TxInterface interface {
	// 获取字段(TN=>TypeName)的值.
	// 函数体实际上是{ return self.TN }.
	GetTN() string

	// 计算类型的名字(calc type name).
	// if  (modifyTN == true) { TN = TypeName }.
	CalcTN(modifyTN bool) string

	// 将自身转成json字符串.
	// 转换失败的话,如果(panicWhenError == true),就panic; 否则返回空字符串.
	ToJSON(panicWhenError bool) string
}

//inner_check_by_compile 在编译时,检查各个结构体是否进行正常书写.
// 每个通信结构体都要写到这里面去,因为解析器还依赖了这个函数.
func inner_check_by_compile() ([]TxInterface, []interface{}) {
	sliceTx := make([]TxInterface, 0)
	sliceOrig := make([]interface{}, 0)
	sliceTx = append(sliceTx, new(BaseDataTx))
	sliceOrig = append(sliceOrig, BaseDataTx{})
	sliceTx = append(sliceTx, new(UnknownNotice))
	sliceOrig = append(sliceOrig, UnknownNotice{})
	sliceTx = append(sliceTx, new(LoginReq))
	sliceOrig = append(sliceOrig, LoginReq{})
	sliceTx = append(sliceTx, new(LoginRsp))
	sliceOrig = append(sliceOrig, LoginRsp{})
	sliceTx = append(sliceTx, new(ReportReq))
	sliceOrig = append(sliceOrig, ReportReq{})
	sliceTx = append(sliceTx, new(ReportRsp))
	sliceOrig = append(sliceOrig, ReportRsp{})
	sliceTx = append(sliceTx, new(ReportData))
	sliceOrig = append(sliceOrig, ReportData{})
	sliceTx = append(sliceTx, new(AddUserReq))
	sliceOrig = append(sliceOrig, AddUserReq{})
	sliceTx = append(sliceTx, new(AddUserRsp))
	sliceOrig = append(sliceOrig, AddUserRsp{})
	sliceTx = append(sliceTx, new(SubscribeReq))
	sliceOrig = append(sliceOrig, SubscribeReq{})
	sliceTx = append(sliceTx, new(SubscribeRsp))
	sliceOrig = append(sliceOrig, SubscribeRsp{})
	sliceTx = append(sliceTx, new(ActionCategoryReq))
	sliceOrig = append(sliceOrig, ActionCategoryReq{})
	sliceTx = append(sliceTx, new(ActionCategoryRsp))
	sliceOrig = append(sliceOrig, ActionCategoryRsp{})
	//
	if len(sliceTx) != len(sliceOrig) {
		panic("使用指针类型的时候,解析器无法工作,待优化")
	}
	for i := 0; i < len(sliceTx); i++ {
		if sliceTx[i].CalcTN(false) != reflect.TypeOf(sliceOrig[i]).Name() {
			panic("使用指针类型的时候,解析器无法工作,待优化")
		}
	}
	return sliceTx, sliceOrig
}

//BaseDataTx 通信结构体的基本数据(每个通信结构体里面都要有它们)
type BaseDataTx struct {
	TN string //([通信]结构体必需)(TN=>TypeName)
}

//BaseDataReq 请求结构体的基本数据(每个请求结构体里面都要有它们)
type BaseDataReq struct {
	InnerID int64  //([请求]结构体必需)内部ID(用户无关)(API同时支持[同步]&&[异步]时需要的字段)
	RefID   int64  //([请求]结构体必需)参考ID(用户填值)
	DU      int64  //([请求]结构体选填)如果未登录,用它执行一次性的校验(disposable user)
	DP      string //([请求]结构体选填)如果未登录,用它执行一次性的校验(disposable pwd)
}

//BaseDataRsp 响应结构体的基本数据(每个响应结构体里面都要有它们)
type BaseDataRsp struct {
	Code    int    //([响应]结构体必需)执行请求结构体的返回值.
	Message string //([响应]结构体必需)执行请求结构体的返回详情.
}

//UnknownNotice 客户端发过来一个消息,服务端解析消息失败,此时也要告诉客户端,就是返回这个消息.
type UnknownNotice struct {
	BaseDataTx
	Message    string //执行情况
	RawMessage string //原始消息
}

//AddUserReq omit
type AddUserReq struct {
	BaseDataTx
	BaseDataReq
}

//AddUserRsp omit
type AddUserRsp struct {
	BaseDataTx
	BaseDataRsp
	ReqData     *AddUserReq
	NewUserID   int64  //新用户的ID
	NewPassword string //新用户的密码
}

//LoginReq omit
type LoginReq struct {
	BaseDataTx
	BaseDataReq
	UserID     int64  //
	Password   string //
	Way        int    //登录方式(手机登录,网页登录,PC登录)
	ForceLogin bool   //强制登录(将已登录的用户踢下线)
	LastMsgID  int64  //全局消息,我接收到哪一个了(-1代表无效字段)
}

//LoginRsp omit
type LoginRsp struct {
	BaseDataTx
	BaseDataRsp
	ReqData *LoginReq
}

//ActionCategoryReq 增删查类别
type ActionCategoryReq struct {
	BaseDataTx
	BaseDataReq
	Action   int //-1:删除,0:查询,1:增加
	Category string
}

//ActionCategoryRsp omit
type ActionCategoryRsp struct {
	BaseDataTx
	BaseDataRsp
	ReqData   *ActionCategoryReq
	QryResult []string //查询结果
}

//SubscribeReq omit
type SubscribeReq struct {
	BaseDataTx
	BaseDataReq
	SubUID  int64  //要订阅的用户([正数]时才有效)
	SubData string //要订阅的类别([非空]时有效)
}

//SubscribeRsp omit
type SubscribeRsp struct {
	BaseDataTx
	BaseDataRsp
	ReqData *SubscribeReq
}

//ReportReq omit
type ReportReq struct {
	BaseDataTx
	BaseDataReq
	UserID      int64
	UserTagID   int64     //rowId
	UserTagTime time.Time //rowUpdateTime
	AttachedI   int64     //随附的一个数字
	Message     string    //消息的详情
	Category    string    //消息所属的类别
}

//ReportRsp omit
type ReportRsp struct {
	BaseDataTx
	BaseDataRsp
	ReqData *ReportReq
}

//ReportData omit
type ReportData struct {
	BaseDataTx
	ID   int64     `xorm:"notnull pk"` //数据库的递增序号.
	Time time.Time //插入数据库的时刻.
	//
	UserID      int64
	UserTagID   int64     //用户在维护这一条数据时,可能会给它贴上一个序号,这个序号的值
	UserTagTime time.Time //用户在维护这一条数据时,可能会给它贴上一个时刻,这个时刻的值
	AttachedI   int64     //随附的一个数字(或许你上报的消息,想表示(成功/失败)或(某某的个数)呢.)
	Message     string    //消息的详情(当然,这里面也可以是json,这样就不需要上面的字段了)
	Category    string    //消息所属的类别
}
