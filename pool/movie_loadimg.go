//影片表检测和下载图片
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
* 检测文件，不存在的去补下
 */
func LoadImg(restart bool, rTime string) {

	//进程开始时间
	startTime := time.Now()

	//定义数据更新对象
	MOVIE := new(models.Movie)
	//定义更新未成功的对象
	CM := new(models.CollectionMovie)

	//每次读取条数
	All := MOVIE.Total(restart, rTime)
	pageSize := 10
	pageAll := int(math.Ceil(float64(All) / float64(pageSize)))

	//数据库中读取数据
	lastId := 0
	v := []*models.Movie{}
	for i := 0; i < pageAll; i++ {

		lastId, v = MOVIE.Lists(lastId, pageSize, restart, rTime)

		//线程池
		lenPool := len(v)
		threadPool := make(chan interface{}, lenPool)
		defer close(threadPool)

		logs.Warning("扫描影片数据", i, "处理线程池", lenPool, len(threadPool))

		for _, obj := range v {

			//定义一个值，用来标记是否存在必填字段缺少图的情况
			chk := true

			obj.Map = strings.ReplaceAll(obj.Map, `\`, "")
			arrMap := []interface{}{}
			jsonByte := []byte(obj.Map)
			jsonErr := json.Unmarshal(jsonByte, &arrMap)

			if jsonErr != nil {
				logs.Error(obj.Id, "影片map的数据解析错误：", jsonErr.Error())
			}

			//定义需要下载的文件列表
			lists := []string{}
			//下载封面图[必填]
			if len(obj.BigCove) > 1 {
				if ImgIsLoad(obj.BigCove) == true {
					_, bimg := formatUrl(obj.BigCove)
					lists = append(lists, bimg)
				}

			} else {
				chk = false
			}
			//下载小封面[必填]
			if len(obj.SmallCover) > 1 {
				if ImgIsLoad(obj.SmallCover) == true {
					_, bimg := formatUrl(obj.SmallCover)
					lists = append(lists, bimg)
				}

			} else {
				chk = false
			}
			//下载预告片
			if len(obj.Trailer) > 1 {
				if ImgIsLoad(obj.Trailer) == true {
					_, bimg := formatUrl(obj.Trailer)
					lists = append(lists, bimg)
				}
			}

			//遍历map，影片组图
			for _, val := range arrMap {
				v1, _ := val.(map[string]interface{})
				simg := common.UnknowToString(v1["img"])
				if len(simg) > 1 {
					if ImgIsLoad(simg) == true {
						_, simg = formatUrl(simg)
						lists = append(lists, simg)
					}
				}
				bimg := common.UnknowToString(v1["big_img"])
				if len(bimg) > 1 {
					if ImgIsLoad(bimg) == true {
						_, bimg = formatUrl(bimg)
						lists = append(lists, bimg)
					}
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
			vv.Number = obj.Number

			//推一个任务进入线程
			threadPool <- vv
			go func(vv *redisQueue) {

				if len(vv.Lists) > 0 {

					logs.Debug("扫描数据需要下载 id:", vv.Sid, len(vv.Lists), "条资源")

					//开启子线程下载
					oDown := new(runtimes.DoRun)
					//下载结果
					done, msg := oDown.Work(vv.Sid, vv.Lists)

					//下载完成，并且必填的也不是空的
					if done == true && chk == true {
						MOVIE.Save(vv.Sid, "1")
					} else {
						MOVIE.Save(vv.Sid, "2")
						CM.SaveWithNumber(vv.Number, "3")
						logs.Info("文件下载失败:", vv.Sid, msg)
						if chk == false {
							logs.Info("缺少必填的封面:", vv.Sid, msg)
						}
					}

					logs.Debug("扫描下载结果:", vv.Sid, done)
				} else {
					logs.Debug("扫描数据不需要下载:", vv.Sid)
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
					logs.Error("扫描当前总任务", All, "当前页数", i, "剩余通道数", len(threadPool), "线程列表数据下载超时", obj)
					break
				}

			}

		}

		//线程池完成一组
		logs.Warning("影片当前总任务", All, "当前页数", i, "完成线程数", lenPool)
	}

	logs.Warning("本次下载完成，影片", All, "条,耗时：", time.Now().Unix()-startTime.Unix())
}
