<h1 align="center">Hopper</h1>

```
                        ██╗  ██╗ ██████╗ ██████╗ ██████╗ ███████╗██████╗ 
                        ██║  ██║██╔═══██╗██╔══██╗██╔══██╗██╔════╝██╔══██╗
                        ███████║██║   ██║██████╔╝██████╔╝█████╗  ██████╔╝
                        ██╔══██║██║   ██║██╔═══╝ ██╔═══╝ ██╔══╝  ██╔══██╗
                        ██║  ██║╚██████╔╝██║     ██║     ███████╗██║  ██║
                        ╚═╝  ╚═╝ ╚═════╝ ╚═╝     ╚═╝     ╚══════╝╚═╝  ╚═╝
```
<p align="center">
Modular Highly Parallel Fuzzer
</p>

## **Inspiration**

Hopper is a composable and distributed Fuzzer written in Golang. It is a Coverage-Guided
Greybox Fuzzer inspired by AFL++, the current state of the art fuzzer. I use a mutation based input generation
and a Priority Queue circular buffer for seed selection with an energy strategy similar to the one implemented
in AFLFast. For binary instrumentation, a strategy similar to AFL++’s instrumentation with afl-clang-fast is
used. Hopper uses LLVM’s SanitizerCoverage to gather coverage, along with LLVM’s built-in Address Sanitizer
to detect crashes. Finally, the RPC communication schema and task scheduling is heavily inspired by Google’s
Map-Reduce

## Usage

Work in progress..
