package main

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"log"
	"math/rand"
	"os"
	"regexp"
	"strings"
	"time"
)

func crearDisco(size int64, path string, name string, unit byte) {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		err = os.MkdirAll(path, 0777)
		if err != nil {
			panic(err)
		}
	}

	temp := strings.Split(name, ".")
	nombre := temp[0]
	match, _ := regexp.MatchString("[A-Za-z][a-zA-Z0-9_]*", nombre)

	if match == true && temp[1] == "dsk" {
		if size > 0 {
			if unit == 0 || (unit == 'm' || unit == 'M') || (unit == 'k' || unit == 'K') {
				f, err := os.Create(path + name)
				//defer f.Close()
				if err != nil {
					panic(err)
				}

				var cero int8 = 0
				var binario bytes.Buffer
				var tam int64

				f.Seek(0, 0)
				err4 := binary.Write(&binario, binary.BigEndian, cero)
				if err4 != nil {
					fmt.Println("binary error ", err4)
				}
				escribirBytes(f, binario.Bytes())

				if unit == 'k' || unit == 'K' {
					tam = int64(1024 * size)
					f.Seek(int64(tam-1), 0)
					err2 := binary.Write(&binario, binary.BigEndian, cero)
					if err2 != nil {
						fmt.Println("binary error ", err2)
					}
					escribirBytes(f, binario.Bytes())

				} else if unit == 0 || (unit == 'm' || unit == 'M') {
					tam = int64(1024 * 1024 * size)
					f.Seek(int64(tam-1), 0)
					err3 := binary.Write(&binario, binary.BigEndian, cero)
					if err3 != nil {
						fmt.Println("binary error ", err3)
					}
					escribirBytes(f, binario.Bytes())
				}

				f.Seek(0, 0)

				r := rand.New(rand.NewSource(99))

				asignarMBR := mbr{Mbrtam: tam, Mbrdisksig: r.Int31()}
				fecha := time.Now().Format("2006-01-02 15:04:05")
				//fmt.Println(fecha)
				copy(asignarMBR.Mbrfecha[:], fecha)

				for i := 0; i < 4; i++ {
					asignarMBR.Prt[i].Partstatus = -1
				}

				s := &asignarMBR
				var binario2 bytes.Buffer
				err1 := binary.Write(&binario2, binary.BigEndian, s)
				if err1 != nil {
					fmt.Println("4- binary error ", err1)
				}
				escribirBytes(f, binario2.Bytes())
				f.Close()
			} else {
				fmt.Println("Error en el tipo de unidad de tamano de disco")
			}
		} else {
			fmt.Println("El tamano definido para el disco debe ser mayor a cero")
		}
	} else {
		fmt.Println("No se puede crear el disco, Extension invalida o Caracter invalido en el nombre")
	}

	obtenerMbr(path + name)

}

func escribirBytes(file *os.File, bytes []byte) {

	_, err := file.Write(bytes)

	if err != nil {
		log.Fatal(err)
	}

}
