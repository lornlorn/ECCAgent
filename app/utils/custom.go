package utils

import (
    "bytes"
    "fmt"
    "github.com/cihub/seelog"
    "os"
    "path/filepath"
    "runtime"
    "time"
)

func WriteRunLog2File(uuid string, content []byte, beginTimeStr string, endTimeStr string) (string, error) {
    // logFileFormat := "20060102150405" // 日志文件命名时间格式化
    // logFileNameTime := nowTime.Format(logFileFormat)
    var UUID string
    if uuid == "" || len(uuid) != 32 {
        UUID = GetUniqueID()
    } else {
        UUID = uuid
    }

    // fileBaseName := filepath.Base(cmd)
    logFileDir := GetConfig("app", "logdir")
    logFilePath := fmt.Sprintf("%v/%v.log", logFileDir, UUID)
    // seelog.Debugf("写入执行结果日志 : %v", logFilePath)

    var head []byte
    var tail []byte
    if runtime.GOOS == "windows" {
        head = []byte(fmt.Sprintf("\r\n===BEGIN===%v===\r\n", beginTimeStr))
        tail = []byte(fmt.Sprintf("\r\n===END===%v===\r\n", endTimeStr))
    } else {
        head = []byte(fmt.Sprintf("\n===BEGIN===%v===\n", beginTimeStr))
        tail = []byte(fmt.Sprintf("\n===END===%v===\n", endTimeStr))
    }

    // 使用Buffer进行byte拼接,string也可用，效率高
    var buffer bytes.Buffer
    buffer.Write(head)
    buffer.Write(content)
    buffer.Write(tail)
    data := buffer.Bytes()

    // err := utils.WriteFile(logFilePath, data)
    err := AppendFile(logFilePath, data)
    if err != nil {
        // seelog.Errorf("执行结果日志写入失败 : %v", err.Error())
        return "", err
    }
    return logFilePath, nil
}

func SendHXMsg(uuid string, title string, tos string, data string) error {
    nowTime := time.Now()
    timeFormat := "20060102150405" // 时间格式化模板
    nowTimeStr := nowTime.Format(timeFormat)

    hxData := fmt.Sprintf("title=%v\ntop=%v\ncontent={\n[%v]\n%v\n}END\n", title, tos, nowTimeStr, data)
    seelog.Debugf("[%v]%v", uuid, hxData)

    filename := fmt.Sprintf("./data/msgfile/%v_%v.dat", RandomString(10), nowTimeStr)
    err := WriteFile(filename, hxData)
    if err != nil {
        seelog.Errorf("[%v]ERROR : %v", uuid, err.Error())
        return err
    }
    seelog.Infof("[%v]通知消息文件已生成[%v] ...", uuid, filename)

    err = ftpMsgFile(filename)
    if err != nil {
        seelog.Errorf("[%v]行信通知发送失败->ERROR:\n%v", uuid, err.Error())
        return err
    }
    seelog.Infof("[%v]行信通知发送成功 ...", uuid)

    return nil
}

func ftpMsgFile(file string) error {
    ftp := struct {
        Host      string
        Username  string
        Password  string
        RemoteDir string
    }{
        Host:      GetConfig("hx", "host"),
        Username:  GetConfig("hx", "username"),
        Password:  GetConfig("hx", "password"),
        RemoteDir: GetConfig("hx", "remotedir"),
    }

    var err error
    client, err := NewFtpClient(ftp.Host, ftp.Username, ftp.Password)
    if err != nil {
        //seelog.Errorf("ERROR : %v", err.Error())
        return err
    }
    err = client.ChangeDir(ftp.RemoteDir)
    if err != nil {
        //seelog.Errorf("ERROR : %v", err.Error())
        return err
    }

    fileHandle, err := os.Open(file)
    if err != nil {
        //seelog.Errorf("ERROR : %v", err.Error())
        return err
    }
    defer fileHandle.Close()

    fileBaseName := filepath.Base(file)
    err = client.Stor(fileBaseName, fileHandle)
    if err != nil {
        //seelog.Errorf("ERROR : %v", err.Error())
        return err
    }

    client.Logout()
    client.Quit()
    return nil
}
