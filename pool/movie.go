package pool

import (
	"encoding/json"
	"hdbdown/common"
	"hdbdown/models"
	"hdbdown/runtimes"
	"math"
	"strconv"
	"strings"
	"time"

	"github.com/beego/beego/v2/core/logs"
)

/**
* 读数据库下载资源文件
 */
func CliDo(restart bool, rTime string) {

	//进程开始时间
	startTime := time.Now()

	//定义数据更新对象
	MOVIE := new(models.CollectionMovie)

	//每次读取条数
	All := MOVIE.Total(restart, rTime)
	pageSize := 10
	pageAll := int(math.Ceil(float64(All) / float64(pageSize)))

	//数据库中读取数据
	lastId := 0
	v := []*models.CollectionMovie{}
	for i := 0; i < pageAll; i++ {

		lastId, v = MOVIE.Lists(lastId, pageSize, restart, rTime)

		//线程池
		lenPool := len(v)
		threadPool := make(chan interface{}, lenPool)
		defer close(threadPool)

		logs.Warning("当前页数", i, "影片实际开启下载处理线程池", lenPool, len(threadPool))

		for _, obj := range v {

			obj.Map = strings.ReplaceAll(obj.Map, `\`, "")
			arrMap := []interface{}{}
			jsonByte := []byte(obj.Map)
			jsonErr := json.Unmarshal(jsonByte, &arrMap)

			if jsonErr != nil {
				logs.Error(obj.Id, "影片map的数据解析错误：", jsonErr.Error())
			}

			//定义需要下载的文件列表
			lists := []string{}
			//下载封面图
			if len(obj.BigCove) > 1 {
				_, bimg := formatUrl(obj.BigCove)
				lists = append(lists, bimg)
			}
			//下载小封面
			if len(obj.SmallCover) > 1 {
				_, bimg := formatUrl(obj.SmallCover)
				lists = append(lists, bimg)
			}
			//下载预告片
			if len(obj.Trailer) > 1 {
				_, bimg := formatUrl(obj.Trailer)
				lists = append(lists, bimg)
			}

			//遍历map，影片组图
			for _, val := range arrMap {
				v1, _ := val.(map[string]interface{})
				simg := common.UnknowToString(v1["img"])
				if len(simg) > 1 {
					_, simg = formatUrl(simg)
					lists = append(lists, simg)
				}
				bimg := common.UnknowToString(v1["big_img"])
				if len(bimg) > 1 {
					_, bimg = formatUrl(bimg)
					lists = append(lists, bimg)
				}
			}

			//设置一个开始时间
			sTime := time.Now()

			//用于下载完成后，写入数据库的数据
			D := new(dbQueue)
			D.Mid = obj.Id
			D.Type = "javdb"
			D.Big_cove = obj.BigCove
			D.Small_cover = obj.SmallCover
			D.Trailer = obj.Trailer
			D.Map = arrMap

			bytes, _ := json.Marshal(D)
			da := strings.ToLower(string(bytes))

			//过滤掉json字符串的域名
			if strings.Contains(da, downdomain) == true {
				da = strings.ReplaceAll(da, downdomain, "/")
			}

			//拼接需要处理的数据
			vv := new(redisQueue)
			vv.Sid = strconv.Itoa(obj.Id)
			vv.Ty = "javdb"
			vv.Txt = da
			vv.Lists = lists

			//推一个任务进入线程
			threadPool <- vv
			go func(vv *redisQueue) {

				if len(vv.Lists) > 0 {
					//开启子线程下载
					oDown := new(runtimes.DoRun)
					//下载结果
					done, msg := oDown.Work(vv.Sid, vv.Lists)

					//只有下载成功后，才会更新状态
					if done == true {
						MOVIE.Save(vv.Sid, "2", vv.Txt, "{}")
					} else {
						MOVIE.Save(vv.Sid, "3", vv.Txt, "{}")
						logs.Error("影片下载记录失败,down error:", vv.Sid, msg, vv.Txt)
					}

					logs.Debug("影片线程下载结果:", vv.Sid, done)
				} else {
					logs.Error("影片获取下载列表失败,list error:", vv.Sid, vv.Txt)
				}

				//跑完一个线程，关闭一条
				<-threadPool
			}(vv)

			//这里等待线程执行完成
			for {
				//执行完成
				if len(threadPool) == 0 {
					break
				}
				//设置一个超时时间，超时时间设置成5分钟
				diff := time.Now().Unix() - sTime.Unix()

				if diff >= 300 {
					logs.Error("影片当前总任务", All, "当前页数", i, "剩余通道数", len(threadPool), "线程列表数据下载超时", obj)
					break
				}

			}

		}

		//线程池完成一组
		logs.Warning("影片当前总任务", All, "当前页数", i, "完成线程数", lenPool)
	}

	logs.Warning("本次下载完成，影片", All, "条,耗时：", time.Now().Unix()-startTime.Unix())
}
