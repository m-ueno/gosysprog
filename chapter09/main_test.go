package main_test

import (
	"fmt"
	"os"
	"path/filepath"
)

// 9.2.11
func ExampleOpenDir() {

	dir, _ := os.Open("/")
	limit := -1
	fileInfos, _ := dir.Readdir(limit)
	for _, fileInfo := range fileInfos {
		if fileInfo.IsDir() {
			fmt.Printf("[Dir]\t%s/\n", fileInfo.Name())
		} else {
			fmt.Printf("[File]\t%s\n", fileInfo.Name())
		}
	}

	//Output:

}

func ExampleWalk() {
	root := ".."
	walkFunc := func(path string, info os.FileInfo, err error) error {
		fmt.Println(path)
		fmt.Println(info.Size())
		return nil
	}

	err := filepath.Walk(root, walkFunc)
	if err != nil {
		panic(err)
	}

	//Output:

}

func ExampleRename() {
	os.Mkdir("olddir")
	os.Mkdir("newdir")
	os.Create("olddir/hoge")
	os.Rename("olddir/hoge", "newdir")

	//Output:
}
