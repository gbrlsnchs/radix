# radix (radix tree)

[![Build Status](https://travis-ci.org/gbrlsnchs/radix.svg?branch=master)](https://travis-ci.org/gbrlsnchs/radix)
[![GoDoc](https://godoc.org/github.com/gbrlsnchs/radix?status.svg)](https://godoc.org/github.com/gbrlsnchs/radix)

<img src="https://upload.wikimedia.org/wikipedia/commons/a/ae/Patricia_trie.svg" align="right">

```javascript
Example
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

## Important
- Until version 1.0 is released, **anything** can change, including names of methods or even their existence.
- Until version [0.3.0], this package was named `patricia`, despite implementing a radix tree. 
If you're looking for a PATRICIA tree implementation, try [this package] instead.

## About
This package is an implementation of a [radix tree] in [Go] (or Golang).  
Some of its features are based on [this awesome package].  

## Features
- No memory allocation for default search.
- Priority sort.
- Named parameter matching.

## Usage
Full documentation [here].  
[HEAD] holds the most recent features.

## Contribution
### How to help:
- Pull Requests
- Issues
- Opinions

[0.3.0]: https://github.com/gbrlsnchs/radix/tree/v0.3.0
[this package]: https://github.com/gbrlsnchs/patricia
[radix tree]: https://en.wikipedia.org/wiki/Radix_tree
[Go]: https://golang.org
[this awesome package]: https://github.com/julienschmidt/httprouter
[here]: https://godoc.org/github.com/gbrlsnchs/radix
[HEAD]: https://github.com/gbrlsnchs/radix/commit/HEAD
