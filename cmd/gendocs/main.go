// create client command markdown docfiles automatically.  

package main

import (
	"gdfs/internal/client/cmd"
	"log"
)

func main() {
	// relative file path, refer to the loaction of excutable file.
	filepath := "/docs/api/cmd/"
	log.Println(cmd.GenDocs(filepath))
}