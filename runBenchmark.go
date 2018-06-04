package main

import (
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/lambda"
	"strings"
	"sort"
	"os"
	"fmt"
)

var lambdaClient *lambda.Lambda

func main() {

	resultFile, err := os.OpenFile("results.csv", os.O_CREATE|os.O_APPEND, 777)
	if err != nil {
		panic(err)
	}

	sess := session.Must(session.NewSession(&aws.Config{
		Region: aws.String("eu-central-1")},
	))

	lambdaClient = lambda.New(sess)
	allFunctions := make([]*lambda.FunctionConfiguration, 0)

	lambdaClient.ListFunctionsPages(&lambda.ListFunctionsInput{}, func(out *lambda.ListFunctionsOutput, isLast bool) bool {

		for _, f := range out.Functions {
			allFunctions = append(allFunctions, f)
		}

		return true
	})

	functionsToBench := make([]*lambda.FunctionConfiguration, 0)
	for _, val := range allFunctions {
		if strings.HasPrefix(*val.FunctionName, "LambdaPerformance") {
			functionsToBench = append(functionsToBench, val)
		}
	}

	if len(functionsToBench) == 0 {
		panic("found 0 lambda functions to benchmark")
	}

	sort.Slice(functionsToBench, func(i, k int) bool {
		return *functionsToBench[i].FunctionName < *functionsToBench[k].FunctionName
	})

	for _, function := range functionsToBench {
		name := *function.FunctionName
		last := strings.LastIndex(name, "dev-")
		resultFile.Write([]byte(name[last+4:] + ";"))
	}

	resultFile.WriteString("\r\n")

	fmt.Println("AWS Lambda Benchmark")
	for i := 0; i < 1; i++ {
		runBenchRound(resultFile, functionsToBench)
	}

}

func runBenchRound(resultFile *os.File, functions []*lambda.FunctionConfiguration) {

	for _, function := range functions {

		req, out := lambdaClient.InvokeRequest(&lambda.InvokeInput{
			FunctionName: function.FunctionName,
		})

		err := req.Send()

		if err != nil {
			panic(err)
		}

		executionTime := out.Payload
		fmt.Println(*function.FunctionName + " : " + strings.Trim(string(executionTime), "\"") + "ms")

		resultFile.Write(out.Payload)
		resultFile.WriteString(";")
	}

	resultFile.WriteString("\r\n")
}
