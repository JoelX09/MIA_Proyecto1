package main

import (
	"container/list"
	"fmt"
	"unsafe"
)

var existePart bool

func adminParticion(fd datoDisco) {
	listaP := list.New()

	type nodoPart struct {
		estado int8
		num    int
		ini    int64
		tam    int64
		nomb   [16]byte
	}

	var primaria, extendida int
	m := mbr{}
	var size int = int(unsafe.Sizeof(m))
	fmt.Println(size)

	for i := 0; i < 4; i++ {
		part := m.Prt[i]
		if part.Partsize > 0 {
			var datosPart nodoPart
			datosPart.estado = 1
			datosPart.num = i
			datosPart.ini = part.Partstart
			datosPart.tam = part.Partsize
			datosPart.nomb = part.Partname

			if part.Parttype == 'P' {
				primaria++
			}
			if part.Parttype == 'E' {
				extendida++
			}

			listaP.PushBack(datosPart)
		}
	}

	if listaP.Len() > 0 {
		existePart = true
		for ele := listaP.Front(); ele != nil; ele = ele.Next() {

		}
	}

}
