<h1 align="center">Hopper</h1>

<div align="center">
<h3>
Coverage-Guided Greybox Distributed Fuzzer inspired
by <a href="https://github.com/AFLplusplus/AFLplusplus">AFL++</a>
</h3>

<h4> Hopper aims to improve performance of Fuzzing in large-scale
distributed environments, it's not meant to replace AFL++ in most cases.
</h4>

<img src="master.png" align="center" alt="Runemaster Icon"/><br>

*Hopper Master*

</div>

## Usage

#### Pre-Reqs:

- [LLVM](https://clang.llvm.org/) toolchain, specifically
  [clang](https://clang.llvm.org/get_started.html) with it's built-in
  [ASAN](https://clang.llvm.org/docs/AddressSanitizer.html)
- [clang-tools](https://clang.llvm.org/docs/ClangTools.html), specifically the
  [SanitizerCoverage](https://clang.llvm.org/docs/SanitizerCoverage.html)
  sancov utility.

#### Instrumentation:

- The [compile](test/compile) script adds all the flags required to compiler
  the target program with clang++.

Ex:
> `./compile target.c`

#### Master:

- <kbd>-I</kbd>: Path to input corpus, directory containing files each being a
  seed
- <kbd>-H</kbd>: Havoc level to use in mutator, defaults to `1` (recommended:
  1-10, for builtin mutator)
- <kbd>-P</kbd>: Port to host Master on, defaults to `:6969`

Ex:
> `go build . && ./hopper -H 5 -I test/in`

#### Node:

- <kbd>-I</kbd>: Node ID, usually just a unique int
- <kbd>-T</kbd>: Path to instrumented target binary
- <kbd>-M</kbd>: IP/address of Master, defaults to `localhost`
- <kbd>-P</kbd>: Port of Master, defaults to `:6969`
- <kbd>--args</kbd>: Args to use against target, ex: `--depth=1 @@`
- <kbd>--env</kbd>: Env variables for target seperated by a `;`, ex:
  `ENV1=foo;ENV2=bar;`
- <kbd>--stdin</kbd>: Should seed be fed as stdin or as an argument, defaults
  to `false`

Ex: 
> Args: `go run node/node.go -I 1 -T test/target --args "--depth=2 @@"` 
>
> Stdin: `go run node/node.go -I 1 -T test/target --stdin`

#### Example:

If you want to run Hopper locally with 10 fuzzing Nodes on a test application
with a known vulnerability you can do the following:

1. Clone project: `git clone https://github.com/Cybergenik/hopper.git`
2. Compile target: `cd hopper/test && ./compile getdomain.c`
3. Run Master: `./run_master.sh`
4. Run Nodes: `./run_node.sh 10` (I'd recommend no more than 1.5x # of logical cores on your machine, any more
nodes on one system and they just get throttled and competing for CPU time)
5. Look at the nice TUI :>

## **Inspiration**

I'll be graduating May 2023, and I had to choose whether to do a Capstone
Project or do a Thesis, I chose Thesis as I felt that I could build something
pretty cool related to Fuzzers. I have been researching Fuzzers for a few years
at the [FLUX](https://www.flux.utah.edu/) research group, from chasing WASM
fuzzing and other compiler based fuzzing approaches, to bootstrapping GDB with
python to get better crash reports on crashes. But I wasn't sure exactly what I
should write my Thesis about. Until, a serendipitously timed assignment to
implement Map-Reduce in my Distributed Systems course taught by Professor [Ryan
Stutsman](https://rstutsman.github.io/), this almost felt like divine
premonition. I began looking at current approaches into the space of parallel
distributed fuzzers, I found that very little had been done in this area. With
this in mind, I set out to build Hopper! Although Google's Map-Reduce was the initial
inspiration and definitely inspired the way Hopper manages the Task Scheduling
and assignment strategy, Hopper’s core infrastructure is substantially
different.

## Architecture

<div align="center"><img src="arch.png" align="center" alt="Runemaster
Icon"/></div><br>

