package main

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"log"
	"os"
	"strings"
	"unsafe"
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
			//binario := bytes.NewBuffer([]byte{})
			var tam int64

			f.Seek(0, 0)
			err4 := binary.Write(&binario, binary.BigEndian, cero)
			if err4 != nil {
				fmt.Println("binary error ", err4)
			}
			writeNextBytes(f, binario.Bytes())

			if unit == 'k' || unit == 'K' {
				tam = int64(1024 * size)
				f.Seek(int64(tam-1), 0)
				err2 := binary.Write(&binario, binary.BigEndian, cero)
				if err2 != nil {
					fmt.Println("binary error ", err2)
				}
				writeNextBytes(f, binario.Bytes())

			} else if unit == 0 || (unit == 'm' || unit == 'M') {
				tam = int64(1024 * 1024 * size)
				f.Seek(int64(tam-1), 0)
				err3 := binary.Write(&binario, binary.BigEndian, cero)
				if err3 != nil {
					fmt.Println("binary error ", err3)
				}
				writeNextBytes(f, binario.Bytes())

			} else {
				fmt.Println("Error en el tipo de unidad de tamano de disco")
			}

			f.Seek(0, 0)

			asignarMBR := mbr{Mbrtam: tam}
			//t := time.Now()
			//t1 := string(t.Year()) + "-" + string(t.Month()) + "-" + string(t.Day()) + " " + string(t.Hour()) + ":" + string(t.Minute()) + " " + string(t.Second())
			//fmt.Println(t1)
			//copy(asignarMBR.mbrfecha[:], t1)

			s := &asignarMBR
			var binario2 bytes.Buffer
			err1 := binary.Write(&binario2, binary.BigEndian, s)
			if err1 != nil {
				fmt.Println("4- binary error ", err1)
			}
			writeNextBytes(f, binario2.Bytes())

			//fmt.Println(asignarMBR)

		} else {
			fmt.Println("El tamano definido para el disco debe ser mayor a cero")
		}
	} else {
		fmt.Println("No se puede crear el disco, Extension invalida.")
	}

	//leer(path + name)

}

func readFile(path string) {

	file, err := os.Open(path)
	defer file.Close()
	if err != nil {
		log.Fatal(err)
	}

	m := mbr{}
	var size int = int(unsafe.Sizeof(m))

	data := readNextBytes(file, size)
	buffer := bytes.NewBuffer(data)

	fmt.Println(data)

	err = binary.Read(buffer, binary.BigEndian, &m)
	if err != nil {
		log.Fatal("binary.Read failed", err)
	}

	fmt.Println(m)

	/*fmt.Printf("Caracter: %c\nCadena: %s\n", m.Caracter, m.Cadena)

	file.Seek(0, 0) // segundo parametro: 0, 1, 2.     0 -> Inicio, 1-> desde donde esta el puntero, 2 -> Del fin para atras
	file.Seek(int64(unsafe.Sizeof(m)), 0)
	//Struct 2
	fmt.Println("Struct 2: ")
	m2 := mbr2{}
	size = int(unsafe.Sizeof(m2))

	data = readNextBytes(file, size)
	buffer = bytes.NewBuffer(data)

	fmt.Println(data)

	err = binary.Read(buffer, binary.BigEndian, &m2)
	if err != nil {
		log.Fatal("binary.Read failed", err)
	}

	fmt.Println(m2)*/

}

func readNextBytes(file *os.File, number int) []byte {
	bytes := make([]byte, number)

	_, err := file.Read(bytes)
	if err != nil {
		log.Fatal(err)
	}

	return bytes
}

func writeNextBytes(file *os.File, bytes []byte) {

	_, err := file.Write(bytes)

	if err != nil {
		log.Fatal(err)
	}

}
