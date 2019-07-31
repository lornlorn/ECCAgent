package scheduler

import (
    "app/utils"
    "fmt"
    "github.com/cihub/seelog"
    "net/http"
    "net/http/cookiejar"
    "strings"
)

/*
CheckUrl func(ip string, data []byte) ([]byte, error)
*/
func CheckUrl(method string, url string, hxTos string) error {
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
        return err
    }

    header := map[string][]string{}
    req.Header = header

    res, err := Client.Do(req)
    if err != nil {
        seelog.Errorf("ERROR : %v", err.Error())
        utils.SendHXMsg("URL检查失败通知", hxTos, url)
        return err
    }
    defer res.Body.Close()

    if res.StatusCode != 200 {
        seelog.Warnf("%v响应码非200", url)
        utils.SendHXMsg("URL响应码异常通知", hxTos, fmt.Sprintf("%v [%v]", url, res.StatusCode))
    }

    /*
       body, err := ioutil.ReadAll(res.Body)
       if err != nil {
           seelog.Errorf("ERROR : %v", err.Error())
           return nil, err
       }

       seelog.Trace(string(body))

    */

    seelog.Infof(">>> Check [%v] Status Code : %v <<<", url, res.StatusCode)

    return nil
}
