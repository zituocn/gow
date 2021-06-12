/*
app_config.go

*/

package gow

import "github.com/zituocn/gow/lib/config"

// AppConfig unified configuration entry
type AppConfig struct {
	AppName       string `json:"app_name,omitempty" yaml:"app_name"`             // app name
	RunMode       string `json:"run_mode,omitempty" yaml:"run_mode"`             // app run mode
	HTTPAddr      string `json:"http_addr,omitempty" yaml:"http_addr"`           // http address
	AutoRender    bool   `json:"auto_render,omitempty" yaml:"auto_render"`       // if true load html template
	Views         string `json:"views,omitempty" yaml:"views"`                   // html template dir
	TemplateLeft  string `json:"template_left,omitempty" yaml:"template_left"`   // html template tag symbol
	TemplateRight string `json:"template_right,omitempty" yaml:"template_right"` // html template tag symbol
	SessionOn     bool   `json:"session_on,omitempty" yaml:"session_on"`         // if true open session
	GzipOn        bool   `json:"gzip_on,omitempty" yaml:"gzip_on"`               // if true:load gzip middleware
	IgnoreCase    bool   `json:"ignore_case" yaml:"ignore_case"`                 // if true ignore case on route
}

// GetAppConfig  return engine config
func GetAppConfig() *AppConfig {
	return &AppConfig{
		AppName:       config.DefaultString("app_name", "gow"),
		RunMode:       config.DefaultString("run_mode", "dev"),
		HTTPAddr:      config.DefaultString("http_addr", ":8080"),
		AutoRender:    config.DefaultBool("auto_render", false),
		Views:         config.DefaultString("views", "views"),
		TemplateLeft:  config.DefaultString("template_left", "{{"),
		TemplateRight: config.DefaultString("template_right", "}}"),
		SessionOn:     config.DefaultBool("session_on", false),
		GzipOn:        config.DefaultBool("gzip_on", false),
		IgnoreCase:    config.DefaultBool("ignore_case", true),
	}
}
