package scheduler

import (
    "app/utils"
    "github.com/cihub/seelog"
    "io/ioutil"
)

func ScanCMDList() error {
    Ftp := struct{
        Host string
        Username string
        Password string
        RemoteDir string
    }{
        Host:      utils.GetConfig("hx", "host"),
        Username:  utils.GetConfig("hx", "username"),
        Password:  utils.GetConfig("hx", "password"),
        RemoteDir: utils.GetConfig("hx", "remotedir"),
    }

    var err error
    client, err := utils.NewFtpClient(Ftp.Host, Ftp.Username, Ftp.Password)
    if err != nil {
        seelog.Errorf("ERROR : %v", err.Error())
        return err
    }
    err = client.ChangeDir(Ftp.RemoteDir)
    if err != nil {
        seelog.Errorf("ERROR : %v", err.Error())
        return err
    }

    res, err := client.Retr("test-file.txt")
    if err != nil {
        panic(err)
    }

    buf, err := ioutil.ReadAll(res)
    println(string(buf))


    client.Logout()
    client.Quit()
    return nil
}
