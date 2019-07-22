package utils

import (
    "fmt"
)

// retCodeDict return code directory
var retCodeDict = map[string]string{
    "0000": "成功",
    "9999": "未知错误",
}

// AjaxReturn struct
type AjaxReturn struct {
    RetCode string `json:"retcode"`
    RetMsg  string `json:"retmsg"`
}

// // AjaxReturnWithData struct
// type AjaxReturnWithData struct {
// 	RetCode string      `json:"retcode"`
// 	RetMsg  string      `json:"retmsg"`
// 	RetData interface{} `json:"retdata"`
// }

// AjaxReturnWithData struct
type AjaxReturnWithData struct {
    AjaxReturn
    RetData interface{} `json:"retdata"`
}

// GetRetMsg func(retcode string) string
// Get Return Message By Return Code
func GetRetMsg(retcode string) string {
    return retCodeDict[retcode]
}

/*
GetAjaxRetObj func(retcode string, errmsg error) AjaxReturn
Get return code and errors
Convert to AjaxReturn struct object
*/
func GetAjaxRetObj(retcode string, errmsg error) AjaxReturn {

    var retobj AjaxReturn

    if errmsg != nil {
        retobj = AjaxReturn{
            RetCode: retcode,
            RetMsg:  fmt.Sprintf("%v, %v", GetRetMsg(retcode), errmsg),
        }
    } else {
        retobj = AjaxReturn{
            RetCode: retcode,
            RetMsg:  GetRetMsg(retcode),
        }
    }

    return retobj

}

/*
GetAjaxRetJSON func(retcode string, errmsg error) []byte
Get return code and errors
Convert to AjaxReturn struct and Marshal to JSON
*/
func GetAjaxRetJSON(retcode string, errmsg error) []byte {

    var retobj AjaxReturn

    if errmsg != nil {
        retobj = AjaxReturn{
            RetCode: retcode,
            RetMsg:  fmt.Sprintf("%v, %v", GetRetMsg(retcode), errmsg.Error()),
        }
    } else {
        retobj = AjaxReturn{
            RetCode: retcode,
            RetMsg:  GetRetMsg(retcode),
        }
    }

    ret := Convert2JSON(retobj)

    return ret

}

/*
GetAjaxRetWithDataObj func(retcode string, errmsg error, data []byte) AjaxReturnWithData
Get return message and data
Convert to AjaxReturnWithData struct object
*/
func GetAjaxRetWithDataObj(retcode string, errmsg error, data interface{}) AjaxReturnWithData {

    var retobj AjaxReturnWithData

    if errmsg != nil {
        retobj.RetCode = retcode
        retobj.RetMsg = fmt.Sprintf("%v, %v", GetRetMsg(retcode), errmsg)
        retobj.RetData = data
        // retobj = AjaxReturnWithData{
        // 	RetCode: retcode,
        // 	RetMsg:  fmt.Sprintf("%v, %v", GetRetMsg(retcode), errmsg),
        // 	RetData: data,
        // }
    } else {
        retobj.RetCode = retcode
        retobj.RetMsg = GetRetMsg(retcode)
        retobj.RetData = data
        // retobj = AjaxReturnWithData{
        // 	RetCode: retcode,
        // 	RetMsg:  GetRetMsg(retcode),
        // 	RetData: data,
        // }
    }

    return retobj

}

/*
GetAjaxRetWithDataJSON func(retcode string, errmsg error, data []byte) []byte
Get return message and data
Convert to AjaxReturnWithData struct and Marshal to JSON
*/
func GetAjaxRetWithDataJSON(retcode string, errmsg error, data interface{}) []byte {

    var retobj AjaxReturnWithData

    if errmsg != nil {
        retobj.RetCode = retcode
        retobj.RetMsg = fmt.Sprintf("%v, %v", GetRetMsg(retcode), errmsg)
        retobj.RetData = data
        // retobj = AjaxReturnWithData{
        // 	RetCode: retcode,
        // 	RetMsg:  fmt.Sprintf("%v, %v", GetRetMsg(retcode), errmsg),
        // 	RetData: data,
        // }
    } else {
        retobj.RetCode = retcode
        retobj.RetMsg = GetRetMsg(retcode)
        retobj.RetData = data
        // retobj = AjaxReturnWithData{
        // 	RetCode: retcode,
        // 	RetMsg:  GetRetMsg(retcode),
        // 	RetData: data,
        // }
    }

    ret := Convert2JSON(retobj)

    return ret

}
