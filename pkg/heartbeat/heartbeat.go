package heartbeat

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/inconshreveable/go-update"
	"io"
	"log"
	"net/http"
	"time"
)

type Config struct {
	Address  string        `yaml:"address"`
	Interval time.Duration `yaml:"interval"`
	NodeName string        `yaml:"nodeName"`
}

func Start(config Config) {
	timer := time.NewTicker(config.Interval * time.Second)
	defer timer.Stop()

	// 循环接收定时器的触发事件
	for range timer.C {
		// 执行发送 POST 请求的操作
		err := sendPostRequest(config)
		if err != nil {
			fmt.Println("send heartbeat post failed:", err)
		}
	}
}

func sendPostRequest(config Config) error {
	fmt.Println("sending post request... v1")
	data := Config{
		NodeName: config.NodeName,
	}

	jsonData, err := json.Marshal(data)
	if err != nil {
		return err
	}

	body := bytes.NewBuffer(jsonData)

	resp, err := http.Post(config.Address, "application/json", body)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		return errors.New(fmt.Sprintf("The request to the beat server: %s received a response code: %d",
			config.Address, resp.StatusCode))
	}
	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Error reading response body:", err)
		return err
	}
	bodyString := string(bodyBytes)
	// 现在你可以使用bodyString了
	fmt.Println(bodyString)

	// TODO 根据心跳返回的内容，配置文件
	// TODO 根据心跳返回的内容，自我更新程序
	doUpdate()
	return nil
}

func doUpdate() {
	// 下载更新文件
	updateURL := "http://127.0.0.1:8000/static/loggie"
	resp, err := http.Get(updateURL)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()

	err = update.Apply(resp.Body, update.Options{})
	if err != nil {
		// 错误处理
	}

	// 在这里写应用程序的逻辑
	fmt.Println("应用程序更新成功！")
}
