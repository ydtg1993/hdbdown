package main

import (
	"fmt"
	"hdbdown/pool"
	"os"
	"runtime"
	"strings"

	//	_ "github.com/CodyGuo/godaemon"
	"github.com/beego/beego/v2/core/logs"
	beego "github.com/beego/beego/v2/server/web"
	"github.com/urfave/cli"
)

func main() {

	//日志天数
	logDay, _ := beego.AppConfig.Int("logday")
	//日志路径
	logpath, _ := beego.AppConfig.String("logpath")
	//日志级别
	loglevel, _ := beego.AppConfig.String("loglevel")

	if logDay < 1 {
		logDay = 7
	}
	if logpath == "." {
		logpath = ""
	}
	logName := fmt.Sprintf("%shdbdown.log", logpath)

	logErr := "error"
	if len(loglevel) > 1 {
		sLevel := strings.Split(loglevel, ",")
		logErr = `"` + strings.Join(sLevel, `","`) + `"`
	}

	level := 2
	if strings.ContainsAny(loglevel, "debug") == true {
		level = 7
	}

	logCfg := fmt.Sprintf(`{"filename":"%s","level":%d,"maxdays":%d,"separate":[%s]}`, logName, level, logDay, logErr)

	//记录日志
	logs.SetLogger(logs.AdapterMultiFile, logCfg)

	// 开始前的线程数
	logs.Debug(("线程数量 starting: %d\n"), runtime.NumGoroutine())

	//如果后面跟了参数，直接进行命令行操作
	if len(os.Args) > 1 {
		cliDo()
		fmt.Println("执行完成")
		return
	}

	go pool.TaskActor()

	pool.TaskDBDown()

}

/**
* 执行命令行操作
 */
func cliDo() {

	app := cli.NewApp()
	app.Version = "1.0.1"
	app.Name = "hdbdown"
	app.Usage = "参数"
	app.UsageText = "本程序专门为下载黄豆瓣的资源数据"
	app.ArgsUsage = ``

	app.Email = "qqc88.abo@gmail.com"
	app.Author = "abo"

	//这个方法就是这个命令已启动会运行什么
	app.Action = func(c *cli.Context) error {
		cli.ShowAppHelp(c) //这个是打印app的help界面
		return nil
	}
	app.Commands = []cli.Command{
		{
			Name:    "down",
			Aliases: []string{"-dn"},
			Usage: `手动更新数据表中未处理完下载的影片
					`,
			Action: func(c *cli.Context) error {
				//参数1

				pool.CliDo(false, "")

				return nil
			},
		},
		{
			Name:    "restart",
			Aliases: []string{"-rt"},
			Usage: `yyyy-mm-dd 4位年-2位月-2位日 重新下载指定时间段影片的图片  
					`,
			Action: func(c *cli.Context) error {
				//参数1
				first := c.Args().Get(0)
				second := c.Args().Get(1)

				if len(first) < 10 {
					fmt.Println("时间格式错误，请使用help查看帮助", first)
				}
				if len(second) > 2 {
					first = first + " " + second
				}

				fmt.Println("输入的时间", first)

				pool.CliDo(true, first)

				return nil
			},
		},
		{
			Name:    "actor",
			Aliases: []string{"-ac"},
			Usage: `手动更新数据表中未处理完下载的演员
					`,
			Action: func(c *cli.Context) error {
				//参数1

				pool.ActorDo(false, "")

				return nil
			},
		},
		{
			Name:    "actor_restart",
			Aliases: []string{"-acr"},
			Usage: `yyyy-mm-dd 4位年-2位月-2位日 重新下载指定时间段演员的图片  
					`,
			Action: func(c *cli.Context) error {
				//参数1
				first := c.Args().Get(0)
				second := c.Args().Get(1)

				if len(first) < 10 {
					fmt.Println("时间格式错误，请使用help查看帮助", first)
				}
				if len(second) > 2 {
					first = first + " " + second
				}

				fmt.Println("输入的时间", first)

				pool.ActorDo(true, first)

				return nil
			},
		},
		{
			Name:    "isload",
			Aliases: []string{"-is"},
			Usage: `yyyy-mm-dd 4位年-2位月-2位日 重新下载指定时间段演员的图片  
					`,
			Action: func(c *cli.Context) error {
				//参数1
				pool.LoadImg(false, "")

				return nil
			},
		},
	}

	app.Run(os.Args)
}
