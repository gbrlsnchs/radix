# patricia (PATRICIA Tree)

[![Build Status](https://travis-ci.org/gbrlsnchs/patricia.svg?branch=master)](https://travis-ci.org/gbrlsnchs/patricia)
[![GoDoc](https://godoc.org/github.com/gbrlsnchs/patricia?status.svg)](https://godoc.org/github.com/gbrlsnchs/patricia)

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
Until version 1.0 is released, **anything** can change, including names of methods or even their existence.

## About
This package is an implementation of a PATRICIA tree in [Go] (or Golang).  
Some of its features are based on [this package].  

## Usage
Full documentation [here].  
[HEAD] holds the most recent features.

## Contribution
### How to help:
- Pull Requests
- Issues
- Opinions

[Go]: https://golang.org
[this package]: https://github.com/julienschmidt/httprouter
[this example]: https://upload.wikimedia.org/wikipedia/commons/a/ae/Patricia_trie.svg
[here]: https://godoc.org/github.com/gbrlsnchs/patricia
[HEAD]: https://github.com/gbrlsnchs/patricia/commit/HEAD
