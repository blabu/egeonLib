package golang

import (
	"fmt"
	"sync"
	"time"
)

var allStores []LocalCache // Все кеши базы данных собраны здесь
var mt sync.Mutex          // Этот мьтекс защищает изменения в allStores

// entry - конкретная запись в кеше (создана для агрегирования вермени жизни записи)
type entry struct {
	value      interface{}
	expireTime time.Time
}

// LocalCache - Предназначен для локального харнения данных с доступом на чтения
// Хранит ряд таблиц распределеных по идентификаторам. В одном LocalCache может быть множество таблиц с разным идентификатором
type LocalCache struct {
	expire time.Duration // Для всех записей этого кеша будет применятся заданный промежуток протухания с даты создания записи
	mt     *sync.Mutex
	store  map[uint32]*sync.Map // список таблиц данных для конкретного кеша
}

//GetNewCache - создание нового кеша (аналог базы данных) для хранения данных в памяти программы
// принимает expire - время жизни каждой записи
// и идентификаторы на основе которых создаются хранилища ключей-значений (аналог таблицы в БД)
func GetNewCache(expire time.Duration, ids ...uint32) *LocalCache {
	cache := LocalCache{
		store:  make(map[uint32]*sync.Map),
		mt:     &sync.Mutex{},
		expire: expire,
	}
	for _, i := range ids {
		cache.store[i] = &sync.Map{}
	}
	mt.Lock()
	defer mt.Unlock()
	allStores = append(allStores, cache)
	return &cache
}

//AddStorage - создает новое хранилище ключ-значение (новую таблицу)
func (lc *LocalCache) AddStorage(id uint32) {
	lc.mt.Lock()
	if _, ok := lc.store[id]; !ok {
		lc.store[id] = &sync.Map{}
	}
	lc.mt.Unlock()
}

//StoreItem - сохраняет ключ значение в заданную таблицу с идентификатором id
func (lc *LocalCache) StoreItem(id uint32, key interface{}, val interface{}) {
	value := entry{
		value:      val,
		expireTime: time.Now().Add(lc.expire),
	}
	cached, ok := lc.store[id]
	if ok && cached != nil {
		cached.Store(key, value)
	} else {
		panic(fmt.Sprintf("Undefine id %d", id))
	}
}

//GetItem - чтение из таблицы с идентификатором id значения по ключу key
func (lc *LocalCache) GetItem(id uint32, key interface{}) interface{} {
	cached, ok := lc.store[id]
	if ok && cached != nil {
		if val, ok := cached.Load(key); ok {
			if v, ok := val.(entry); ok {
				if lc.expire != 0 && v.expireTime.Before(time.Now()) {
					return nil
				}
				return v.value
			}
		}
	} else {
		cached = &sync.Map{}
		lc.store[id] = cached
	}
	return nil
}
