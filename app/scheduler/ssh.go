package scheduler

import (
    "app/models"
    "app/utils"
    "github.com/cihub/seelog"
    "golang.org/x/crypto/ssh"
    "io/ioutil"
    "net"
    "strings"
)

func SSHConnect(host string, user string, password string, key string) (*ssh.Session, error) {
    var (
        auth []ssh.AuthMethod
        //addr         string
        clientConfig *ssh.ClientConfig
        client       *ssh.Client
        config       ssh.Config
        session      *ssh.Session
        err          error
    )
    auth = make([]ssh.AuthMethod, 0)
    if key == "" {
        auth = append(auth, ssh.Password(password))
    } else {
        pemBytes, err := ioutil.ReadFile(key)
        if err != nil {
            return nil, err
        }
        var signer ssh.Signer
        if password == "" {
            signer, err = ssh.ParsePrivateKey(pemBytes)
        } else {
            signer, err = ssh.ParsePrivateKeyWithPassphrase(pemBytes, []byte(password))
        }
        if err != nil {
            return nil, err
        }
        auth = append(auth, ssh.PublicKeys(signer))
    }

    clientConfig = &ssh.ClientConfig{
        User: user,
        Auth: auth,
        //Timeout: 60 * time.Second,
        Timeout: 0,
        Config:  config,
        HostKeyCallback: func(hostname string, remote net.Addr, key ssh.PublicKey) error {
            return nil
        },
        //HostKeyCallback: ssh.InsecureIgnoreHostKey(),
    }

    //addr = fmt.Sprintf("%s:%d", host, port)

    if client, err = ssh.Dial("tcp", host, clientConfig); err != nil {
        return nil, err
    }

    if session, err = client.NewSession(); err != nil {
        return nil, err
    }

    /*
       // Get term size
       fd := int(os.Stdin.Fd())
       oldState, err := terminal.MakeRaw(fd)
       if err != nil {
           panic(err)
       }
       defer terminal.Restore(fd, oldState)
       termWidth, termHeight, err := terminal.GetSize(fd)
       if err != nil {
           panic(err)
       }

    */
    modes := ssh.TerminalModes{
        ssh.ECHO:          1,
        ssh.TTY_OP_ISPEED: 14400,
        ssh.TTY_OP_OSPEED: 14400,
    }
    termHeight := 80
    termWidth := 80
    if err := session.RequestPty("xterm", termHeight, termWidth, modes); err != nil {
        return nil, err
    }

    return session, nil
}

//func SSHRun(host string, auth string, privkey string, cmdstr string, hxTos string, cronName string) error {
func CronSSHRun(cron models.SysCron) error {
    UUID := utils.GetUniqueID()
    seelog.Infof("[%v]SSH->[%v@%v]->[%v] Begin ...", UUID, cron.CronAuth, cron.CronHost, cron.CronCmd)

    var user string
    var password string

    userInfo := strings.Split(cron.CronAuth, "/")
    if len(userInfo) == 1 {
        user = userInfo[0]
    } else if len(userInfo) == 2 {
        user = userInfo[0]
        password = userInfo[1]
    } else {
        seelog.Errorf("[%v]SSH Auth ERROR : [%v@%v]", UUID, cron.CronAuth, cron.CronHost)
    }

    session, err := SSHConnect(cron.CronHost, user, password, cron.CronPrivkey)
    if err != nil {
        seelog.Errorf("[%v]ERROR : %v", UUID, err.Error())
        utils.SendHXMsg(cron.CronName, cron.CronHx, err.Error())
        return err
    }
    defer session.Close()

    buf, err := session.CombinedOutput(cron.CronCmd)
    if err != nil {
        seelog.Errorf("[%v]ERROR : %v", UUID, err.Error())
        utils.SendHXMsg(cron.CronName, cron.CronHx, err.Error())
        return err
    }
    seelog.Debug(string(buf))

    /*
       cmdlist := strings.Split(cmdstr, ";")
       stdinBuf, err := session.StdinPipe()
       if err != nil {
           seelog.Error(err.Error())
           return
       }

       var outbt, errbt bytes.Buffer
       session.Stdout = &outbt
       session.Stderr = &errbt

       err = session.Shell()
       if err != nil {
           seelog.Error(err.Error())
           return
       }

       for _, cmd := range cmdlist {
           stdinBuf.Write([]byte(cmd + "\n"))
       }
       stdinBuf.Close()
       session.Wait()

       seelog.Trace(outbt.String())
       seelog.Trace(errbt.String())
       //return

    */

    utils.SendHXMsg(cron.CronName, cron.CronHx, string(buf))

    seelog.Infof("[%v]SSH->[%v@%v]->[%v] Finish ...", UUID, cron.CronAuth, cron.CronHost, cron.CronCmd)

    return nil
}
