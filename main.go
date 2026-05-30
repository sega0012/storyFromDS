package main

import (
	"fmt"
	"os"
	"strings"
	"sync"

	"github.com/xuri/excelize/v2"
)

var (
	wg           sync.WaitGroup
	prompts      string
	apiKey       string
	model        string
	worldSys     string
	chapterOrder string
)

func main() {
	f, err := excelize.OpenFile("大纲.xlsx")
	if err != nil {
		fmt.Println(err)
		return
	}
	configRows, _ := f.GetRows("配置")
	apiKey = strings.TrimSpace(configRows[1][0])
	model = strings.TrimSpace(configRows[1][1])
	worldSys = strings.TrimSpace(configRows[1][2])
	chapterOrder = strings.TrimSpace(configRows[1][3])
	rows, _ := f.GetRows("细纲")
	for i, row := range rows {
		if i == 0 {
			continue
		}
		prompts += row[1] + row[2] + "\n"
		if row[3] == "否" {
			continue
		}
		wg.Add(1)
		zjprompt := fmt.Sprintf("%s,这是小说大纲，现在写%s的内容，%s。", prompts, row[1], chapterOrder)
		go CallDeepSeekAPI(zjprompt, row[0], row[1], i)

	}
	// for j, r := range rows {
	// 	if j == 0 {
	// 		continue
	// 	}
	// 	if r[3] == "否" {
	// 		continue
	// 	}
	// 	wg.Add(1)
	// 	zjprompt := fmt.Sprintf("%s,这是小说大纲，现在写%s的内容，写5000字的小说，要有详细的性爱描写。", prompts, r[1])
	// 	go CallDeepSeekAPI(zjprompt, r[0], r[1], j)
	// }
	wg.Wait()
	title := rows[1][0]
	xsList, err := getXs(title)
	if err != nil {
		fmt.Println("获取数据错误:", err)
		return
	}
	filename := fmt.Sprintf("%s.txt", title)
	filename = strings.ReplaceAll(filename, "/", "_")
	filename = strings.ReplaceAll(filename, "\\", "_")
	filename = strings.ReplaceAll(filename, ":", "_")
	filename = strings.ReplaceAll(filename, "*", "_")
	filename = strings.ReplaceAll(filename, "?", "_")
	filename = strings.ReplaceAll(filename, "\"", "_")
	filename = strings.ReplaceAll(filename, "<", "_")
	filename = strings.ReplaceAll(filename, ">", "_")
	filename = strings.ReplaceAll(filename, "|", "_")
	var data string
	for _, xs := range xsList {

		data += fmt.Sprintf("Chapter: %s\nContent:\n%s \n", xs.Chapter, xs.Content)

	}
	err = os.WriteFile(filename, []byte(data), 0644)
	if err != nil {
		fmt.Printf("写入文件失败 %s: %v\n", filename, err)
	}
}
