package esops

import (
	"context"
	"errors"
	"es-backup/utils"
	"github.com/olivere/elastic/v7"
	"strings"
)

func FetchIndices(esAddress, esPort string) ([]string,error) {
	var result []string
	ready := utils.CheckServiceReady(esAddress, esPort)
	if ready {
		esURL := "http://" + esAddress + ":" + esPort
		esClient, err := elastic.NewClient(
			elastic.SetURL(esURL),
			//sniff是提供给客户端用来维护集群状态的，现在客户端是在容器里运行
			// 而es集群在k8s集群外面
			//这样就无法sniff，所以把这个sniff给关闭掉
			elastic.SetSniff(false),
			//如果es开启了用户认证，需要添加对应的用户名和密码设置
			//elastic.SetBasicAuth("elastic","root123"),
		)
		if err != nil {
			panic(err)
		}
		response, err := esClient.IndexNames()
		if err != nil {
			return result,err
		}
		for _,v := range response {
			if !utils.IsDirNameStartWithDot(v) && strings.Contains(v,"-20") {
				result = append(result,v)
			}
		}
		return result,nil
	}
	return result,errors.New("es service not ready")
}

func DeleteIndex (esAddress, esPort string, indices []string) error {
	ready := utils.CheckServiceReady(esAddress, esPort)
	if ready {
		esURL := "http://" + esAddress + ":" + esPort
		esClient, err := elastic.NewClient(
			elastic.SetURL(esURL),
			//sniff是提供给客户端用来维护集群状态的，现在客户端是在容器里运行
			// 而es集群在k8s集群外面
			//这样就无法sniff，所以把这个sniff给关闭掉
			elastic.SetSniff(false),
			//如果es开启了用户认证，需要添加对应的用户名和密码设置
			//elastic.SetBasicAuth("elastic","root123"),
		)
		if err != nil {
			return err
		}
		//通过待删除的索引的数组构造出删除索引的字符串
		//构造的结果是：index1,index2,index3
		deleteIndexStr := ""
		for _,index := range indices {
			deleteIndexStr += index + ","
		}
		deleteIndexStr = strings.TrimSuffix(deleteIndexStr,",")

		//执行索引删除
		_, err = esClient.DeleteIndex(deleteIndexStr).Do(context.Background())
		if err != nil {
			return err
		}
		return nil
	}
	return errors.New("delete index failed")
}


