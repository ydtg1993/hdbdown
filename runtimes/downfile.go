package runtimes

import (
	"fmt"
	"hdbdown/common"
	"net/url"
	"os"
	"os/user"
	"runtime"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/beego/beego/v2/core/logs"
	beego "github.com/beego/beego/v2/server/web"
)

//定义一个并发的协程池通道，用来约束发送请求的并发数量
var maxthreads, _ = beego.AppConfig.Int("maxthreads")
var chans = make(chan string, maxthreads)

//定义的下载文件目录
var downpath, _ = beego.AppConfig.String("downpath")

//定义协程超时时间
var downTimeOut, _ = beego.AppConfig.Int64("downtimeout")

//所有者
var ownUser, _ = beego.AppConfig.String("ownuser")
var ownGroup, _ = beego.AppConfig.String("owngroup")

type DoRun struct {
	Sid       string     //协程的id
	Count     int        //协程完成的计数器
	Err       bool       //错误状态
	countLock sync.Mutex //并发锁
	Msg       string     //错误内容
}

func (this *DoRun) process(sUrl string) {

	//切割url，得到最后的文件名
	arrUrl := strings.Split(sUrl, `/`)
	fileName := arrUrl[len(arrUrl)-1]

	//获取中间的目录
	oUrl, oErr := url.Parse(sUrl)
	if oErr != nil {
		logs.Error("下载地址不能解析", sUrl, oErr.Error())
		this.Err = true
		return
	}
	//获取除掉域名后的地址
	sPath := oUrl.EscapedPath()
	//切割并去掉结尾的名称
	aPath := strings.SplitAfter(sPath, "/")
	a1Path := aPath[0 : len(aPath)-1]
	for i := 0; i < len(a1Path); i++ {
		a1Path[i] = strings.ReplaceAll(a1Path[i], "/", "")
	}

	filePath := downpath + strings.Join(a1Path, "/")

	//logs.Info("调试", sUrl, filePath)

	//下载文件
	_, down := common.DownFileFor5(sUrl, filePath, fileName)

	//异常情况处理,将协程状态变成异常
	if down == false {
		logs.Error("连续5次尝试下载文件失败", sUrl)
		this.Err = true
		this.Msg = fmt.Sprintf("下载失败,%s", sUrl)
		return
	}

	this.countLock.Lock()
	defer this.countLock.Unlock()

	//下载成功，更新计数器,防止并发计数，这里使用协程锁
	this.Count = this.Count + 1

	//修改文件夹授权
	oUser, oUErr := user.Lookup(ownUser)
	oGroup, oGErr := user.LookupGroup(ownGroup)

	if oUErr != nil {
		logs.Error("获取授权用户失败,请检测配置文件ownUser", oUErr.Error())
		return
	}
	if oGErr != nil {
		logs.Error("获取授权组失败,请检测配置文件ownUser", oGErr.Error())
		return
	}
	//获取用户id和组id
	Uid := oUser.Uid
	iUid, _ := strconv.Atoi(Uid)
	Gid := oGroup.Gid
	iGid, _ := strconv.Atoi(Gid)

	//修改文件夹所有者
	cErr := os.Chown(filePath, iUid, iGid)
	if cErr != nil {
		logs.Error("对文件夹授权失败", filePath, cErr.Error())
	}

	logs.Info("授权", filePath, iUid, iGid)
}

func (this *DoRun) handle(url string) {

	this.process(url)
	// 信号完成：开始启用下一个请求

	// 将缓冲区释放一个容量
	<-chans
}

/**
* 多协程池入口
* sId  		需要处理数据的id
* lists		需要下载的图片数组
 */
func (this *DoRun) Work(sId string, lists []string) (bool, string) {
	// 当通道已满的时候将被阻塞
	// 所以停在这里等待，直到有容量（被释放），才能继续去处理请求

	startTime := time.Now()

	this.Sid = sId

	for k, v := range lists {
		//开启协程，占用一个缓冲区容量
		chans <- fmt.Sprintf("%s:%d", sId, k)

		//对象赋值
		go this.handle(v)

		logs.Debug(("当前协程数量 runing: %d, 下载任务id %s 第%d个资源,%s 进入下载协程池 \n"), runtime.NumGoroutine(), sId, k+1, v)
	}

	res := false
	msg := "下载完成"
	//协程等待，直到计数器完成或者超时
	for {
		//协程执行完成
		if this.Count == len(lists) && this.Count > 0 {
			res = true
			break
		}

		//协程中存在错误
		if this.Err == true {
			logs.Debug("文件下载错误，请查看错误日志,查看相关id", sId)
			msg = this.Msg
			break
		}

		//多协程执行超时，设置一个超时时间
		diff := time.Now().Unix() - startTime.Unix()
		//协程超时最低不能小于10秒
		if downTimeOut < 10 {
			downTimeOut = 10
		}

		if diff >= downTimeOut {
			logs.Error("下载任务协程执行超时", sId)
			msg = fmt.Sprintf("下载任务协程执行超时,%s", sId)
			break
		}
	}

	logs.Debug("协程任务id", sId, "成功下载", this.Count, "条,协程耗时:", time.Now().Unix()-startTime.Unix())

	return res, msg
}
