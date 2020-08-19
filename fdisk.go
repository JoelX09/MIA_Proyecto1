package main

import (
	"container/list"
	"fmt"
	"unsafe"
)

func adminParticion(fd datoDisco, fl banderaParam) {

	type nodoPart struct {
		Estado     int
		Partstatus int8
		Parttype   byte
		Partfit    [2]byte
		Partstart  int64
		Partsize   int64
		Partname   [16]byte
	}

	listaP := list.New()
	existePart := false
	var datosPart nodoPart
	primaria, extendida := 0, 0
	m := obtenerMbr(fd.path)

	//tam := m.Mbrtam
	//fmt.Println("-------------------------\n\n" + strconv.FormatInt(tam, 10) + "\n\n-------------------------")

	var size int = int(unsafe.Sizeof(m))
	//fmt.Println(size)

	for i := 0; i < 4; i++ { // Recorro el arreglo de particiones y construyo la lista con las existentes
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

	if listaP.Len() > 0 { // Recorro la lista para crear los espacios disponibles
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
	}

	fmt.Println("Contenido despues de generar los nodos ocupados y disponibles:")
	for element := listaP.Front(); element != nil; element = element.Next() {
		temp := element.Value.(nodoPart)
		if temp.Estado == 1 {
			/*var nombre [16]byte
			copy(nombre[:], fd.name)
			if bytes.Compare(element.Value.(nodoPart).Partname, nombre) == 0 {

			}*/
			fmt.Println(string(temp.Partname[:]))
		}
		fmt.Println(element.Value)
	}

	if fl.deleteY == false && fl.addY == false { // Si son falsos es porque se va crear una nueva

		unidad, tipoPart, tipoFit := true, true, false
		var tam int64
		if fd.unit == 'K' || fd.unit == 0 || fd.unit == 'k' {
			tam = fd.size * 1024
		} else if fd.unit == 'M' || fd.unit == 'm' {
			tam = fd.size * 1024 * 1024
		} else if fd.unit == 'B' || fd.unit == 'b' {
			tam = fd.size
		} else {
			unidad = false
			fmt.Println("No se puede crear la Particion, Tipo de unidad errorneo.")
		}

		if tam < 0 {
			fmt.Println("El tamano de la particion debe ser mayor a cero")
		}

		if extendida+primaria < 4 {

			if dato.typeP == 'E' {
				if extendida >= 1 {
					fmt.Println("Ya existe una particion primaria.")
					tipoPart = false
				}
			} else if dato.typeP == 'L' {
				if extendida == 0 {
					fmt.Println("No existe una particion Extendida para crear la particion logica.")
					tipoPart = false
				}
			}
		} else {
			fmt.Println("Se alcanzo el limite de particiones que puede crear")
			tipoPart = false
		}

		if dato.fit == "BF" || dato.fit == "FF" || dato.fit == "WF" || dato.fit == "" {
			tipoFit = true
		} else {
			fmt.Println("El tipo de ajuste es incorrecto")
		}

		if unidad == true && tipoPart == true && tipoFit == true && tam > 0 {
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
						if temp.Partsize >= tam {
							temp2 := ele.Prev()

							temp.Estado = 1
							temp.Partstatus = 0
							temp.Parttype = fd.typeP
							copy(temp.Partfit[:], fd.fit)
							temp.Partsize = tam
							copy(temp.Partname[:], fd.name)
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
	}

	fmt.Println("Contenido despues de insertar un particion")
	for element := listaP.Front(); element != nil; element = element.Next() {
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
