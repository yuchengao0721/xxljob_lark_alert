package main

import (
	"fmt"
	"strconv"
	"time"
	"xxl_job_alert/alertinit"
	"xxl_job_alert/alertsender"

	"github.com/gofiber/fiber/v2"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func IsNotFound(err error) bool {
	s, ok := status.FromError(err)
	if !ok {
		return ok
	}
	if s.Code() == codes.NotFound {
		return true
	}

	return false
}
func main() {
	fmt.Println("程序启动")
	// 全局设置
	alertinit.Init()
	sender := new(alertsender.FeishuSender)
	alertsender.InitializeConnectionPools()
	go startScheduledTask(sender)
	app := fiber.New(fiber.Config{
		Prefork:       false, //docker环境下千万别开，会导致程序执行闪退
		CaseSensitive: true,
		StrictRouting: true,
		ServerHeader:  "Fiber",
		AppName:       alertinit.Conf.Application.Name})
	// 测试是否成功的接口
	app.Get("/healthy", func(c *fiber.Ctx) error {
		return c.SendString("true")
	})
	// 测试是否成功的接口
	app.Get("/test", func(c *fiber.Ctx) error {
		var result = sender.SendMsg()
		// 将对象序列化为 JSON 字节
		return c.JSON(result)
	})
	port := alertinit.Conf.Application.Port
	fmt.Println("程序启动成功")
	app.Listen(":" + strconv.Itoa(port))

}
func startScheduledTask(sender *alertsender.FeishuSender) {
	// 创建一个定时任务，每隔一段时间执行一次
	ticker := time.NewTicker(time.Duration(alertinit.Conf.Interval) * time.Second)
	defer ticker.Stop()
	for {
		select {
		case <-ticker.C:
			// 在这里执行你的定时任务逻辑
			fmt.Println("开始执行", time.Now())
			sender.SendMsg()
		}
	}
}
