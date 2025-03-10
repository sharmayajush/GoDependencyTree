package dependency

import (
	"fmt"
	"os/exec"
	"slices"
	"strings"
)

type Artifact struct {
	Name         string
	Version      string
	Dependencies []*Artifact
}

func GetGoModGraph(tmpDir string) ([]string, error) {
	cmd := exec.Command("go", "mod", "graph")
	cmd.Dir = tmpDir

	output, err := cmd.CombinedOutput()
	if err != nil {
		fmt.Printf("Failed to run `go mod graph`: %s\n", err)
		return []string{}, err
	}

	lines := strings.Split(string(output), "\n")
	return lines, nil
}

func MakeDependencyGraph(lines *[]string) *Artifact {
	var ans *Artifact
	artifactBuffer := make(map[string]*Artifact)

	for i := len(*lines) - 1; i >= 0; i-- {
		line := (*lines)[i]

		if line == "" {
			continue
		}

		relations := strings.Split(line, " ")
		if len(relations) < 2 {
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
	return ans
}
