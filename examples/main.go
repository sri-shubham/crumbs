package main

import (
	"fmt"

	"github.com/sri-shubham/crumbs/examples/basic_example"
	"github.com/sri-shubham/crumbs/examples/context_example"
	"github.com/sri-shubham/crumbs/examples/logging_example"
	"github.com/sri-shubham/crumbs/examples/middleware_example"
	"github.com/sri-shubham/crumbs/examples/stack_example"
	"github.com/sri-shubham/crumbs/examples/std_errors_example"
)

func main() {
	fmt.Println("====================================")
	fmt.Println("CRUMBS LIBRARY EXAMPLES")
	fmt.Println("====================================")

	fmt.Println("\n\n====================================")
	fmt.Println("BASIC USAGE EXAMPLE")
	fmt.Println("====================================")
	basic_example.DemonstrateBasicUsage()

	fmt.Println("\n\n====================================")
	fmt.Println("CONTEXT AND CRUMBS EXAMPLE")
	fmt.Println("====================================")
	context_example.RunExample()

	fmt.Println("\n\n====================================")
	fmt.Println("STACK TRACES EXAMPLE")
	fmt.Println("====================================")
	stack_example.DemonstrateStackTraces()

	fmt.Println("\n\n====================================")
	fmt.Println("LOGGING INTEGRATION EXAMPLE")
	fmt.Println("====================================")
	logging_example.DemonstrateLoggingIntegration()

	fmt.Println("\n\n====================================")
	fmt.Println("STANDARD ERRORS INTEGRATION EXAMPLE")
	fmt.Println("====================================")
	std_errors_example.DemonstrateStandardErrorsMethods()

	fmt.Println("\n\n====================================")
	fmt.Println("MIDDLEWARE INTEGRATION EXAMPLE")
	fmt.Println("====================================")
	middleware_example.RunExample()

	fmt.Println("\n\nAll examples completed!")
}
