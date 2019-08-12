package scheduler

import (
    "app/models"
    "app/utils"
    "errors"
    "github.com/cihub/seelog"
    "net"
    "time"
)

type destPort struct {
    addr  string
    hxTos string
}

func (p destPort) portScan() error {

    conn, err := net.DialTimeout("tcp", p.addr, 5*time.Second)
    if err != nil {
        return err
    }
    defer conn.Close()

    return nil
}

func PortScan(obj interface{}) error {
    var dest destPort

    switch obj.(type) {
    case models.SysCron:
        data := obj.(models.SysCron)
        dest = destPort{
            addr:  data.CronCmd,
            hxTos: data.CronHx,
        }
    case string:

    default:
        return errors.New("Wrong Arg Type ...")
    }

    UUID := utils.GetUniqueID()
    seelog.Infof("[%v]PORT->[%v] Begin ...", UUID, dest.addr)

    err := dest.portScan()
    if err != nil {
        seelog.Errorf("[%v]PORT->ERROR:\n%v", UUID, err.Error())
        utils.SendHXMsg(UUID, "端口探测失败通知", dest.hxTos, dest.addr)
        return err
    }

    seelog.Infof("[%v]PORT->>>> Check [%v] OK! <<<", UUID, dest.addr)
    seelog.Infof("[%v]PORT->[%v] Finish ...", UUID, dest.addr)

    return nil
}
