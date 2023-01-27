package mongo

import (
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"hdbdown/global/orm"
	"strconv"
)

/**
domain_id:2,
domain_name:www.javlibrary.com,
group:,
release_time:2021-05-07,
video_time:240分钟,
producer:プラム,
publisher:素人オンリープラム,
series:,
score:0,
score_man:0,
video_sort:[
  "女同性恋",
...
],
actor:[],
direct:null,
sell:,
uid:RS-061,
video_url:/cn/?v=javme5lsqy,
small_cover:www.javlibrary.com/60a1d9a13b7c7dc3639e6ce6/small_cover.jpg,
video_title:LGBTびっくり！6人のレズビアン（後編）「わたしが女の素敵なお下品を教えてあげる…」,
actual_source:http://www.javlibrary.com/cn/?v=javme5lsqy,
big_cover:www.javlibrary.com/60a1d9a13b7c7dc3639e6ce6/big_cover.jpg,
comments:[],
preview_img:[
  	"javdb.com/62a0ffd9a79fcf81649ed991/preview_img_1.jpg"
  ...
],
preview_big_img:[
	"javdb.com/62a0ffd9a79fcf81649ed991/preview_big_img_1.jpg"
  ...
]
trailer:www.javlibrary.com/60a1d9a13b7c7dc3639e6ce6/trailer.mp4,
ctime:2021-05-17 02:55:38
*/
type MovieBase struct {
	//Id            MongoID      `json:"_id" bson:"_id"`
	Uid           string       `json:"uid" bson:"uid"`
	DomainName    string       `json:"domain_name" bson:"domain_name"`
	Group         string       `json:"group" bson:"group"`
	ReleaseTime   string       `json:"release_time" bson:"release_time"`
	VideoTime     string       `json:"video_time" bson:"video_time"`
	VideoTitle    string       `json:"video_title" bson:"video_title"`
	Producer      string       `json:"producer" bson:"producer"`
	Publisher     string       `json:"publisher" bson:"publisher"`
	Series        []string     `json:"series" bson:"series"`
	Score         interface{}  `json:"score" bson:"score"`
	ScoreMan      interface{}  `json:"score_man" bson:"score_man"`
	VideoSort     interface{}  `json:"video_sort" bson:"video_sort"`
	Actor         interface{}  `json:"actor" bson:"actor"`
	Direct        string       `json:"direct" bson:"direct"`
	Sell          string       `json:"sell" bson:"sell"`
	Comments      interface{}  `json:"comments" bson:"comments"`
	VideoUrl      string       `json:"video_url" bson:"video_url"`
	Magnet        []MagnetMode `json:"magnet" bson:"magnet"`
	Ctime         string       `json:"ctime" bson:"ctime"`
	Utime         string       `json:"utime" bson:"utime"`
	ActualSource  string       `json:"actual_source" bson:"actual_source"`
	PreviewImg    []string     `json:"preview_img" bson:"preview_img"`
	PreviewBigImg []string     `json:"preview_big_img" bson:"preview_big_img"`
	SmallCover    string       `json:"small_cover" bson:"small_cover"`
	BigCover      string       `json:"big_cover" bson:"big_cover"`
	Trailer       string       `json:"trailer" bson:"trailer"`
}

// mongodb  返回的数据据类型不确定，正确的数据类型为 float32
func (mv *MovieBase) ScoreTypeChange() (op float32) {
	switch mv.Score.(type) {
	case float32:
		op, _ = mv.Score.(float32)
		break
	case string:
		scores := mv.Score.(string)
		sc, err := strconv.ParseFloat(scores, 32)
		if err != nil {
			op = 0
		}
		op = float32(sc)
		break
	default:
		op = 0
		break
	}

	return
}

// mongodb  返回的数据据类型不确定，正确的数据类型为 float32
func (mv *MovieBase) ScoreManTypeChange() (op float32) {
	switch mv.ScoreMan.(type) {
	case float32:
		op, _ = mv.ScoreMan.(float32)
		break
	case string:
		scores := mv.ScoreMan.(string)
		sc, err := strconv.ParseFloat(scores, 32)
		if err != nil {
			op = 0
		}
		op = float32(sc)
		break
	default:
		op = 0
		break
	}
	return
}

// mongodb  返回的数据据类型不确定，正确的数据类型为 []string
func (mv *MovieBase) ActorTypeChange() (op []map[string]string) {
	if pa, ok := mv.Actor.(primitive.A); ok {
		valueMSI := []interface{}(pa)
		for _, v := range valueMSI {
			var actor = map[string]string{}

			switch v.(type) {
			case primitive.A:
				actorWithGender := []interface{}(v.(primitive.A))
				if actorWithGender[0] != nil {
					actor["name"] = actorWithGender[0].(string)
					actor["sex"] = "♀"
				}
				if actorWithGender[1] != nil {
					actor["sex"] = actorWithGender[1].(string)
				}
				break
			case string:
				actor["name"] = v.(string)
				actor["sex"] = "♀"
				break
			default:
				return
			}
			op = append(op, actor)
		}
		return
	}
	return
}

// mongodb  返回的数据据类型不确定，正确的数据类型为 []string
func (mv *MovieBase) VideoSortTypeChange() (op []string) {
	if pa, ok := mv.VideoSort.(primitive.A); ok {
		valueMSI := []interface{}(pa)

		for _, v := range valueMSI {
			op = append(op, v.(string))
		}
		return
	}
	return
}

/**
{
  "commentator": "hyhaiml",
  "comment_time": "2021-05-15 16:29:56",
  "comment_text": "纯粹就是喜剧片",
  "score": null
}
*/
type Comments struct {
	Commentator string      `json:"commentator" bson:"commentator"`
	CommentTime string      `json:"comment_time" bson:"comment_time"`
	CommentText string      `json:"comment_text" bson:"comment_text"`
	Score       interface{} `json:"score" bson:"score"`
}

type MagnetMode struct {
	Name      string      `json:"name" bson:"name"`
	Time      string      `json:"time" bson:"time"`
	Url       string      `json:"url" bson:"url"`
	Meta      string      `json:"meta" bson:"meta"`
	IsSmall   interface{} `json:"is-small" bson:"is-small"`
	IsWarning interface{} `json:"is-warning" bson:"is-warning"`
}

/**
通过番号获取数据
*/
func GetMovieDataByNumber(dbName string, number string) (err error, mv *MovieBase) {
	findOptions := options.FindOne()

	err = orm.DBClient.Collection(dbName).FindOne(context.TODO(), bson.M{"uid": number}, findOptions).Decode(&mv)
	if err != nil && err != mongo.ErrNoDocuments {
		return
	}
	return
}

/**
统计  ctime > lastTime 的数据总和
*/
func Count(dbName string, lastTime string) (err error, num int64) {
	filter := bson.D{{"ctime", bson.D{{"$gt", lastTime}}}}
	num, err = orm.DBClient.Collection(dbName).CountDocuments(context.TODO(), filter)
	if err != nil {
		return
	}
	return
}

/**
获取 ctime > lastTime 的数据
*/
func Find(dbName string, lastTime string, currentPage int64, limit int64) (err error, list []*MovieBase) {
	findOptions := options.Find()
	skip := currentPage*limit - limit
	findOptions.SetLimit(limit)
	findOptions.SetSkip(skip)

	projection := bson.D{
		{"ctime", 1},
		{"utime", 1},
		{"utime", 1},
		{"uid", 1},
		{"big_cover", 1},
		{"small_cover", 1},
	}
	findOptions.SetProjection(projection)

	sortMap := make(map[string]interface{}, 0)
	sortMap["ctime"] = 1
	findOptions.SetSort(sortMap)

	filter := bson.D{{"ctime", bson.D{{"$gt", lastTime}}}}
	cur, err := orm.DBClient.Collection(dbName).Find(context.TODO(), filter, findOptions)
	if err != nil {
		return
	}

	for cur.Next(context.TODO()) {
		//raw := cur.Current
		//fmt.Println(string(raw))
		var tmp MovieBase
		if err = cur.Decode(&tmp); err != nil {
			return
		}
		list = append(list, &tmp)
	}

	return
}

/**
projection := bson.D{
		{"video_sort", 1},
		{"uid", 1},
	}
*/
func FindWithCondition(dbName string, lastTime string, currentPage int64, limit int64, fields interface{}) (err error, list []*MovieBase) {
	findOptions := options.Find()
	findOptions.SetLimit(limit)
	findOptions.SetSkip(currentPage * (limit - 1))

	if fields != nil {
		findOptions.SetProjection(fields)
	}

	sortMap := make(map[string]interface{}, 0)
	sortMap["ctime"] = 1
	findOptions.SetSort(sortMap)

	filter := bson.D{{"ctime", bson.D{{"$gt", lastTime}}}}
	cur, err := orm.DBClient.Collection(dbName).Find(context.TODO(), filter, findOptions)
	if err != nil {
		return
	}

	for cur.Next(context.TODO()) {
		var tmp MovieBase
		if err = cur.Decode(&tmp); err != nil {
			return
		}
		list = append(list, &tmp)
	}

	return
}
