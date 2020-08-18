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
		estado     int
		Partstatus byte
		Parttype   byte
		Partfit    byte
		Partstart  int64
		Partsize   int64
		Partname   [16]byte
	}

	primaria, extendida := 0, 0
	m := obtenerMbr(fd.path)
	var size int = int(unsafe.Sizeof(m))
	fmt.Println(size)

	for i := 0; i < 4; i++ {
		part := m.Prt[i]
		if part.Partsize > 0 {
			var datosPart nodoPart
			datosPart.estado = 1
			datosPart.Partstatus = part.Partstatus
			datosPart.Parttype = part.Parttype
			datosPart.Partfit = part.Partfit
			datosPart.Partstart = part.Partstart
			datosPart.Partsize = part.Partsize
			datosPart.Partname = part.Partname

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
			var temp nodoPart = ele.Value.(nodoPart)
			if temp.Partstart == int64(size) {
			}
		}
	}

}
