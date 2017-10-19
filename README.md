# tea-go [![Build Status](https://travis-ci.org/tealang/tea-go.svg?branch=master)](https://travis-ci.org/tealang/tea-go) [![Go Report Card](https://goreportcard.com/badge/github.com/tealang/tea-go)](https://goreportcard.com/report/github.com/tealang/tea-go)  [![codecov](https://codecov.io/gh/tealang/tea-go/branch/master/graph/badge.svg)](https://codecov.io/gh/tealang/tea-go)

Welcome to the repository of *tea-go*, the Go implementation of the Tealang runtime. **It is NOT compatible to the Python implementation in some aspects due to further language changes. This project is heavily work-in-progress, do not use in any production environment.**

## How does FizzBuzz look?
```tea
use io;

operator /?(a, b: int): bool {
    return a % b == 0;
}

for (var i = 0; i < 100; i++) {
    var a, b: string;
    if (i /? 3) {
        a = "Fizz";
    } else if (i /? 5) {
        b = "Buzz";
    }
    match (var v = a + b) {
    case "" => {
        io.println(i);
    }
    default => {
        io.println(v);
    }
    }
}
```
