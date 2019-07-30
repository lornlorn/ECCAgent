package scheduler

import (
    "app/udfuncs"
    "github.com/cihub/seelog"
    "io/ioutil"
    "net/http"
    "net/http/cookiejar"
    "strings"
)

/*
CheckUrl func(ip string, data []byte) ([]byte, error)
*/
func CheckUrl(method string, url string, hxTos string) ([]byte, error) {
    // Client http.Client
    var Client *http.Client
    //seelog.Info("InitClient begin ...")
    Client = &http.Client{}
    jar, _ := cookiejar.New(nil)
    Client.Jar = jar

    var data string

    req, err := http.NewRequest(method, url, strings.NewReader(data))
    if err != nil {
        seelog.Errorf("ERROR : %v", err.Error())
        return nil, err
    }

    header := map[string][]string{}
    req.Header = header

    res, err := Client.Do(req)
    if err != nil {
        seelog.Errorf("ERROR : %v", err.Error())
        // HX
        //tos := utils.GetConfig("hx", "tos")
        udfuncs.SendHXMsg("URL检查失败通知", hxTos, url)
        return nil, err
    }
    defer res.Body.Close()

    body, err := ioutil.ReadAll(res.Body)
    if err != nil {
        seelog.Errorf("ERROR : %v", err.Error())
        return nil, err
    }

    //seelog.Trace(string(body))

    seelog.Infof(">>> Check [%v] Status Code : %v <<<", url, res.StatusCode)

    return body, nil
}
