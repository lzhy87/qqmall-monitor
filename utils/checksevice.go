package utils

import (
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/astaxie/beego/config"
	"gopkg.in/gomail.v2"
)

var (
	// mailTo1 = []string{
	// 	"1610469455@qq.com", //杜潇
	// 	"55900695@qq.com",   //李梓仪
	// }

	mailTo []string
	//邮件主题
	subject = "qqcmall后端服务出现故障"
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
)

func init() {

	hostname, _ = os.Hostname()
	conf, err := config.NewConfig("ini", "app.conf")
	if err != nil {
		Logs.Error("config配置件读取失败%v", err)
		log.Fatalf("config配置件读取失败%v", err)
	}

	apiurl = conf.String("monitor::checkapiurl")
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
func MonitorPort() {

	//file, err := os.Create("./monitorPort.pprof")
	//if err != nil {
	//	fmt.Printf("create cpu pprof failed, err:%v\n", err)
	//	return
	//}
	//pprof.WriteHeapProfile(file)
	//
	//defer file.Close()

	var num int
	sendFlag := false

	for {
		time.Sleep(time.Second * time.Duration(apitime))
		l, err = net.Listen("tcp", ":"+serverport)

		//服务正常绑定端口，服务启动成功
		if err != nil {
			Logs.Info(serverport + "已经绑定")
			num = 0
			continue
		}
		//检测3次，每次间隔时间为serverrestarttime；如果3次检测依然为成功，则发送邮件
		num++
		err = l.Close()
		if err != nil {
			fmt.Println("network Listen关闭错误", err)
		}
		time.Sleep(time.Second * time.Duration(serverrestarttime))
		Logs.Warning("服务端口检测绑定异常，正在重新检测中,第%v次检测中。。。", num)

		if num >= 3 && !sendFlag {
			num = 0
			sendFlag = true
			body = fmt.Sprintf("（可能正在重启服务）qqcmall服务未启动，端口监听未启用,端口号为：%v，重新检测时间为:%v秒，检查次数为3次", serverport, serverrestarttime)
			Logs.Error(body)
			err = SendMail(mailTo, subject, body)
			if err != nil {
				Logs.Error("邮件发送失败：%v", err)
				return
			}

		} else if sendFlag {
			num = -10
			sendFlag = false
		}

	}

}

//CheckService 调用API接口获取服务信息情况
func CheckService() {
	//file, err := os.Create("./checkSerice.pprof")
	//if err != nil {
	//	fmt.Printf("create cpu pprof failed, err:%v\n", err)
	//	return
	//}
	//pprof.WriteHeapProfile(file)
	//defer file.Close()

	var num1 int
	sendFlag := false
	for {
		time.Sleep(time.Second * time.Duration(apitime))
		req, err := http.Get(apiurl)
		if err != nil {
			Logs.Warning("服务已经关闭，无法连接上" + apiurl)
			continue
		}
		req.Body.Close()

		if req.StatusCode == 200 {
			num1 = 0

			Logs.Info("服务正常运行中")
		} else {
			num1++
			time.Sleep(time.Second * time.Duration(serverrestarttime))
			Logs.Warning("服务接口检测失败，正在重新检测中,第%v次检测中。。。", num1)

			if num1 >= 3 && !sendFlag {
				num1 = 0
				body = "服务未调用成功,错误代码:" + req.Status
				err := SendMail(mailTo, subject, body)
				if err != nil {
					Logs.Error("邮件发送失败：%v", err)
					return
				}
				Logs.Error("******邮件已发送*********服务未调用成功,错误代码:%v", req.Status)
				sendFlag = true
			} else if sendFlag {
				num1 = -10
				sendFlag = false
			}
		}

	}
}

//SendMail 发送邮件
func SendMail(mailTo []string, subject string, body string) error {
	//定义邮箱服务器连接信息，如果是阿里邮箱 pass填密码，qq邮箱填授权码
	mailConn := map[string]string{
		// "user": "55900695@qq.com",
		// "pass": "ynxxuwgmpeovbgbd",
		// "host": "smtp.qq.com",
		// "port": "465",
		"user": user,
		"pass": pass,
		"host": host,
		"port": port,
	}

	port, _ := strconv.Atoi(mailConn["port"]) //转换端口类型为int

	m := gomail.NewMessage()
	m.SetHeader("From", "qqcmall server "+hostname+" <"+mailConn["user"]+">") //这种方式可以添加别名，即“XD Game”， 也可以直接用<code>m.SetHeader("From",mailConn["user"])</code> 读者可以自行实验下效果
	m.SetHeader("To", mailTo...)                                              //发送给多个用户
	m.SetHeader("Subject", subject)                                           //设置邮件主题
	m.SetBody("text/html", body)                                              //设置邮件正文

	d := gomail.NewDialer(mailConn["host"], port, mailConn["user"], mailConn["pass"])

	err := d.DialAndSend(m)
	return err

}

/*  //定义收件人
  mailTo := []string {
 "zhangqiang@xxx.com",
 "abc@qq.com",
"sssdd@qq.com",
 }
//邮件主题为"Hello"
 subject := "Hello"
// 邮件正文
 body := "Good"
 SendMail(mailTo, subject, body)*/
