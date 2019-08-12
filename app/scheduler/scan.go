package scheduler

import (
    "app/models"
    "app/utils"
    "errors"
    "fmt"
    "github.com/cihub/seelog"
    "io/ioutil"
    "path/filepath"
    "strings"
)

type destScan struct {
    name     string
    host     string
    user     string
    password string
    file     string
    hxTos    string
}

func (p destScan) scanFtpFile2Line() ([]byte, error) {

    var err error
    client, err := utils.NewFtpClient(p.host, p.user, p.password)
    if err != nil {
        //seelog.Errorf("ERROR : %v", err.Error())
        return nil, err
    }

    filePath, fileName := filepath.Split(p.file)

    err = client.ChangeDir(filePath)
    if err != nil {
        //seelog.Errorf("ERROR : %v", err.Error())
        return nil, err
    }

    res, err := client.Retr(fileName)
    if err != nil {
        return nil, err
    }

    buf, _ := ioutil.ReadAll(res)
    //println(string(buf))

    client.Logout()
    client.Quit()

    return buf, nil
}

func ScanDaemon(obj interface{}) error {
    var dest destScan

    switch obj.(type) {
    case models.SysCron:
        data := obj.(models.SysCron)
        var user string
        var password string

        userInfo := strings.Split(data.CronAuth, "/")
        if len(userInfo) == 2 {
            user = userInfo[0]
            password = userInfo[1]
        } else {
            seelog.Errorf("SCAN->Wrong Auth Config : [%v@%v]", data.CronAuth, data.CronHost)
            return errors.New("Wrong Auth Config")
        }
        dest = destScan{
            name:     data.CronName,
            host:     data.CronHost,
            user:     user,
            password: password,
            file:     data.CronCmd,
            hxTos:    data.CronHx,
        }
    case string:

    default:
        return errors.New("Wrong Arg Type ...")
    }

    UUID := utils.GetUniqueID()
    seelog.Infof("[%v]SCAN->[%v@%v]->[%v] Begin ...", UUID, dest.user, dest.host, dest.file)

    ret, err := dest.scanFtpFile2Line()
    if err != nil {
        seelog.Errorf("[%v]SCAN->ERROR:\n%v", UUID, err.Error())
        utils.SendHXMsg(UUID, "接收命令文件失败", dest.hxTos, fmt.Sprintf("HOST:%v@%v\nFILE:%v\nRESULT:%v", dest.user, dest.host, dest.file, err.Error()))
        return err
    }

    seelog.Debugf("[%v]SCAN->\n%v", UUID, string(ret))
    seelog.Infof("[%v]SCAN->[%v@%v]->[%v] Finish ...", UUID, dest.user, dest.host, dest.file)

    return nil
}
