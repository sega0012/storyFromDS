package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
)

type Message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type Request struct {
	Model    string    `json:"model"`
	Messages []Message `json:"messages"`
	Stream   bool      `json:"stream"`
}

type Response struct {
	Choices []struct {
		Message Message `json:"message"`
	} `json:"choices"`
}

func CallDeepSeekAPI(prompt, title, chapter string, idindex int) {
	var res string

	url := "https://api.deepseek.com/chat/completions"

	reqBody := Request{
		Model: model,
		Messages: []Message{
			{Role: "system", Content: worldSys},
			{Role: "user", Content: strings.TrimSpace(prompt)},
		},
		Stream: false,
	}

	jsonData, err := json.Marshal(reqBody)
	if err != nil {
		fmt.Println("JSON编码错误:", err)
		return
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		fmt.Println("创建请求错误:", err)
		return
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+apiKey)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("发送请求错误:", err)
		return
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("读取响应错误:", err)
		return
	}

	if resp.StatusCode != http.StatusOK {
		fmt.Printf("API错误: %s\n", body)
		return
	}

	var response Response
	err = json.Unmarshal(body, &response)
	if err != nil {
		fmt.Println("JSON解码错误:", err)
		return
	}

	if len(response.Choices) > 0 {
		res = response.Choices[0].Message.Content
	} else {
		fmt.Println("无响应内容")
		res = ""
	}
	xs := XS{
		Title:   title,
		Chapter: chapter,
		Content: res,
		Prompt:  prompt,
		IDindex: idindex,
	}
	err = insertXS(xs)
	fmt.Println(chapter + "完成")
	wg.Done()
}
