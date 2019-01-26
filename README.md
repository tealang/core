# ![tea core](https://raw.githubusercontent.com/tealang/core/master/docs/logo.png)
[![Build Status](https://travis-ci.org/tealang/core.svg?branch=master)](https://travis-ci.org/tealang/core) [![Go Report Card](https://goreportcard.com/badge/github.com/tealang/core)](https://goreportcard.com/report/github.com/tealang/core)  [![codecov](https://codecov.io/gh/tealang/core/branch/master/graph/badge.svg)](https://codecov.io/gh/tealang/core)

Welcome to the repository of *core*, the core Tealang runtime. **It is NOT compatible to the Python implementation in some aspects due to further language changes. This project is heavily work-in-progress, do not use in any production environment.**

## How does FizzBuzz look?
```tea
operator /?(a, b: int): bool {
    return a % b == 0;
}

for var i = 0; i < 100; i = i + 1 {
    var a, b: string;
    if i /? 3 {
        a = "Fizz";
    }
    if i /? 5 {
        b = "Buzz";
    }
    let v = a + b;
    match v {
        case "" {
            print(i);
        }
        default {
            print(v);
        }
    }
}
```
