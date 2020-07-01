package config

import (
	"github.com/go-redis/redis"
	"go.mongodb.org/mongo-driver/mongo"
)

var (
	ConsulAddress    string //consul地址
	EtcdAddress      string //etcd地址
	LogPath          string //日志地址
	LogLevel         int    //日志级别
	LogMaxAge        int    //日志保留时长
	LogMaxSize       int    //日志最大size
	LogMaxBackups    int    //日志最多留备份数
	MongoReplicas    string //mongdb地址
	MongoMaxPoolSize int64  //最大连接池数
	MongoMinPoolSize int64  //最小连接池数
)

var (
	// 声明一个全局的reidsclient变量
	RedisClient *redis.Client
	// Pool MongoDB连接池
	Pool   *mongo.Client
	Client *mongo.Client
)

type Node struct {
	Data int
	Next *Node
}

func CreateNode(head *Node, len int) *Node {
	if head == nil {
		head = &Node{Data: -1}
	}
	tNode := head
	for i := 0; i < len; i++ {
		node := &Node{
			Data: i,
		}

		tNode.Next = node
		tNode = node
	}
	return head
}

func Reverse(node *Node) *Node {
	for node.Next == nil || node == nil {
		return nil
	}
	var pre *Node
	var cur *Node
	next := node.Next

	for next != nil {
		cur = next.Next
		next.Next = pre
		pre = next
		next = cur
	}
	return pre
}
