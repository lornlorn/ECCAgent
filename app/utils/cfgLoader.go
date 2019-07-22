package utils

import (
    "github.com/cihub/seelog"
    "github.com/robfig/config"
)

/*
Config Global Configurations
*/
var cfg *config.Config

/*
InitConfig func(path string) error
Initialize The Config Global Variable
*/
func InitConfig(path string) error {
    conf, err := config.ReadDefault(path)
    if err != nil {
        // seelog.Errorf("Read Config File [%v] Fail : %v", path, err)
        return err
    }

    cfg = conf

    return nil
}

/*
GetConfig func(section string, option string) string
Return String
*/
func GetConfig(section string, option string) string {
    value, err := cfg.String(section, option)
    if err != nil {
        seelog.Errorf("Get Config [%v].[%v] Fail : %v", section, option, err)
        return ""
    }

    return value
    // result is string "http://www.example.com/some/path"

    // c.Int("service-1", "maxclients")
    // // result is int 200

    // c.Bool("service-1", "delegation")
    // // result is bool true

    // c.String("service-1", "comments")
    // // result is string "This is a multi-line\nentry"
}

/*
GetConfigInt func(section string, option string) int
Return String
*/
func GetConfigInt(section string, option string) int {
    value, err := cfg.Int(section, option)
    if err != nil {
        seelog.Errorf("Get Config [%v].[%v] Fail : %v", section, option, err)
        return 0
    }

    return value
    // result is string "http://www.example.com/some/path"

    // c.Int("service-1", "maxclients")
    // // result is int 200

    // c.Bool("service-1", "delegation")
    // // result is bool true

    // c.String("service-1", "comments")
    // // result is string "This is a multi-line\nentry"
}
