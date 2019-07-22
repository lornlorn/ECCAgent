package udfuncs

import (
    "app/utils"
    "bytes"
    "fmt"
    "runtime"
)

// WriteRunLog2File func(nowTime time.Time, cmd string) error
// Write Execute Result To Log File
func WriteRunLog2File(uuid string, content []byte, beginTimeStr string, endTimeStr string) (string, error) {
    // logFileFormat := "20060102150405" // 日志文件命名时间格式化
    // logFileNameTime := nowTime.Format(logFileFormat)
    var UUID string
    if uuid == "" || len(uuid) != 32 {
        UUID = utils.GetUniqueID()
    } else {
        UUID = uuid
    }

    // fileBaseName := filepath.Base(cmd)
    logFileDir := utils.GetConfig("app", "logdir")
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
    err := utils.AppendFile(logFilePath, data)
    if err != nil {
        // seelog.Errorf("执行结果日志写入失败 : %v", err.Error())
        return "", err
    }
    return logFilePath, nil
}
