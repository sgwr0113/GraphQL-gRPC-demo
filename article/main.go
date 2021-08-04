// gRPCサーバーの動作確認用のファイル
package main

import (
	"context"
	"fmt"
	"io"
	"log"

	"github.com/k88t76/GraphQL-gRPC-demo/article/client"
	"github.com/k88t76/GraphQL-gRPC-demo/article/pb"
)

func main() {
	// clientを生成
	c, _ := client.NewClient("localhost:50051")
	create(c)
	// read(c)
	// update(c)
	// delete(c)
	// list(c)
}

func create(c *client.Client) {
	// 広告をCREATE
	input := &pb.ArticleInput{
		DisplayAppName: "テストアプリ",
		IconSrc:        "testSrc",
		CvCondition:    "インストールで獲得",
		BasePoint:      1000,
	}
	res, err := c.Service.CreateArticle(context.Background(), &pb.CreateArticleRequest{ArticleInput: input})
	if err != nil {
		log.Fatalf("Failed to CreateArticle: %v\n", err)
	}
	fmt.Printf("CreateArticle Response: %v\n", res)
}

func read(c *client.Client) {
	// 広告をREAD
	var id int64 = 16
	res, err := c.Service.ReadArticle(context.Background(), &pb.ReadArticleRequest{Id: id})
	if err != nil {
		log.Fatalf("Failed to ReadArticle: %v\n", err)
	}
	fmt.Printf("ReadArticle Response: %v\n", res)
}

func update(c *client.Client) {
	// 広告をUPDATE
	var id int64 = 16
	input := &pb.ArticleInput{
		DisplayAppName: "テストアプリ改",
		IconSrc:        "newSrc",
		CvCondition:    "Lv.100到達で獲得",
		BasePoint:      2500,
	}
	res, err := c.Service.UpdateArticle(context.Background(), &pb.UpdateArticleRequest{Id: id, ArticleInput: input})
	if err != nil {
		log.Fatalf("Failed to UpdateArticle: %v\n", err)
	}
	fmt.Printf("UpdateArticle Response: %v\n", res)
}

func delete(c *client.Client) {
	// 広告をDELETE
	var id int64 = 13
	res, err := c.Service.DeleteArticle(context.Background(), &pb.DeleteArticleRequest{Id: id})
	if err != nil {
		log.Fatalf("Failed to UpdateArticle: %v\n", err)
	}
	fmt.Printf("The article has been deleted (%v)\n", res)
}

func list(c *client.Client) {
	// 広告を全取得
	stream, err := c.Service.ListArticle(context.Background(), &pb.ListArticleRequest{})
	if err != nil {
		log.Fatalf("Failed to ListArticle: %v\n", err)
	}

	// Server Streamingで渡されたレスポンスを１つ１つ受け取る
	for {
		res, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatalf("Failed to Server Streaming: %v\n", err)
		}
		fmt.Println(res)
	}
}
