package mongo

import (
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

/**
actor_name:華原永遠,
group:有碼,
actual_source:https://www.javbus.com/star/afu,
act:,
gender:2,
interflow:{},
ctime:2022-07-14 16:41:43,
birthday:,
age:,
height:,
cup:D,
bust:87cm,
waistline:62cm,
hip:86cm,
hobby:,
birthplace:,
avatar:javbus.com/actor/62cfd6c72467fa8f8256100d/actors.jpg
*/
type ActorBase struct {
	ActorName    string            `json:"actor_name" bson:"actor_name"`
	Group        string            `json:"group" bson:"group"`
	ActualSource string            `json:"actual_source" bson:"actual_source"`
	Act          string            `json:"act" bson:"act"`
	Gender       int               `json:"gender" bson:"gender"`
	Interflow    map[string]string `json:"interflow" bson:"interflow"`
	Ctime        string            `json:"ctime" bson:"ctime"`
	Birthday     string            `json:"birthday" bson:"birthday"`
	Age          interface{}       `json:"age" bson:"age"`
	Height       string            `json:"height" bson:"height"`
	Cup          string            `json:"cup" bson:"cup"`
	Bust         string            `json:"bust" bson:"bust"`
	Waistline    string            `json:"waistline" bson:"waistline"`
	Hip          string            `json:"hip" bson:"hip"`
	Birthplace   string            `json:"birthplace" bson:"birthplace"`
	Avatar       string            `json:"avatar" bson:"avatar"`
}

func FindActor(name string) (err error, act *ActorBase) {
	findOptions := options.FindOne()

	rest := DBClient.Collection("javdb_actor").FindOne(context.TODO(), bson.M{"actor_name": name}, findOptions).Decode(&act)
	if rest != nil && rest != mongo.ErrNoDocuments {
		err = rest
		return
	}

	return
}
