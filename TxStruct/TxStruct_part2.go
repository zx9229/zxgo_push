package txstruct

import (
	"encoding/json"
	"reflect"
)

//////////////////////////////////////////////////////////////////////

func (thls *BaseDataTx) GET_TN() string {
	return thls.TN
}

func (thls *BaseDataTx) CALC_TN(modifyTN bool) string {
	TypeName := reflect.ValueOf(*thls).Type().Name()
	if modifyTN {
		thls.TN = TypeName
	}
	return TypeName
}

func (thls *BaseDataTx) TO_JSON(panicWhenError bool) string {
	if bytes, err := json.Marshal(thls); err != nil {
		if panicWhenError {
			panic(err)
		}
		return ""
	} else {
		return string(bytes)
	}
}

//////////////////////////////////////////////////////////////////////

func (thls *UnknownNotice) GET_TN() string {
	return thls.TN
}

func (thls *UnknownNotice) CALC_TN(modifyTN bool) string {
	TypeName := reflect.ValueOf(*thls).Type().Name()
	if modifyTN {
		thls.TN = TypeName
	}
	return TypeName
}

func (thls *UnknownNotice) TO_JSON(panicWhenError bool) string {
	if bytes, err := json.Marshal(thls); err != nil {
		if panicWhenError {
			panic(err)
		}
		return ""
	} else {
		return string(bytes)
	}
}

//////////////////////////////////////////////////////////////////////

func (thls *LoginReq) GET_TN() string {
	return thls.TN
}

func (thls *LoginReq) CALC_TN(modifyTN bool) string {
	TypeName := reflect.ValueOf(*thls).Type().Name()
	if modifyTN {
		thls.TN = TypeName
	}
	return TypeName
}

func (thls *LoginReq) TO_JSON(panicWhenError bool) string {
	if bytes, err := json.Marshal(thls); err != nil {
		if panicWhenError {
			panic(err)
		}
		return ""
	} else {
		return string(bytes)
	}
}

//////////////////////////////////////////////////////////////////////

func (thls *LoginRsp) GET_TN() string {
	return thls.TN
}

func (thls *LoginRsp) CALC_TN(modifyTN bool) string {
	TypeName := reflect.ValueOf(*thls).Type().Name()
	if modifyTN {
		thls.TN = TypeName
	}
	return TypeName
}

func (thls *LoginRsp) TO_JSON(panicWhenError bool) string {
	if bytes, err := json.Marshal(thls); err != nil {
		if panicWhenError {
			panic(err)
		}
		return ""
	} else {
		return string(bytes)
	}
}

//////////////////////////////////////////////////////////////////////
func (thls *ReportReq) GET_TN() string {
	return thls.TN
}

func (thls *ReportReq) CALC_TN(modifyTN bool) string {
	TypeName := reflect.ValueOf(*thls).Type().Name()
	if modifyTN {
		thls.TN = TypeName
	}
	return TypeName
}

func (thls *ReportReq) TO_JSON(panicWhenError bool) string {
	if bytes, err := json.Marshal(thls); err != nil {
		if panicWhenError {
			panic(err)
		}
		return ""
	} else {
		return string(bytes)
	}
}

//////////////////////////////////////////////////////////////////////

func (thls *ReportRsp) GET_TN() string {
	return thls.TN
}

func (thls *ReportRsp) CALC_TN(modifyTN bool) string {
	TypeName := reflect.ValueOf(*thls).Type().Name()
	if modifyTN {
		thls.TN = TypeName
	}
	return TypeName
}

func (thls *ReportRsp) TO_JSON(panicWhenError bool) string {
	if bytes, err := json.Marshal(thls); err != nil {
		if panicWhenError {
			panic(err)
		}
		return ""
	} else {
		return string(bytes)
	}
}

//////////////////////////////////////////////////////////////////////