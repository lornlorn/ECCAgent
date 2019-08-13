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

func (p destScan) scanFtpFile2Byte() ([]byte, error) {

    var buf []byte
    //var lines []string

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

    //lst, _ := client.List(filePath)
    //seelog.Trace(lst)

    res, err := client.Retr(fileName)
    if err != nil {
        //seelog.Errorf("ERROR : %v", err.Error())
        if err.Error() == "550 Failed to open file." {
            return nil, errors.New("未找到目标文件")
        } else {
            return nil, err
        }
    } else {
        buf, _ = ioutil.ReadAll(res)
        //buff := bufio.NewReader(res)
        //for {
        //    line, err := buff.ReadString('\n')
        //    if err != nil || io.EOF == err {
        //        break
        //    }
        //    seelog.Tracef(">>>%v<<<", line)
        //    lines = append(lines, line)
        //}

        err = client.Delete(fileName)
        if err != nil {
            //seelog.Errorf("ERROR : %v", err.Error())
            if err.Error() != "226 Transfer complete." {
                return nil, err
            }
        }
    }

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

    ret, err := dest.scanFtpFile2Byte()
    if err != nil {
        if err.Error() != "未找到目标文件" {
            seelog.Errorf("[%v]SCAN->ERROR:\n%v", UUID, err.Error())
            utils.SendHXMsg(UUID, "接收命令文件失败", dest.hxTos, fmt.Sprintf("HOST:%v@%v\nFILE:%v\nRESULT:%v", dest.user, dest.host, dest.file, err.Error()))
            return err
        } else {
            seelog.Infof("[%v]SCAN->无待处理命令文件 ...", UUID)
        }
    } else {
        // something to run
        seelog.Debugf("[%v]SCAN->命令文件内容:\n%v", UUID, string(ret))
        lines := strings.Split(string(ret), "\n")
        for idx, line := range lines {
            if line != "" {
                seelog.Tracef("No.%v : %v", idx, line)
            }
        }
    }

    seelog.Infof("[%v]SCAN->[%v@%v]->[%v] Finish ...", UUID, dest.user, dest.host, dest.file)

    return nil
}
