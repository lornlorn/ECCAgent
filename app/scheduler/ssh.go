package scheduler

import (
    "app/models"
    "app/utils"
    "errors"
    "fmt"
    "github.com/cihub/seelog"
    "golang.org/x/crypto/ssh"
    "io/ioutil"
    "net"
    "strings"
)

type destSsh struct {
    name     string
    host     string
    user     string
    password string
    privkey  string
    command  string
    hxTos    string
}

func sshConnect(host string, user string, password string, key string) (*ssh.Session, error) {
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
func (s destSsh) sshRun() ([]byte, error) {

    session, err := sshConnect(s.host, s.user, s.password, s.privkey)
    if err != nil {
        return nil, err
    }
    defer session.Close()

    buf, err := session.CombinedOutput(s.command)
    if err != nil {
        return nil, err
    }

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

    return buf, nil
}

func SSHRun(obj interface{}) error {

    var dest destSsh

    switch obj.(type) {
    case models.SysCron:
        data := obj.(models.SysCron)
        var user string
        var password string

        userInfo := strings.Split(data.CronAuth, "/")
        if len(userInfo) == 1 {
            user = userInfo[0]
        } else if len(userInfo) == 2 {
            user = userInfo[0]
            password = userInfo[1]
        } else {
            seelog.Errorf("SSH->Wrong Auth Config : [%v@%v]", data.CronAuth, data.CronHost)
            return errors.New("Wrong Auth Config")
        }
        dest = destSsh{
            name:     data.CronName,
            host:     data.CronHost,
            user:     user,
            password: password,
            privkey:  data.CronPrivkey,
            command:  data.CronCmd,
            hxTos:    data.CronHx,
        }
    case string:

    default:
        return errors.New("Wrong Arg Type ...")
    }

    UUID := utils.GetUniqueID()
    seelog.Infof("[%v]SSH->[%v@%v]->[%v] Begin ...", UUID, dest.user, dest.host, dest.command)

    ret, err := dest.sshRun()
    if err != nil {
        seelog.Errorf("[%v]SSH->ERROR:\n%v", UUID, err.Error())
        utils.SendHXMsg(UUID, dest.name, dest.hxTos, fmt.Sprintf("IP:%v\nCMD:%v\nRESULT:%v", dest.host,dest.command, err.Error()))
        return err
    }
    seelog.Debugf("[%v]SSH->\n%v", UUID, string(ret))

    utils.SendHXMsg(UUID, dest.name, dest.hxTos, fmt.Sprintf("IP:%v\nCMD:%v\nRESULT:%v", dest.host,dest.command, string(ret)))

    seelog.Infof("[%v]SSH->[%v@%v]->[%v] Finish ...", UUID, dest.user, dest.host, dest.command)

    return nil
}
