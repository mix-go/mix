package xutil

import (
	"context"
	"errors"
	"fmt"
	"github.com/patrickmn/go-cache"
	"runtime"
	"runtime/debug"
	"time"
)

var alerterSingleton *Alerter

type Alerter struct {
	RateLimitInterval time.Duration
	Notifier          *Notifier
	Logger            Logger
	cache             *cache.Cache
	ch                chan *AlertMsg
}

type AlertMsg struct {
	Context    context.Context
	UUID       string
	Title      string
	Content    string
	MentionAll bool
}

func NewAlert(ntype NotifierType, credential, titlePrefix string, rateLimitInterval time.Duration, mqSize int, logger Logger) (*Alerter, error) {
	n, err := NewNotifier(ntype, credential, titlePrefix)
	if err != nil {
		return nil, err
	}
	a := &Alerter{
		RateLimitInterval: rateLimitInterval,
		Notifier:          n,
		Logger:            logger,
		cache:             cache.New(rateLimitInterval, rateLimitInterval*2),
		ch:                make(chan *AlertMsg, mqSize),
	}

	go func() {
		for {
			msg := <-a.ch

			// 降频
			if msg.UUID != "" {
				_, ok := a.cache.Get(msg.UUID)
				if ok {
					continue
				}
			}

			if err := a.Notifier.SendAlert(msg.Context, msg.Title, msg.Content, msg.MentionAll); err != nil {
				a.Logger.Errorw("alert: send alert msg fail", "uuid", msg.UUID, TraceIDKey, msg.Context.Value(TraceIDKey))
			}

			if msg.UUID != "" {
				a.cache.SetDefault(msg.UUID, true)
			}
		}
	}()

	return a, nil
}

// StartAlerter 启动报警器
func StartAlerter(ntype NotifierType, credential, titlePrefix string, rateLimitInterval time.Duration, mqSize int, logger Logger) error {
	alert, err := NewAlert(ntype, credential, titlePrefix, rateLimitInterval, mqSize, logger)
	if err != nil {
		return err
	}
	alerterSingleton = alert
	return nil
}

// PushAlert 队列限速发送
func PushAlert(ctx context.Context, title, content, uuid string, mentionAll bool) error {
	if alerterSingleton == nil {
		return errors.New("alert: instance not initialized")
	}

	select {
	case alerterSingleton.ch <- &AlertMsg{
		Context:    ctx,
		UUID:       uuid,
		Title:      title,
		Content:    content,
		MentionAll: mentionAll,
	}:
	default:
		// 队列满就丢弃
	}
	return nil
}

// SendAlert 直接发送
func SendAlert(ctx context.Context, title, content string, mentionAll bool) error {
	if alerterSingleton == nil {
		return errors.New("alert: instance not initialized")
	}

	if err := alerterSingleton.Notifier.SendAlert(ctx, title, content, mentionAll); err != nil {
		return fmt.Errorf("alert: send alert msg fail: %s", err)
	}

	return nil
}

// ErrorAlert 错误警报
func ErrorAlert(err error, needStack bool) {
	if err == nil {
		return
	}

	_, file, line, _ := runtime.Caller(1)
	fileLine := fmt.Sprintf("%s:%d", file, line)

	title := fmt.Sprintf("ERROR: %s", err)
	content := fileLine
	if needStack {
		content = fmt.Sprintf("%s\n%s", content, string(debug.Stack()))
	}

	uuid := fileLine

	_ = PushAlert(context.Background(), title, content, uuid, true)
}
