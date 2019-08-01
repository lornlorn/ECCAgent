package scheduler

import (
    "app/models"
    "app/utils"
    "errors"
    "github.com/cihub/seelog"
    "net"
    "time"
)

type dest struct {
    addr  string
    hxTos string
}

func (s dest) portScan() error {

    conn, err := net.DialTimeout("tcp", s.addr, 5*time.Second)
    if err != nil {
        return err
    }
    defer conn.Close()

    return nil
}

func PortScan(obj interface{}) error {
    var scanner dest

    switch obj.(type) {
    case models.SysCron:
        data := obj.(models.SysCron)
        scanner = dest{
            addr:  data.CronCmd,
            hxTos: data.CronHx,
        }
    case string:

    default:
        return errors.New("Wrong Arg Type ...")
    }

    UUID := utils.GetUniqueID()
    seelog.Infof("[%v]PORT->[%v] Begin ...", UUID, scanner.addr)

    err := scanner.portScan()
    if err != nil {
        seelog.Errorf("[%v]PORT->ERROR : %v", UUID, err.Error())
        utils.SendHXMsg("端口探测失败通知", scanner.hxTos, scanner.addr)
        return err
    }

    seelog.Infof("[%v]PORT->>>> Check [%v] OK! <<<", UUID, scanner.addr)
    seelog.Infof("[%v]PORT->[%v] Finish ...", UUID, scanner.addr)

    return nil
}
