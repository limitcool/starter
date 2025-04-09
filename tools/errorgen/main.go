package main

import (
	"bufio"
	"bytes"
	"fmt"
	"os"
	"regexp"
	"strconv"
	"strings"
	"text/template"
)

// 错误定义结构
type ErrorDef struct {
	Code        int
	Name        string
	Message     string
	HTTPStatus  int
	Description string // 可选的描述信息
}

// 错误组定义
type ErrorGroup struct {
	Name        string
	Description string
	BaseCode    int
	Range       string
	Errors      []ErrorDef
}

// 解析Markdown表格行
func parseTableRow(line string) (*ErrorDef, error) {
	parts := strings.Split(line, "|")
	if len(parts) < 5 {
		return nil, fmt.Errorf("无效的表格行: %s", line)
	}

	// 清理每个字段的空白
	for i := range parts {
		parts[i] = strings.TrimSpace(parts[i])
	}

	// 解析错误码
	code, err := strconv.Atoi(parts[1])
	if err != nil {
		return nil, fmt.Errorf("错误码解析失败: %s", parts[1])
	}

	// 解析HTTP状态码
	httpStatus, err := strconv.Atoi(parts[4])
	if err != nil {
		return nil, fmt.Errorf("HTTP状态码解析失败: %s", parts[4])
	}

	return &ErrorDef{
		Code:       code,
		Name:       parts[2],
		Message:    parts[3],
		HTTPStatus: httpStatus,
	}, nil
}

// 解析Markdown文件中的错误定义
func parseMarkdownFile(filePath string) ([]ErrorGroup, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var errorGroups []ErrorGroup
	var currentGroup *ErrorGroup

	scanner := bufio.NewScanner(file)
	inTable := false
	headerPattern := regexp.MustCompile(`^## (.+) \((\d+)-(\d+)\)$`)

	for scanner.Scan() {
		line := scanner.Text()

		// 检查是否是组标题行
		if matches := headerPattern.FindStringSubmatch(line); len(matches) > 0 {
			groupName := matches[1]
			baseCodeStr := matches[2]
			rangeStr := fmt.Sprintf("%s-%s", matches[2], matches[3])

			baseCode, _ := strconv.Atoi(baseCodeStr)

			// 保存当前组并创建新组
			if currentGroup != nil && len(currentGroup.Errors) > 0 {
				errorGroups = append(errorGroups, *currentGroup)
			}

			currentGroup = &ErrorGroup{
				Name:     groupName,
				BaseCode: baseCode,
				Range:    rangeStr,
				Errors:   []ErrorDef{},
			}
			inTable = false
			continue
		}

		// 基础错误码特殊处理 - 跳过，因为我们在模板中直接定义SuccessCode
		if line == "## 基础错误码" {
			inTable = false
			// 跳过基础错误码部分，直到下一个标题
			for scanner.Scan() {
				line = scanner.Text()
				if strings.HasPrefix(line, "## ") {
					// 找到下一个标题，处理这个新标题
					if matches := headerPattern.FindStringSubmatch(line); len(matches) > 0 {
						groupName := matches[1]
						baseCodeStr := matches[2]
						rangeStr := fmt.Sprintf("%s-%s", matches[2], matches[3])

						baseCode, _ := strconv.Atoi(baseCodeStr)

						// 保存当前组并创建新组
						if currentGroup != nil && len(currentGroup.Errors) > 0 {
							errorGroups = append(errorGroups, *currentGroup)
						}

						currentGroup = &ErrorGroup{
							Name:     groupName,
							BaseCode: baseCode,
							Range:    rangeStr,
							Errors:   []ErrorDef{},
						}
					}
					break
				}
			}
			continue
		}

		// 检查是否是表格分隔行
		if strings.Contains(line, "|------") {
			inTable = true
			continue
		}

		// 跳过表头行
		if strings.Contains(line, "| 错误码 | 名称 | 错误消息 | HTTP状态码 |") {
			continue
		}

		// 解析表格数据行
		if inTable && strings.HasPrefix(line, "|") && currentGroup != nil {
			errorDef, err := parseTableRow(line)
			if err != nil {
				fmt.Printf("警告: %v\n", err)
				continue
			}
			currentGroup.Errors = append(currentGroup.Errors, *errorDef)
		}
	}

	// 不要忘记添加最后一个组
	if currentGroup != nil && len(currentGroup.Errors) > 0 {
		errorGroups = append(errorGroups, *currentGroup)
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return errorGroups, nil
}

// 生成Go代码的模板
const codeTemplate = `// 本文件由工具自动生成，请勿手动修改
// 生成命令: go run tools/errorgen/main.go

package errorx

import "net/http"

// 错误码基础值定义
const (
	// 错误码基础值 (按模块划分)
	CommonErrorBase  = 10000 // 通用错误码 (10000-19999)
	{{- range .ErrorGroups }}
	{{- if and (ne .BaseCode 0) (ne .BaseCode 10000) }}
	{{ .GroupConstName }}ErrorBase = {{ .BaseCode }} // {{ .Name }} ({{ .Range }})
	{{- end }}
	{{- end }}
)

// 错误码定义
const (
	// 基础错误码
	SuccessCode = 0 // 成功

	{{- range .ErrorGroups }}
	// {{ .Name }}
	{{- $baseConst := .GroupConstName }}
	{{- range $index, $error := .Errors }}
	{{ $error.Name }} = {{ $baseConst }}ErrorBase + {{ $index }} // {{ $error.Message }}
	{{- end }}
	{{- end }}
)

// 预定义错误实例
var (
	// 基础错误
	ErrSuccess = NewAppError(SuccessCode, "成功", http.StatusOK)

	{{- range .ErrorGroups }}
	// {{ .Name }}实例
	{{- range .Errors }}
	Err{{ .InstanceName }} = NewAppError({{ .Name }}, "{{ .Message }}", http.Status{{ .HTTPStatusName }})
	{{- end }}
	{{- end }}
)

// GetError 根据错误码获取预定义错误
func GetError(code int) *AppError {
	switch code {
	case SuccessCode:
		return ErrSuccess
	{{- range .ErrorGroups }}
	{{- range .Errors }}
	case {{ .Name }}:
		return Err{{ .InstanceName }}
	{{- end }}
	{{- end }}
	default:
		return ErrUnknown // 默认返回未知错误
	}
}`

// 模板数据结构
type TemplateData struct {
	ErrorGroups []ErrorGroupTemplate
}

// 错误组模板数据
type ErrorGroupTemplate struct {
	Name           string
	GroupConstName string
	BaseCode       int
	Range          string
	Errors         []ErrorDefTemplate
}

// 错误定义模板数据
type ErrorDefTemplate struct {
	Code           int
	Name           string
	Message        string
	HTTPStatus     int
	HTTPStatusName string
	InstanceName   string
}

// HTTP状态码转换为Go常量名称
func httpStatusToName(code int) string {
	statusMap := map[int]string{
		200: "OK",
		201: "Created",
		204: "NoContent",
		400: "BadRequest",
		401: "Unauthorized",
		403: "Forbidden",
		404: "NotFound",
		409: "Conflict",
		500: "InternalServerError",
	}

	if name, ok := statusMap[code]; ok {
		return name
	}
	return fmt.Sprintf("%d", code)
}

// 将组名称转为常量名称
func groupNameToConstName(name string) string {
	// 处理特殊情况
	nameMap := map[string]string{
		"通用错误码":   "Common",
		"数据库错误码":  "Database",
		"用户相关错误码": "User",
		"权限相关错误码": "Auth",
		"缓存相关错误码": "Cache",
		"文件相关错误码": "File",
	}

	if constName, ok := nameMap[name]; ok {
		return constName
	}

	// 移除非字母字符
	reg := regexp.MustCompile(`[^a-zA-Z]`)
	name = reg.ReplaceAllString(name, "")

	// 首字母大写，其余小写
	if len(name) > 0 {
		name = strings.ToUpper(name[:1]) + strings.ToLower(name[1:])
	}

	return name
}

// 生成错误实例名称
func generateInstanceName(errorName string) string {
	// 移除Code后缀
	instanceName := strings.TrimSuffix(errorName, "Code")

	// 已有良好命名的情况，直接返回
	if !strings.HasPrefix(instanceName, "Error") {
		return instanceName
	}

	// 处理Error前缀的特殊情况
	if len(instanceName) > 5 { // "Error" 长度为5
		return strings.TrimPrefix(instanceName, "Error")
	}

	return instanceName
}

func main() {
	if len(os.Args) < 3 {
		fmt.Println("用法: go run main.go <markdown文件路径> <输出Go文件路径>")
		os.Exit(1)
	}

	markdownFile := os.Args[1]
	outputFile := os.Args[2]

	// 解析Markdown文件
	errorGroups, err := parseMarkdownFile(markdownFile)
	if err != nil {
		fmt.Printf("解析Markdown文件失败: %v\n", err)
		os.Exit(1)
	}

	// 准备模板数据
	templateData := TemplateData{
		ErrorGroups: make([]ErrorGroupTemplate, 0, len(errorGroups)),
	}

	for _, group := range errorGroups {
		templateGroup := ErrorGroupTemplate{
			Name:           group.Name,
			GroupConstName: groupNameToConstName(group.Name),
			BaseCode:       group.BaseCode,
			Range:          group.Range,
			Errors:         make([]ErrorDefTemplate, 0, len(group.Errors)),
		}

		for _, errDef := range group.Errors {
			templateGroup.Errors = append(templateGroup.Errors, ErrorDefTemplate{
				Code:           errDef.Code,
				Name:           errDef.Name,
				Message:        errDef.Message,
				HTTPStatus:     errDef.HTTPStatus,
				HTTPStatusName: httpStatusToName(errDef.HTTPStatus),
				InstanceName:   generateInstanceName(errDef.Name),
			})
		}

		templateData.ErrorGroups = append(templateData.ErrorGroups, templateGroup)
	}

	// 执行模板
	tmpl, err := template.New("codeTemplate").Parse(codeTemplate)
	if err != nil {
		fmt.Printf("模板解析失败: %v\n", err)
		os.Exit(1)
	}

	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, templateData); err != nil {
		fmt.Printf("模板执行失败: %v\n", err)
		os.Exit(1)
	}

	// 写入输出文件
	if err := os.WriteFile(outputFile, buf.Bytes(), 0644); err != nil {
		fmt.Printf("写入输出文件失败: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("成功生成错误码文件: %s\n", outputFile)
}
