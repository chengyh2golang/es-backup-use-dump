package main

import (
	"es-backup/esops"
	"es-backup/utils"
	"flag"
	"fmt"
	"log"
	"os"
	"os/exec"
	"time"
)

var (
	help bool
	backupDestBaseDir string // "/tmp/es-backup/"
	expired int // 7
	logDestDir  string // "/tmp/es-task-log"
	esAddr string //192.268.250.9
	esPort string //9200
)

func init()  {
	flag.BoolVar(&help, "help", false, "es backup help usage")
	flag.StringVar(&backupDestBaseDir, "dst", "/tmp/es-backup","index backup destination directory")
	flag.IntVar(&expired, "expired", 7,"expired time, day")
	flag.StringVar(&logDestDir, "log", "/tmp/es-task-log","log file directory")
	flag.StringVar(&esAddr, "es", "192.168.250.9","es server ip address")
	flag.StringVar(&esPort, "port", "9200","es server port")
	flag.Usage = usage
}

func usage() {
	_, _ = fmt.Fprintf(os.Stderr, `backup version: 1.0
Options:
`)
	flag.PrintDefaults()
}

func esIndexBackup(logDestDir,backupDestBaseDir,esAddr,esPort string,expired int) {
	//检查指定的log目录是否存在，如果不存在，就创建出这个目录
	logDirExists, err := utils.CheckExists(logDestDir)
	if err != nil {
		log.Fatal(err)
	}
	if !logDirExists {
		err := os.MkdirAll(logDestDir, os.ModePerm)
		if err != nil {
			log.Fatal(err)
		}
	}

	dayFormat := time.Now().Format("2006-01-02")
	logFileName := logDestDir + "es-backup-" + dayFormat + ".log"

	exists, err := utils.CheckExists(logFileName)
	if err != nil {
		log.Fatal(err)
	}
	if !exists {
		file, err := os.Create(logFileName)

		defer func() {
			if err != nil {
				file.Close()
			} else {
				err = file.Close()
			}
		}()
		if err != nil {
			panic(err)
		}
	}

	f, err := os.OpenFile(logFileName, os.O_RDWR | os.O_CREATE | os.O_APPEND, 0666)
	defer f.Close()
	log.SetOutput(f)

	//获取es中的日志相关的索引
	indexArr, err := esops.FetchIndices(esAddr, esPort)

	if err != nil {
		log.Printf("从es获取索引失败：%v\n",err)
	}

	var backupIndicesArr []string

	//遍历索引slice，判断是否需要备份，如果满足备份条件，执行备份
	for _,indexName := range indexArr {

		expiredBool := utils.IsNeedBackup(indexName, expired)
		if expiredBool {

			//判断备份目标目录是否存在，如果不存在，创建出来
			//目标目录的格式：/tmp/es-backup/2020.04.01
			indexDay := utils.FetchDay(indexName)
			destPath := backupDestBaseDir + indexDay
			exists, _ := utils.CheckExists( destPath )
			if !exists {
				err := os.MkdirAll(destPath, os.ModePerm)
				if err != nil {
					log.Printf("创建备份文件夹: %v 失败：%v\n",destPath,err)
				}
			}
			destPath = utils.PathWrapper(destPath)
			destFileName := destPath + indexName + ".json"

			//检查要生成的索引的json文件是否存在，如果已经存在，将其改名为xx.bak
			exists, _ = utils.CheckExists( destFileName )
			if exists {
				err := os.Rename(destFileName, destFileName+".bak")
				if err != nil {
					log.Printf("创建备份文件: %v 失败：%v\n",destFileName +".bak",err)
				}
			}

			//执行备份
			//构建docker的执行命令string
			//dockerCmd := "docker run --rm -v " + destPath + ":/tmp/es-backup-data " +
			//	"docker.io/taskrabbit/elasticsearch-dump:latest" +
			//	" --input=http://" + esAddr + ":" + esPort + "/" + indexName +
			//	" --output=" + "/tmp/es-backup-data/" + indexName + ".json" +
			//	" --type=data"

			//在docker容器内部执行，使用该命令
			dockerCmd := "elasticdump " +
				" --input=http://" + esAddr + ":" + esPort + "/" + indexName +
				" --output=" + destFileName +
				" --type=data"


			cmd := exec.Command("/bin/sh", "-c", dockerCmd)
			err = cmd.Run()
			if err == nil {
				backupIndicesArr = append(backupIndicesArr,indexName)
			} else {
				log.Printf("运行elasticdump报错：%v\n",err)
			}

		}
	}

	//删除过期的索引
	if len(backupIndicesArr) > 0 {
		err = esops.DeleteIndex(esAddr, esPort, backupIndicesArr)
		if err == nil {
			log.Printf("本次备份任务备份了%v个索引：\n",len(backupIndicesArr))
			for _,v := range backupIndicesArr {
				log.Println(v)
			}
	} else {
			log.Printf("删除索引失败：%v\n",err)
		}
	}
}

func main() {

	flag.Parse()
	if help {
		flag.Usage()
	} else {
		backupDestBaseDir = utils.PathWrapper(backupDestBaseDir)
		logDestDir = utils.PathWrapper(logDestDir)
		esIndexBackup(logDestDir,backupDestBaseDir,esAddr,esPort,expired)
	}


}
