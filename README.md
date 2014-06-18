# Lightstore

*NO LONGER ACTIVELY MAINTAINED* - this project should be considered Archived. If you are interested in maintaining an official fork, please contact me directly.

Minimal native-hash indexing in-memory data storage. Define index functions and
query this structure based on those indexes. 

## Installation

```bash
go get github.com/jkassemi/lightstore
```

## Usage

Import the library

```go
import (
  "log"
  "github.com/jkassemi/lightstore"
)
```

Create a record structure

```go
type struct MyRecord {
  Name        string
  Value       string
}
```

Initialize a data storage instance and define some indexes 

```go
var store *lighstore.LightStore

func init(){
  store = lightstore.NewStore()

  store.DefineIndex(&Index{
    Name: "byName",
    Fn: func(v interface{}) []interface{} {
      return []interface{}{v.(*MyRecord).Name}
    },
  })

  store.DefineIndex(&Index{
    Name: "byValue",
    Fn: func(v interface{}) []interface{} {
      return []interface{}{v.(*MyRecord).Value}
    },
  })
}
```

Adds some records and perform some queries

```go
func main(){
  records := []*MyRecord{
    &MyRecord{Name: "Hello", Value: "Mars"},
    &MyRecord{Name: "Hello", Value: "Colonization"},
    &MyRecord{Name: "Goodbye", Value: "Earth"},
    &MyRecord{Name: "Goodbye", Value: "Imaginatory Stagnation"},
  }

  for _, record := range records {
    store.AddRecord(record)
  }

  hellos := store.Query("byName", "Hello")
  // hellos == [records[0], records[1]] 

  goodbyes := store.Query("byName", "Goodbye")
  // goodbyes == [records[2], records[3]]

  mars := store.Query("byValue", "Mars")
  // mars == [records[0]]
}
```
