package main

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"log"
	"os"
	"strings"
)

func crearDisco(size int, path string, name string, unit byte) {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		err = os.MkdirAll(path, 0777)
		if err != nil {
			panic(err)
		}
	}

	temp := strings.Split(name, ".")
	if temp[1] == "dsk" {
		if size > 0 {
			f, err := os.Create(path + name)
			defer f.Close()
			if err != nil {
				panic(err)
			}
			var cero int8 = 0
			var binario bytes.Buffer

			f.Seek(0, 0)
			binary.Write(&binario, binary.BigEndian, cero)
			writeNextBytes(f, binario.Bytes())

			if unit == 'k' || unit == 'K' {
				tam := 1024 * size
				f.Seek(int64(tam-1), 0)
				binary.Write(&binario, binary.BigEndian, cero)
				writeNextBytes(f, binario.Bytes())

			} else if unit == 0 || (unit == 'm' || unit == 'M') {
				tam := 1024 * 1024 * size
				f.Seek(int64(tam-1), 0)
				binary.Write(&binario, binary.BigEndian, cero)
				writeNextBytes(f, binario.Bytes())

			} else {
				fmt.Println("Error en el tipo de unidad de tamano de disco")
			}

		} else {
			fmt.Println("El tamano definido para el disco debe ser mayor a cero")
		}
	} else {
		fmt.Println("No se puede crear el disco, Extension invalida.")
	}

}

func writeNextBytes(file *os.File, bytes []byte) {

	_, err := file.Write(bytes)

	if err != nil {
		log.Fatal(err)
	}

}
