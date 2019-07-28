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
    conn, err := ftp.DialTimeout(host, time.Second*30)
    if err != nil {
        return &ftp.ServerConn{}, err
    }
    err = conn.Login(username, password)
    if err != nil {
        return &ftp.ServerConn{}, err
    }
    return conn, nil
}