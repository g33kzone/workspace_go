package main

import (
	"fmt"
	"strings"
)

func main() {
	repoURL := "https://g33kzone@bitbucket.org/g33kzone/sb-hello-world.git"

	repoArray := strings.Split(repoURL, "/")

	fmt.Println(repoArray[4])

}
