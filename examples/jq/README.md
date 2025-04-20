# Fuzzing jq JSON Processor with Hopper

This example demonstrates how to fuzz the popular [jq](https://stedolan.github.io/jq/) JSON processor using Hopper.

## Steps to Run

1. Clone hopper: 
```bash
git clone https://github.com/Cybergenik/hopper.git && cd hopper
```

2. Build base Hopper image: 
```bash
docker build -t hopper-node .
```
3. Build instrumented jq fuzzing image:

```bash
cd examples/jq
docker build -t hopper-jq .
```
4. Prepare a corpus directory with JSON seeds (already included under `examples/jq/corpus`).

5. Run Hopper master, pointing to corpus, with havoc level 2 (example):
```bash
./master_docker.sh
```

6. In a seperate shell, run one or more (ex: 10) Hopper fuzz nodes with the jq binary:

```bash
./node_docker.sh 1 10
```

*Uses --stdin to feed seeds through jq's stdin.*

To see how the nodes and master are being run, look in the scripts
`master_docker.sh` and `node_docker.sh`. Try changing the havoc level on the
master, see how that affects fuzzing.

## Observe

1. Monitor the TUI from the Hopper master.

2. Review crashes saved by Hopper under the `HOPPER_OUT` directory.
