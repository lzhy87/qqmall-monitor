package utils

import (
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/astaxie/beego/config"
)

var (
	serviceName string
	mailTo      []string
	//邮件主题
	subject = "后端服务出现故障"
	// 邮件正文
	body string
	err  error
	l    net.Listener
	//配置文件内容
	hostname   string
	apiurl     string
	serverport string
	//Apitime 配置文件中的检测时间间隔
	apitime           int64
	serverrestarttime int64
	user              string
	pass              string
	host              string
	port              string
	mailto            string
	serviceNames      []string
)

func init() {

	hostname, _ = os.Hostname()
	conf, err := config.NewConfig("ini", "app.conf")
	if err != nil {
		Logs.Error("config配置件读取失败%v", err)
		log.Fatalf("config配置件读取失败%v", err)
	}

	serviceName = conf.String("monitor::checkapiurl")
	serviceNames = strings.Fields(serviceName)
	serverport = conf.String("monitor::serverport")
	apitime, _ = conf.Int64("monitor::apitime")
	serverrestarttime, _ = conf.Int64("monitor::serverrestarttime")
	user = conf.String("mail::user")
	pass = conf.String("mail::pass")
	host = conf.String("mail::host")
	port = conf.String("mail::port")
	mailto = conf.String("mail::mailto")
	mailTo = strings.Fields(mailto)

	// if err != nil {
	// 	Logs.Error("获取hostname失败：", err)

	// }
}

//MonitorPort 监控指定端口
// func MonitorPort() {

// 	//file, err := os.Create("./monitorPort.pprof")
// 	//if err != nil {
// 	//	fmt.Printf("create cpu pprof failed, err:%v\n", err)
// 	//	return
// 	//}
// 	//pprof.WriteHeapProfile(file)
// 	//
// 	//defer file.Close()

// 	var num int
// 	sendFlag := false

// 	for {
// 		time.Sleep(time.Second * time.Duration(apitime))
// 		l, err = net.Listen("tcp", ":"+serverport)

// 		//服务正常绑定端口，服务启动成功
// 		if err != nil {
// 			Logs.Info(serverport + "已经绑定")
// 			num = 0
// 			continue
// 		}
// 		//检测3次，每次间隔时间为serverrestarttime；如果3次检测依然为成功，则发送邮件
// 		num++
// 		err = l.Close()
// 		if err != nil {
// 			fmt.Println("network Listen关闭错误", err)
// 		}
// 		time.Sleep(time.Second * time.Duration(serverrestarttime))
// 		Logs.Warning("服务端口检测绑定异常，正在重新检测中,第%v次检测中。。。", num)

// 		if num >= 3 && !sendFlag {
// 			num = 0
// 			sendFlag = true
// 			body = fmt.Sprintf("（可能正在重启服务）%v服务未启动，端口监听未启用,端口号为：%v，重新检测时间为:%v秒，检查次数为3次", serviceNames[0], serverport, serverrestarttime)
// 			Logs.Error(body)
// 			err = SendMail(mailTo, subject, body)
// 			if err != nil {
// 				Logs.Error("邮件发送失败：%v", err)
// 				return
// 			}

// 		} else if sendFlag {
// 			num = -10
// 			sendFlag = false
// 		}

// 	}

// }

//CheckService 调用API接口获取服务信息情况
func CheckService() {
	//file, err := os.Create("./checkSerice.pprof")
	//if err != nil {
	//	fmt.Printf("create cpu pprof failed, err:%v\n", err)
	//	return
	//}
	//pprof.WriteHeapProfile(file)
	//defer file.Close()
	num1 := 1
	sendFlag := false
	for {
		time.Sleep(time.Second * time.Duration(apitime))
		for _, s := range serviceNames {
			req, err := http.Get(s)
			if err != nil {
				Logs.Error("服务已经关闭，无法连接上" + s)
				num1++

				time.Sleep(time.Second * time.Duration(serverrestarttime))
				Logs.Warning("服务接口检测失败，正在重新检测中,第%v次检测中。。。", num1)
				fmt.Println(num1)
				if num1 >= 3 && !sendFlag {
					num1 = 0

					err := DingToInfo("告警：" + s + "服务出现问题，请及时处理")

					if err != true {
						Logs.Error("钉钉告警信息发送失败！")
						return
					}
					// Logs.Error("******邮件已发送*********服务未调用成功,错误代码:%v", req.Status)
					Logs.Error("******钉钉告警信息已发送*********服务未调用成功,错误代码:%v", req.Status)
					sendFlag = true
					continue
				} else if sendFlag {
					num1 = -10
					sendFlag = false
				}
				continue
			}

			if req.StatusCode == 200 {
				// num1 = 0
				Logs.Info(s + "服务正常运行中")

			}
			// req.Body.Close()
		}
	}
}
