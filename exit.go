package main

import (
	"fmt"
	"github.com/astaxie/beego/logs"
	"os"
	"path"
	"time"
)

// 超时或者日志上传完毕关闭管道，退出各个goroutine
func goroutineExit() {
	defer waitGroup.Done()
	defer close(hadpChan)
	defer close(lzopChan)
	defer close(rsyncChan)

	ticker := time.NewTicker(5 * time.Second)
	delPath := path.Dir(path.Dir(appConf.localAddr))
	delCmdFmt := `/usr/bin/find %s -name "*.log.*"|grep %s |xargs rm -f`
	delCmd := fmt.Sprintf(delCmdFmt, delPath, timeStamp)
	for range ticker.C {

		if int(time.Now().Unix()-timeStart) > (appConf.timeout * 60) {
			_, err := ExecCmdLocal(delCmd)
			if err != nil {
				logs.Error("execute cmd:%s error:%v", delCmd, err)
				return
			}
			logs.Error("timeout: more than %s mins. delete local logs success! del cmd: %s.", appConf.timeout, delCmd)
			os.Exit(0)
			return
		}

		if sucessNum == int32(len(appConf.hadoopClients)*len(hostMap)) {
			_, err := ExecCmdLocal(delCmd)
			if err != nil {
				logs.Error("execute cmd:%s error:%v", delCmd, err)
				return
			}
			logs.Info("Push all log success! delete local logs success! del cmd: %s.", delCmd)
			return
		}
	}
}
