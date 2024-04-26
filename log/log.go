package log

import (
	"encoding/json"
	"time"

	"code.gitee.cn/databridge/common"

	"go.opentelemetry.io/otel/attribute"
)

const (
	FatalLevel = "Fatal" // 致命错误，严重影响用户使用
	ErrorLevel = "Error" // 发生严重问题，部分功能受到影响
	WarnLevel  = "Warn"  // 可能存在潜在问题，但当前不影响用户体验和使用
	InfoLevel  = "Info"  // 普通信息输出
	DebugLevel = "Debug" // 调试信息
)

const (
	LogTypeServer = "server" // 后端
	LogTypeApp    = "app"    // 前端
	LogTypeSlow   = "slow"   // 系统缓慢
	LogTypeGc     = "gc"     // 垃圾回收信息
	LogTypeAudit  = "audit"  // 审计日志
)

type OTELLog struct {
	Type      string `json:"type"`
	Timestamp string `json:"timestamp"`
	Level     string `json:"level"`
	Component string `json:"component"`
	Tenant    string `json:"tenant"`
	RequestID string `json:"requestId"`
	Message   string `json:"message"`
}

func (o *OTELLog) SetTimestamp() {
	o.Timestamp = time.Now().Format("2006-01-02T15:04:05-07:00")
}

func (o *OTELLog) SetComponent() {
	o.Component = common.ModuleName
}

func (o *OTELLog) MakeupLogAttr() attribute.KeyValue {
	if len(o.Component) == 0 {
		o.SetComponent()
	}
	if len(o.Timestamp) == 0 {
		o.SetTimestamp()
	}
	ol, _ := json.Marshal(o)
	return attribute.String("log", string(ol))
}
