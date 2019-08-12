package scheduler

import (
    "app/models"
    "app/utils"
    "errors"
    "github.com/cihub/seelog"
    "net/http"
    "net/http/cookiejar"
    "strings"
)

type destUrl struct {
    url   string
    hxTos string
}

/*
CheckUrl func(ip string, data []byte) ([]byte, error)
*/
func (u destUrl) checkUrl() error {

    // Client http.Client
    var Client *http.Client
    //seelog.Info("InitClient begin ...")
    Client = &http.Client{}
    jar, _ := cookiejar.New(nil)
    Client.Jar = jar

    var data string

    req, err := http.NewRequest("GET", u.url, strings.NewReader(data))
    if err != nil {
        return err
    }

    header := map[string][]string{}
    req.Header = header

    res, err := Client.Do(req)
    if err != nil {
        return err
    }
    defer res.Body.Close()

    if res.StatusCode != 200 {
        return errors.New("响应码非200")
    }

    /*
       body, err := ioutil.ReadAll(res.Body)
       if err != nil {
           seelog.Errorf("ERROR : %v", err.Error())
           return nil, err
       }

       seelog.Trace(string(body))

    */

    return nil
}

func CheckUrl(obj interface{}) error {
    var dest destUrl

    switch obj.(type) {
    case models.SysCron:
        data := obj.(models.SysCron)
        dest = destUrl{
            url:   data.CronCmd,
            hxTos: data.CronHx,
        }
    case string:

    default:
        return errors.New("Wrong Arg Type ...")
    }

    UUID := utils.GetUniqueID()
    seelog.Infof("[%v]URL->[%v] Begin ...", UUID, dest.url)

    err := dest.checkUrl()
    if err != nil {
        seelog.Errorf("[%v]URL->ERROR:\n%v", UUID, err.Error())
        utils.SendHXMsg(UUID, "URL检查失败通知", dest.hxTos, dest.url)
        return err
    }

    seelog.Infof("[%v]URL->>>> Check [%v] OK! <<<", UUID, dest.url)
    seelog.Infof("[%v]URL->[%v] Finish ...", UUID, dest.url)

    return nil
}
