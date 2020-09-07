package main

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
	"time"
	"unsafe"
)

func crearCarpeta(vd string, path string, p bool) {
	var idDisco byte
	idDisco = vd[2]
	idDisco2 := idDisco - 97
	idP, _ := strconv.Atoi(vd[3:])
	idP--

	if arregloMount[idDisco2].estado == 1 {
		if arregloMount[idDisco2].discos[idP].estado == 1 {
			inicioPart := arregloMount[idDisco2].discos[idP].Partstart
			tamPart := arregloMount[idDisco2].discos[idP].Partsize
			rutaDisco := arregloMount[idDisco2].Ruta
			superBloque := obtenerSB(rutaDisco, inicioPart)

			if path == "/" {
				fmt.Println("Se va a crear la carpeta Raiz")
				/*listadobitmap(rutaDisco, superBloque.SBapBAVD, int(superBloque.SBsizeStructAVD))
				fmt.Println("\nEjecuacion pausada... Presione enter para continuar")
				fmt.Scanln()*/

				nuevoAVD := avd{}
				for j := 0; j < 6; j++ {
					nuevoAVD.AVDapArraySub[j] = -1
				}
				nuevoAVD.AVDapDetalleDir = -1
				nuevoAVD.AVDapAVD = -1
				fecha := time.Now().Format("2006-01-02 15:04:05")
				copy(nuevoAVD.AVDfechaCreacion[:], fecha)
				copy(nuevoAVD.AVDnombreDirectorio[:], path)
				nuevoAVD.AVDproper = 1
				/*fmt.Println("Avd raiz")
				fmt.Println(nuevoAVD)
				fmt.Println("En la posicion")
				fmt.Println(superBloque.SBapAVD)
				fmt.Println("EN la ruta")
				fmt.Println(rutaDisco)
				fmt.Println("Posicion bitmap")
				fmt.Println(superBloque.SBfirstFreeBitAVD)
				fmt.Println("posicion en el espacio para avd")
				fmt.Println(superBloque.SBapBAVD)*/

				escribirStructAVD(rutaDisco, superBloque.SBapAVD, nuevoAVD)
				superBloque.SBavdFree--
				actualizarValorBitmap(rutaDisco, superBloque.SBapBAVD, '1')
				nuevoFFBAVD := obtenerFirstFreeBit(rutaDisco, superBloque.SBapBAVD, int(tamPart))
				/*fmt.Println(nuevoFFBAVD)
				fmt.Println("\nEjecuacion pausada... Presione enter para continuar")
				fmt.Scanln()*/
				superBloque.SBfirstFreeBitAVD = nuevoFFBAVD
				escribirSuperBloque(rutaDisco, inicioPart, superBloque)
				/*listadobitmap(rutaDisco, superBloque.SBapBAVD, int(superBloque.SBsizeStructAVD))
				fmt.Println("\nEjecuacion pausada... Presione enter para continuar")
				fmt.Scanln()*/
			} else {
				fmt.Println("Crear carpetas en raiz")
				path1 := strings.TrimPrefix(path, "/")
				path2 := strings.TrimSuffix(path1, "/")
				pathPart := strings.Split(path2, "/")

				//i := 0
				encontrado := false
				posEncontrado := superBloque.SBapAVD
				t := 0
				/*for i := 0; i < len(pathPart); i++ {
					fmt.Println(pathPart[i])
				}*/
				for i := 0; i < len(pathPart); i++ {
					t = i
					fmt.Println(pathPart[i])
					encontrado, posEncontrado = buscarDir(posEncontrado, pathPart[i], rutaDisco)
					if encontrado == false {
						break
					}
				}

				if encontrado == false {
					if t == len(pathPart)-1 {
						crearDir(rutaDisco, superBloque, path2, inicioPart)
					} else {
						if p == true {
							for i := 0; i < len(pathPart); i++ {
								/*fmt.Println("Entro")
								fmt.Println("\nEjecuacion pausada... Presione enter para continuar")
								fmt.Scanln()*/
								superBloque = obtenerSB(rutaDisco, inicioPart)
								crearDir(rutaDisco, superBloque, path2, inicioPart)
								/*fmt.Println("Salgo")
								fmt.Println("\nEjecuacion pausada... Presione enter para continuar")
								fmt.Scanln()*/
							}
						} else {
							fmt.Println("No se pueden crear las carpetas padres, falta de parametro")
						}
					}
				}
			}

		} else {
			fmt.Println("La particion indicada no esta montada")
		}
	} else {
		fmt.Println("EL disco proporcionado no esta montado")
	}
}

func crearDir(rutaDisco string, superBloque sb, path2 string, inicioPart int64) {
	var posAnterior int64 = 0
	posEncontrado := superBloque.SBapAVD
	encontrado := false
	pathPart := strings.Split(path2, "/")
	for i := 0; i < len(pathPart); i++ {
		posAnterior = posEncontrado
		/*fmt.Println("BUsco en la pos")
		fmt.Println(posEncontrado)
		fmt.Println("EL directotiro")
		fmt.Println(pathPart[i])*/
		encontrado, posEncontrado = buscarDir(posEncontrado, pathPart[i], rutaDisco)
		if encontrado == false {
			nuevoAVD := avd{}
			for j := 0; j < 6; j++ {
				nuevoAVD.AVDapArraySub[j] = -1
			}
			nuevoAVD.AVDapDetalleDir = -1
			nuevoAVD.AVDapAVD = -1
			fecha := time.Now().Format("2006-01-02 15:04:05")
			copy(nuevoAVD.AVDfechaCreacion[:], fecha)
			copy(nuevoAVD.AVDnombreDirectorio[:], pathPart[i])
			nuevoAVD.AVDproper = 1

			/*fmt.Println("Nuevo avd")
			fmt.Println(nuevoAVD)
			fmt.Println("creo directorio")
			fmt.Println(pathPart[i])*/

			posFirstFreeBit := superBloque.SBfirstFreeBitAVD //obtenerFirstFreeBit(rutaDisco, superBloque.SBapBAVD, int(superBloque.SBsizeStructAVD))
			//listadobitmap(rutaDisco, superBloque.SBapBAVD, int(superBloque.SBsizeStructAVD))
			/*fmt.Println("POsicion en bitmap avd libre")
			fmt.Println(posFirstFreeBit)
			fmt.Println("\nEjecuacion pausada... Presione enter para continuar")
			fmt.Scanln()*/

			posNuevoAVD := superBloque.SBapAVD + (posFirstFreeBit * superBloque.SBsizeStructAVD)
			/*fmt.Println("Posicion en la que insertare avd")
			fmt.Println(posNuevoAVD)
			fmt.Println("\nEjecuacion pausada... Presione enter para continuar")
			fmt.Scanln()*/

			avdPadre := obtenerAVD(rutaDisco, posAnterior)

			for j := 0; j < len(avdPadre.AVDapArraySub); j++ {
				if avdPadre.AVDapArraySub[j] == -1 {
					avdPadre.AVDapArraySub[j] = posNuevoAVD
					break
				}
			}
			/*fmt.Println("Actualizo datos del padre, esta en la pos")
			fmt.Println(posAnterior)
			fmt.Println("\nEjecuacion pausada... Presione enter para continuar")
			fmt.Scanln()*/
			escribirStructAVD(rutaDisco, posAnterior, avdPadre)
			/*fmt.Println("Ingreso el nuevo avd en la pos")
			fmt.Println(posNuevoAVD)
			fmt.Println("\nEjecuacion pausada... Presione enter para continuar")
			fmt.Scanln()*/
			escribirStructAVD(rutaDisco, posNuevoAVD, nuevoAVD)
			superBloque.SBavdFree--
			/*fmt.Println("****************************************************************")
			fmt.Println("Valor 1 al bitmap avd en la pos")
			fmt.Println(superBloque.SBapBAVD + posFirstFreeBit)
			fmt.Println("**********************************************************************")*/
			actualizarValorBitmap(rutaDisco, superBloque.SBapBAVD+posFirstFreeBit, '1')
			posFirstFreeBit = obtenerFirstFreeBit(rutaDisco, superBloque.SBapBAVD, int(superBloque.SBsizeStructAVD))
			/*fmt.Println("*************************************************************")
			fmt.Println("Nuevo primer bit libre")
			fmt.Println(posFirstFreeBit)
			fmt.Println("*************************************************************")*/
			superBloque.SBfirstFreeBitAVD = posFirstFreeBit

			escribirSuperBloque(rutaDisco, inicioPart, superBloque)

			fmt.Println("Se creo el directorio: " + pathPart[i])
			/*fmt.Println("\nEjecuacion pausada... Presione enter para continuar")
			fmt.Scanln()*/
			break
		}
	}
	/*fmt.Println("Saliendo ptm....")
	fmt.Println("\nEjecuacion pausada... Presione enter para continuar")
	fmt.Scanln()
	fmt.Println("CSTM")*/
}

func buscarDir(pos int64, nombre string, ruta string) (bool, int64) {

	arbol := obtenerAVD(ruta, pos)
	encontrado := false
	/*fmt.Println("Buscar en el avd:")
	fmt.Println(arbol)
	fmt.Println("Presione enter para continuar")
	fmt.Scanln()
	fmt.Println("El nombre")
	fmt.Println(nombre)
	fmt.Println("Presione enter para continuar")
	fmt.Scanln()*/
	var posEncontrado int64 = 0
	for i := 0; i < len(arbol.AVDapArraySub); i++ {
		/*fmt.Println("Analizando pos")
		fmt.Println(i)
		fmt.Println("Presione enter para continuar")
		fmt.Scanln()*/

		if arbol.AVDapArraySub[i] != -1 {
			/*fmt.Println("Las pos de abajo esta ocupada")
			fmt.Println(i)
			fmt.Println("Presione enter para continuar")
			fmt.Scanln()*/
			temp := obtenerAVD(ruta, arbol.AVDapArraySub[i])
			nomb := ""
			for j := 0; j < len(temp.AVDnombreDirectorio); j++ {
				if temp.AVDnombreDirectorio[j] != 0 {
					nomb += string(temp.AVDnombreDirectorio[j])
				}
			}
			if nombre == nomb {
				encontrado = true
				posEncontrado = arbol.AVDapArraySub[i]
				fmt.Println("Se encontro el directorio: " + nombre)
				break
			}
		}
	}
	if encontrado == false {
		if arbol.AVDapAVD != -1 {
			encontrado, posEncontrado = buscarDir(arbol.AVDapAVD, nombre, ruta)
		}
	}
	return encontrado, posEncontrado
}

func obtenerFirstFreeBit(ruta string, posIni int64, tam int) int64 {

	file, err := os.OpenFile(ruta, os.O_RDWR, 0777)
	if err != nil {
		log.Fatal(err)
	}

	var car byte
	var size int = int(unsafe.Sizeof(car))
	pos := 0
	for i := 0; i < tam; i++ {
		file.Seek(posIni, 0)

		data := obtenerBytesEbr(file, size)
		buffer := bytes.NewBuffer(data)

		err = binary.Read(buffer, binary.BigEndian, &car)
		if err != nil {
			log.Fatal("binary.Read failed", err)
		}
		/*fmt.Println("Posicion bitmap obteniendo First free")
		fmt.Println(posIni)
		fmt.Println("Valor en la posicion")
		fmt.Println(car)
		fmt.Println("valor i")
		fmt.Println(i)*/
		if car != 49 {
			pos = i
			break
		}
		posIni++
	}

	file.Close()
	return int64(pos)
}

/*
func listadobitmap(ruta string, posIni int64, tam int) int64 {

	file, err := os.OpenFile(ruta, os.O_RDWR, 0777)
	if err != nil {
		log.Fatal(err)
	}

	var car byte
	var size int = int(unsafe.Sizeof(car))
	pos := 0
	for i := 0; i < tam; i++ {
		file.Seek(posIni, 0)

		data := obtenerBytesEbr(file, size)
		buffer := bytes.NewBuffer(data)

		err = binary.Read(buffer, binary.BigEndian, &car)
		if err != nil {
			log.Fatal("binary.Read failed", err)
		}

		fmt.Print("valor i ")
		fmt.Println(car)
		posIni++
	}

	file.Close()
	return int64(pos)
}
*/
