/**
屏蔽IP

- 实现思路：保存IP的最近x次访问的时间，最新与最后一次对比，如果时间间隔少于n秒，则视为恶意IP(一秒钟超过10个请求则屏蔽)
- 目的：屏蔽一定程度的CC攻击
*/
package libs

import (
	"container/list"
	"sync"
	"time"
)

const (
	MAX_BETWEEN = 1  // 最后一次请求和最近的请求的间隔多少秒
	MAX_LIMIT   = 20 // MAX_BETWEEN时间内的请求量

	MAX_FILTER_LIMIT = 200
)

type RequestAt struct {
	Ip string
	At time.Time
}

var Filter = make(chan RequestAt, MAX_FILTER_LIMIT)

type BlackWall struct {
	blacks   map[string]bool
	runtimes map[string]*list.List

	max_limit   int
	max_between int64

	done chan bool

	mtx sync.Mutex
}

var Wall *BlackWall

func init() {
	if Wall == nil {
		Wall = &BlackWall{
			max_limit:   MAX_LIMIT,
			max_between: MAX_BETWEEN,
			runtimes:    make(map[string]*list.List),
			blacks:      make(map[string]bool),
		}
		Wall.Start()
	}
}

func (r *BlackWall) Start() {
	go func() {
		for {
			select {
			case req := <-Filter:
				r.AddIp(req)
			case <-r.done:
				return
			}
		}
	}()
}

func (r *BlackWall) Stop() {
	r.done <- true
}

func (r *BlackWall) AddIp(ip RequestAt) {
	if r.IsBlack(ip.Ip) {
		return
	}

	r.mtx.Lock()
	defer r.mtx.Unlock()

	if r.runtimes[ip.Ip] == nil {
		r.runtimes[ip.Ip] = list.New()
	}
	r.runtimes[ip.Ip].PushFront(ip.At)

	r.check(ip)
}

func (r *BlackWall) Black(ip string) {
	r.blacks[ip] = true
}

func (r *BlackWall) IsBlack(ip string) bool {
	r.mtx.Lock()
	defer r.mtx.Unlock()

	return r.blacks[ip]
}

func (r *BlackWall) check(ip RequestAt) {
	if r.runtimes[ip.Ip].Len() < r.max_limit {
		return
	}
	back := r.runtimes[ip.Ip].Back().Value.(time.Time)
	front := r.runtimes[ip.Ip].Front().Value.(time.Time)
	if (front.Unix() - back.Unix()) <= r.max_between { // 最新和最后一次请求时间小于等于MAX_BETWEEN则为非法请求
		r.Black(ip.Ip)
	}
	r.runtimes[ip.Ip].Init() // 超过MAX_LIMIT个请求清空
}

func (r *BlackWall) Remove(ip string) {
	r.mtx.Lock()
	defer r.mtx.Unlock()

	delete(r.blacks, ip)
}
