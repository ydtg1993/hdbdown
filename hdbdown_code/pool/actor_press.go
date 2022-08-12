package pool

import (
	"encoding/json"
	"fmt"
	"github.com/beego/beego/v2/core/logs"
	"hdbdown/models"
	"hdbdown/mongo"
	"hdbdown/rd"
	"hdbdown/runtimes"
	"strconv"
)

func ImportUpdateActorData()  {
	if checkDisk() == false {
		return
	}

	total, err  := rd.LLen(models.MovieActorPress)
	if err != nil {
		logs.Error("ImportUpdateActorData 进程","redis 读取错误:", err.Error())
		return
	}



	var chans = make(chan int, 5)
	for i := 0; i< int(total); i ++  {
		if rd.CheckLock() {
			return
		}
		chans <- i
		
		name, err := rd.LPop(models.MovieActorPress)
		if err != nil {
			logs.Error("ImportUpdateActorData 进程","redis 读取错误:", err.Error())
			<- chans
			continue
		}

		err, actMongo := mongo.FindActor(name)
		if err != nil {
			logs.Error("ImportUpdateActorData 进程","mongo 数据读取失败:", err.Error(), name)
			continue
		}

		if actMongo == nil{
			<- chans
			continue
		}

		var actor models.MovieActor
		if err := models.GetGormDb().Where("name = ?", name).First(&actor).Error ; err != nil{
			logs.Error("ImportUpdateActorData 进程","movie_actor 读取错误:", err)
			<- chans
			continue
		}

		if actMongo.Gender == 2 {
			actor.Sex = "♀"
		}else {
			actor.Sex = "♂"
		}

		social, err := json.Marshal(actMongo.Interflow)
		if err != nil {
			logs.Error("ImportUpdateActorData 进程","movie_actor 数据解析错误:", err, actMongo.Interflow)
			<- chans
			continue
		}

		actor.SocialAccounts = string(social)
		actor.Photo = actMongo.Avatar

		rest := models.GetGormDb().Model(&actor).Where("id = ?", actor.Id).Updates(models.MovieActor{
			Photo:         actor.Photo,
			Sex:           actor.Sex,
			SocialAccounts: actor.SocialAccounts,
		})
		if rest.Error != nil {
			logs.Error("ImportUpdateActorData 进程","movie_actor 数据存储错误:", err, actor)
			<- chans
			continue
		}


		//开启子线程下载
		oDown := new(runtimes.DoRun)
		//下载结果
		var list []string
		_, uri := formatUrl(actor.Photo)
		list = append(list, uri)
		done, msg := oDown.Work(fmt.Sprintf("%s:%s:%s", actor.Name, "actor" , strconv.Itoa(actor.Id)), list)
		if done != true {
			logs.Error("演员下载记录失败,actor error:", msg)
		}

		<- chans
	}

}
