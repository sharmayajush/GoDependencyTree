package main

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/sharmayajush/lineajeassignment/dependency"
	"github.com/sharmayajush/lineajeassignment/repository"
)

func main() {

	if len(os.Args) < 3 {
		fmt.Println("Usage: go run main.go <githubURL> <tag>")
		return
	}

	githubURL := os.Args[1]
	tag := os.Args[2]

	tmpDir := "./tmp/repo-clone"
	defer os.RemoveAll(tmpDir)

	repo, err := repository.CloneRepo(githubURL, tmpDir)
	if err != nil {
		fmt.Println("unable to clone repository")
		return
	}

	if err = repository.CheckoutToTag(repo, tag); err != nil {
		fmt.Println("unable to checkout to tag ", tag)
		return
	}

	// get go mod graph to identify all the dependencies
	lines, err := dependency.GetGoModGraph(tmpDir)
	if err != nil {
		fmt.Println("error in executing go mod graph for the given directory")
		return
	}

	ans := dependency.MakeDependencyGraph(&lines)

	fmt.Println("marshalling and printing result to output.json...")

	result, err := json.MarshalIndent(*ans, "", "	")
	if err != nil {
		fmt.Println("unable to marshal the result", err)
		return
	}

	// Create a new JSON file
	file, err := os.Create("output.json")
	if err != nil {
		fmt.Println("Unable to create file:", err)
		return
	}
	defer file.Close() // Ensure the file is closed after writing

	// Write the marshalled JSON data to the file
	_, err = file.Write(result)
	if err != nil {
		fmt.Println("Unable to write to file:", err)
		return
	}

	fmt.Println("JSON data successfully written to output.json")

}
