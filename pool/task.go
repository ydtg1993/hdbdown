package pool

import (
	"time"
	//	beego "github.com/beego/beego/v2/server/web"
)

//定义数据结构
type redisQueue struct {
	Sid    string
	Ty     string
	Lists  []string
	Txt    string
	Number string
}

type dbQueue struct {
	Mid         int
	Type        string
	Small_cover string
	Big_cove    string
	Trailer     string
	Map         []interface{}
}

//每60分钟（从队列读取）
func TaskDBDown() {
	//队列key

	d := time.Duration(time.Second * time.Duration(3600))
	t := time.NewTicker(d)
	defer t.Stop()

	//启动更新数据到es
	CliDo(false, "")

	for {
		<-t.C
		CliDo(false, "")
	}
}

//每60分钟（从队列读取）
func TaskActor() {
	//队列key

	d := time.Duration(time.Second * time.Duration(3600))
	t := time.NewTicker(d)
	defer t.Stop()

	//启动更新数据到es
	ActorDo(false, "")

	for {
		<-t.C
		ActorDo(false, "")
	}
}
