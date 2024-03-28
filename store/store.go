package storages

import (
	"sync"
	"time"
)

// Эта структура описывает состояние сущности host
// Status = 0 - отвечает на пинг
// Status = 1 - не ответил на пинг в течении Time_response >= 2 сек
// Date_time - время события
type Host struct {
	Ip            string
	Time_response float64
	Status        int
	Date_time     string
	Hostname      string
	Descriptor    string
	TimeStamp     time.Time
}

type DataStore interface {
	WriteData(key string, data Host)
	ReadData(key string) (Host, bool)
	ReadAll()(map[string]Host, bool)
}

type InMemoryDataStore struct {
	data  map[string]Host
	mutex sync.RWMutex
}

func NewInMemoryDataStore() *InMemoryDataStore {
	return &InMemoryDataStore{
		data:  make(map[string]Host),
		mutex: sync.RWMutex{},
	}
}

func (ds *InMemoryDataStore) WriteData(key string, data Host) {
	ds.mutex.Lock()
	defer ds.mutex.Unlock()
	ds.data[key] = data
}

func (ds *InMemoryDataStore) ReadData(key string) (Host, bool) {
	ds.mutex.RLock()
	defer ds.mutex.RUnlock()
	data, ok := ds.data[key]
	return data, ok
}

func (ds *InMemoryDataStore) ReadAll() (map[string]Host, bool) {
	ds.mutex.Lock()
	defer ds.mutex.Unlock()
	return ds.data, true
}



func Version() string {

	return "Vers1.0"
}

//func main() {
// Пример использования
//	dataStore := NewInMemoryDataStore()
//	dataStore.WriteData("host1", Host{Ip: "192.168.1.1", Tr: 4.5, St: 1, Dt: "2021-12-31"})
//	hostData, ok := dataStore.ReadData("host1")
//	if ok {
//		fmt.Println("Data found: ", hostData)
//	} else {
//		fmt.Println("Data not found")
//	}
//}
