package main

import (
	"fmt"
	"log"
	"rabbit-mq-Publish-Subscribe/config"
	"time"

	"github.com/streadway/amqp"
)

func main() {
	author := "author"
	users := []string{"user1", "user2", "user3"}
	// 启动发布者
	go startPublisher(author)

	//启动用户订阅服务
	for i := range users {
		go startMessageProcessor(author, users[i])
	}

	// 阻塞主线程
	select {}
}

func startPublisher(author string) {
	conn, err := amqp.Dial(config.RabbitMQURL)
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	ch, err := conn.Channel()
	if err != nil {
		log.Fatal(err)
	}
	defer ch.Close()

	authorExchange := config.ArticleExchange + "_" + author

	// 声明扇出交换机
	err = ch.ExchangeDeclare(authorExchange, "fanout", true, false, false, false, nil)
	if err != nil {
		log.Fatal(err)
	}

	// 模拟发布文章
	for i := 1; i <= 5; i++ {
		article := fmt.Sprintf("Article %d", i)
		err = ch.Publish(authorExchange, "", false, false, amqp.Publishing{
			ContentType: "text/plain",
			Body:        []byte(article),
		})
		if err != nil {
			log.Println("Failed to publish article:", err)
		}
		time.Sleep(time.Second) // 模拟发布间隔
	}
}

func startMessageProcessor(author, user string) {
	conn, err := amqp.Dial(config.RabbitMQURL)
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	ch, err := conn.Channel()
	if err != nil {
		log.Fatal(err)
	}
	defer ch.Close()

	authorExchange := config.ArticleExchange + "_" + author
	userQueue := config.ArticleQueue + "_" + user

	// 声明队列
	_, err = ch.QueueDeclare(userQueue, true, false, false, false, nil)
	if err != nil {
		log.Fatal(err)
	}

	// 绑定队列到交换机
	err = ch.QueueBind(userQueue, "", authorExchange, false, nil)
	if err != nil {
		log.Fatal(err)
	}

	// 消费消息
	msgs, err := ch.Consume(userQueue, "", true, false, false, false, nil)
	if err != nil {
		log.Fatal(err)
	}

	// 处理消息
	for msg := range msgs {
		processArticle(msg.Body)
	}
}

//
func processArticle(article []byte) {
	// 模拟异步处理文章的操作
	log.Printf("Processing article: %s\n", article)
	time.Sleep(time.Second * 2)
	log.Println("Article processed.")
}
