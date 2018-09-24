# radix (Radix tree implementation in Go)

[![Build status](https://travis-ci.org/gbrlsnchs/radix.svg?branch=master)](https://travis-ci.org/gbrlsnchs/radix)
[![Build status](https://ci.appveyor.com/api/projects/status/0eyx2cdvefhx0xo5/branch/master?svg=true)](https://ci.appveyor.com/project/gbrlsnchs/radix/branch/master)
[![Sourcegraph](https://sourcegraph.com/github.com/gbrlsnchs/radix/-/badge.svg)](https://sourcegraph.com/github.com/gbrlsnchs/radix?badge)
[![GoDoc](https://godoc.org/github.com/gbrlsnchs/radix?status.svg)](https://godoc.org/github.com/gbrlsnchs/radix)
[![Minimal version](https://img.shields.io/badge/minimal%20version-go1.10%2B-5272b4.svg)](https://golang.org/doc/go1.10)

## About
This package is an implementation of a [radix tree](https://en.wikipedia.org/wiki/Radix_tree) in [Go](https://golang.org) (or Golang).  

Searching for static values in the tree doesn't allocate memory on the heap, what makes it pretty fast.  
It can also sort nodes by priority, therefore traversing nodes that hold more non-nil values first.

## Usage
Full documentation [here](https://godoc.org/github.com/gbrlsnchs/radix).  

### Installing
#### Go 1.10
`vgo get -u github.com/gbrlsnchs/radix`
#### Go 1.11
`go get -u github.com/gbrlsnchs/radix`

### Importing
```go
import (
	// ...

	"github.com/gbrlsnchs/radix"
)
```

### Building [this example from Wikipedia](https://upload.wikimedia.org/wikipedia/commons/a/ae/Patricia_trie.svg)
```go
tr := radix.New(radix.Tdebug)
tr.Add("romane", 1)
tr.Add("romanus", 2)
tr.Add("romulus", 3)
tr.Add("rubens", 4)
tr.Add("ruber", 5)
tr.Add("rubicon", 6)
tr.Add("rubicundus", 7)
tr.Sort(radix.PrioritySort) // optional step
log.Print(tr.String())
```

#### The code above will print this
```
. (14 nodes)
└── 7↑ r → <nil>
    ├── 4↑ ub → <nil>
    │   ├── 2↑ ic → <nil>
    │   │   ├── 1↑ undus 🍂 → 7
    │   │   └── 1↑ on 🍂 → 6
    │   └── 2↑ e → <nil>
    │       ├── 1↑ r 🍂 → 5
    │       └── 1↑ ns 🍂 → 4
    └── 3↑ om → <nil>
        ├── 2↑ an → <nil>
        │   ├── 1↑ us 🍂 → 2
        │   └── 1↑ e 🍂 → 1
        └── 1↑ ulus 🍂 → 3
```

### Retrieving a value from the tree
```go
n, _ := tr.Get("rubicon") // zero-allocation search
log.Print(n.Value)        // prints "6"
```

### Building a dynamic tree
A dynamic tree is a tree that can match labels based on a placeholder and a demiliter (e.g. an HTTP router that accepts dynamic routes).  
Note that this only works with prefix trees, not binary ones.

```go
tr := radix.New(0) // passing 0 means passing no flags
tr.Add("/dynamic/path/@id", 1)
tr.Add("/dynamic/path/@id/subpath/@name", 2)
tr.Add("/static/path", 3)
tr.SetBoundaries('@', '/')

var (
	n *radix.Node
	p map[string]string
)
n, p = tr.Get("/dynamic/path/123")
log.Print(n.Value) // prints "1"
log.Print(p["id"]) // prints "123"

n, p = tr.Get("/dynamic/path/456/subpath/foobar")
log.Print(n.Value)   // prints "2"
log.Print(p["id"])   // prints "456"
log.Print(p["name"]) // prints "foobar"

n, _ = tr.Get("/static/path") // p would be nil
log.Print(n.Value)            // prints "3"
```

### Building a binary tree
```go
tr := radix.New(radix.Tdebug | radix.Tbinary)
tr.Add("deck", 1)
tr.Add("did", 2)
tr.Add("doe", 3)
tr.Add("dog", 4)
tr.Add("doge", 5)
tr.Add("dogs", 6)
```

#### The code above will print this
```
. (71 nodes)
01100100011001010110001101101011 🍂 → 1
011001000110100101100100 🍂 → 2
011001000110111101100101 🍂 → 3
011001000110111101100111 → 4
01100100011011110110011101100101 🍂 → 5
01100100011011110110011101110011 🍂 → 6
```

## Contributing
### How to help
- For bugs and opinions, please [open an issue](https://github.com/gbrlsnchs/radix/issues/new)
- For pushing changes, please [open a pull request](https://github.com/gbrlsnchs/radix/compare)
