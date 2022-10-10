package main

import (
	"errors"
	"os"
	"ply2octree/pkg"
)

func main() {
	// ToDo - Add cli lib
	pkg.PrintInfo("Ply file to octree point cloud chunk converter")
	args := os.Args
	if len(args) < 3 {
		pkg.PrintError(errors.New("arguments not valid"))
	} else {
		converter := pkg.NewConverter(args[1], args[2], 5)
		err := converter.Convert()
		if err != nil {
			pkg.PrintError(err)
		}
	}
}
