package alertsender

import (
	"bytes"
	"html/template"
	"io/ioutil"
	"xxl_job_alert/alertinit"
	"xxl_job_alert/alertmodel"

	"fmt"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/imroc/req/v3"
	"github.com/rs/zerolog/log"
)

const feishu_tplPath = "./etc/xxl_job_alert/conf/lark.tpl"

type FeishuSender struct{}

// 飞书发送消息
func (s *FeishuSender) SendMsg() []alertmodel.XXLAlertResult {
	//配置项里面的飞书token必填,或者填写grafana里面通知媒介里面Message内的fs_rebot_token值
	fs_tokens := make(alertmodel.Set)
	if strings.TrimSpace(alertinit.Conf.LarkToken) != "" {
		var arr = strings.Split(alertinit.Conf.LarkToken, ",")
		fs_tokens.AddArr(arr)
	}
	if len(fs_tokens) == 0 {
		return nil
	}
	var list = QueryList()
	if len(list) == 0 {
		return nil
	}
	// 常规的告警
	client := req.C().DevMode()
	var feishu_card alertmodel.FeishuCard
	content, _ := create_common_alert_content(list)
	if len(content) == 0 {
		return nil
	}

	feishu_card = alertmodel.CreateFsCard("⚠️  告警通知", content, "S1")
	for token := range fs_tokens {
		if len(token) > 0 {
			feishu_url := fmt.Sprintf("https://open.feishu.cn/open-apis/bot/v2/hook/%s", token)
			resp, err := client.R().
				SetHeader(fiber.HeaderContentType, fiber.MIMEApplicationJSONCharsetUTF8).
				SetHeader(fiber.HeaderHost, "open.feishu.cn").
				SetBody(feishu_card).
				Post(feishu_url)
			if err != nil {
				log.Error().Msgf("飞书通知异常了: %v", err)
			}
			if !resp.IsSuccessState() {
				log.Error().Msgf("飞书通知失败了: %v", err)
			}
		}

	}
	return list
}

// 拼接通知的消息内容
func create_common_alert_content(list []alertmodel.XXLAlertResult) (string, error) {
	content, err := ioutil.ReadFile(feishu_tplPath)
	if err != nil {
		fmt.Printf("读取文件失败：%v\n", err)
		return "", err
	}
	tpl := string(content)
	t, err := template.New("feishu").Parse(tpl)
	if err != nil {
		log.Error().Msgf("模板加载错误:%v", err)
		return "", err
	}
	// 解析模板
	var buf bytes.Buffer
	// 应用模板并输出结果
	err = t.Execute(&buf, list)
	if err != nil {
		log.Error().Msgf("模板执行错误:%v", err)
		return "", err
	}
	return buf.String(), nil
}
