package scheduler

import (
    "app/models"
    "github.com/cihub/seelog"
    "github.com/rfyiamcool/cronlib"
)

// Cron global CronScheduler
var Cron *cronlib.CronSchduler

// InitCron func()
func InitCron() error {
    Cron = cronlib.New()

    crons, err := models.GetCrons()
    if err != nil {
        seelog.Errorf("Get Cron List Error : %v", err.Error())
        return err
    }

    // seelog.Debug(crons)
    for idx, cron := range crons {
        seelog.Debugf("Cron Job %v : %v", idx, cron)

        // 复制对象
        cronObj := cron

        switch cronObj.CronType { //finger is declared in switch
        case "SSH":
            //seelog.Trace("SSH")
            job, err := cronlib.NewJobModel(
                cronObj.CronSpec,
                func() {
                    //Execute("cron", "", cronCmd, cronEnvs, cronArgs...)
                    CronSSHRun(cronObj)
                },
                cronlib.AsyncMode(),
            )
            if err != nil {
                seelog.Errorf("Cron Set Fail : [%v]", cronObj.CronName)
                return err
            }

            err = Cron.Register(cronObj.CronName, job)
            if err != nil {
                seelog.Errorf("Cron Register Error : %v", err.Error())
                return err
            }
        case "SQL":
            seelog.Trace("SQL")
        case "URL":
            //seelog.Trace("URL")
            job, err := cronlib.NewJobModel(
                cronObj.CronSpec,
                func() {
                    CronCheckUrl(cronObj)
                },
                cronlib.AsyncMode(),
            )
            if err != nil {
                seelog.Errorf("Cron Set Fail : [%v]", cronObj.CronName)
                return err
            }

            err = Cron.Register(cronObj.CronName, job)
            if err != nil {
                seelog.Errorf("Cron Register Error : %v", err.Error())
                return err
            }
        case "PORT":
            //seelog.Trace("PORT")
            job, err := cronlib.NewJobModel(
                cronObj.CronSpec,
                func() {
                    PortScan(cronObj)
                },
                cronlib.AsyncMode(),
            )
            if err != nil {
                seelog.Errorf("Cron Set Fail : [%v]", cronObj.CronName)
                return err
            }

            err = Cron.Register(cronObj.CronName, job)
            if err != nil {
                seelog.Errorf("Cron Register Error : %v", err.Error())
                return err
            }
        case "SCAN":
            seelog.Trace("SCAN")

        default: //default case
            seelog.Warn("incorrect cron type")
        }
    }

    Cron.Start()
    //Cron.Wait()
    // cron.Join()
    return nil
}

// AddCronJob func(cron models.NewCron) error
func AddCronJob(cron models.NewCron) error {
    seelog.Debugf("Set New Job : %v", cron)
    cronName := cron.CronName
    cronSpec := cron.CronSpec
    //cronEnvs := strings.Split(cron.CronEnvs, " ")
    var cronEnvs []string
    cronCmd := cron.CronCmd
    //cronArgs := strings.Split(cron.CronArgs, " ")
    var cronArgs []string
    //cronUuid := cron.CronUuid

    job, err := cronlib.NewJobModel(
        cronSpec,
        func() {
            Execute("cron", "", cronCmd, cronEnvs, cronArgs...)
        },
        cronlib.AsyncMode(),
    )
    if err != nil {
        seelog.Errorf("Cron Set Fail : [%v]", cronName)
        return err
    }

    err = Cron.UpdateJobModel(cronName, job)
    if err != nil {
        seelog.Errorf("Cron Register Error : %v", err.Error())
        return err
    }

    return nil
}

// DelCronJob func(cron models.NewCron)
func DelCronJob(cron models.NewCron) {
    StopCronJob(cron.CronName)
    UnregisterCronJob(cron.CronName)
}

// StopCronJob func(cronName string)
func StopCronJob(cronName string) {
    Cron.StopService(cronName)
}

// UnregisterCronJob func(cronName string)
func UnregisterCronJob(cronName string) {
    Cron.UnRegister(cronName)
}

// UpdateCronJob func(cron models.NewCron) error
func UpdateCronJob(cron models.NewCron) error {
    seelog.Debugf("Update Job : %v", cron)
    cronName := cron.CronName
    cronSpec := cron.CronSpec
    var cronEnvs []string
    cronCmd := cron.CronCmd
    //cronArgs := strings.Split(cron.CronArgs, " ")
    var cronArgs []string
    //cronUuid := cron.CronUuid

    job, err := cronlib.NewJobModel(
        cronSpec,
        func() {
            Execute("cron", "", cronCmd, cronEnvs, cronArgs...)
        },
        cronlib.AsyncMode(),
    )
    if err != nil {
        seelog.Errorf("Cron Update Fail : [%v]", cronName)
        return err
    }

    err = Cron.UpdateJobModel(cronName, job)
    if err != nil {
        seelog.Errorf("Cron Register Error : %v", err.Error())
        return err
    }

    return nil
}
