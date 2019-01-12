# go-recyclebin [![Build Status](https://travis-ci.com/trubitsyn/go-recyclebin.svg?branch=master)](https://travis-ci.com/trubitsyn/go-recyclebin) [![GoDoc](https://godoc.org/github.com/trubitsyn/go-recyclebin?status.svg)](https://godoc.org/github.com/trubitsyn/go-recyclebin) [![Go Report Card](https://goreportcard.com/badge/github.com/trubitsyn/go-recyclebin)](https://goreportcard.com/report/github.com/trubitsyn/go-recyclebin)
Cross-platform way to use Trash or Recycle Bin from Go.

**Currently under development.**

## Installation
`go get github.com/trubitsyn/go-recyclebin`

## Usage
```
package main

import (
	"github.com/trubitsyn/go-recyclebin"
	"fmt"
)

func main() {
    bin, err := recyclebin.ForLocation("/home/user")
    if err != nil {
        fmt.Println(err)
    }
    if err := bin.Empty(); err != nil {
    	fmt.Println(err)
    }
    fmt.Println("Trash is empty now.")
}
```

## Testing
```
go get -t github.com/trubitsyn/go-recyclebin
go test github.com/trubitsyn/go-recyclebin
```

## References
* [The FreeDesktop.org Trash specification](https://standards.freedesktop.org/trash-spec/trashspec-1.0.html)

## LICENSE
```
Copyright 2018 Nikola Trubitsyn

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
```
