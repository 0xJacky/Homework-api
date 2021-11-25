package pkg

import (
	"log"
	"os"
)

func ExistsOrCreate(path string) {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		err = os.Mkdir(path, 0755)
		if err != nil {
			log.Fatal("[ExistsOrCreate] fail to create", path)
		}
	}
}
