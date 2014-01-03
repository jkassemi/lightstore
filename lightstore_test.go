package lightstore

import (
  "testing"
)

var indexQueryTests = []struct {
  indexFn func(interface{}) []interface{}
  records []interface{}
  queries []interface{}
  results [][]interface{}
}{
  {
    indexFn: func(v interface{}) []interface{} { return []interface{}{v} },
    records: []interface{}{"hello", "world"},
    queries: []interface{}{"hello", "world"},
    results: [][]interface{}{[]interface{}{"hello"}, []interface{}{"world"}},
  },
  {
    indexFn: func(v interface{}) []interface{} { return []interface{}{"hello"} },
    records: []interface{}{"hello", "world"},
    queries: []interface{}{"hello"},
    results: [][]interface{}{[]interface{}{"hello", "world"}},
  },
  {
    indexFn: func(v interface{}) []interface{} { return []interface{}{"hello", v} },
    records: []interface{}{"r1", "r2"},
    queries: []interface{}{"r1", "r2", "hello"},
    results: [][]interface{}{[]interface{}{"r1"}, []interface{}{"r2"}, []interface{}{"r1", "r2"}},
  },
}

func comp(v1 []interface{}, v2 []interface{}) bool {
  if len(v1) != len(v2) {
    return false
  }

  for i, v := range v1 {
    if v2[i] != v {
      return false
    }
  }

  return true
}

func TestIndexQueries(t *testing.T) {
  for i, iq := range indexQueryTests {
    ls := NewStore()

    ls.DefineIndex(&Index{
      Name: "idx",
      Fn:   iq.indexFn,
    })

    for _, r := range iq.records {
      ls.AddRecord(r)
    }

    for j, query := range iq.queries {
      results := ls.Query("idx", query)

      if !comp(results, iq.results[j]) {
        t.Errorf("%d. .Query(%q, %q) => %q, want %q", i, iq.records, query, results, iq.results[j])
      }
    }
  }
}

func TestMultiIndex(t *testing.T) {
  ls := NewStore()

  ls.DefineIndex(&Index{
    Name: "idx1",
    Fn: func(v interface{}) []interface{} {
      if v.(string) == "r1" {
        return []interface{}{v}
      } else {
        return []interface{}{}
      }
    },
  })

  ls.DefineIndex(&Index{
    Name: "idx2",
    Fn: func(v interface{}) []interface{} {
      if v.(string) == "r2" {
        return []interface{}{v}
      } else {
        return []interface{}{}
      }
    },
  })

  records := []interface{}{"r1", "r2"}

  for _, r := range records {
    ls.AddRecord(r)
  }

  if !comp(ls.Query("idx1", "r1"), []interface{}{"r1"}) ||
    !comp(ls.Query("idx2", "r2"), []interface{}{"r2"}) {
    t.Errorf("Unexpected output")
  }
}
