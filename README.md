# recyclebin [![Build Status](https://travis-ci.org/trubitsyn/recyclebin.svg?branch=master)](https://travis-ci.org/trubitsyn/recyclebin) [![GoDoc](https://godoc.org/github.com/trubitsyn/recyclebin?status.svg)](https://godoc.org/github.com/trubitsyn/recyclebin)
Cross-platform way to use Trash from Go.

## Installation
`go get github.com/trubitsyn/recyclebin`

## Usage
```
package main

import (
	"github.com/trubitsyn/recyclebin"
	"fmt"
)

func main() {
	recyclebin.EmptyTrash()
	fmt.Println("Trash is empty now.")
}
```

## Testing
```
go get -t github.com/trubitsyn/recyclebin
go test github.com/trubitsyn/recyclebin
```

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
