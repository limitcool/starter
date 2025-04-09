package main

import (
	"bytes"
	"fmt"
	"io"
	"log"

	"github.com/limitcool/starter/internal/pkg/storage"
)

func main() {
	// 创建本地存储
	localConfig := storage.Config{
		Type: storage.StorageTypeLocal,
		Path: "./uploads",
		URL:  "http://localhost:8080/uploads", // 本地开发URL
	}

	localStorage, err := storage.New(localConfig)
	if err != nil {
		log.Fatalf("创建本地存储失败: %v", err)
	}

	// 上传文件示例
	fileContent := []byte("这是一个测试文件内容")
	reader := bytes.NewReader(fileContent)

	filePath := storage.GeneratePath("avatars", "user1.jpg")
	err = localStorage.Put(filePath, reader)
	if err != nil {
		log.Fatalf("上传文件失败: %v", err)
	}

	fmt.Println("文件上传成功!")

	// 获取文件URL
	url, err := localStorage.GetURL(filePath)
	if err != nil {
		log.Fatalf("获取文件URL失败: %v", err)
	}

	fmt.Printf("文件访问URL: %s\n", url)

	// 读取文件
	file, err := localStorage.Get(filePath)
	if err != nil {
		log.Fatalf("获取文件失败: %v", err)
	}
	defer file.Close()

	content, err := io.ReadAll(file)
	if err != nil {
		log.Fatalf("读取文件内容失败: %v", err)
	}

	fmt.Printf("文件内容: %s\n", content)

	// 列出目录文件
	objects, err := localStorage.List("avatars")
	if err != nil {
		log.Fatalf("列出文件失败: %v", err)
	}

	fmt.Println("目录下的文件:")
	for _, obj := range objects {
		fmt.Printf(" - %s\n", obj.Name)
	}

	// 获取文件MIME类型
	mime := storage.GetMimeType("user1.jpg")
	fmt.Printf("文件MIME类型: %s\n", mime)

	// S3示例配置
	/*
		s3Config := storage.Config{
			Type:      storage.StorageTypeS3,
			AccessKey: "your-access-key",
			SecretKey: "your-secret-key",
			Region:    "us-west-2",
			Bucket:    "your-bucket-name",
			Endpoint:  "https://s3.us-west-2.amazonaws.com",
		}

		s3Storage, err := storage.New(s3Config)
		if err != nil {
			log.Fatalf("创建S3存储失败: %v", err)
		}

		// 使用与本地存储相同的API
	*/
}
