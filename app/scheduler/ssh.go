package scheduler

import (
    "bytes"
    "fmt"
    "github.com/cihub/seelog"
    "golang.org/x/crypto/ssh"
    "io/ioutil"
    "net"
    "strings"
    "time"
)

func SSHConnect(user string, password string, host string, port int, key string) (*ssh.Session, error) {
    var (
        auth         []ssh.AuthMethod
        addr         string
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
        User:    user,
        Auth:    auth,
        Timeout: 10 * time.Second,
        Config:  config,
        HostKeyCallback: func(hostname string, remote net.Addr, key ssh.PublicKey) error {
            return nil
        },
        //HostKeyCallback: ssh.InsecureIgnoreHostKey(),
    }

    addr = fmt.Sprintf("%s:%d", host, port)

    if client, err = ssh.Dial("tcp", addr, clientConfig); err != nil {
        return nil, err
    }

    if session, err = client.NewSession(); err != nil {
        return nil, err
    }

    modes := ssh.TerminalModes{
        ssh.ECHO:          0,
        ssh.TTY_OP_ISPEED: 14400,
        ssh.TTY_OP_OSPEED: 14400,
    }
    if err := session.RequestPty("xterm", 80, 40, modes); err != nil {
        return nil, err
    }

    return session, nil
}

func SSHRun(cmdstr string) {
    seelog.Debugf("SSHRun...%v", cmdstr)
    session, err := SSHConnect("root", "", "ecsgn.raydio.site", 22, "F:\\WorkSpace\\id_rsa")
    if err != nil {
        seelog.Error(err.Error())
        return
    }
    defer session.Close()

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
        cmd = cmd + "\n"
        stdinBuf.Write([]byte(cmd))
    }
    session.Wait()

    seelog.Trace(outbt.String())
    seelog.Trace(errbt.String())
    //return
    seelog.Trace("finish...")
}
