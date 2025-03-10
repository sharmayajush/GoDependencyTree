package main

import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"slices"
	"strings"

	"github.com/go-git/go-git/v5"
)

type Artifact struct {
	Name         string
	Version      string
	Dependencies []*Artifact
}


func main() {
	githubURL := "https://github.com/etcd-io/etcd"
	tag := "v3.5.5"
	// cloneDir := "./etcd-repo"
	artifactBuffer := make(map[string]*Artifact)

	tmpDir := "/tmp/repo-clone"
	defer os.RemoveAll(tmpDir)

	repo, err := git.PlainClone(tmpDir, false, &git.CloneOptions{
		URL:      githubURL,
		Progress: os.Stdout,
	})
	if err != nil {
		fmt.Println("unable to clone repository")
		return
	}
	fmt.Println("cloned repository")

	tagRef, err := repo.Tag(tag)
	if err != nil {
		return
	}

	worktree, err := repo.Worktree()
	if err!=nil{
		fmt.Println("")
	}

	worktree.Checkout(&git.CheckoutOptions{
		Hash: tagRef.Hash(),
	})

	cmd := exec.Command("go", "mod", "graph")
	cmd.Dir = tmpDir

	output, err := cmd.CombinedOutput()
	if err != nil {
		fmt.Printf("Failed to run `go mod graph`: %s\n", err)
		return
	}

	lines := strings.Split(string(output), "\n")

	var ans *Artifact

	for i := len(lines) - 1; i >= 0; i-- {
		line := lines[i]
		fmt.Println(line)

		if line == "" {
			continue
		}

		relations := strings.Split(line, " ")
		if len(relations) < 2 {
			fmt.Println(relations)
			continue
		}
		if _, ok := artifactBuffer[relations[1]]; !ok {
			art := strings.Split(relations[1], "@")
			if len(art) < 2 {
				art = append(art, "v0")
			}
			artifactBuffer[relations[1]] = &Artifact{
				Name:         art[0],
				Version:      art[1],
				Dependencies: nil,
			}
		}

		if _, ok := artifactBuffer[relations[0]]; !ok {
			art := strings.Split(relations[0], "@")
			if len(art) < 2 {
				art = append(art, "v0")
			}
			artifactBuffer[relations[0]] = &Artifact{
				Name:         art[0],
				Version:      art[1],
				Dependencies: []*Artifact{},
			}
		}
		if slices.Contains(artifactBuffer[relations[1]].Dependencies, artifactBuffer[relations[0]]) {
			continue
		}

		artifactBuffer[relations[0]].Dependencies = append(artifactBuffer[relations[0]].Dependencies, artifactBuffer[relations[1]])
		ans = artifactBuffer[relations[0]]

	}

	// Print the output of `go mod graph`
	fmt.Println("Output of `go mod graph`:")
	fmt.Println(ans)
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

	// fmt.Println(string(output))
}
