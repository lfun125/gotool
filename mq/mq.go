package mq

import (
	"strings"
	"sync"
	"time"

	"github.com/kr/beanstalk"
	"github.com/lfun125/gotool/logger"
	"github.com/lfun125/gotool/run"
)

var log logger.Interface

type Data struct {
	Error chan error
	Item  chan *Item
}

type Item struct {
	Body []byte
	ID   uint64
	Conn *beanstalk.Conn
	Wait sync.WaitGroup
}

func newData() *Data {
	data := &Data{}
	data.Error = make(chan error)
	data.Item = make(chan *Item)
	return data
}

func SetLogger(logger logger.Interface) {
	log = logger
}

func Put(addr string, data []byte, key string, pri uint32, delay, trr time.Duration) (id uint64, err error) {
	var conn *beanstalk.Conn
	if conn, err = beanstalk.Dial("tcp", addr); err != nil {
		return
	}
	defer func() {
		if err := conn.Close(); err != nil {
			log.Error(err)
		}
	}()
	tube := beanstalk.Tube{Conn: conn, Name: key}
	id, err = tube.Put(data, pri, delay, trr)
	return
}

func PutAt(addr string, data []byte, key string, pri uint32, t time.Time, trr time.Duration) (id uint64, err error) {
	delay := t.Sub(time.Now())
	if delay < 0 {
		delay = 0
	}
	return Put(addr, data, key, pri, delay, trr)
}

func Subscribe(addr, tube, tag string, do ProcessFunc) (err error) {
	data := newData()
	if err = watch(addr, tube, tag, data); err != nil {
		return
	}
	for {
		select {
		case err = <-data.Error:
			return
		case item := <-data.Item:
			process(item, do)
		}
	}
}

func watch(addr, tube, tag string, data *Data) (err error) {
	var conn *beanstalk.Conn
	if conn, err = beanstalk.Dial("tcp", addr); err != nil {
		return
	}
	tubeSet := beanstalk.NewTubeSet(conn, tube)
	run.GO(log, func() {
		defer func() {
			if err := conn.Close(); err != nil {
				log.Error(err)
			}
		}()
		for {
			id, body, err := tubeSet.Reserve(30 * time.Minute)
			if err != nil && strings.HasPrefix(err.Error(), "reserve-with-timeout") {
				continue
			} else if err != nil {
				log.With("tag", tag).With("err", err).Error("beanstalk reserve error")
				data.Error <- err
				return
			}
			bodyData := map[string]interface{}{
				tag: string(body),
			}
			log.With("tag", tag).With("id", id, "body_data", bodyData).Info("reserve new message")
			item := &Item{
				Body: body,
				ID:   id,
				Conn: conn,
				Wait: sync.WaitGroup{},
			}
			item.Wait.Add(1)
			data.Item <- item
			item.Wait.Wait()
		}
	}, 1)
	return
}

func beanstalkDelete(conn *beanstalk.Conn, id uint64) {
	log.With("id", id).Info("beanstalk delete")
	if err := conn.Delete(id); err != nil {
		log.With("err", err, "id", id).Error("beanstalk.Delete")
	}
}

func beanstalkRelease(conn *beanstalk.Conn, id uint64, delay time.Duration) {
	log.With("id", id, "delay", delay).Info("beanstalk release")
	if err := conn.Release(id, 1024, delay); err != nil {
		log.With("err", err, "id", id).Error("beanstalk.Release")
	}
}

type ProcessFunc func(item *Item) (delay time.Duration, isDel bool, err error)

func process(item *Item, f ProcessFunc) {
	defer func() {
		item.Wait.Done()
	}()

	delay, isDel, err := f(item)
	if err != nil {
		log.With("beanstalk.id", item.ID).Error(err)
	}
	if isDel {
		beanstalkDelete(item.Conn, item.ID)
	} else {
		if err != nil && delay == 0 {
			delay = 15 * time.Second
		}
		beanstalkRelease(item.Conn, item.ID, delay)
	}
}
