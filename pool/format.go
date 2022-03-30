package pool

import (
	"encoding/json"
	"hdbdown/common"
	"strings"

	"github.com/beego/beego/v2/core/logs"
	beego "github.com/beego/beego/v2/server/web"
)

var downdomain, _ = beego.AppConfig.String("downdomain")
var downpath, _ = beego.AppConfig.String("downpath")

/**
* 判断文件是否需要下载
* @param	filename	文件目录+名称
 */
func ImgIsLoad(filename string) bool {
	//必须是本地文件
	if strings.Count(filename, "http") >= 1 {
		return false
	}

	allPath := downpath + "/" + filename

	//判断文件是否存在，不存在的文件需要下载
	chk := common.IsExist(allPath)
	if chk == false {
		return true
	}

	//文件大小小于1也加入队列
	size, _ := common.GetFileSize(allPath)
	if size < 1 {
		return true
	}

	return false
}

/**
* 针对url进行加工
 */
func formatUrl(sUrl string) (bool, string) {
	//如果url有问题,数据里面存在2次http
	if strings.Count(sUrl, "http") >= 2 {

		//包含http或者https，这个可能是重复的http头,去掉重复的http头
		if strings.LastIndex(sUrl, "http://") > 1 || strings.LastIndex(sUrl, "https://") > 1 {
			sUrl = strings.ReplaceAll(sUrl, downdomain, "")
		}
	}
	//不包含http头的，补上
	if strings.Count(sUrl, "http") <= 0 {
		sUrl = downdomain + sUrl
	}
	return true, sUrl
}

/**
* 针对json数据的解析
* @param 	string		str
* return	string		表id
* return	string		表类型
* return	[]string    需要下载的资源数组
 */
func decodeJson(str string) (string, string, []string) {

	sId := ""
	sTy := ""
	lists := []string{}

	jOut := map[string]interface{}{}
	//解析json数据
	err := json.Unmarshal([]byte(str), &jOut)
	if err != nil {
		logs.Error("解析json错误", str, err.Error())
		return "", "", lists
	}

	//获取结果
	sId = common.UnknowToString(jOut["mid"])
	sTy = common.UnknowToString(jOut["type"])

	smallCober := common.UnknowToString(jOut["small_cover"])
	if len(smallCober) > 1 {
		_, smallCober = formatUrl(smallCober)
		lists = append(lists, smallCober)
	}
	bigCober := common.UnknowToString(jOut["big_cove"])
	if len(bigCober) > 1 {
		_, bigCober = formatUrl(bigCober)
		lists = append(lists, bigCober)
	}
	trailer := common.UnknowToString(jOut["trailer"])
	if len(trailer) > 1 {
		_, trailer = formatUrl(trailer)
		lists = append(lists, trailer)
	}

	//针对演员
	avatar := common.UnknowToString(jOut["avatar"])
	if len(avatar) > 1 {
		_, avatar = formatUrl(avatar)
		lists = append(lists, avatar)
	}
	photo := common.UnknowToString(jOut["photo"])
	if len(photo) > 1 {
		_, photo = formatUrl(photo)
		lists = append(lists, photo)
	}
	aId := common.UnknowToString(jOut["aid"])
	if len(aId) > 0 {
		sId = aId
	}

	//处理map
	maps, _ := jOut["map"].([]interface{})
	//遍历map
	for _, val := range maps {
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

	//异常过滤,如果缺少大图或者小图，过滤掉这条记录的下载
	if sTy == "javdb" && (len(smallCober) < 1 || len(bigCober) < 1) {
		lists = []string{}
	}

	//返回最终结果
	return sId, sTy, lists
}
