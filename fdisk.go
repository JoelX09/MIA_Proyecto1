package main

import (
	"bufio"
	"bytes"
	"container/list"
	"encoding/binary"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"strings"
	"unsafe"
)

type estructEBR struct {
	EstadoL     int8
	PartstatusL int8
	PartfitL    byte
	PartstartL  int64
	PartsizeL   int64
	PartnextL   int64
	PartnameL   [16]byte
}

type nodoPart struct {
	Estado     int
	Partstatus int8
	Parttype   byte
	Partfit    byte
	Partstart  int64
	Partsize   int64
	Partname   [16]byte
}

type valExt struct {
	inicioE int64
	tamE    int64
}

func listaInicialPE(path string) (*list.List, [2]int) {
	var listaP = list.New()
	listaP.Init()
	m := obtenerMbr(path)
	var valores [2]int
	primaria, extendida := 0, 0
	var datosPart nodoPart

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
	valores[0] = primaria
	valores[1] = extendida
	return listaP, valores
}

func crearParticion(fd datoDisco) {
	m := obtenerMbr(fd.path)
	sizeMBR := int(unsafe.Sizeof(m))
	var existePart bool
	var datosPart nodoPart

	listaP, numPart := listaInicialPE(fd.path)
	var listaPtemp = list.New()
	listaPtemp.PushFrontList(listaP)
	listaP.Init()
	existePart, listaP = espaciosPEdisp(sizeMBR, m, listaPtemp)
	existeNombrePE, valoresExt, _ := imprimirListaPE(fd.name, false, true, listaP)

	if fd.typeP == 'L' {
		if numPart[1] == 1 {

			fmt.Println("----------\nSe va a crear una Logica\n----------")

			unidad, tipoFit, tam, fit := validarValores(fd.unit, fd.size, fd.fit)
			fd.fit = fit

			listaL := listaInicialL(fd.path, valoresExt.inicioE, valoresExt.tamE, fd.name, valoresExt.inicioE)

			fmt.Println("Lista con los ebr que existen")
			imprimirListaL(fd.name, true, false, listaL)
			var listaLtemp = list.New()
			listaLtemp.PushFrontList(listaL)
			listaL.Init()
			listaL = espaciosLL(valoresExt.inicioE, valoresExt.tamE, listaLtemp)
			fmt.Println("Lista con los ebr y espacios disponibles")
			existeNombreL, _ := imprimirListaL(fd.name, true, true, listaL)

			if existeNombreL == false {
				if unidad == true && tipoFit == true && fd.size > 0 {
					if listaL.Len() == 1 {
						datos := listaL.Front().Value.(estructEBR)
						if datos.EstadoL == 0 {
							if datos.PartsizeL >= tam {
								datos.EstadoL = 1
								datos.PartstatusL = 0
								datos.PartfitL = fd.fit[0]
								datos.PartsizeL = tam
								datos.PartnextL = -1
								copy(datos.PartnameL[:], fd.name)
								listaL.Remove(listaL.Front())
								listaL.PushFront(datos)
							}
						}
					} else if listaL.Len() == 2 {
						actual := listaL.Front()
						datosActual := actual.Value.(estructEBR)
						siguiente := listaL.Back()
						datosSiguiente := siguiente.Value.(estructEBR)

						if datosSiguiente.EstadoL == 0 {
							if datosSiguiente.PartsizeL >= tam {
								datosSiguiente.EstadoL = 1
								datosSiguiente.PartstatusL = 0
								datosSiguiente.PartfitL = fd.fit[0]
								datosSiguiente.PartsizeL = tam
								datosSiguiente.PartnextL = -1
								copy(datosSiguiente.PartnameL[:], fd.name)

								datosActual.PartnextL = datosSiguiente.PartstartL

								listaL.Remove(actual)
								listaL.PushFront(datosActual)
								listaL.Remove(siguiente)
								listaL.PushBack(datosSiguiente)
							}
						}
					} else {
						for ele := listaL.Front(); ele != nil; ele = ele.Next() {
							actual := ele.Value.(estructEBR)
							if actual.EstadoL == 0 {
								if actual.PartsizeL >= tam {
									if actual.PartstartL == valoresExt.inicioE {
										actual.EstadoL = 1
										actual.PartstatusL = 0
										actual.PartfitL = fd.fit[0]
										actual.PartsizeL = tam
										copy(actual.PartnameL[:], fd.name)

										listaL.Remove(ele)
										listaL.PushFront(actual)
										break
									} else {
										actual.EstadoL = 1
										actual.PartstatusL = 0
										actual.PartfitL = fd.fit[0]
										actual.PartsizeL = tam
										copy(actual.PartnameL[:], fd.name)

										anterior := ele.Prev()
										datosAnterior := anterior.Value.(estructEBR)
										datosAnterior.PartnextL = actual.PartstartL

										listaL.Remove(anterior)
										apuntadorTemp := listaL.InsertBefore(datosAnterior, ele)

										listaL.Remove(ele)
										listaL.InsertAfter(actual, apuntadorTemp)
										break
									}
								}
							}
						}
					}
					fmt.Println("Contenido despues de insertar una particion Logica")
					imprimirListaL(fd.name, true, false, listaL)
					fmt.Println("------------------------------------------------")
					fmt.Println("Escribiendo EBR's")
					fmt.Println("------------------------------------------------")
					escribirListaEbr(fd.path, listaL)

				}

			} else {
				fmt.Println("NO se puede crear la particion logica, ya existe una con ese nombre")
			}

		} else {
			fmt.Println("----------\nNo existe una particion extendida para crear logicas\n----------")
		}

	} else if existeNombrePE == false && fd.typeP != 'L' {
		fmt.Println("Contenido de los nodos P Y E ocupados y disponibles:")
		imprimirListaPE(fd.name, true, false, listaP)

		tipoPart := true
		unidad, tipoFit, tam, fit := validarValores(fd.unit, fd.size, fd.fit)
		fd.fit = fit

		if numPart[0]+numPart[1] < 4 {
			if fd.typeP == 'E' {
				if numPart[1] == 1 {
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

		if unidad == true && tipoPart == true && tipoFit == true && fd.size > 0 {
			if existePart == false {
				datosPart.Estado = 1
				datosPart.Partstatus = 0
				datosPart.Parttype = fd.typeP
				datosPart.Partfit = fd.fit[0]
				datosPart.Partstart = int64(sizeMBR)
				datosPart.Partsize = tam
				copy(datosPart.Partname[:], fd.name)
				listaP.PushFront(datosPart)

				if fd.typeP == 'E' {
					fmt.Println("-------------------")
					fmt.Println(valoresExt.inicioE)
					fmt.Println(datosPart.Partstart)
					valoresExt.inicioE = datosPart.Partstart
					valoresExt.tamE = tam
					fmt.Println("-------------------")
					asignarebr := ebr{Partstatus: -1, Partstart: datosPart.Partstart, Partnext: -1}
					var sizebr int = int(unsafe.Sizeof(asignarebr))
					if datosPart.Partsize > int64(sizebr) {
						escribirEbr(fd.path, asignarebr, datosPart.Partstart)
						fmt.Println("Particion agregada exitosamente")
					} else {
						fmt.Println("La particion no se puede crear. Tamano insuficiente para EBR")
					}
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
							temp.Partfit = fd.fit[0]
							temp.Partsize = tam
							copy(temp.Partname[:], fd.name)
							done = true
							if fd.typeP == 'E' {
								fmt.Println("-------------------")
								fmt.Println(valoresExt.inicioE)
								fmt.Println(temp.Partstart)
								valoresExt.inicioE = datosPart.Partstart
								valoresExt.tamE = tam
								fmt.Println("-------------------")
								asignarebr := ebr{Partstatus: -1, Partstart: temp.Partstart, Partnext: -1}
								var sizebr int = int(unsafe.Sizeof(asignarebr))
								if temp.Partsize > int64(sizebr) {
									escribirEbr(fd.path, asignarebr, temp.Partstart)
									fmt.Println("Particion agregada exitosamente")
									done = true
								} else {
									done = false
									fmt.Println("La particion no se puede crear. Tamano insuficiente para EBR")
								}
							}
							if done == true {
								listaP.Remove(ele)
								listaP.InsertAfter(temp, temp2)
							}
							break
						}
					}
				}
				if done == false {
					fmt.Println("No se pudo crear la particion, no hay espacios de disco disponibles o el tamano disponible es insuficiente")
				}
			}
		}
		fmt.Println("Contenido despues de insertar una particion")
		imprimirListaPE(fd.name, true, false, listaP)
		actualizarMBR(fd.path, listaP)

	} else if existeNombrePE == true && fd.typeP != 'L' {
		fmt.Println("Ya existe una particion P o E con ese nombre")
		fmt.Println("Contenido sin modificar")
		imprimirListaPE(fd.name, true, false, listaP)
	}
}

func eliminarParticion(fd datoDisco) {

	listaP, _ := listaInicialPE(fd.path)
	fmt.Println("Contenido de los nodos P Y E:")
	existeNombrePE, valoresExt, _ := imprimirListaPE(fd.name, true, true, listaP)

	econtrado := false
	for ele := listaP.Front(); ele != nil; ele = ele.Next() {
		temp := ele.Value.(nodoPart)
		if temp.Estado == 1 {

			var tempcomp [16]byte
			copy(tempcomp[:], fd.name)

			if temp.Parttype == 'E' && existeNombrePE == false {

				fmt.Println("Recorrer Logicas para ver si es la que se elimina")

				listaL := listaInicialL(fd.path, valoresExt.inicioE, valoresExt.tamE, fd.name, valoresExt.inicioE)

				fmt.Println("Lista con los ebr que existen")
				existeNombreL, _ := imprimirListaL(fd.name, true, true, listaL)

				if existeNombreL == true {
					for eleL := listaL.Front(); eleL != nil; eleL = eleL.Next() {
						tempL := eleL.Value.(estructEBR)
						fmt.Println(tempL)
						fmt.Println("\nEjecuacion pausada... Presione enter para continuar")
						fmt.Scanln()
						if tempL.EstadoL == 1 {

							if tempL.PartnameL == tempcomp {
								if strings.ToLower(fd.deleteP) == "fast" {
									if confirmarEliminacion() == true {
										if tempL.PartstartL == valoresExt.inicioE {
											tempL.PartstatusL = -1
											tempL.PartfitL = 0
											tempL.PartsizeL = tempL.PartnextL - tempL.PartstartL
											for j := 0; j < len(tempL.PartnameL); j++ {
												tempL.PartnameL[j] = 0
											}
											listaL.Remove(eleL)
											listaL.PushFront(tempL)
											fmt.Println("Particion eliminada correctamente")
										} else {
											tempLAnt := eleL.Prev()
											valtempLAnt := tempLAnt.Value.(estructEBR)
											valtempLAnt.PartnextL = tempL.PartnextL
											listaL.Remove(tempLAnt)
											listaL.InsertBefore(valtempLAnt, eleL)
											listaL.Remove(eleL)
											fmt.Println("Particion eliminada correctamente")
										}
										econtrado = true
										break
									}
								} else if strings.ToLower(fd.deleteP) == "full" {
									if confirmarEliminacion() == true {
										if tempL.PartstartL == valoresExt.inicioE {
											tempL.PartstatusL = -1
											tempL.PartfitL = 0
											tempL.PartsizeL = tempL.PartnextL - tempL.PartstartL
											for j := 0; j < len(tempL.PartnameL); j++ {
												tempL.PartnameL[j] = 0
											}
											var tamEBR int64
											tamEBR = int64(unsafe.Sizeof(ebr{}))
											deleteFull(fd.path, tempL.PartstartL+tamEBR, tempL.PartsizeL-tempL.PartstartL+tamEBR)
											listaL.Remove(eleL)
											listaL.PushFront(tempL)
											fmt.Println("Particion eliminada correctamente")
										} else {
											tempLAnt := eleL.Prev()
											valtempLAnt := tempLAnt.Value.(estructEBR)
											valtempLAnt.PartnextL = tempL.PartnextL
											deleteFull(fd.path, tempL.PartstartL, tempL.PartsizeL)
											listaL.Remove(tempLAnt)
											listaL.InsertBefore(valtempLAnt, eleL)
											listaL.Remove(eleL)
											fmt.Println("Particion eliminada correctamente")
										}
										econtrado = true
										break
									}
								} else {
									fmt.Println("Valor del delete incorrecto")
									break
								}

							}
						}
					}

					fmt.Println("Contenido despues de eliminar una particion Logica")
					imprimirListaL(fd.name, true, false, listaL)
					fmt.Println("------------------------------------------------")
					fmt.Println("Escribiendo EBR's")
					fmt.Println("------------------------------------------------")
					escribirListaEbr(fd.path, listaL)

					break
				}

			}

			if econtrado == false && temp.Partname == tempcomp {
				if strings.ToLower(fd.deleteP) == "fast" {
					if confirmarEliminacion() == true {
						listaP.Remove(ele)
						fmt.Println("Particion eliminada correctamente")
						econtrado = true

						fmt.Println("Contenido despues de eliminar una particion")
						imprimirListaPE(fd.name, true, false, listaP)
						actualizarMBR(fd.path, listaP)

						break
					}
				} else if strings.ToLower(fd.deleteP) == "full" {
					if confirmarEliminacion() == true {
						deleteFull(fd.path, temp.Partstart, temp.Partsize)
						listaP.Remove(ele)
						fmt.Println("Particion eliminada correctamente")
						econtrado = true

						fmt.Println("Contenido despues de eliminar una particion")
						imprimirListaPE(fd.name, true, false, listaP)
						actualizarMBR(fd.path, listaP)

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
}

func aumentarParticion(fd datoDisco) {
	m := obtenerMbr(fd.path)
	sizeMBR := int(unsafe.Sizeof(m))

	listaP, _ := listaInicialPE(fd.path)
	var listaPtemp = list.New()
	listaPtemp.PushFrontList(listaP)
	listaP.Init()
	_, listaP = espaciosPEdisp(sizeMBR, m, listaPtemp)
	fmt.Println("Contenido de los nodos P Y E ocupados y disponibles:")
	existeNombrePE, valoresExt, _ := imprimirListaPE(fd.name, false, true, listaP)

	econtrado := false
	for ele := listaP.Front(); ele != nil; ele = ele.Next() {
		temp := ele.Value.(nodoPart)

		unidad, tam := validarValorAdd(fd.unit, fd.add)

		if unidad == false {
			break
		}

		if temp.Estado == 1 {

			var tempcomp [16]byte
			copy(tempcomp[:], fd.name)

			if temp.Parttype == 'E' && existeNombrePE == false {
				fmt.Println("Recorrer Logicas para ver si es la que se aumenta")

				listaL := listaInicialL(fd.path, valoresExt.inicioE, valoresExt.tamE, fd.name, valoresExt.inicioE)

				fmt.Println("Lista con los ebr que existen")
				imprimirListaL(fd.name, true, false, listaL)
				var listaLtemp = list.New()
				listaLtemp.PushFrontList(listaL)
				listaL.Init()
				listaL = espaciosLL(valoresExt.inicioE, valoresExt.tamE, listaLtemp)
				fmt.Println("Lista con los ebr y espacios disponibles")
				existeNombreL, _ := imprimirListaL(fd.name, true, true, listaL)

				if existeNombreL == true {
					for eleL := listaL.Front(); eleL != nil; eleL = eleL.Next() {
						tempL := eleL.Value.(estructEBR)
						if tempL.PartnameL == tempcomp {
							fmt.Println("-----------------------------")
							fmt.Println("Lo encontro")
							var cero int64
							cero = 0
							fmt.Println(cero)
							fmt.Println("-----------------------------")

							tempLSig := eleL.Next()
							tempLSigVal := tempLSig.Value.(estructEBR)
							if fd.add >= cero {
								fmt.Println("-----------------------------")
								fmt.Println("Tam positivo")
								fmt.Println("-----------------------------")
								if tempLSigVal.EstadoL == 0 {
									fmt.Println("-----------------------------")
									fmt.Println("Espacio siguiente disponible")
									fmt.Println("-----------------------------")
									if tempL.PartstartL+tempL.PartsizeL+tam-1 < tempLSigVal.PartstartL+tempLSigVal.PartsizeL {
										tempL.PartsizeL = tempL.PartsizeL + tam
										listaL.Remove(eleL)
										listaL.InsertBefore(tempL, tempLSig)
										fmt.Println("Particion aumentada exitosamente")
										econtrado = true
										break
									} else {
										fmt.Println("NO hay espacio libre suficiente despues de la particion para aumentar tamano")
										econtrado = true
										break
									}
								}
							} else {
								if tempL.PartsizeL+tam > 0 {
									tempL.PartsizeL = tempL.PartsizeL + tam
									listaL.Remove(eleL)
									listaL.InsertBefore(tempL, tempLSig)
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

					fmt.Println("Contenido despues de aumentar o disminuir una particion Logica")
					imprimirListaL(fd.name, true, false, listaL)
					fmt.Println("------------------------------------------------")
					fmt.Println("Escribiendo EBR's")
					fmt.Println("------------------------------------------------")
					escribirListaEbr(fd.path, listaL)

					break
				}
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

							fmt.Println("Contenido despues modificar una particion")
							imprimirListaPE(fd.name, true, false, listaP)
							actualizarMBR(fd.path, listaP)

							break
						} else {
							fmt.Println("NO hay espacio libre suficiente despues de la particion para aumentar tamano")
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

						fmt.Println("Contenido despues modificar una particion")
						imprimirListaPE(fd.name, true, false, listaP)
						actualizarMBR(fd.path, listaP)

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

func actualizarMBR(path string, listaP *list.List) {
	m := obtenerMbr(path)
	for i := 0; i < 4; i++ { //Vaciar arreglo de particiones
		m.Prt[i].Partstatus = -1
		m.Prt[i].Parttype = 0
		m.Prt[i].Partfit = 0
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
	escribirMbr(path, m)
}

func adminParticion(fd datoDisco, fl banderaParam) {

	if fl.deleteY == false && fl.addY == false { // Si son falsos es porque se va crear una nueva
		crearParticion(fd)
	}

	if fl.deleteY == true {
		eliminarParticion(fd)

	} else if fl.addY == true {
		aumentarParticion(fd)
	}
}

func crearDot() {
	contenido := "digraph g{\n" +
		"15->7->-3\n" +
		"15->22->17\n" +
		"7->10\n" +
		"10->8\n" +
		"22->35\n" +
		"}"

	f, err := os.Create("/home/joel/Escritorio/Prueba.dot")
	if err != nil {
		panic(err)
	}
	f.Close()

	err = ioutil.WriteFile("/home/joel/Escritorio/Prueba.dot", []byte(contenido), 0777)
	if err != nil {
		panic(err)
	}

	cmd := exec.Command("dot", "-Tpng", "/home/joel/Escritorio/Prueba.dot", "-o", "/home/joel/Escritorio/Prueba.png")
	cmd.Run()
	//log.Printf("Command finished with error: %v", erro)

	/*
		cmd := exec.Command("comando)
			cmd.Stdout = os.Stdout
			cmd.Run()
	*/
}

func imprimirListaPE(name string, imprimir bool, buscarNombre bool, listaP *list.List) (bool, valExt, nodoPart) {
	var valoresExt valExt
	var nodoReturn nodoPart

	existeNombrePE := false
	for element := listaP.Front(); element != nil; element = element.Next() {
		if buscarNombre == true {
			temp := element.Value.(nodoPart)
			if temp.Estado == 1 {

				var tempcomp [16]byte
				copy(tempcomp[:], name)

				if temp.Parttype == 'E' {
					valoresExt.inicioE = temp.Partstart
					valoresExt.tamE = temp.Partsize
				}

				if temp.Partname == tempcomp {
					existeNombrePE = true
					nodoReturn = temp
					fmt.Println("El nombre de la particion se econtro")
				}
			}
		}
		if imprimir == true {
			fmt.Println(element.Value)
		}
	}
	return existeNombrePE, valoresExt, nodoReturn
}

func espaciosPEdisp(size int, m mbr, listaPa *list.List) (bool, *list.List) {
	existePart := false
	pos := 0
	listaP := listaPa
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
			pos++
		}
	}
	return existePart, listaP
}

func validarValores(unit byte, size int64, fit string) (bool, bool, int64, string) {
	unidad, tipoFit := true, false
	var tam int64
	if unit == 'K' || unit == 0 || unit == 'k' {
		tam = size * 1024
	} else if unit == 'M' || unit == 'm' {
		tam = size * 1024 * 1024
	} else if unit == 'B' || unit == 'b' {
		tam = size
	} else {
		unidad = false
		fmt.Println("No se puede crear la Particion, Tipo de unidad erroneo.")
	}

	if size < 0 {
		fmt.Println("El valor de size debe ser mayor a cero")
	}

	if fit == "BF" || fit == "FF" || fit == "WF" || fit == "" {
		tipoFit = true
		if fit == "" {
			fit = "WF"
		}
	} else {
		fmt.Println("El tipo de ajuste es incorrecto")
	}
	return unidad, tipoFit, tam, fit
}

func validarValorAdd(unit byte, add int64) (bool, int64) {
	unidad := true
	var tam int64
	if unit == 'K' || unit == 0 || unit == 'k' {
		tam = add * 1024
	} else if unit == 'M' || unit == 'm' {
		tam = add * 1024 * 1024
	} else if unit == 'B' || unit == 'b' {
		tam = add
	} else {
		unidad = false
		fmt.Println("No se puede aumentar o disminuir la Particion, Tipo de unidad errorneo.")
	}
	return unidad, tam
}

func deleteFull(path string, inicio int64, tam int64) {
	f, err := os.OpenFile(path, os.O_RDWR, 0777)
	defer f.Close()
	if err != nil {
		log.Fatal(err)
	}

	var binario bytes.Buffer

	f.Seek(inicio, 0)

	temptam := make([]byte, tam)

	err4 := binary.Write(&binario, binary.BigEndian, temptam)
	if err4 != nil {
		fmt.Println("binary error ", err4)
	}
	escribirBytesDelete(f, binario.Bytes())
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

var listaNL = list.New()

func listaInicialL(path string, inicioE int64, tamE int64, name string, iniciEBR int64) *list.List {
	var datosEBR estructEBR
	contenidoEBR := obtenerEbr(path, iniciEBR)
	if contenidoEBR.Partstart == inicioE {
		if contenidoEBR.Partstatus == -1 {
			if contenidoEBR.Partnext == -1 {
				datosEBR.EstadoL = 0
				contenidoEBR.Partstatus = -1 // <------------ Confirmar
				contenidoEBR.Partsize = tamE
			} else if contenidoEBR.Partnext != -1 {
				datosEBR.EstadoL = 0
				contenidoEBR.Partstatus = -1
				contenidoEBR.Partsize = contenidoEBR.Partnext - contenidoEBR.Partstart
			}

		}
	}

	if contenidoEBR.Partstatus == 0 || contenidoEBR.Partstatus == 1 {
		datosEBR.EstadoL = 1
	}
	datosEBR.PartstatusL = contenidoEBR.Partstatus
	datosEBR.PartfitL = contenidoEBR.Partfit
	datosEBR.PartstartL = contenidoEBR.Partstart
	datosEBR.PartsizeL = contenidoEBR.Partsize
	datosEBR.PartnextL = contenidoEBR.Partnext
	datosEBR.PartnameL = contenidoEBR.Partname
	listaNL.PushBack(datosEBR)
	if contenidoEBR.Partnext != -1 {
		listaInicialL(path, inicioE, tamE, name, contenidoEBR.Partnext)
	}
	return listaNL
}

func espaciosLL(inicioE int64, tamE int64, listaL *list.List) *list.List {
	for ele := listaL.Front(); ele != nil; ele = ele.Next() {
		actual := ele.Value.(estructEBR)
		if listaL.Len() == 1 {
			if actual.EstadoL == 1 {
				var nuevo estructEBR
				nuevo.EstadoL = 0
				nuevo.PartstartL = actual.PartstartL + actual.PartsizeL
				nuevo.PartsizeL = inicioE + tamE - nuevo.PartstartL
				nuevo.PartnextL = -1
				listaL.PushBack(nuevo)
			}
			break
		} else {
			if actual.PartnextL != -1 && actual.EstadoL == 1 {
				if actual.PartstartL+actual.PartsizeL < actual.PartstartL {
					var nuevo estructEBR
					nuevo.EstadoL = 0
					nuevo.PartstatusL = 0
					nuevo.PartfitL = 0
					nuevo.PartstartL = actual.PartstartL + actual.PartsizeL
					nuevo.PartsizeL = actual.PartnextL - nuevo.PartstartL
					nuevo.PartnextL = actual.PartnextL
					for j := 0; j < len(nuevo.PartnameL); j++ {
						nuevo.PartnameL[j] = 0
					}
					listaL.InsertAfter(nuevo, ele)
				}
			}
			if actual.PartnextL == -1 && actual.EstadoL == 1 {
				if actual.PartstartL+actual.PartsizeL < inicioE+tamE {
					var nuevo estructEBR
					nuevo.EstadoL = 0
					nuevo.PartstatusL = 0
					nuevo.PartfitL = 0
					nuevo.PartstartL = actual.PartstartL + actual.PartsizeL
					nuevo.PartsizeL = inicioE + tamE - nuevo.PartstartL
					nuevo.PartnextL = -1
					for j := 0; j < len(nuevo.PartnameL); j++ {
						nuevo.PartnameL[j] = 0
					}
					listaL.InsertAfter(nuevo, ele)
				}
			}

		}
	}
	return listaL
}

func imprimirListaL(name string, imprimir bool, buscarNombre bool, listaL *list.List) (bool, estructEBR) {
	encontrado := false
	var nodoReturn estructEBR

	for ele := listaL.Front(); ele != nil; ele = ele.Next() {
		temp := ele.Value.(estructEBR)
		if temp.EstadoL == 1 {

			var tempcomp [16]byte
			copy(tempcomp[:], name)

			if temp.PartnameL == tempcomp {
				if buscarNombre == true {
					fmt.Println("Nombres iguales Logicas") //<--------------Cambiar
					nodoReturn = temp
					encontrado = true
				}
			} /*else {
				fmt.Println("No son iguales los nombres Logicas") //<--------------Cambiar
			}*/
		}

		if imprimir == true {
			fmt.Println(ele.Value)
		}
	}
	return encontrado, nodoReturn
}

func escribirListaEbr(path string, listaL *list.List) {
	for ele := listaL.Front(); ele != nil; ele = ele.Next() {
		temp := ele.Value.(estructEBR)
		if temp.EstadoL == 1 {

			asignarebr := ebr{}
			asignarebr.Partstatus = temp.PartstatusL
			asignarebr.Partfit = temp.PartfitL
			asignarebr.Partstart = temp.PartstartL
			asignarebr.Partsize = temp.PartsizeL
			asignarebr.Partnext = temp.PartnextL
			asignarebr.Partname = temp.PartnameL
			escribirEbr(path, asignarebr, temp.PartstartL)

		}
	}
}
