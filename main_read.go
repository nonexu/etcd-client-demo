package main

import (
	"context"
	"log"
	"strconv"
	"time"

	"go.etcd.io/etcd/client/v3"
)

const (
	ETCD_TEST_KEY = "/dev/game/opentime"
)

var open_time int64
var open_time_str string

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

	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	getResp, err := cli.Get(ctx, ETCD_TEST_KEY, clientv3.WithPrefix())
	cancel()
	if err != nil {
		log.Printf("get from etcd failed, err:%v\n", err)
		return
	}

	log.Println("getResp:", getResp)

	for _, ev := range getResp.Kvs {
		log.Printf("%s:%s\n", ev.Key, ev.Value)
	}
	go watchKey(cli, ETCD_TEST_KEY, &open_time, sto64)
	go watchKey(cli, ETCD_TEST_KEY, &open_time_str, stos)
	go debugInfo()

	rch := cli.Watch(context.Background(), ETCD_TEST_KEY)
	for wresp := range rch {
		for _, _ = range wresp.Events {
			//log.Printf("%s %q : %q\n", ev.Type, ev.Kv.Key, ev.Kv.Value)
		}
	}
}

func debugInfo() {
	for {
		time.Sleep(time.Second)
		log.Println("golbal key update key", open_time, open_time_str)
	}
}

func sto64(data interface{}, str string) {
	log.Println("sto64", str)
	value, _ := strconv.ParseInt(str, 10, 64)
	*data.(*int64) = value
}

func stos(data interface{}, str string) {
	log.Println("stos", str)
	*data.(*string) = str
}

func preWatch(cli *clientv3.Client, key string, data interface{}, convert func(interface{}, string)) {

	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	getResp, err := cli.Get(ctx, ETCD_TEST_KEY, clientv3.WithPrefix())
	cancel()
	if err != nil {
		log.Printf("get from etcd failed, err:%v\n", err)
		return
	}
	for _, env := range getResp.Kvs {
		convert(data, string(env.Value))
	}
}

func watchKey(cli *clientv3.Client, key string, data interface{}, convert func(interface{}, string)) {
	log.Println("start etcd watch key :", key)
	rch := cli.Watch(context.Background(), key, clientv3.WithPrevKV())
	preWatch(cli, key, data, convert)
	for wresp := range rch {
		for _, ev := range wresp.Events {
			switch ev.Type {
			case clientv3.EventTypePut:
				convert(data, string(ev.Kv.Value))
				/* code */
			case clientv3.EventTypeDelete:
				convert(data, string(ev.Kv.Value))
			}
		}
	}

}
