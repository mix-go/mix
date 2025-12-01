package xutil

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/go-lark/lark/v2"
	"github.com/slack-go/slack"
)

// NotifierType 定义了支持的通知平台
type NotifierType string

const (
	NotifierTypeSlack NotifierType = "slack"
	NotifierTypeLark  NotifierType = "lark" // 飞书/Lark
)

// Notifier 是透明发送消息的结构体
type Notifier struct {
	Type NotifierType

	TitlePrefix string

	// Slack Webhook 专用的 URL
	SlackWebhookURL string

	// Lark 机器人专用的客户端（通常通过 Webhook URL 初始化）
	LarkBot *lark.Bot
}

// NewNotifier 根据传入的类型和凭证创建一个新的 Notifier 实例。
// 凭证通常是 Webhook URL，但可以根据平台需求传入其他 Token。
func NewNotifier(ntype NotifierType, credential, titlePrefix string) (*Notifier, error) {
	if credential == "" {
		return nil, errors.New("notifier: credential (Webhook URL) cannot be empty")
	}

	n := &Notifier{
		Type:        ntype,
		TitlePrefix: titlePrefix,
	}

	switch ntype {
	case NotifierTypeSlack:
		// 对于 Slack，直接存储 Webhook URL
		if !strings.HasPrefix(credential, "https://hooks.slack.com/") {
			return nil, errors.New("notifier: slack webhook url format seems incorrect")
		}
		n.SlackWebhookURL = credential
	case NotifierTypeLark:
		// 对于 Lark (飞书)，使用 go-lark 库初始化一个 NotificationBot
		// 飞书 Webhook 格式通常为: https://open.larksuite.com/open-apis/bot/v2/hook/...
		if !strings.Contains(credential, "/bot/v2/hook/") {
			return nil, errors.New("notifier: lark webhook url format seems incorrect")
		}
		// go-lark/lark 库提供了 NotificationBot 来封装 Webhook 发送
		n.LarkBot = lark.NewNotificationBot(credential)
	default:
		return nil, fmt.Errorf("notifier: unsupported notifier type: %s", ntype)
	}

	return n, nil
}

// SendAlert 根据 Notifier 的类型，将告警消息发送到目标平台。
//
// title: 消息标题 (会被加粗或用于Block)
// content: 消息主要内容
// mentionAll: 是否 @所有人 (在 Slack 中是 <!channel> 或 <!here>, 在 Lark 中是 @all)
func (n *Notifier) SendAlert(ctx context.Context, title, content string, mentionAll bool) error {
	if n.TitlePrefix != "" {
		title = fmt.Sprintf("[%s] %s", n.TitlePrefix, title)
	}

	switch n.Type {
	case NotifierTypeSlack:
		fullMessage := fmt.Sprintf("*%s*\n%s", title, content)
		return n.sendToSlack(ctx, fullMessage, mentionAll)
	case NotifierTypeLark:
		return n.sendToLark(ctx, title, content, mentionAll)
	default:
		return fmt.Errorf("notifier: unsupported notifier type: %s", n.Type)
	}
}

// sendToSlack 使用 slack-go/slack 库发送消息。
func (n *Notifier) sendToSlack(ctx context.Context, message string, mentionAll bool) error {
	// 在 Slack 中，@all 是 <!channel> 或 <!here>
	prefix := ""
	if mentionAll {
		// 使用 <!here> 只通知当前活跃的用户，更友好一些
		prefix = "<!here>\n"
	}

	// Slack 使用 Markdown/Block Kit 格式，我们构建一个简单的附件（Attachment）来突出标题。
	// 当然也可以使用 PostWebhook 函数发送简单的 WebhookMessage。这里采用 PostWebhook。

	msg := slack.WebhookMessage{
		Text: prefix + message, // 直接发送带有 Markdown 的文本
		// Block 也可以，但 Text 是最简单通用的 Webhook 方式
	}

	// slack.PostWebhook 是一个封装了 HTTP POST 请求的简便函数
	return slack.PostWebhook(n.SlackWebhookURL, &msg)
}

// sendToLark 使用 go-lark/lark 库发送消息。
func (n *Notifier) sendToLark(ctx context.Context, title, content string, mentionAll bool) error {
	// Lark (飞书) 推荐使用 RichText (Post) 或 Card 消息来发送结构化信息。
	// 这里我们使用最简单的 Post 富文本格式，它支持标题和 @all。

	// 创建一个富文本构建器
	postBuilder := lark.NewPostBuilder().
		Title(title).
		TextTag(content, 1, true).
		Render()

	// 增加 @所有人
	if mentionAll {
		// 该方法已经废弃
		// postBuilder.MentionAll()
	}

	// 构建整个消息体
	msg := lark.NewMsgBuffer(lark.MsgPost).
		Post(postBuilder).
		Build()

	// PostNotification 是 lark bot 发送消息的标准方法
	// 注意: lark 库返回的错误可能包含响应体，需要检查
	resp, err := n.LarkBot.PostNotification(ctx, msg)

	if err != nil {
		return fmt.Errorf("notifier: lark: bot post error: %w", err)
	}

	// 检查 Lark 响应的 code 是否成功 (code=0 为成功)
	if resp.Code != 0 {
		return fmt.Errorf("notifier: lark: API error code %d: %s", resp.Code, resp.Msg)
	}

	return nil
}
