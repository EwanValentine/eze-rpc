package main

import (
	"fmt"
	"github.com/EwanValentine/eze-rpc/generator"
	"log"
	"os"
)

func main() {
	content, err := os.ReadFile("./example/example.eze")
	if err != nil {
		log.Panic(err)
	}

	dsl := generator.ParseDSL(string(content))

	output := generator.GenerateCode(dsl)

	log.Println(output)

	log.Println(dsl.Package)

	os.Mkdir(dsl.Package, 0755)

	if err := os.WriteFile(fmt.Sprintf("%s/%s", dsl.Package, "eze_output.go"), []byte(output), 0644); err != nil {
		log.Panic(err)
	}
}
