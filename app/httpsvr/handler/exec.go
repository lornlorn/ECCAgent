package handler

import (
    "app/scheduler"
    "app/utils"
    "net/http"

    "github.com/cihub/seelog"
)

// ExecuteHandler func(res http.ResponseWriter, req *http.Request)
func ExecuteHandler(res http.ResponseWriter, req *http.Request) {

    seelog.Infof("Route Execute : %v", req.URL)
    // key := mux.Vars(req)["key"]

    reqBody := utils.ReadRequestBody2JSON(req.Body)
    seelog.Debugf("Request Body : %v", string(reqBody))

    reqURL := req.URL.Query()
    seelog.Debugf("Request Params : %v", reqURL)

    envs := utils.ReadJSONData2Array(reqBody, "data.envs")
    cmd := utils.GetJSONResultFromRequestBody(reqBody, "data.cmd")
    args := utils.ReadJSONData2Array(reqBody, "data.args")

    var es []string
    for _, env := range envs {
        es = append(es, env.String())
    }
    var as []string
    for _, arg := range args {
        as = append(as, arg.String())
    }

    seelog.Debugf("[%v][%v][%v]", es, cmd, as)

    var ret []byte
    output, err := scheduler.Execute("http", "", cmd.String(), es, as...)
    if err != nil {
        seelog.Errorf("Command Run Error : %v", err.Error())
        ret = utils.GetAjaxRetJSON("9999", err)
    } else {
        seelog.Debugf("执行结果 : %v", string(output))
        ret = utils.GetAjaxRetWithDataJSON("0000", nil, string(output))
    }

    res.Write(ret)
    return
}
