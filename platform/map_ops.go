package platform

import (
	"fmt"
	"sync"
)

// AddChanToChanMap function to write the channel name to lock
func AddChanToChanMap(url string, channelMap map[string]chan bool, lock *sync.RWMutex, channel chan bool) {
	lock.Lock()
	defer lock.Unlock()
	channelMap[url] = channel
	fmt.Println("Added channel for ", url)
}

// ReadChanFromChanMap function to read the channel name from map
func ReadChanFromChanMap(url string, channelMap map[string]chan bool, lock *sync.RWMutex) (chan bool, bool) {
	lock.RLock()
	defer lock.RUnlock()
	val, ok := channelMap[url]
	if ok {
		return val, true
	}
	return val, false
}

// DeleteChanFromChannelMap to delete the channel from channel map
func DeleteChanFromChannelMap(url string, channelMap map[string]chan bool, lock *sync.RWMutex) bool {
	lock.Lock()
	defer lock.Unlock()
	_, ok := channelMap[url]
	if ok {
		delete(channelMap, url)
		return true
	}
	return false
}
