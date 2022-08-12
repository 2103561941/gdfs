// create client command markdown docfiles automatically.

package main

import (
	"log"

	"github.com/cyb0225/gdfs/internal/client/cmd"
)

func main() {
	// relative file path, refer to the loaction of excutable file.
	filepath := "/docs/api/cmd/"
	log.Println(cmd.GenDocs(filepath))
}
