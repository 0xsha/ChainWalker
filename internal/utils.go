package internal

import (
	"os"
)

func WriteHexToFile(address string, bytecode string) {

	file, err := os.OpenFile("output/"+address+".evm", os.O_RDWR|os.O_CREATE, 0755)
	if err != nil {
		return
	}
	_, err = file.Write([]byte(bytecode))
	if err != nil {
		return
	}

}
