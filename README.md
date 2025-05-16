# LineajeAssignment
This project is a Go-based tool that clones a Git repository, checks out a specific tag, and generates a dependency graph using the go mod graph command. The dependency graph is then saved as a JSON file.

## Prerequisites
Go 1.20 or higher

Git installed on your system

Access to a Git repository with a valid Go module

## Installation

```bash
git clone https://github.com/sharmayajush/GoDependencyTree.git
cd GoDependencyTree
```

## Command
```bash
go run main.go githubURL tag
```
### Arguments
```bash
<githubURL>: The URL of the Git repository (e.g., https://github.com/etcd-io/etcd).
<tag>: The tag or version to checkout (e.g., v3.5.5).
```
### Example
```bash
go run main.go https://github.com/etcd-io/etcd v3.5.5
```
## Output

The dependency graph will be saved in a file named output.json in the project root directory.

