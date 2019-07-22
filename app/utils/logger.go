package utils

import (
    "github.com/cihub/seelog"
)

/*
InitLogger initial a logger by seelog
Default config file SeelogCfg = "./config/seelog.xml"
*/
func InitLogger(path string) error {
    defer seelog.Flush()

    logger, err := seelog.LoggerFromConfigAsFile(path)
    if err != nil {
        return err
    }
    seelog.ReplaceLogger(logger)

    return nil
}
