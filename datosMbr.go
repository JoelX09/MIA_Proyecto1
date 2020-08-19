package main

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"log"
	"os"
	"unsafe"
)

func obtenerMbr(path string) mbr {
	file, err := os.Open(path)
	defer file.Close()
	if err != nil {
		log.Fatal(err)
	}

	m := mbr{}
	var size int = int(unsafe.Sizeof(m))

	data := obtenerBytes(file, size)
	buffer := bytes.NewBuffer(data)

	//fmt.Println(data)

	err = binary.Read(buffer, binary.BigEndian, &m)
	if err != nil {
		log.Fatal("binary.Read failed", err)
	}

	fmt.Println(m)
	return m
}

func obtenerBytes(file *os.File, number int) []byte {
	bytes := make([]byte, number)

	_, err := file.Read(bytes)
	if err != nil {
		log.Fatal(err)
	}

	return bytes
}

func escribirMbr(path string, m mbr) {
	file, err := os.OpenFile(path, os.O_RDWR, 0777)
	defer file.Close()
	if err != nil {
		log.Fatal(err)
	}

	temp := m
	fmt.Println("para escribir:")
	fmt.Println(temp)
	file.Seek(0, 0)
	s := &temp
	var binario2 bytes.Buffer
	err1 := binary.Write(&binario2, binary.BigEndian, s)
	if err1 != nil {
		fmt.Println("4- binary error ", err1)
	}
	escribirBytesMBR(file, binario2.Bytes())
}

func escribirBytesMBR(file *os.File, bytes []byte) {
	_, err := file.Write(bytes)
	if err != nil {
		log.Fatal(err)
	}
}
