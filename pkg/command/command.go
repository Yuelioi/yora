package main

import (
	"fmt"
	"strconv"
	"strings"
)

// Arg 位置参数
type Arg struct {
	Name     string // 参数名
	Type     string // 类型: "string", "int", "bool", "float"
	Optional bool   // 是否可选
	Help     string // 帮助信息
}

// Opt 选项参数
type Opt struct {
	Short    string      // 短选项 如: "f"
	Long     string      // 长选项 如: "force"
	Type     string      // 类型
	Default  interface{} // 默认值
	Required bool        // 是否必需
	Help     string      // 帮助信息
}

// Command 命令定义
type Command struct {
	Name        string              // 命令名
	Help        string              // 帮助信息
	Aliases     []string            // 别名
	Args        []Arg               // 位置参数
	Options     []Opt               // 选项参数
	SubCommands map[string]*Command // 子命令
	Handler     func(*ParseResult)  // 处理函数
}

// ParseResult 解析结果
type ParseResult struct {
	Command   *Command
	Args      map[string]interface{}
	Options   map[string]interface{}
	RawArgs   []string
	Remaining []string
}

// Parser 命令解析器
type Parser struct {
	Commands map[string]*Command
	Prefix   string
}

// NewParser 创建解析器
func NewParser(prefix string) *Parser {
	return &Parser{
		Commands: make(map[string]*Command),
		Prefix:   prefix,
	}
}

// Register 注册命令
func (p *Parser) Register(cmd *Command) {
	p.Commands[cmd.Name] = cmd
}

// Parse 解析命令
func (p *Parser) Parse(input string) (*ParseResult, error) {
	input = strings.TrimSpace(input)

	// 检查前缀
	if p.Prefix != "" && !strings.HasPrefix(input, p.Prefix) {
		return nil, fmt.Errorf("command must start with prefix: %s", p.Prefix)
	}

	if p.Prefix != "" {
		input = strings.TrimPrefix(input, p.Prefix)
		input = strings.TrimSpace(input)
	}

	if input == "" {
		return nil, fmt.Errorf("empty command")
	}

	tokens := tokenize(input)
	if len(tokens) == 0 {
		return nil, fmt.Errorf("no command provided")
	}

	// 查找主命令
	cmdName := tokens[0]
	cmd, exists := p.Commands[cmdName]
	if !exists {
		return nil, fmt.Errorf("unknown command: %s", cmdName)
	}

	tokens = tokens[1:]

	// 检查子命令
	if len(tokens) > 0 && cmd.SubCommands != nil {
		if subCmd, exists := cmd.SubCommands[tokens[0]]; exists {
			cmd = subCmd
			tokens = tokens[1:]
		}
	}

	return p.parseCommand(cmd, tokens)
}

// parseCommand 解析具体命令
func (p *Parser) parseCommand(cmd *Command, tokens []string) (*ParseResult, error) {
	result := &ParseResult{
		Command: cmd,
		Args:    make(map[string]interface{}),
		Options: make(map[string]interface{}),
		RawArgs: tokens,
	}

	// 初始化选项默认值
	for _, opt := range cmd.Options {
		if opt.Default != nil {
			key := opt.Long
			if key == "" {
				key = opt.Short
			}
			result.Options[key] = opt.Default
		}
	}

	i := 0
	argIndex := 0

	for i < len(tokens) {
		token := tokens[i]

		if strings.HasPrefix(token, "--") {
			// 长选项
			optName := strings.TrimPrefix(token, "--")
			opt := findOptionByLong(cmd, optName)
			if opt == nil {
				return nil, fmt.Errorf("unknown option: --%s", optName)
			}

			if opt.Type == "bool" {
				result.Options[opt.Long] = true
				i++
			} else {
				if i+1 >= len(tokens) {
					return nil, fmt.Errorf("option --%s requires a value", optName)
				}
				value, err := convertValue(tokens[i+1], opt.Type)
				if err != nil {
					return nil, fmt.Errorf("invalid value for option --%s: %s", optName, err)
				}
				result.Options[opt.Long] = value
				i += 2
			}
		} else if strings.HasPrefix(token, "-") && len(token) > 1 {
			// 短选项
			optName := strings.TrimPrefix(token, "-")
			opt := findOptionByShort(cmd, optName)
			if opt == nil {
				return nil, fmt.Errorf("unknown option: -%s", optName)
			}

			key := opt.Long
			if key == "" {
				key = opt.Short
			}

			if opt.Type == "bool" {
				result.Options[key] = true
				i++
			} else {
				if i+1 >= len(tokens) {
					return nil, fmt.Errorf("option -%s requires a value", optName)
				}
				value, err := convertValue(tokens[i+1], opt.Type)
				if err != nil {
					return nil, fmt.Errorf("invalid value for option -%s: %s", optName, err)
				}
				result.Options[key] = value
				i += 2
			}
		} else {
			// 位置参数
			if argIndex < len(cmd.Args) {
				arg := cmd.Args[argIndex]
				value, err := convertValue(token, arg.Type)
				if err != nil {
					return nil, fmt.Errorf("invalid value for argument %s: %s", arg.Name, err)
				}
				result.Args[arg.Name] = value
				argIndex++
			} else {
				result.Remaining = append(result.Remaining, token)
			}
			i++
		}
	}

	// 检查必需的参数
	for i, arg := range cmd.Args {
		if i >= argIndex && !arg.Optional {
			return nil, fmt.Errorf("missing required argument: %s", arg.Name)
		}
	}

	// 检查必需的选项
	for _, opt := range cmd.Options {
		if opt.Required {
			key := opt.Long
			if key == "" {
				key = opt.Short
			}
			if _, exists := result.Options[key]; !exists {
				return nil, fmt.Errorf("missing required option: %s", key)
			}
		}
	}

	return result, nil
}

// Execute 执行命令
func (result *ParseResult) Execute() {
	if result.Command.Handler != nil {
		result.Command.Handler(result)
	}
}

// 辅助函数
func findOptionByLong(cmd *Command, name string) *Opt {
	for _, opt := range cmd.Options {
		if opt.Long == name {
			return &opt
		}
	}
	return nil
}

func findOptionByShort(cmd *Command, name string) *Opt {
	for _, opt := range cmd.Options {
		if opt.Short == name {
			return &opt
		}
	}
	return nil
}

type Bool bool

func convertValue(value string, targetType string) (interface{}, error) {
	switch targetType {
	case "string":
		return value, nil
	case "int":
		return strconv.Atoi(value)
	case "float":
		return strconv.ParseFloat(value, 64)
	case "bool":
		return strconv.ParseBool(value)
	case "[]string":
		return strings.Split(value, ","), nil
	case "[]int":
		values := strings.Split(value, ",")
		result := make([]int, len(values))
		for i, v := range values {
			n, err := strconv.Atoi(v)
			if err != nil {
				return nil, err
			}
			result[i] = n
		}
		return result, nil
	case "url":
		if !strings.HasPrefix(value, "http://") && !strings.HasPrefix(value, "https://") {
			return "", fmt.Errorf("invalid url: %s", value)
		}
		return value, nil
	case "image":
		if !strings.HasPrefix(value, "http://") && !strings.HasPrefix(value, "https://") {
			return "", fmt.Errorf("invalid image url: %s", value)
		}
		return value, nil

	default:
		return value, nil
	}
}

func tokenize(input string) []string {
	var tokens []string
	var current strings.Builder
	inQuotes := false
	quoteChar := byte(0)

	for i := 0; i < len(input); i++ {
		char := input[i]

		if inQuotes {
			if char == quoteChar {
				inQuotes = false
				quoteChar = 0
			} else {
				current.WriteByte(char)
			}
		} else {
			if char == '"' || char == '\'' {
				inQuotes = true
				quoteChar = char
			} else if char == ' ' || char == '\t' {
				if current.Len() > 0 {
					tokens = append(tokens, current.String())
					current.Reset()
				}
			} else {
				current.WriteByte(char)
			}
		}
	}

	if current.Len() > 0 {
		tokens = append(tokens, current.String())
	}

	return tokens
}

// 示例用法
func main() {
	// 创建解析器
	parser := NewParser("/")

	// 直接定义命令结构
	userCmd := &Command{
		Name: "user",
		Help: "用户管理命令",
		SubCommands: map[string]*Command{
			"add": {
				Name: "add",
				Help: "添加用户",
				Args: []Arg{
					{Name: "username", Type: "string", Optional: false, Help: "用户名"},
					{Name: "email", Type: "string", Optional: true, Help: "邮箱地址"},
				},
				Options: []Opt{
					{Short: "a", Long: "admin", Type: "bool", Default: false, Help: "是否为管理员"},
					{Short: "g", Long: "group", Type: "string", Default: "default", Help: "用户组"},
				},
				Handler: func(result *ParseResult) {
					username := result.Args["username"].(string)
					email, hasEmail := result.Args["email"]
					isAdmin := result.Options["admin"].(bool)
					group := result.Options["group"].(string)

					fmt.Printf("添加用户: %s\n", username)
					if hasEmail {
						fmt.Printf("邮箱: %s\n", email.(string))
					}
					fmt.Printf("管理员: %t\n", isAdmin)
					fmt.Printf("用户组: %s\n", group)
				},
			},
			"list": {
				Name: "list",
				Help: "列出用户",
				Options: []Opt{
					{Short: "l", Long: "limit", Type: "int", Default: 10, Help: "限制数量"},
					{Short: "f", Long: "filter", Type: "string", Help: "过滤条件"},
				},
				Handler: func(result *ParseResult) {
					limit := result.Options["limit"].(int)
					filter, hasFilter := result.Options["filter"]

					fmt.Printf("列出用户，限制: %d\n", limit)
					if hasFilter {
						fmt.Printf("过滤条件: %s\n", filter.(string))
					}
				},
			},
		},
	}

	// 系统命令
	sysCmd := &Command{
		Name: "sys",
		Help: "系统管理命令",
		Args: []Arg{
			{Name: "action", Type: "string", Optional: false, Help: "操作类型"},
		},
		Options: []Opt{
			{Short: "f", Long: "force", Type: "bool", Default: false, Help: "强制执行"},
			{Short: "t", Long: "timeout", Type: "int", Default: 30, Help: "超时时间(秒)"},
		},
		Handler: func(result *ParseResult) {
			action := result.Args["action"].(string)
			force := result.Options["force"].(bool)
			timeout := result.Options["timeout"].(int)

			fmt.Printf("系统操作: %s\n", action)
			fmt.Printf("强制执行: %t\n", force)
			fmt.Printf("超时时间: %d秒\n", timeout)
		},
	}

	// 注册命令
	parser.Register(userCmd)
	parser.Register(sysCmd)

	// 测试命令解析
	testCommands := []string{
		"/user add john john@example.com --admin --group staff",
		"/user list --limit 20 --filter active",
		"/sys restart --force --timeout 60",
		"/user add alice 123",
	}

	for _, cmd := range testCommands {
		fmt.Printf("\n解析命令: %s\n", cmd)
		fmt.Println(strings.Repeat("-", 50))

		result, err := parser.Parse(cmd)
		if err != nil {
			fmt.Printf("错误: %s\n", err)
			continue
		}

		result.Execute()
	}
}
