package scheduler

import (
    "app/models"
    "app/udfuncs"
    "app/utils"
    "fmt"
    "os/exec"
    "runtime"
    "strings"
    "time"

    "github.com/cihub/seelog"
)

// Execute func(src string, command string, envs []string, args ...string) ([]byte, error)
func Execute(src string, uuid string, command string, envs []string, args ...string) ([]byte, error) {
    seelog.Debug("Execute Job ...")
    nowTime := time.Now()
    timeFormat := "2006-01-02 15:04:05" // 时间格式化模板
    nowTimeStr := nowTime.Format(timeFormat)
    var sysLog = models.NewLog{
        RunTime: nowTimeStr,
        // RunEnvs: strings.Replace(strings.Trim(fmt.Sprint(es), "[]"), " ", ",", -1),
        RunEnvs: strings.Trim(fmt.Sprint(envs), "[]"),
        RunCmd:  command,
        RunArgs: strings.Trim(fmt.Sprint(args), "[]"),
        ReqSrc:  src,
    }
    output, beginTimeStr, endTimeStr, err := Run(command, envs, args...)
    if err != nil {
        sysLog.RunStatus = "失败"
        sysLog.RunMsg = err.Error()
        return nil, err
    }
    var result []byte
    if runtime.GOOS == "windows" {
        seelog.Debug("Decode GBK to UTF-8 ...")
        result, err = utils.DecodeGBK2UTF8(output)
        if err != nil {
            seelog.Errorf("Decode GBK to UTF-8 Error : %v", err.Error())
            result = output
        }
    } else {
        result = output
    }
    seelog.Trace(string(result))

    logFilePath, err := udfuncs.WriteRunLog2File(uuid, result, beginTimeStr, endTimeStr)
    if err != nil {
        seelog.Errorf("Write Run Log Fail : %v", err.Error())
        sysLog.LogfilePath = fmt.Sprintf("Write Run Log Fail : %v", err.Error())
    } else {
        sysLog.LogfilePath = logFilePath
    }
    seelog.Debugf("Write Run Log : %v", logFilePath)

    sysLog.RunStatus = "成功"
    sysLog.RunMsg = string(result)

    err = sysLog.Save()
    if err != nil {
        seelog.Errorf("Write DB Run Log Fail : %v", err.Error())
    }

    return result, nil

}

/*
Run func(command string, args ...string) ([]byte, string, string, error)
*/
func Run(command string, envs []string, args ...string) ([]byte, string, string, error) {
    timeFormat := "2006-01-02 15:04:05" // 时间格式化模板
    beginTime := time.Now()
    beginTimeStr := beginTime.Format(timeFormat)

    // 执行
    cmd := exec.Command(command, args...)
    cmd.Env = envs

    output, err := cmd.StdoutPipe()
    // output, err := cmd.CombinedOutput()
    if err != nil {
        return nil, "", "", err
    }

    if err = cmd.Start(); err != nil {
        return nil, "", "", err
    }

    var out = make([]byte, 0, 1024)
    for {
        tmp := make([]byte, 128)
        n, err := output.Read(tmp)
        out = append(out, tmp[:n]...)
        if err != nil {
            break
        }
    }

    if err = cmd.Wait(); err != nil {
        return nil, "", "", err
    }

    endTime := time.Now()
    endTimeStr := endTime.Format(timeFormat)

    return out, beginTimeStr, endTimeStr, nil
}
