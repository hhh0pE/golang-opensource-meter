package main

import (
	"errors"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

var GOPATH string
var GOROOT string

func main() {
	if len(os.Args) == 1 {
		panic("not found package arg")
	}

	var packageSrc string
	if param_parts := strings.Split(os.Args[1], "="); len(param_parts) == 2 {
		if param_parts[0] == "-package" {
			packageSrc = param_parts[1]
		}
	}

	if packageSrc == "" {
		panic("cannot get package src from args")
	}

	fmt.Println(packageSrc)
	GOPATH = os.Getenv("GOPATH")
	GOROOT = os.Getenv("GOROOT")

	//goExePath := goroot+string(os.PathSeparator)+"bin"+string(os.PathSeparator)+"go.exe"
	os.Chdir(fullPathToPackage(packageSrc))
	goListCmd := exec.Command("cmd", "/C", "go", "list", "-f", `"{{join .Deps "\n"}}"`)
	listBytes, cmd_err := goListCmd.CombinedOutput()
	fmt.Println(string(listBytes))
	if cmd_err != nil {
		panic(cmd_err)
	}

	var pkgSizeMap = make(map[string]int64)

	deps := strings.Split(string(listBytes), "\n")
	var openSourceSize, stdLibSize, restSize int64
	for _, depPkg := range deps {
		depPkg = strings.Trim(depPkg, `"`)
		if depPkg == "" {
			continue
		}
		if pkgSize, err := pkgSize(depPkg); err != nil {
			log.Println("Error when getting size of package: ", err.Error())
			continue
		} else {
			pkgSizeMap[depPkg] = pkgSize
		}

		if isLibPkg(depPkg) {
			stdLibSize += pkgSizeMap[depPkg]
		} else if isPkgOpensource(depPkg) {
			openSourceSize += pkgSizeMap[depPkg]
		} else {
			restSize += pkgSizeMap[depPkg]
		}

		//fmt.Println(depPkg, isPkgOpensource(depPkg), pkgSizeMap[depPkg])
	}

	var totalSize = openSourceSize + restSize + stdLibSize

	//fmt.Println(openSourceSize, restSize)
	fmt.Println("StdLib %", (float64(stdLibSize)/float64(totalSize))*100)
	fmt.Println("OpenSource %", (float64(openSourceSize)/float64(totalSize))*100)
	fmt.Println("Own code %", (float64(restSize)/float64(totalSize))*100)

	//fmt.Println(pkgSizeMap)

	//import

}

func fullPathToPackage(packageName string) string {
	return GOPATH + string(os.PathSeparator) + "src" + string(os.PathSeparator) + packageName
}

func pkgSize(pkgSrc string) (int64, error) {
	//var dirFileInfo os.FileInfo
	var dirFullSrc string
	if _, err := os.Stat(srcPath(GOPATH) + pkgSrc); os.IsNotExist(err) {
		if _, err2 := os.Stat(srcPath(GOROOT) + pkgSrc); os.IsNotExist(err2) {
			return 0, errors.New("Cannot find path's directory " + pkgSrc)
		} else {
			//dirFileInfo = dirStat2
			dirFullSrc = srcPath(GOROOT) + pkgSrc
		}
	} else {
		//dirFileInfo = dirStat
		dirFullSrc = srcPath(GOPATH) + pkgSrc
	}

	//fmt.Println("dirFullSrc", dirFullSrc)

	var dirSize int64
	filepath.Walk(dirFullSrc, func(path string, info os.FileInfo, err error) error {
		if !info.IsDir() {
			dirSize += info.Size()
		}
		return nil
	})
	return dirSize, nil
}

func srcPath(goBasePath string) string {
	return goBasePath + string(os.PathSeparator) + "src" + string(os.PathSeparator)
}

func isPkgOpensource(pkg string) bool {
	var opensourceSites = []string{"github.com", "google.golang.org", "golang.org", "gopkg.in"}
	if isLibPkg(pkg) {
		return true
	}

	for _, oss := range opensourceSites {
		if strings.HasPrefix(pkg, oss) {
			return true
		}
	}
	return false
}

func isLibPkg(pkg string) bool {
	if parts := strings.Split(pkg, "/"); len(parts) > 1 {
		if foundDotIndex := strings.Index(parts[0], "."); foundDotIndex == -1 {
			return true
		}
	} else {
		return true
	}
	return false
}
