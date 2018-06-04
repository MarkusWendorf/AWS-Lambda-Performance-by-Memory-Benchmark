# AWS-Lambda-Performance-by-Memory-Benchmark

A performance benchmark to illustrate the relationship of available memory and CPU performance on AWS Lambda.

## Requirements
* Golang Runtime
* AWS CLI
* Serverless Framework

## How to

1. Run "go run build.go"
2. Run "sls deploy"
3. Run "go run runBenckmark.go"

The benchmark results are saved to "results.csv".
