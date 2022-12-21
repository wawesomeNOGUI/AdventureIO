package main

import (
	"sync"

	"github.com/pion/webrtc/v3"
)

// Map of data channels so we can broadcast messages to everyone
// with mutex so different goroutines can use with no contention
type DataChannelContainer struct {
	mu sync.Mutex
	chans map[string]*webrtc.DataChannel
}

func (c *DataChannelContainer) Broadcast(msg string) {
	c.mu.Lock()
    defer c.mu.Unlock()

	for _, v := range c.chans {
		v.SendText(msg)
	}
}

func (c *DataChannelContainer) SendToPlayer(tag string, msg string) {
	c.mu.Lock()
    defer c.mu.Unlock()

	_, ok := c.chans[tag]
	if ok {
		c.chans[tag].SendText(msg)
	}
}

func (c *DataChannelContainer) DeletePlayerChan(tag string) {
	c.mu.Lock()
    defer c.mu.Unlock()

	delete(c.chans, tag)
}

func (c *DataChannelContainer) AddPlayerChan(tag string, d *webrtc.DataChannel) {
	c.mu.Lock()
    defer c.mu.Unlock()

	c.chans[tag] = d
}