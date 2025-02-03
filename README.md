<div align="center" style="display:grid;place-items:center;">
<p>
    <img width="100" src="./shark_file_icon.svg" alt="Statusify Logo">
</p>
<h1>The Shark Programming Language</h1>

<h4>Shark is a programming language with a language server, compiler and virtual machine</h4>
</div>

## About

Shark is written in Go aiming to be a simple, dynamically typed language with a focus on simplicity and ease of use. The language is compiled to bytecode that a virtual machine can run. It is inspired by languages like TypeScript and Dart.

## Key Features of SharkLang

- Dynamically typed
- Optional type annotations
- Compiles to bytecode that a virtual machine can run
- Garbage collected
- VS Code extension for syntax highlighting
- Language server
- Caching

> [!NOTE]
> Shark is currently in development and there is no release yet.

## Example

```shark
let reduce = (arr: array<i64>, initial: i64, f: func<(i64, i64)->i64>) => {
    let iter = (arr: array<i64>, result: i64) => {
        if (len(arr) == 0) {
            result
        } else {
            iter( rest(arr), f(result, first(arr)) ) 
        }
    } iter(arr, initial)
};

let sum = (arr: array<i64>) => {
    reduce(arr, 0, (initial: i64, el: i64): i64 => {
        initial + el
    })
};

let result = sum(1..5);

puts(result);
```

For more examples, look at the [examples](./examples) directory.

## Benchmarks

The following Shark code is executed:

```shark
let fibonacci = (x) => {
    if (x == 0) {
        return 0;
    } else {
        if (x == 1) {
            return 1;
        } else {
            fibonacci(x - 1) + fibonacci(x - 2);
        }
    }
};
fibonacci(50);
```

### Results

The results were recorded at `February 3rd, 2025`. The results include VM caching enabled with the cache size of `1024`.

#### Time

```
BenchmarkRecursiveFibonacci-10    	1000000000	         0.0003148 ns/op	       0 B/op	       0 allocs/op
```

#### Memory Usage

```
Showing nodes accounting for 5034.07kB, 100% of 5034.07kB total
Showing top 10 nodes out of 16
      flat  flat%   sum%        cum   cum%
 1762.94kB 35.02% 35.02%  1762.94kB 35.02%  runtime/trace.Start
 1184.27kB 23.53% 58.55%  1184.27kB 23.53%  runtime/pprof.StartCPUProfile
 1184.27kB 23.53% 82.07%  1184.27kB 23.53%  shark/vm.New
  902.59kB 17.93%   100%   902.59kB 17.93%  compress/flate.NewWriter (inline)
         0     0%   100%   902.59kB 17.93%  compress/gzip.(*Writer).Write
         0     0%   100%  2947.21kB 58.55%  main.main
         0     0%   100%  2947.21kB 58.55%  runtime.main
         0     0%   100%   902.59kB 17.93%  runtime/pprof.(*profileBuilder).build
         0     0%   100%   902.59kB 17.93%  runtime/pprof.profileWriter
         0     0%   100%  1184.27kB 23.53%  shark/vm.NewDefault
```
