package main

import (
	"container/list"
	"fmt"
	"strconv"
	"unsafe"
)

func adminParticion(fd datoDisco, fl banderaParam) {
	listaP := list.New()
	existePart := false

	type nodoPart struct {
		Estado     int
		Partstatus int8
		Parttype   byte
		Partfit    [2]byte
		Partstart  int64
		Partsize   int64
		Partname   [16]byte
	}

	primaria, extendida := 0, 0
	m := obtenerMbr(fd.path)

	tam := m.Mbrtam
	fmt.Println("-------------------------\n\n" + strconv.FormatInt(tam, 10) + "\n\n-------------------------")

	var size int = int(unsafe.Sizeof(m))
	fmt.Println(size)

	var datosPart nodoPart

	for i := 0; i < 4; i++ {
		part := m.Prt[i]
		if part.Partstatus != -1 {

			datosPart.Estado = 1
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
			var tempVacio nodoPart
			if ele.Prev() == nil {
				if temp.Partstart != int64(size) {
					tempVacio.Estado = 0
					tempVacio.Partstart = int64(size)
					tempVacio.Partsize = temp.Partstart - int64(size)

					listaP.InsertBefore(tempVacio, ele)
				}
				if ele.Next() == nil {
					tempVacio.Estado = 0
					tempVacio.Partstart = temp.Partstart + temp.Partsize
					tempVacio.Partsize = m.Mbrtam - tempVacio.Partstart
					listaP.InsertAfter(tempVacio, ele)
				}
			} else if ele.Next() == nil {
				var tempAnt nodoPart = ele.Prev().Value.(nodoPart)
				if temp.Partstart != (tempAnt.Partstart + tempAnt.Partsize) {
					tempVacio.Estado = 0
					tempVacio.Partstart = tempAnt.Partstart + tempAnt.Partsize
					tempVacio.Partsize = temp.Partstart - tempVacio.Partstart
					listaP.InsertBefore(tempVacio, ele)
				}
				if (temp.Partstart + temp.Partsize - 1) != m.Mbrtam-1 {
					tempVacio.Estado = 0
					tempVacio.Partstart = temp.Partstart + temp.Partsize
					tempVacio.Partsize = m.Mbrtam - tempVacio.Partstart

					listaP.InsertAfter(tempVacio, ele)
				}
			} else {
				var tempAnt nodoPart = ele.Prev().Value.(nodoPart)
				if temp.Partstart != (tempAnt.Partstart + tempAnt.Partsize) {
					tempVacio.Estado = 0
					tempVacio.Partstart = tempAnt.Partstart + tempAnt.Partsize
					tempVacio.Partsize = temp.Partstart - tempVacio.Partstart
					listaP.InsertBefore(tempVacio, ele)
				}
			}

		}
	} //-------------------------------------------

	fmt.Println("Contenido despues de generar los nodos ocupados y disponibles:")
	for element := listaP.Front(); element != nil; element = element.Next() {
		// do something with element.Value
		fmt.Println(element.Value)
	}

	//-----------------

	if fl.deleteY == false && fl.addY == false {

		var tam int64
		if fd.unit == 'K' {
			tam = fd.size * 1024
		}

		if existePart == false {
			datosPart.Estado = 1
			datosPart.Partstatus = 0
			datosPart.Parttype = fd.typeP
			copy(datosPart.Partfit[:], fd.fit)
			datosPart.Partstart = int64(size)
			datosPart.Partsize = tam
			copy(datosPart.Partname[:], fd.name)
			listaP.PushFront(datosPart)
		} else {
			done := false
			for ele := listaP.Front(); ele != nil; ele = ele.Next() {
				temp := ele.Value.(nodoPart)
				if temp.Estado == 0 {
					if temp.Partsize >= fd.size {
						temp2 := ele.Prev()

						temp.Estado = 1
						temp.Partstatus = 0
						temp.Parttype = fd.typeP
						copy(temp.Partfit[:], fd.fit)
						temp.Partsize = tam
						copy(datosPart.Partname[:], fd.name)
						listaP.Remove(ele)
						listaP.InsertAfter(temp, temp2)
						done = true
						break
					}
				}
			}
			if done == false {
				fmt.Println("No se pudo crear la particion")
			}
		}
	}

	fmt.Println("Contenido despues de insertar un particion")
	for element := listaP.Front(); element != nil; element = element.Next() {
		// do something with element.Value
		fmt.Println(element.Value)
	}

	for i := 0; i < 4; i++ {
		m.Prt[i].Partstatus = -1
		m.Prt[i].Parttype = 0
		copy(m.Prt[i].Partfit[:], "")
		m.Prt[i].Partstart = 0
		m.Prt[i].Partsize = 0
		copy(m.Prt[i].Partname[:], "")
	}

	pos := 0
	for element := listaP.Front(); element != nil; element = element.Next() {
		// do something with element.Value
		valPart := element.Value.(nodoPart)

		if valPart.Estado == 1 {
			m.Prt[pos].Partstatus = valPart.Partstatus
			m.Prt[pos].Parttype = valPart.Parttype
			m.Prt[pos].Partfit = valPart.Partfit
			m.Prt[pos].Partstart = valPart.Partstart
			m.Prt[pos].Partsize = valPart.Partsize
			m.Prt[pos].Partname = valPart.Partname
			pos++
		}

	}

	escribirMbr(fd.path, m)

}
