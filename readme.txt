appname = hdbdown
httpport = 8080
runmode = pro
autorender = false
copyrequestbody = true
EnableDocs = true

mysqlhost = "127.0.0.1"         //数据库地址
mysqluser = "yellowDouban"      //数据库用户名
mysqlpass = "yellowDouban"      //数据库密码
mysqldbname   = "yellowDoubanDB"  //数据库名
mysqllifetime = 15              //连接空闲断开时间（秒）
mysqlidletime = 120             //连接生命周期（秒）
mysqlmaxconn = 200              //连接池最大连接数

redishost = "127.0.0.1"        //redis地址
redisport = 6379               //redis端口
redispass = "root"				//redis的密码
redisdb = "1"                  //redis选择的数据库
rediskey = "collectionresources:info:queue"    //队列的键值

downpath = "./img/"				//资源文件图片等下载地址
downsleep = 120					//计划任务间隔时间（秒）
downtimeout = 60               //每一个线程下载超时时间（秒）
downlimit = 10000              //每次读取队列的长度（条）

maxthreads = 300				//线程池最大并发数

logdays = 1						//日志保留时间（秒）
logpath = /Users/jack/			//日志目录，以/结尾
loglevel = error				//日志错误等级，error只记录错误，debug记录调试数据