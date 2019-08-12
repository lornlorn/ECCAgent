package utils

import (
    "time"

    "github.com/jlaffaye/ftp"
)

/*
NewFtpClient func(host string, username string, password string) (*ftp.ServerConn, error)
*/
func NewFtpClient(host string, username string, password string) (*ftp.ServerConn, error) {
    var err error
    conn, err := ftp.DialTimeout(host, 5*time.Second)
    if err != nil {
        return nil, err
    }
    err = conn.Login(username, password)
    if err != nil {
        return nil, err
    }
    err = conn.NoOp()
    if err != nil {
        return nil, err
    }
    return conn, nil
}
