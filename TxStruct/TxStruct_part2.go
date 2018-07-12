package txstruct

import (
	"encoding/json"
	"reflect"
)

//////////////////////////////////////////////////////////////////////

//GetTN omit
func (thls *BaseDataTx) GetTN() string {
	return thls.TN
}

//CalcTN omit
func (thls *BaseDataTx) CalcTN(modifyTN bool) string {
	TypeName := reflect.ValueOf(*thls).Type().Name()
	if modifyTN {
		thls.TN = TypeName
	}
	return TypeName
}

//ToJSON omit
func (thls *BaseDataTx) ToJSON(panicWhenError bool) string {
	if bytes, err := json.Marshal(thls); err == nil {
		return string(bytes)
	} else if panicWhenError {
		panic(err)
	}
	return ""
}

//////////////////////////////////////////////////////////////////////

//GetTN omit
func (thls *UnknownNotice) GetTN() string {
	return thls.TN
}

//CalcTN omit
func (thls *UnknownNotice) CalcTN(modifyTN bool) string {
	TypeName := reflect.ValueOf(*thls).Type().Name()
	if modifyTN {
		thls.TN = TypeName
	}
	return TypeName
}

//ToJSON omit
func (thls *UnknownNotice) ToJSON(panicWhenError bool) string {
	if bytes, err := json.Marshal(thls); err == nil {
		return string(bytes)
	} else if panicWhenError {
		panic(err)
	}
	return ""
}

//////////////////////////////////////////////////////////////////////

//GetTN omit
func (thls *LoginReq) GetTN() string {
	return thls.TN
}

//CalcTN omit
func (thls *LoginReq) CalcTN(modifyTN bool) string {
	TypeName := reflect.ValueOf(*thls).Type().Name()
	if modifyTN {
		thls.TN = TypeName
	}
	return TypeName
}

//ToJSON omit
func (thls *LoginReq) ToJSON(panicWhenError bool) string {
	if bytes, err := json.Marshal(thls); err == nil {
		return string(bytes)
	} else if panicWhenError {
		panic(err)
	}
	return ""
}

//////////////////////////////////////////////////////////////////////

//GetTN omit
func (thls *LoginRsp) GetTN() string {
	return thls.TN
}

//CalcTN omit
func (thls *LoginRsp) CalcTN(modifyTN bool) string {
	TypeName := reflect.ValueOf(*thls).Type().Name()
	if modifyTN {
		thls.TN = TypeName
	}
	return TypeName
}

//ToJSON omit
func (thls *LoginRsp) ToJSON(panicWhenError bool) string {
	if bytes, err := json.Marshal(thls); err == nil {
		return string(bytes)
	} else if panicWhenError {
		panic(err)
	}
	return ""
}

//////////////////////////////////////////////////////////////////////

//GetTN omit
func (thls *ReportReq) GetTN() string {
	return thls.TN
}

//CalcTN omit
func (thls *ReportReq) CalcTN(modifyTN bool) string {
	TypeName := reflect.ValueOf(*thls).Type().Name()
	if modifyTN {
		thls.TN = TypeName
	}
	return TypeName
}

//ToJSON omit
func (thls *ReportReq) ToJSON(panicWhenError bool) string {
	if bytes, err := json.Marshal(thls); err == nil {
		return string(bytes)
	} else if panicWhenError {
		panic(err)
	}
	return ""
}

//////////////////////////////////////////////////////////////////////

//GetTN omit
func (thls *ReportRsp) GetTN() string {
	return thls.TN
}

//CalcTN omit
func (thls *ReportRsp) CalcTN(modifyTN bool) string {
	TypeName := reflect.ValueOf(*thls).Type().Name()
	if modifyTN {
		thls.TN = TypeName
	}
	return TypeName
}

//ToJSON omit
func (thls *ReportRsp) ToJSON(panicWhenError bool) string {
	if bytes, err := json.Marshal(thls); err == nil {
		return string(bytes)
	} else if panicWhenError {
		panic(err)
	}
	return ""
}

//////////////////////////////////////////////////////////////////////

//GetTN omit
func (thls *ReportData) GetTN() string {
	return thls.TN
}

//CalcTN omit
func (thls *ReportData) CalcTN(modifyTN bool) string {
	TypeName := reflect.ValueOf(*thls).Type().Name()
	if modifyTN {
		thls.TN = TypeName
	}
	return TypeName
}

//ToJSON omit
func (thls *ReportData) ToJSON(panicWhenError bool) string {
	if bytes, err := json.Marshal(thls); err == nil {
		return string(bytes)
	} else if panicWhenError {
		panic(err)
	}
	return ""
}

//////////////////////////////////////////////////////////////////////

//GetTN omit
func (thls *AddUserReq) GetTN() string {
	return thls.TN
}

//CalcTN omit
func (thls *AddUserReq) CalcTN(modifyTN bool) string {
	TypeName := reflect.ValueOf(*thls).Type().Name()
	if modifyTN {
		thls.TN = TypeName
	}
	return TypeName
}

//ToJSON omit
func (thls *AddUserReq) ToJSON(panicWhenError bool) string {
	if bytes, err := json.Marshal(thls); err == nil {
		return string(bytes)
	} else if panicWhenError {
		panic(err)
	}
	return ""
}

//////////////////////////////////////////////////////////////////////

//GetTN omit
func (thls *AddUserRsp) GetTN() string {
	return thls.TN
}

//CalcTN omit
func (thls *AddUserRsp) CalcTN(modifyTN bool) string {
	TypeName := reflect.ValueOf(*thls).Type().Name()
	if modifyTN {
		thls.TN = TypeName
	}
	return TypeName
}

//ToJSON omit
func (thls *AddUserRsp) ToJSON(panicWhenError bool) string {
	if bytes, err := json.Marshal(thls); err == nil {
		return string(bytes)
	} else if panicWhenError {
		panic(err)
	}
	return ""
}

//////////////////////////////////////////////////////////////////////

//GetTN omit
func (thls *SubscribeReq) GetTN() string {
	return thls.TN
}

//CalcTN omit
func (thls *SubscribeReq) CalcTN(modifyTN bool) string {
	TypeName := reflect.ValueOf(*thls).Type().Name()
	if modifyTN {
		thls.TN = TypeName
	}
	return TypeName
}

//ToJSON omit
func (thls *SubscribeReq) ToJSON(panicWhenError bool) string {
	if bytes, err := json.Marshal(thls); err == nil {
		return string(bytes)
	} else if panicWhenError {
		panic(err)
	}
	return ""
}

//////////////////////////////////////////////////////////////////////

//GetTN omit
func (thls *SubscribeRsp) GetTN() string {
	return thls.TN
}

//CalcTN omit
func (thls *SubscribeRsp) CalcTN(modifyTN bool) string {
	TypeName := reflect.ValueOf(*thls).Type().Name()
	if modifyTN {
		thls.TN = TypeName
	}
	return TypeName
}

//ToJSON omit
func (thls *SubscribeRsp) ToJSON(panicWhenError bool) string {
	if bytes, err := json.Marshal(thls); err == nil {
		return string(bytes)
	} else if panicWhenError {
		panic(err)
	}
	return ""
}

//////////////////////////////////////////////////////////////////////

//GetTN omit
func (thls *ActionCategoryReq) GetTN() string {
	return thls.TN
}

//CalcTN omit
func (thls *ActionCategoryReq) CalcTN(modifyTN bool) string {
	TypeName := reflect.ValueOf(*thls).Type().Name()
	if modifyTN {
		thls.TN = TypeName
	}
	return TypeName
}

//ToJSON omit
func (thls *ActionCategoryReq) ToJSON(panicWhenError bool) string {
	if bytes, err := json.Marshal(thls); err == nil {
		return string(bytes)
	} else if panicWhenError {
		panic(err)
	}
	return ""
}

//////////////////////////////////////////////////////////////////////

//GetTN omit
func (thls *ActionCategoryRsp) GetTN() string {
	return thls.TN
}

//CalcTN omit
func (thls *ActionCategoryRsp) CalcTN(modifyTN bool) string {
	TypeName := reflect.ValueOf(*thls).Type().Name()
	if modifyTN {
		thls.TN = TypeName
	}
	return TypeName
}

//ToJSON omit
func (thls *ActionCategoryRsp) ToJSON(panicWhenError bool) string {
	if bytes, err := json.Marshal(thls); err == nil {
		return string(bytes)
	} else if panicWhenError {
		panic(err)
	}
	return ""
}
