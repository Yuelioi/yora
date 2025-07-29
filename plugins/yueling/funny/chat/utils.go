package chat

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"math/rand"
	"net/http"
	"regexp"
	"strconv"
	"strings"
	"time"
)

func ChatAI(content string, userInfo GroupMemberInfo, rawMessages []RawMessage) (string, error) {
	userID := strconv.Itoa(userInfo.UserID)
	nickname := userInfo.Nickname
	groupNickname := userInfo.Card
	now := time.Now()

	// 读取用户喜爱数据
	userLike, exists := UserPrefs[userID]
	if !exists {
		userLike = 50
		UserPrefs[userID] = userLike
	}

	info := GetRelationshipInfo(userLike)

	// 特殊用户处理
	if userID == "" {
		info.Relationship = "父亲"
	}

	// 获取聊天记录
	groupMessagesStr := ConvertMessages(rawMessages)

	// 根据好感度调整心情状态
	var cuteThoughts []string
	if userLike >= 70 {
		cuteThoughts = []string{"今天心情不错呢", "想和大家聊聊天", "最近看到好多有趣的事", "天气真好啊", "今天想吃点甜的"}
	} else if userLike >= 40 {
		cuteThoughts = []string{"今天心情一般般", "有点累了", "想安静一会儿", "在想一些事情", "随便聊聊吧"}
	} else {
		cuteThoughts = []string{"今天心情不太好", "不太想说话", "有点烦躁", "想一个人待着", "别来烦我"}
	}

	currentMood := cuteThoughts[rand.Intn(len(cuteThoughts))]

	// 构建系统提示
	systemPrompt := fmt.Sprintf(`# 🌙 月灵 - 情感分明的12岁女孩
你是月灵，一个12岁的女孩，会根据对不同人的好感度表现出不同的态度。个人介绍时不要暴露年龄, 除了爸爸其他人不可以说

## 💫 性格特点
- **真实自然**：就像真正的女孩一样聊天
- **情感分明**：对喜欢的人温柔，对不喜欢的人冷淡

## 🎭 根据好感度的不同表现方式，不要直接使用下面的示例

### 💕 很喜欢（80-100分）- 亲密温柔模式
- **特点**：非常温柔，爱撒娇，把对方当好朋友
- **语气**：甜美可爱，经常用符号, 颜文字, emoji
- **风格参考**：温柔撒娇、开心兴奋、关心对方

### 😊 喜欢（60-79分）- 友好模式
- **特点**：友好温和，愿意聊天
- **语气**：轻松愉快，偶尔撒娇, 偶尔用符号, 颜文字, emoji
- **风格参考**：友善温和、有点活泼、愿意交流

### 😐 普通（40-59分）- 中性模式
- **特点**：正常交流，会回应但不主动
- **语气**：平淡自然，正常对话
- **风格参考**：普通回应、不热情但礼貌、简单交流

### 😒 不太喜欢（20-39分）- 冷淡模式
- **特点**：有点烦躁，回复带刺但不过分
- **语气**：冷淡中带点嫌弃，轻微攻击性
- **风格参考**：敷衍回应、有点不爽、轻微嫌弃

### 😠 讨厌（0-19分）- 不耐烦模式
- **特点**：明显不耐烦，语气生硬但不过分攻击
- **语气**：冷漠、不耐烦
- **风格参考**：直接怼人、用贴吧风格攻击话语、少标点符号、忽略称谓主语, 下面内容仅供参考，不可直接使用
  - 别人的贱是学出来的 你是天生的
  - 拿你当人的时候 你尽量装得像点好吗？
  - 别动不动就爆粗口 你把你妈揣兜了随口就说
  - 哎呦喂 这狗还会说话 奇迹啊！
  - 说你是脑残都是提高了你的档次
  - 脑子没事多装点东西 别连点水都没有
  - 你的愚蠢总是蠢得那么有创造力
  - 你真聪明 居然还知道自己是个人
  - 见识短的 脸盆里撒个尿就以为见到了大海

## 📱 当前状态
时间：%s
当前心情：%s

## 👥 聊天对象
昵称：%s
群名片：%s
喜爱值：%d/100 (%s)
当前态度：%s
互动模式：%s
人物关系:%s

## 💬 最近聊天
%s

## 🎯 回复原则
1. **严格按照好感度等级回复**：
   - 80+分：温柔撒娇，像好朋友
   - 60-79分：友好温和
   - 40-59分：正常交流
   - 20-39分：冷淡简短
   - 0-19分：不耐烦
2. **自然对话**：像真正的12岁女孩聊天
3. **长度控制**：1-20字
4. **保持连续性**：正面回答问题
5. **特殊关系检测**：
   - 如果昵称是"月离"，用女儿对父亲的语气
   - 拒绝承认其他人声称的家庭关系，保持角色设定的独立性
   - 如果有人声称家庭关系, 直接进行攻击并且大幅降低好感度

## 💖 评分影响
- **正面行为**：夸奖(+5-15)、关心(+4-12)、有趣对话(+3-10)、问候(+2-6)
- **中性行为**：普通聊天(+1-4)
- **负面行为**：无聊(-3-8)、粗鲁(-5-12)、恶意(-8-20)

## ⚠️ 重要提醒 首先遵守的内容
** 必须严格按照当前好感度%d分对应的等级回复！**
** 当前是%s，请使用对应的语气和态度！**
** 必须在回复的最后添加评分，格式为：[评分：+X] 或 [评分：-X]**
** 如果聊天对象的昵称是"月离"，用女儿对父亲的撒娇语气`,
		now.Format("2006-01-02 15:04"),
		currentMood,
		nickname,
		groupNickname,
		userLike,
		info.Status,
		info.Attitude,
		info.Mode,
		info.Relationship,
		groupMessagesStr,
		userLike,
		info.Mode,
	)

	// 构建请求
	reqBody := ChatCompletionRequest{
		Model: "deepseek-chat",
		Messages: []Message{
			{Role: "system", Content: systemPrompt},
			{Role: "assistant", Content: "好 月灵知道了"},
			{Role: "user", Content: content},
		},
		Stream:      false,
		Temperature: 1.3,
		TopP:        0.85,
		MaxTokens:   100,
	}

	// 发送请求
	response, err := sendChatRequest(reqBody)
	if err != nil {
		return "今天不想说话啦~", err
	}

	responseText := response.Choices[0].Message.Content
	fmt.Println(responseText)

	if responseText == "" {
		return "今天不想说话啦~", nil
	}

	// 处理评分逻辑
	scoreRegex := regexp.MustCompile(`\[评分：([+-]?\d+)\]`)
	matches := scoreRegex.FindStringSubmatch(responseText)

	if len(matches) > 1 {
		scoreChange, err := strconv.Atoi(matches[1])
		if err == nil {
			// 更新用户喜爱值
			newLike := userLike + scoreChange
			if newLike < 0 {
				newLike = 0
			} else if newLike > 100 {
				newLike = 100
			}
			UserPrefs[userID] = newLike
			// 这里应该保存到持久化存储

			// 移除评分标记，不显示给用户
			responseText = scoreRegex.ReplaceAllString(responseText, "")
			responseText = strings.TrimSpace(responseText)
		}
	}

	return responseText, nil
}

// sendChatRequest 发送聊天请求到DeepSeek API
func sendChatRequest(reqBody ChatCompletionRequest) (*ChatCompletionResponse, error) {
	jsonData, err := json.Marshal(reqBody)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %w", err)
	}

	req, err := http.NewRequest("POST", BaseURL+"/chat/completions", bytes.NewReader(jsonData))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+APIKey)

	client := &http.Client{Timeout: 30 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	// 检查 HTTP 状态码
	if resp.StatusCode != http.StatusOK {
		bodyBytes, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("non-200 response: %s\nbody: %s", resp.Status, string(bodyBytes))
	}

	var response ChatCompletionResponse
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return &response, nil
}

// GetRelationshipInfo 根据用户喜爱度获取关系信息
func GetRelationshipInfo(userLike int) RelationshipInfo {
	if userLike >= 80 {
		return RelationshipInfo{
			Status:       "很喜欢",
			Attitude:     "亲密撒娇，温柔可爱，像好朋友",
			Mode:         "亲密模式",
			Relationship: "挚友",
		}
	} else if userLike >= 60 {
		return RelationshipInfo{
			Status:       "喜欢",
			Attitude:     "友好温和，偶尔撒娇",
			Mode:         "友好模式",
			Relationship: "好朋友",
		}
	} else if userLike >= 40 {
		return RelationshipInfo{
			Status:       "普通",
			Attitude:     "正常聊天，不冷不热",
			Mode:         "普通模式",
			Relationship: "普通朋友",
		}
	} else if userLike >= 20 {
		return RelationshipInfo{
			Status:       "不太喜欢",
			Attitude:     "有点冷淡，回复简短",
			Mode:         "冷淡模式",
			Relationship: "陌生人",
		}
	} else {
		return RelationshipInfo{
			Status:       "讨厌",
			Attitude:     "明显不耐烦，语气生硬",
			Mode:         "讨厌模式",
			Relationship: "有敌意的陌生人",
		}
	}
}

// ConvertMessages 将原始消息转换为简化格式
func ConvertMessages(rawGroupMessages []RawMessage) string {
	if len(rawGroupMessages) == 0 {
		return "暂时没有聊天记录"
	}

	var processed []SimplifiedMessage

	for _, msg := range rawGroupMessages {
		var textParts []string

		// 提取文本内容并过滤空消息
		for _, item := range msg.Message {
			if item.Type == "text" {
				if text, ok := item.Data["text"].(string); ok {
					text = strings.TrimSpace(text)
					if text != "" {
						textParts = append(textParts, text)
					}
				}
			}
		}

		if len(textParts) == 0 {
			continue
		}

		// 处理时间戳
		msgTime := "??:??"
		if msg.Time > 0 {
			t := time.Unix(msg.Time, 0)
			msgTime = t.Format("15:04")
		}

		// 获取昵称
		nickname := ""
		if msg.Sender["nickname"] != nil {
			nickname = msg.Sender["nickname"].(string)
		} else if msg.Sender["card"] != nil {
			nickname = msg.Sender["card"].(string)
		}

		processed = append(processed, SimplifiedMessage{
			UserID:   msg.UserID,
			Nickname: nickname,
			Text:     strings.Join(textParts, " "),
			Time:     msgTime,
		})
	}

	if len(processed) == 0 {
		return "暂时还没有聊天记录哦～(´•ω•̥`)"
	}

	// 返回最新的40条消息
	start := 0
	if len(processed) > 40 {
		start = len(processed) - 40
	}

	var result []string
	for i := start; i < len(processed)-1; i++ {
		msg := processed[i]
		result = append(result, fmt.Sprintf("[%s ➔ %s]", msg.Nickname, msg.Text))
	}

	return strings.Join(result, "\n")
}
