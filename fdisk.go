package main

import (
	"bufio"
	"bytes"
	"container/list"
	"encoding/binary"
	"fmt"
	"log"
	"os"
	"strings"
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
	existePart, existeNombrePE := false, false
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

			var tempcomp [16]byte
			copy(tempcomp[:], fd.name)

			if temp.Partname == tempcomp {
				fmt.Println("Nombres iguales") //<--------------Cambiar
				existeNombrePE = true
			} else {
				fmt.Println("No son iguales los nombres") //<--------------Cambiar
			}
		}

		fmt.Println(element.Value)
	}

	if fl.deleteY == false && fl.addY == false { // Si son falsos es porque se va crear una nueva

		if fd.typeP == 'L' {
			if extendida == 1 {
				fmt.Println("----------\nSe va a crear una Logica\n----------")
			} else {
				fmt.Println("----------\nNo existe una particion extendida para crear logicas\n----------")
			}

		} else if existeNombrePE == false {
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

			if fd.size < 0 {
				fmt.Println("El valor de size debe ser mayor a cero")
			}

			if extendida+primaria < 4 {
				if fd.typeP == 'E' {
					if extendida == 1 {
						fmt.Println("Ya existe una particion Extendida.")
						tipoPart = false
					}
				} else if fd.typeP == 0 {
					fmt.Println("No se declaro tipo particion")
					fd.typeP = 'P'
					fmt.Println(fd.typeP)
				} else if fd.typeP != 'P' {
					tipoPart = false
					fmt.Println("Tipo de particion incorrecto")
				}
			} else {
				fmt.Println("Se alcanzo el limite de particiones que puede crear")
				tipoPart = false
			}

			if fd.fit == "BF" || fd.fit == "FF" || fd.fit == "WF" || fd.fit == "" {
				tipoFit = true
				if fd.fit == "" {
					fd.fit = "WF"
				}
			} else {
				fmt.Println("El tipo de ajuste es incorrecto")
			}

			if unidad == true && tipoPart == true && tipoFit == true && fd.size > 0 {
				if existePart == false {
					datosPart.Estado = 1
					datosPart.Partstatus = 0
					datosPart.Parttype = fd.typeP
					copy(datosPart.Partfit[:], fd.fit)
					datosPart.Partstart = int64(size)
					datosPart.Partsize = tam
					copy(datosPart.Partname[:], fd.name)
					listaP.PushFront(datosPart)
					fmt.Println("Particion agregada exitosamente")
					if fd.typeP == 'E' {
						asignarebr := ebr{Partstatus: -1, Partnext: -1}
						escribirEbr(fd.path, asignarebr, datosPart.Partstart)
					}
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
								fmt.Println("Particion agregada exitosamente")
								if fd.typeP == 'E' {
									asignarebr := ebr{Partstatus: -1, Partnext: -1}
									escribirEbr(fd.path, asignarebr, datosPart.Partstart)
								}
								done = true
								break
							}
						}
					}
					if done == false {
						fmt.Println("No se pudo crear la particion, no hay espacios de disco disponibles o el tamano disponible es insuficiente")
					}
				}
			}
		} else {
			fmt.Println("Ya existe una particion P o E con ese nombre")
		}
	}

	if fl.deleteY == true {
		econtrado := false
		for ele := listaP.Front(); ele != nil; ele = ele.Next() {
			temp := ele.Value.(nodoPart)
			if temp.Estado == 1 {

				var tempcomp [16]byte
				copy(tempcomp[:], fd.name)

				if temp.Parttype == 'E' {
					fmt.Println("Recorrer Logicas para ver si es la que se elimina")
				}

				if temp.Partname == tempcomp {
					if strings.ToLower(fd.deleteP) == "fast" {
						if confirmarEliminacion() == true {
							listaP.Remove(ele)
							fmt.Println("Particion eliminada correctamente")
							econtrado = true
							break
						}
					} else if strings.ToLower(fd.deleteP) == "full" {
						if confirmarEliminacion() == true {
							deleteFull(fd.path, temp.Partstart, temp.Partsize)
							listaP.Remove(ele)
							fmt.Println("Particion eliminada correctamente")
							econtrado = true
							break
						}
					} else {
						fmt.Println("Valor del delete incorrecto")
						econtrado = true
						break
					}
				}
			}
		}
		if econtrado == false {
			fmt.Println("No se encontro ningua particion para eliminar con el nombre: " + fd.name)
		}
	} else if fl.addY == true {
		econtrado := false
		for ele := listaP.Front(); ele != nil; ele = ele.Next() {
			temp := ele.Value.(nodoPart)

			unidad := true
			var tam int64
			if fd.unit == 'K' || fd.unit == 0 || fd.unit == 'k' {
				tam = fd.add * 1024
			} else if fd.unit == 'M' || fd.unit == 'm' {
				tam = fd.add * 1024 * 1024
			} else if fd.unit == 'B' || fd.unit == 'b' {
				tam = fd.add
			} else {
				unidad = false
				fmt.Println("No se puede aumentar o disminuir la Particion, Tipo de unidad errorneo.")
				break
			}

			if temp.Estado == 1 {

				var tempcomp [16]byte
				copy(tempcomp[:], fd.name)

				if temp.Parttype == 'E' {
					fmt.Println("Recorrer Logicas para ver si es la que se elimina")
				}

				if temp.Partname == tempcomp && unidad == true {
					fmt.Println("-----------------------------")
					fmt.Println("Lo encontro")
					var cero int64
					cero = 0
					fmt.Println(cero)
					fmt.Println("-----------------------------")

					tempSig := ele.Next()
					tempSigVal := tempSig.Value.(nodoPart)
					if fd.add >= cero {
						fmt.Println("-----------------------------")
						fmt.Println("Tam positivo")
						fmt.Println("-----------------------------")
						if tempSigVal.Estado == 0 {
							fmt.Println("-----------------------------")
							fmt.Println("Espacio siguiente disponible")
							fmt.Println("-----------------------------")
							if temp.Partstart+temp.Partsize+tam-1 < tempSigVal.Partstart+tempSigVal.Partsize {
								temp.Partsize = temp.Partsize + tam
								listaP.Remove(ele)
								listaP.InsertBefore(temp, tempSig)
								fmt.Println("Particion aumentada exitosamente")
								econtrado = true
								break
							} else {
								fmt.Println("NO hay espacio libre despues de la particion para aumentar tamano")
								econtrado = true
								break
							}
						}
					} else {
						if temp.Partsize+tam > 0 {
							temp.Partsize = temp.Partsize + tam
							listaP.Remove(ele)
							listaP.InsertBefore(temp, tempSig)
							fmt.Println("La particion se redujo")
							econtrado = true
							break
						} else {
							fmt.Println("NO se puede reducir la particion a un espacio negativo")
							econtrado = true
							break
						}
					}
				}
			}
		}
		if econtrado == false {
			fmt.Println("No se encontro ninguna particion para aumentar con ese nombre")
		}
	}

	fmt.Println("Contenido despues de insertar o modificar o eliminar una particion")
	for element := listaP.Front(); element != nil; element = element.Next() {
		fmt.Println(element.Value)
	}

	for i := 0; i < 4; i++ { //Vaciar arreglo de particiones
		m.Prt[i].Partstatus = -1
		m.Prt[i].Parttype = 0
		for j := 0; j < len(m.Prt[i].Partfit); j++ {
			m.Prt[i].Partfit[j] = 0
		}
		m.Prt[i].Partstart = 0
		m.Prt[i].Partsize = 0
		for j := 0; j < len(m.Prt[i].Partname); j++ {
			m.Prt[i].Partname[j] = 0
		}
	}

	pos := 0
	for element := listaP.Front(); element != nil; element = element.Next() { //insertar particiones modificadas en mbr
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

func deleteFull(path string, inicio int64, tam int64) {
	f, err := os.OpenFile(path, os.O_RDWR, 0777)
	defer f.Close()
	if err != nil {
		log.Fatal(err)
	}

	var binario bytes.Buffer

	f.Seek(inicio, 0)
	//for i := inicio; i < fin; i++ {
	//var temptam [tam]byte

	temptam := make([]byte, tam)

	err4 := binary.Write(&binario, binary.BigEndian, temptam)
	if err4 != nil {
		fmt.Println("binary error ", err4)
	}
	escribirBytesDelete(f, binario.Bytes())
	//}
}

func escribirBytesDelete(file *os.File, bytes []byte) {

	_, err := file.Write(bytes)

	if err != nil {
		log.Fatal(err)
	}

}

func confirmarEliminacion() bool {
	fmt.Println("Desea remover la particion [y/n]")
	reader := bufio.NewReader(os.Stdin)
	lectura, _ := reader.ReadString('\n')
	eleccion := strings.TrimRight(lectura, "\n")
	if eleccion == "y" {
		return true
	} else if eleccion == "n" {
		fmt.Println("No se eliminara la particion")
		return false
	} else {
		fmt.Println("Confirmacion invalida. No se realizara la eliminacion")
		return false
	}
}
