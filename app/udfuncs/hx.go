package udfuncs

import (
	"app/utils"
    "fmt"
    "github.com/cihub/seelog"
	"os"
	"path/filepath"
)

type FTP struct {
	Host      string
	Username  string
	Password  string
	RemoteDir string
}

func SendHXMsg(title string, tos string, data string) error {

    hxData := fmt.Sprintf("title=%v\ntop=%v\ncontent={\n%v\n}END\n", title, tos, data)
    seelog.Debug(hxData)

    filename := fmt.Sprintf("./data/msgfile/%v_%v.dat", "CHKURL", nowDataStr)
    err := utils.WriteFile(filename, hxData)
    if err != nil {
        seelog.Errorf("ERROR : %v", err.Error())
        return err
    }
    seelog.Debugf("通知消息文件已生成[%v] ...", filename)


    err = ftpMsgFile(filename)
	if err != nil {
		seelog.Errorf("ERROR : %v", err.Error())
		return err
	}
	return nil
}

func ftpMsgFile(file string) error {
	Ftp := FTP{
		Host:      utils.GetConfig("ftp", "host"),
		Username:  utils.GetConfig("ftp", "username"),
		Password:  utils.GetConfig("ftp", "password"),
		RemoteDir: utils.GetConfig("ftp", "remotedir"),
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

	fileHandle, err := os.Open(file)
	if err != nil {
		seelog.Errorf("ERROR : %v", err.Error())
		return err
	}
	defer fileHandle.Close()

	fileBaseName := filepath.Base(file)
	err = client.Stor(fileBaseName, fileHandle)
	if err != nil {
		seelog.Errorf("ERROR : %v", err.Error())
		return err
	}

	client.Logout()
	client.Quit()
	return nil
}
