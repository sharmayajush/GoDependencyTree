package repository

import (
	"fmt"
	"os"

	"github.com/go-git/go-git/v5"
)

func CloneRepo(githubURL, tmpDir string) (*git.Repository, error) {
	fmt.Println("cloning repository...")
	repo, err := git.PlainClone(tmpDir, false, &git.CloneOptions{
		URL:      githubURL,
		Progress: os.Stdout,
	})
	if err != nil {
		fmt.Println("err:", err.Error())
		return nil, err
	}
	fmt.Println("successfully cloned repository...")
	return repo, nil
}

func CheckoutToTag(repo *git.Repository, tag string) error {
	fmt.Printf("checking out to tag %s...\n", tag)
	tagRef, err := repo.Tag(tag)
	if err != nil {
		return err
	}

	worktree, err := repo.Worktree()
	if err != nil {
		fmt.Println("error in getting worktree for the repo")
		return err
	}

	err = worktree.Checkout(&git.CheckoutOptions{
		Hash: tagRef.Hash(),
	})
	if err != nil {
		fmt.Println("error in checkout to tag")
		return err
	}
	fmt.Printf("successfully checked out to tag %s...\n", tag)
	return nil
}
