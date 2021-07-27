package main

import (
	"crypto/tls"
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/yockii/cloud-ip-access/service"
	"github.com/yockii/cloud-ip-access/util"
	"io"
	"net"
	"net/http"
	"strings"
	"time"
)

func main() {
	app := fiber.New(fiber.Config{
		DisableStartupMessage: true,
		//Views:                 engine,
	})
	app.Use(recover.New(recover.Config{
		StackTraceHandler: func(e interface{}) {
			util.Log.Error(e)
		},
	}))

	var externalIP = ""

	app.Get("/", func(ctx *fiber.Ctx) error {
		eip := GetExternalIP()
		if eip == "" {
			return ctx.SendString("无法自动获取公网IP，请稍后再试")
		}
		if eip == externalIP {
			return ctx.SendString("IP地址已设置过，无需重新设置")
		}

		if err := service.SecurityGroupService.UpdateSecurityGroupManageIP(eip); err != nil {
			util.Log.Error(err)
			return ctx.SendString("设置ecs失败，请稍后再试!" + err.Error())
		}
		if err := service.RdsGroupService.UpdateRdsDBInstancesWhiteIps(eip); err != nil {
			util.Log.Error(err)
			return ctx.SendString("设置rds失败，请稍后再试!" + err.Error())
		}
		externalIP = eip
		return ctx.SendString("阿里云IP白名单已设置成功!")
	})

	app.Listen(fmt.Sprintf(":%d", util.Config.GetUint("server.port")))
}

func GetExternalIP() string {
	client := &http.Client{
		Timeout:   3 * time.Second,
		Transport: &http.Transport{TLSClientConfig: &tls.Config{InsecureSkipVerify: true}},
	}

	resp, err := client.Get("http://ip.cip.cc/")
	if err != nil {
		fmt.Println(err)
		resp, err = client.Get("https://ifconfig.co/ip")
		if err != nil {
			fmt.Println(err)
			resp, err = client.Get("http://ifconfig.me/ip")
			if err != nil {
				fmt.Println(err)
				resp, err = client.Get("https://www.taobao.com/help/getip.php")
				if err != nil {
					fmt.Println(err)
					return ""
				}
			}
		}
	}
	defer resp.Body.Close()
	ip, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println(err)
		return ""
	}
	ipstr := strings.ReplaceAll(string(ip), "ipCallback({ip:\"", "")
	ipstr = strings.ReplaceAll(ipstr, "\"})", "")
	return strings.TrimSpace(string(ipstr))
}

func GetLocalIP() string {
	conn, _ := net.Dial("udp", "8.8.8.8:80")
	defer conn.Close()
	localAddr := conn.LocalAddr().String()
	fmt.Println(localAddr)
	idx := strings.LastIndex(localAddr, ":")
	return localAddr[:idx]
}
