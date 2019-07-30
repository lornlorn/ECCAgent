package scheduler

import (
    "app/udfuncs"
    "github.com/cihub/seelog"
    "net"
    "time"
)

func PortScan(dest string, hxTos string) error {
    conn, err := net.DialTimeout("tcp", dest, 5*time.Second)
    if err != nil {
        seelog.Errorf("ERROR : %v", err.Error())
        udfuncs.SendHXMsg("端口探测失败通知", hxTos, dest)
        return err
    }
    defer conn.Close()

    seelog.Infof(">>> Check [%v] OK! <<<", dest)

    return nil
}
