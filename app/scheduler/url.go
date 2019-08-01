package scheduler

import (
    "app/models"
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
func CronCheckUrl(cron models.SysCron) error {
    UUID := utils.GetUniqueID()
    seelog.Infof("[%v]URL->[%v] Begin ...", UUID, cron.CronCmd)

    // Client http.Client
    var Client *http.Client
    //seelog.Info("InitClient begin ...")
    Client = &http.Client{}
    jar, _ := cookiejar.New(nil)
    Client.Jar = jar

    var data string

    req, err := http.NewRequest("GET", cron.CronCmd, strings.NewReader(data))
    if err != nil {
        seelog.Errorf("[%v]ERROR : %v", UUID, err.Error())
        return err
    }

    header := map[string][]string{}
    req.Header = header

    res, err := Client.Do(req)
    if err != nil {
        seelog.Errorf("[%v]ERROR : %v", UUID, err.Error())
        utils.SendHXMsg("URL检查失败通知", cron.CronHx, cron.CronCmd)
        return err
    }
    defer res.Body.Close()

    if res.StatusCode != 200 {
        seelog.Warnf("[%v]%v响应码非200", UUID, cron.CronCmd)
        utils.SendHXMsg("URL响应码异常通知", cron.CronHx, fmt.Sprintf("%v [%v]", cron.CronCmd, res.StatusCode))
    }

    /*
       body, err := ioutil.ReadAll(res.Body)
       if err != nil {
           seelog.Errorf("ERROR : %v", err.Error())
           return nil, err
       }

       seelog.Trace(string(body))

    */

    seelog.Infof("[%v]>>> Check [%v] Status Code : %v <<<", UUID, cron.CronCmd, res.StatusCode)
    seelog.Infof("[%v]URL->[%v] Finish ...", UUID, cron.CronCmd)

    return nil
}
