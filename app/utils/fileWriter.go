package utils

import (
    "errors"
    "io/ioutil"
    "os"

    "github.com/cihub/seelog"
)

/*
WriteFile func(filename string, content string)
*/
func WriteFile(filename string, content string) error {
    data := []byte(content)
    err := ioutil.WriteFile(filename, data, 0644)
    if err != nil {
        return err
    }
    return nil
}

/*
AppendFile func(filepath string, content []byte) error
*/
func AppendFile(filepath string, content []byte) error {
    fileHandler, err := os.OpenFile(filepath, os.O_RDWR|os.O_APPEND|os.O_CREATE, 0644)
    if err != nil {
        return err
    }
    defer fileHandler.Close()

    n, err := fileHandler.Write(content)
    if err != nil {
        return err
    } else if n != len(content) {
        return errors.New("Writen bytes Not Equal The Length Of Content")
    }
    seelog.Tracef("n:%v,len:%v,%v", n, len(content), err)
    return nil
}
