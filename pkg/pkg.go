package pkg

import (
	"log"
	"os"
	"strconv"
)

func ExistsOrCreate(path string) {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		err = os.Mkdir(path, 0755)
		if err != nil {
			log.Fatal("[ExistsOrCreate] fail to create", path)
		}
	}
}

func StrToUInt(str string) uint {
	i, e := strconv.Atoi(str)
	if e != nil {
		return 0
	}
	return uint(i)
}