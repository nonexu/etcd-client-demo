package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"go.etcd.io/etcd/client/v3"
)

const (
	ETCD_TEST_KEY = "/dev/game/opentime"
)

func main() {
	cli, err := clientv3.New(clientv3.Config{
		Endpoints:   []string{"localhost:2379"},
		DialTimeout: 5 * time.Second,
	})

	if err != nil {
		log.Println("connect to etcd failed, err:", err)
		return
	}
	defer cli.Close()

	for {
		time.Sleep(time.Second / 1000)
		timeStr := fmt.Sprintf("%d", time.Now().Unix())

		ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
		putResp, err := cli.Put(ctx, ETCD_TEST_KEY, timeStr)
		cancel()

		if err != nil {
			log.Println("put to ectcd failed:", err)
			return
		}
		log.Println("putResp:", putResp)
		continue
		time.Sleep(time.Second)

		ctx, cancel = context.WithTimeout(context.Background(), 1*time.Second)
		deleteResp, err := cli.Delete(ctx, ETCD_TEST_KEY)
		cancel()

		if err != nil {
			log.Println("put to ectcd failed:", err)
			return
		}
		log.Println("deleteResp:", deleteResp)

	}

}
