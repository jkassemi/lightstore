package lightstore

import (
  "sync"
)

type Index struct {
  Fn     func(interface{}) []interface{}
  Name   string
  Unique bool
  data   map[interface{}][]interface{}
  mu     sync.Mutex
}

type LightStore struct {
  indexes map[string]*Index
  data    []interface{}
  mu      sync.Mutex
}

func NewStore() *LightStore {
  return &LightStore{
    indexes: make(map[string]*Index),
    data:    make([]interface{}, 0),
  }
}

func rm(haystack []interface{}, needle interface{}) []interface{} {
  found := -1

  for i, v := range haystack {
    if v == needle {
      found = i
      break
    }
  }

  if found == -1 {
    return haystack
  }

  return append(haystack[:found], haystack[found+1:]...)
}

func (l *LightStore) DefineIndex(indexDefinition *Index) {
  l.mu.Lock()
  defer l.mu.Unlock()

  l.indexes[indexDefinition.Name] = indexDefinition
}

func (l *LightStore) AddRecord(r interface{}) {
  l.mu.Lock()
  l.data = append(l.data, r)
  l.mu.Unlock()

  ch := make(chan bool, len(l.indexes))

  for _, index := range l.indexes {
    go func() {
      index.mu.Lock()
      defer index.mu.Unlock()

      if index.data == nil {
        index.data = make(map[interface{}][]interface{})
      }

      indexKeys := index.Fn(r)

      for _, indexKey := range indexKeys {
        if index.Unique == true || index.data[indexKey] == nil {
          index.data[indexKey] = []interface{}{r}
        } else {
          index.data[indexKey] = append(index.data[indexKey], r)
        }
      }

      ch <- true
    }()
  }

  <-ch
}

func (l *LightStore) RemoveRecord(r interface{}) {
  l.mu.Lock()
  l.data = rm(l.data, r)
  l.mu.Unlock()

  for _, index := range l.indexes {
    indexKeys := index.Fn(r)

    for _, indexKey := range indexKeys {
      if index.data[indexKey] != nil {
        index.mu.Lock()
        index.data[indexKey] = rm(index.data[indexKey], r)
        index.mu.Unlock()
      }
    }
  }
}

func (l *LightStore) Query(indexName string, key interface{}) []interface{} {
  index := l.indexes[indexName]

  if index.data == nil || index.data[key] == nil {
    return make([]interface{}, 0)
  }

  return index.data[key]
}
