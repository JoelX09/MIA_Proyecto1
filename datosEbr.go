package main

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"log"
	"os"
	"unsafe"
)

func obtenerEbr(path string, pos int64) ebr {
	file, err := os.Open(path)
	defer file.Close()
	if err != nil {
		log.Fatal(err)
	}

	e := ebr{}
	var size int = int(unsafe.Sizeof(e))

	file.Seek(pos, 0)
	data := obtenerBytesEbr(file, size)
	buffer := bytes.NewBuffer(data)

	//fmt.Println(data)

	err = binary.Read(buffer, binary.BigEndian, &e)
	if err != nil {
		log.Fatal("binary.Read failed", err)
	}

	fmt.Println(e)
	return e
}

func obtenerBytesEbr(file *os.File, number int) []byte {
	bytes := make([]byte, number)

	_, err := file.Read(bytes)
	if err != nil {
		log.Fatal(err)
	}

	return bytes
}

func escribirEbr(path string, e ebr, pos int64) {
	file, err := os.OpenFile(path, os.O_RDWR, 0777)
	defer file.Close()
	if err != nil {
		log.Fatal(err)
	}

	temp := e
	fmt.Println("para escribir:")
	fmt.Println(temp)
	file.Seek(pos, 0)
	s := &temp
	var binario2 bytes.Buffer
	err1 := binary.Write(&binario2, binary.BigEndian, s)
	if err1 != nil {
		fmt.Println("4- binary error ", err1)
	}
	escribirBytesEBR(file, binario2.Bytes())
}

func escribirBytesEBR(file *os.File, bytes []byte) {
	_, err := file.Write(bytes)
	if err != nil {
		log.Fatal(err)
	}
}
