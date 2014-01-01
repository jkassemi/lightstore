package lightstore

type Index struct {
  Fn     func(interface{}) interface{}
  Name   string
  Unique bool
  data   map[interface{}][]interface{}
}

type LightStore struct {
  indexes map[string]*Index
  data    []interface{}
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
  l.indexes[indexDefinition.Name] = indexDefinition
}

func (l *LightStore) AddRecord(r interface{}) {
  l.data = append(l.data, r)

  ch := make(chan bool, len(l.indexes))

  for _, index := range l.indexes {
    go func() {
      indexKey := index.Fn(r)

      if index.data == nil {
        index.data = make(map[interface{}][]interface{})
      }

      if index.Unique == true || index.data[indexKey] == nil {
        index.data[indexKey] = []interface{}{r}
      } else {
        index.data[indexKey] = append(index.data[indexKey], r)
      }

      ch <- true
    }()
  }

  <-ch
}

func (l *LightStore) RemoveRecord(r interface{}) {
  l.data = rm(l.data, r)

  for _, index := range l.indexes {
    indexKey := index.Fn(r)

    if index.data[indexKey] != nil {
      index.data[indexKey] = rm(index.data[indexKey], r)
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
