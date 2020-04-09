package utils

import (
	"net"
	"os"
	"strings"
	"time"
)

const (
	BaseDayTime = "2006.01.02"
)

func IsDirNameStartWithDot(name string) bool {
	return strings.HasPrefix(name,".")
}

//检查端口是否监听
//ipAndPort: 192.168.1.1:3306
func CheckPortReady(ipAndPort string) bool {
	_, err := net.Dial("tcp", ipAndPort)
	if err != nil {
		return false
	}
	return true
}

//检查文件或者文件夹是否存在
func CheckExists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}

//检查服务是否就绪，如果不就绪，重复检查5次，间隔5秒
func CheckServiceReady(ip,port string) bool {
	service := ip + ":" + port
	for i:=0;i<5;i++ {
		if CheckPortReady(service) {
			return true
		} else {
			i += 1
			time.Sleep(5 * time.Second)
		}
	}
	return false
}

//根据超时时间设置：比如7天
//再根据日志文件名：env6-default-nginx-base-67bd6995f6-2020.04.03,解析出该文件生成的日期
//判断该文件是否需要备份
func IsNeedBackup(fileName string, expiredDay int) bool {
	//获取该文件的日期
	fileDaytimeStrArr:=  strings.Split(fileName,"-")
	fileDaytimeStr := fileDaytimeStrArr[len(fileDaytimeStrArr)-1]

	fileDaytime, _ := time.Parse(BaseDayTime, fileDaytimeStr)

	nowDaytimeStr := time.Now().Format(BaseDayTime)
	nowDaytime, _ := time.Parse(BaseDayTime, nowDaytimeStr)

	if (fileDaytime.Unix() + int64(86400 * expiredDay) ) <= nowDaytime.Unix() {
		return true
	} else {
		return false
	}
}

//根据索引文件名，解析出索引创建的日期
//env6-default-nginx-base-67bd6995f6-2020.04.03
//解析结果：2020.04.03
func FetchDay(fileName string) string {
	//获取该文件的日期
	fileDaytimeStrArr:=  strings.Split(fileName,"-")
	return fileDaytimeStrArr[len(fileDaytimeStrArr)-1]

}

//判断字符串是否是以某个字符为结尾
//用于判断用户输入的路径是/tmp/back还是/tmp/back/
//如果不是以"/"结尾，加上
func PathWrapper(inputStr string) string {
	if strings.LastIndex(inputStr, "/") == len(inputStr) -1 {
		return inputStr
	} else {
		return inputStr + "/"
	}
}
