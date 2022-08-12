package command

import (
	"fmt"
	_ "fmt"
	"github.com/360EntSecGroup-Skylar/excelize"
	_ "github.com/360EntSecGroup-Skylar/excelize"
	"github.com/golang-module/carbon"
	"go.mongodb.org/mongo-driver/bson"
	"hdbdown/models"
	"hdbdown/mongo"
	"hdbdown/rd"
	"math"
	"strings"
	"time"
)

const MovieTemporaryList = "movie_temporary_list"
const RedisKeyWithAllMongo = "redis_key_with_all_mongo"

func Run() {
	startTime := time.Now()
	f, err := excelize.OpenFile("./command/123.xlsx")
	if err != nil {
		fmt.Println(err)
		return
	}

	//redisTempSave()
	//redisMongoSave()

	rows := f.GetRows("qz-0730-2")
	totalExcel := len(rows)
	totalPress := 0

	for index, row := range rows {
		fmt.Printf("当前行：%d, 总行数:%d \n",index, totalExcel)
		var fh int
		if index == 0 {
			for k, v := range row {
				if v == "番号" {
					fh = k
				}
			}
			continue
		}
		number := row[fh]

		label , err := rd.HGet(RedisKeyWithAllMongo, number)
		if err != nil {
			continue
		}

		if label == "" {
			continue
		}

		ro := fmt.Sprintf("E%d", index + 1)
		f.SetCellValue("qz-0730-2",ro , label)
		totalPress ++
	}

	fmt.Printf("处理完成：%d, 总数：%d, 不存在:%d ,总耗时:%d \n",totalPress, totalExcel, totalExcel - totalPress, time.Now().Unix()- startTime.Unix())

	if err = f.Save(); err != nil {
		fmt.Println(err)
	}
}

func formatLabel(op []string) string {
	var labels []string
	var labelMap = make(map[string]string)
	for _, v := range op {
		_, ok := labelMap[v]
		if ok == false {
			v = strings.TrimSpace(v)
			labelMap[v] = v
			labels = append(labels, v)
		}
	}

	var str string
	for _, v := range labels {
		if str == "" {
			str = v
		}else {
			str = fmt.Sprintf("%s,%s", str, v)
		}
	}

	return str
}

func redisTempSave()  {
	var total int64
	var pageNum int64 = 1000
	err := models.GetGormDb().Model(models.TemporaryMovie{}).Count(&total).Error
	if err != nil {
		panic(err)
	}
	redisLen := rd.HLen(MovieTemporaryList)
	if redisLen == total {
		return
	}

	pageAll := math.Ceil(float64(total) / float64(pageNum))
	var lastId = 0
	var temp models.TemporaryMovie
	var temps []models.TemporaryMovie

	for i:=0; i < int(pageAll); i ++  {

		_, lastId ,temps = temp.ListOfAll(lastId, int(pageNum))
		for _, v := range temps {
			rd.HSet(MovieTemporaryList, v.Number,  v.DbName)
		}
	}
}

func redisMongoSave()  {
	javlist := []string{"javlibrary", "javbus", "javdb"}
	lastTime := carbon.CreateFromDateTime(2021, 1, 1, 0, 0, 0).ToDateTimeString()
	for _, dbName := range javlist {
		var pageNum int64 = 5000
		_ , total := mongo.Count(dbName, lastTime)
		pageAll := math.Ceil(float64(total) / float64(pageNum))

		for i:=0 ; i < int(pageAll); i ++  {
			fmt.Println("数据查询中....")
			fields := bson.D{
				{"video_sort", 1},
				{"uid", 1},
			}
			_, lists := mongo.FindWithCondition(dbName, lastTime, int64(i + 1), pageNum, fields)
			//fmt.Println(dbName, i, pageAll)
			for k, val := range lists {
				label := formatLabel(val.VideoSortTypeChange())
				fmt.Println(dbName,k, i, pageAll)
				if label == "" {
					fmt.Println("数据为空:" + label)
					continue
				}
				rd.HSet(RedisKeyWithAllMongo, val.Uid, label)
			}
		}
	}

}