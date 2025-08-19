// Package formatter 提供了一个基于 text/template 的、带缓存的高性能字符串格式化工具。
package formatter

import (
	"bytes"
	"fmt"
	"sync"
	"text/template"
)

// KV 表示一个键值对 (Key-Value pair)。
type KV struct {
	Key   string
	Value interface{}
}

// D 是一个用于构建格式化参数的便捷类型，其灵感来源于 bson.D。
// 它是一个 KV 类型的切片。
//
// 使用示例:
//
//	formatter.D{
//	    {"name", "Alice"},
//	    {"age", 30},
//	}
type D []KV

// toMap 将 formatter.D 转换为 map[string]interface{}。
func (d D) toMap() map[string]interface{} {
	// 预分配 map 容量以提高性能
	m := make(map[string]interface{}, len(d))
	for _, pair := range d {
		m[pair.Key] = pair.Value
	}
	return m
}

var (
	templateCache = make(map[string]*template.Template)
	cacheMutex    = &sync.RWMutex{}
)

// Format 使用命名参数格式化字符串。
// 它内部使用带缓存的 text/template 引擎，在重复使用相同模板时性能很高。
//
// templateStr: 带有占位符的模板字符串，例如 "Hello, {{.Name}}!"
// params: 参数，可以是 map[string]interface{} 类型，也可以是 formatter.D 类型。
//
// 返回格式化后的字符串和一个 error。
func Format(templateStr string, params interface{}) (string, error) {
	// 步骤 1: 将传入的 params 转换为 map[string]interface{}
	paramMap, err := convertParamsToMap(params)
	if err != nil {
		return "", err
	}

	// 步骤 2: 尝试从缓存中获取模板 (使用读锁)
	cacheMutex.RLock()
	tpl, found := templateCache[templateStr]
	cacheMutex.RUnlock()

	if !found {
		// 步骤 3: 如果缓存中没有，则解析并存入缓存 (使用写锁)
		cacheMutex.Lock()
		if tpl, found = templateCache[templateStr]; !found {
			tpl, err = template.New("formatter").Parse(templateStr)
			if err != nil {
				cacheMutex.Unlock()
				return "", err
			}
			templateCache[templateStr] = tpl
		}
		cacheMutex.Unlock()
	}

	// 步骤 4: 使用获取到的模板执行格式化
	var buf bytes.Buffer
	err = tpl.Execute(&buf, paramMap)
	if err != nil {
		return "", err
	}

	return buf.String(), nil
}

// MustFormat 是 Format 的一个变体，它在发生错误时会 panic。
// 这在模板和参数都由程序控制，确信不会出错的场景下可以简化代码。
func MustFormat(templateStr string, params interface{}) string {
	result, err := Format(templateStr, params)
	if err != nil {
		panic(err)
	}
	return result
}

// convertParamsToMap 是一个内部辅助函数，用于处理不同类型的参数。
func convertParamsToMap(params interface{}) (map[string]interface{}, error) {
	switch p := params.(type) {
	case map[string]interface{}:
		return p, nil // 如果已经是 map，直接返回
	case D:
		return p.toMap(), nil // 如果是 D 类型，转换为 map
	default:
		// 返回一个明确的错误，告知用户支持的类型
		return nil, fmt.Errorf("unsupported params type: %T. must be map[string]interface{} or formatter.D", p)
	}
}
