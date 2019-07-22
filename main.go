package main

import (
    "app/httpsvr"
    "app/scheduler"
    "github.com/cihub/seelog"
    "log"
    "net/http"
    _ "net/http/pprof"

    "app/utils"
)

const (
    // SeelogCfg seelog config file path
    SeelogCfg = "./config/seelog.xml"
    // AppCfg app config file path
    AppCfg = "./config/app.conf"
)

func main() {
    //这里实现了远程获取pprof数据的接口
    go func() {
        log.Println(http.ListenAndServe(":9999", nil))
    }()

    var err error
    var msg string

    // Initialize Logger
    msg = "Initialize Logger ..."
    err = utils.InitLogger(SeelogCfg)
    if err != nil {
        panic("Exit!")
    }
    defer seelog.Flush()
    seelog.Infof("%v Success !", msg)

    // Read Configuration
    msg = "Load Configuration ..."
    err = utils.InitConfig(AppCfg)
    if err != nil {
        seelog.Criticalf("%v Error : %v", msg, err.Error())
        panic("Exit!")
    }
    seelog.Infof("%v Success !", msg)

    // Init DB
    msg = "Connect Database ..."
    dbtype := utils.GetConfig("db", "dbtype")
    dbstr := utils.GetConfig("db", "dbstr")
    err = utils.InitDB(dbtype, dbstr)
    if err != nil {
        seelog.Criticalf("%v Error : %v", msg, err.Error())
        panic("Exit!")
    }
    defer utils.Engine.Close()
    seelog.Infof("%v Success !", msg)

    // Start Cron
    msg = "Start Cron ..."
    err = scheduler.InitCron()
    if err != nil {
        seelog.Criticalf("%v Error : %v", msg, err.Error())
        panic("Exit!")
    }
    seelog.Infof("%v Success !", msg)
    //scheduler.Cron.Wait()

    // Start HTTP Server
    msg = "5 -> Starting HTTP Server"
    seelog.Infof("%v !", msg)
    seelog.Info("***Everything is OK !***")
    // log.Fatalln(httpsvr.StartHTTP())
    httpPort := utils.GetConfigInt("http", "httpport")
    writeTimeout := utils.GetConfigInt("http", "writetimeout")
    readTimeout := utils.GetConfigInt("http", "readtimeout")
    err = httpsvr.StartHTTP(httpPort, writeTimeout, readTimeout)
    if err != nil {
        seelog.Criticalf("%v Error : %v", msg, err)
        panic("Exit!")
    }
    seelog.Infof("%v Success !", msg)
}
