package gtp

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strings"
)

const BASEURL = "https://api.openai.com/v1/"

type ChatMsg struct {
	Role            string `json:"role"`
	ID              string `json:"id"`
	ParentMessageId string `json:"parentMessageId"`

	ConversationId string              `json:"conversationId"`
	Text           string              `json:"text"`
	Detail         ChatGPTResponseBody `json:"detail"`
}

// ChatGPTResponseBody 请求体
type ChatGPTResponseBody struct {
	ID      string                   `json:"id"`
	Object  string                   `json:"object"`
	Created int                      `json:"created"`
	Model   string                   `json:"model"`
	Choices []map[string]interface{} `json:"choices"`
	Usage   map[string]interface{}   `json:"usage"`
}

type ChoiceItem struct {
}

// ChatGPTRequestBody 响应体
type ChatGPTRequestBody struct {
	Model            string  `json:"model"`
	Prompt           string  `json:"prompt"`
	MaxTokens        int     `json:"max_tokens"`
	Temperature      float32 `json:"temperature"`
	TopP             int     `json:"top_p"`
	FrequencyPenalty int     `json:"frequency_penalty"`
	PresencePenalty  int     `json:"presence_penalty"`
}

// Completions gtp文本模型回复
// curl https://api.openai.com/v1/completions
// -H "Content-Type: application/json"
// -H "Authorization: Bearer your chatGPT key"
// -d '{"model": "text-davinci-003", "prompt": "give me good song", "temperature": 0, "max_tokens": 7}'
func Completions(version int, conv_id string, who string, msg string) (string, error) {
	//requestBody := ChatGPTRequestBody{
	//	Model:            "text-davinci-003",
	//	Prompt:           msg,
	//	MaxTokens:        2048,
	//	Temperature:      0.7,
	//	TopP:             1,
	//	FrequencyPenalty: 0,
	//	PresencePenalty:  0,
	//}
	var P = make(map[string]string)
	P["qureyT"] = msg

	requestData, err := json.Marshal(P)

	if err != nil {
		return "", err
	}
	log.Printf("request gtp json string : %v", string(requestData))
	//req, err := http.NewRequest("POST", BASEURL+"completions", bytes.NewBuffer(requestData))
	//var param = url.Values{}
	//param.Add("qureyT", msg)

	//var qureyT = base64.StdEncoding.EncodeToString([]byte(msg))
	//var qureyT = msg
	//fmt.Println("发送查询参数", qureyT)
	var urlP = url.Values{}
	urlP.Set("who", who)
	urlP.Set("message", msg)
	urlP.Set("conv_id", conv_id)

	urlReq := "http://192.168.228.129:8001/ask/v1?" + urlP.Encode()
	if version == 3 {
		//urlReq = "http://192.168.228.129:8001/ask/v3?" + urlP.Encode()
		urlReq = "http://127.0.0.1:8080/ask/v3?" + urlP.Encode()
	}

	req, err := http.NewRequest("GET", urlReq, nil)
	if err != nil {
		return "", err
	}

	//apiKey := config.LoadConfig().ApiKey
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("charset", "UTF-8")
	//req.Header.Set("Authorization", "Bearer "+apiKey)
	client := &http.Client{}
	response, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer response.Body.Close()

	body, err := ioutil.ReadAll(response.Body)

	fmt.Println("收到回复内容是", string(body))
	fmt.Println("收到回复异常err", err)
	if err != nil {
		return "", err
	}

	//gptResponseBody := &ChatGPTResponseBody{}
	//log.Println(string(body))
	//err = json.Unmarshal(body, gptResponseBody)
	//if err != nil {
	//	return "", err
	//}
	//var reply string
	//if len(gptResponseBody.Choices) > 0 {
	//	for _, v := range gptResponseBody.Choices {
	//		reply = v["text"].(string)
	//		break
	//	}
	//}
	reply := string(body)
	//fmt.Println("接收内容中包含的换行符个数：", strings.Count(reply,"\n"))
	log.Printf("gpt response text: %s \n", reply)

	reply = strings.ReplaceAll(reply, "机器人神了", "我晕了不知道你在说什么")
	reply = strings.ReplaceAll(reply, "机器人", "小米粒")

	reply = strings.ReplaceAll(reply, "\"", "")

	fmt.Println("返回输出内容", reply)
	fmt.Println("返回的内容长度是:", len(reply))
	if len(reply) == 0 {
		reply = "小机灵出了点小问题。先休息了，一会回来。"
	}
	return reply, nil
}
