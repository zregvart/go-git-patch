package main

import (
	"fmt"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing/format/diff"
	"github.com/go-git/go-git/v5/plumbing/object"
)

func result[T any](t T, err error) T {
	if err != nil {
		panic(err)
	}

	return t
}

func printFile(f *object.File) error {
	fmt.Printf("+ %s\n", f.Name)

	return nil
}

func printFiles(from, to diff.File) {
	// If the patch creates a new file, "from" will be nil. If the patch deletes a file, "to" will be nil.
	op := "-"
	f := from
	if f == nil {
		op = "+"
		f = to
	} else {
		op = "~"
	}
	fmt.Printf("%s %s\n", op, f.Path())
}

func main() {
	r := result(git.PlainOpen("."))

	commits := result(r.Log(&git.LogOptions{}))
	defer commits.Close()

	commits.ForEach(func(c *object.Commit) error {
		fmt.Println(c.ID())

		switch c.NumParents() {
		case 0:
			result(c.Files()).ForEach(printFile)
		case 1:
			to := result(c.Parent(0))
			p := result(c.Patch(to))
			for _, f := range p.FilePatches() {
				from, to := f.Files()
				printFiles(from, to)
			}
		}

		return nil
	})
}
