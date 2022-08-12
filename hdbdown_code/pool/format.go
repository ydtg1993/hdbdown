package pool

import (
	_ "github.com/joho/godotenv/autoload"
	"hdbdown/common"
	"os"
	"strings"
)

var downdomain string
var downpath string

func init() {
	downdomain= os.Getenv("downdomain")
	downpath = os.Getenv("downpath")
}

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
获取 arr 中 index = n 的值，如果不存在返回空
*/
func getValue(index int, arr []string) string  {
	for k, v := range arr {
		if k == index {
			return v
		}
	}
	return ""
}