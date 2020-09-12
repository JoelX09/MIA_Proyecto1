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

func crearCarpeta(vd string, path string, p bool, registro bool) {
	var idDisco byte
	idDisco = vd[2]
	idDisco2 := idDisco - 97
	idP, _ := strconv.Atoi(vd[3:])
	idP--

	if arregloMount[idDisco2].estado == 1 {
		if arregloMount[idDisco2].discos[idP].estado == 1 {
			inicioPart := arregloMount[idDisco2].discos[idP].Partstart
			rutaDisco := arregloMount[idDisco2].Ruta
			superBloque := obtenerSB(rutaDisco, inicioPart)

			if path == "/" {
				fmt.Println("Se va a crear la carpeta Raiz")

				nuevoAVD := avd{}

				nuevoDD := dd{}
				for i := 0; i < 5; i++ {
					nuevoDD.DDarrayFiles[i].DDfileApInodo = -1
				}
				nuevoDD.DDapDD = -1

				posDD := superBloque.SBapDD + superBloque.SBfirstFreeBitDD*superBloque.SBsizeStructDD

				superBloque.SBddFree--
				actualizarValorBitmap(rutaDisco, superBloque.SBapBDD+superBloque.SBfirstFreeBitDD, '1')
				nuevoFFBDD := obtenerFirstFreeBit(rutaDisco, superBloque.SBapBDD, int(superBloque.SBddCount))

				superBloque.SBfirstFreeBitDD = nuevoFFBDD

				escribirStructDD(rutaDisco, posDD, nuevoDD)

				for j := 0; j < 6; j++ {
					nuevoAVD.AVDapArraySub[j] = -1
				}
				nuevoAVD.AVDapDetalleDir = posDD
				nuevoAVD.AVDapAVD = -1
				fecha := time.Now().Format("2006-01-02 15:04:05")
				copy(nuevoAVD.AVDfechaCreacion[:], fecha)
				copy(nuevoAVD.AVDnombreDirectorio[:], path)
				nuevoAVD.AVDproper = 1

				escribirStructAVD(rutaDisco, superBloque.SBapAVD, nuevoAVD)
				superBloque.SBavdFree--
				actualizarValorBitmap(rutaDisco, superBloque.SBapBAVD, '1')
				nuevoFFBAVD := obtenerFirstFreeBit(rutaDisco, superBloque.SBapBAVD, int(superBloque.SBsizeStructAVD))

				superBloque.SBfirstFreeBitAVD = nuevoFFBAVD
				escribirSuperBloque(rutaDisco, inicioPart, superBloque)
				var sizeBitacora int64 = int64(unsafe.Sizeof(bitacora{}))
				posCopiaSB := superBloque.SBapLOG + (superBloque.SBavdCount * sizeBitacora)
				escribirSuperBloque(rutaDisco, posCopiaSB, superBloque)

			} else {
				fmt.Println("Crear carpetas en raiz")
				path1 := strings.TrimPrefix(path, "/")
				path2 := strings.TrimSuffix(path1, "/")
				pathPart := strings.Split(path2, "/")

				encontrado := false
				posEncontrado := superBloque.SBapAVD
				t := 0

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
								superBloque = obtenerSB(rutaDisco, inicioPart)
								crearDir(rutaDisco, superBloque, path2, inicioPart)
							}
						} else {
							fmt.Println("No se pueden crear las carpetas padres, falta de parametro de permiso")
						}
					}
				}
			}
			if registro == true {
				bit := bitacora{}
				copy(bit.LOGtipoOperacion[:], "mkdir")
				bit.LOGtipo = '0'
				copy(bit.LOGnombre[:], path)
				fecha := time.Now().Format("2006-01-02 15:04:05")
				copy(bit.LOGfecha[:], fecha)
				var sizeBitacora int64 = int64(unsafe.Sizeof(bitacora{}))
				insertaBitacora(rutaDisco, superBloque.SBapLOG, bit, superBloque.SBavdCount, sizeBitacora)
			}

		} else {
			fmt.Println("La particion indicada no esta montada")
		}
	} else {
		fmt.Println("EL disco proporcionado no esta montado")
	}
}

func crearDir(rutaDisco string, superBloque sb, path2 string, inicioPart int64) int64 {
	var posAnterior int64 = 0
	var posNuevo int64 = 0
	posEncontrado := superBloque.SBapAVD
	encontrado := false
	pathPart := strings.Split(path2, "/")
	for i := 0; i < len(pathPart); i++ {
		posAnterior = posEncontrado
		encontrado, posEncontrado = buscarDir(posEncontrado, pathPart[i], rutaDisco)
		if encontrado == false {
			posEncontrado = nuevoAVDindirecto(posAnterior, rutaDisco, superBloque, inicioPart)
			posAnterior = posEncontrado
			superBloque = obtenerSB(rutaDisco, inicioPart)

			nuevoAVD := avd{}

			nuevoDD := dd{}
			for i := 0; i < 5; i++ {
				nuevoDD.DDarrayFiles[i].DDfileApInodo = -1
			}
			nuevoDD.DDapDD = -1

			posDD := superBloque.SBapDD + superBloque.SBfirstFreeBitDD*superBloque.SBsizeStructDD

			superBloque.SBddFree--
			actualizarValorBitmap(rutaDisco, superBloque.SBapBDD+superBloque.SBfirstFreeBitDD, '1')
			nuevoFFBDD := obtenerFirstFreeBit(rutaDisco, superBloque.SBapBDD, int(superBloque.SBddCount))

			superBloque.SBfirstFreeBitDD = nuevoFFBDD

			escribirStructDD(rutaDisco, posDD, nuevoDD)

			for j := 0; j < 6; j++ {
				nuevoAVD.AVDapArraySub[j] = -1
			}
			nuevoAVD.AVDapDetalleDir = posDD
			nuevoAVD.AVDapAVD = -1
			fecha := time.Now().Format("2006-01-02 15:04:05")
			copy(nuevoAVD.AVDfechaCreacion[:], fecha)
			copy(nuevoAVD.AVDnombreDirectorio[:], pathPart[i])
			nuevoAVD.AVDproper = 1

			posFirstFreeBit := superBloque.SBfirstFreeBitAVD

			posNuevoAVD := superBloque.SBapAVD + (posFirstFreeBit * superBloque.SBsizeStructAVD)
			posNuevo = posNuevoAVD
			avdPadre := obtenerAVD(rutaDisco, posAnterior)

			for j := 0; j < len(avdPadre.AVDapArraySub); j++ {
				if avdPadre.AVDapArraySub[j] == -1 {
					avdPadre.AVDapArraySub[j] = posNuevoAVD
					break
				}
			}

			escribirStructAVD(rutaDisco, posAnterior, avdPadre)
			escribirStructAVD(rutaDisco, posNuevoAVD, nuevoAVD)
			superBloque.SBavdFree--
			actualizarValorBitmap(rutaDisco, superBloque.SBapBAVD+posFirstFreeBit, '1')
			posFirstFreeBit = obtenerFirstFreeBit(rutaDisco, superBloque.SBapBAVD, int(superBloque.SBsizeStructAVD))
			superBloque.SBfirstFreeBitAVD = posFirstFreeBit

			escribirSuperBloque(rutaDisco, inicioPart, superBloque)
			var sizeBitacora int64 = int64(unsafe.Sizeof(bitacora{}))
			posCopiaSB := superBloque.SBapLOG + (superBloque.SBavdCount * sizeBitacora)
			escribirSuperBloque(rutaDisco, posCopiaSB, superBloque)

			fmt.Println("Se creo el directorio: " + pathPart[i])
			break
		}
	}
	return posNuevo

}

func nuevoAVDindirecto(pos int64, ruta string, superBloque sb, inicioPart int64) int64 {

	arbol := obtenerAVD(ruta, pos)
	var posEncontrado int64 = 0
	cantidad := 0
	for i := 0; i < len(arbol.AVDapArraySub); i++ {
		if arbol.AVDapArraySub[i] != -1 {
			cantidad++
		} else if arbol.AVDapArraySub[i] == -1 {
			posEncontrado = pos
			break
		}
	}

	if cantidad == 6 {
		if arbol.AVDapAVD != -1 {
			posEncontrado = nuevoAVDindirecto(arbol.AVDapAVD, ruta, superBloque, inicioPart)
		} else {
			nuevoAVD := avd{}

			nuevoDD := dd{}
			for i := 0; i < 5; i++ {
				nuevoDD.DDarrayFiles[i].DDfileApInodo = -1
			}
			nuevoDD.DDapDD = -1

			posDD := superBloque.SBapDD + superBloque.SBfirstFreeBitDD*superBloque.SBsizeStructDD

			superBloque.SBddFree--
			actualizarValorBitmap(ruta, superBloque.SBapBDD+superBloque.SBfirstFreeBitDD, '1')
			nuevoFFBDD := obtenerFirstFreeBit(ruta, superBloque.SBapBDD, int(superBloque.SBddCount))

			superBloque.SBfirstFreeBitDD = nuevoFFBDD

			escribirStructDD(ruta, posDD, nuevoDD)

			for j := 0; j < 6; j++ {
				nuevoAVD.AVDapArraySub[j] = -1
			}
			nuevoAVD.AVDapDetalleDir = posDD
			nuevoAVD.AVDapAVD = -1
			fecha := time.Now().Format("2006-01-02 15:04:05")
			copy(nuevoAVD.AVDfechaCreacion[:], fecha)
			nuevoAVD.AVDnombreDirectorio = arbol.AVDnombreDirectorio
			nuevoAVD.AVDproper = 1

			posFirstFreeBit := superBloque.SBfirstFreeBitAVD

			posNuevoAVD := superBloque.SBapAVD + (posFirstFreeBit * superBloque.SBsizeStructAVD)

			arbol.AVDapAVD = posNuevoAVD

			escribirStructAVD(ruta, pos, arbol)
			escribirStructAVD(ruta, posNuevoAVD, nuevoAVD)
			superBloque.SBavdFree--
			actualizarValorBitmap(ruta, superBloque.SBapBAVD+posFirstFreeBit, '1')
			posFirstFreeBit = obtenerFirstFreeBit(ruta, superBloque.SBapBAVD, int(superBloque.SBsizeStructAVD))
			superBloque.SBfirstFreeBitAVD = posFirstFreeBit

			escribirSuperBloque(ruta, inicioPart, superBloque)
			var sizeBitacora int64 = int64(unsafe.Sizeof(bitacora{}))
			posCopiaSB := superBloque.SBapLOG + (superBloque.SBavdCount * sizeBitacora)
			escribirSuperBloque(ruta, posCopiaSB, superBloque)

			tempNombre := ""

			for i := 0; i < len(arbol.AVDnombreDirectorio); i++ {
				if arbol.AVDnombreDirectorio[i] != 0 {
					tempNombre += string(arbol.AVDnombreDirectorio[i])
				}
			}

			posEncontrado = posNuevoAVD
		}
	}
	return posEncontrado
}

func buscarDir(pos int64, nombre string, ruta string) (bool, int64) {
	arbol := obtenerAVD(ruta, pos)
	encontrado := false
	var posEncontrado int64 = 0
	for i := 0; i < len(arbol.AVDapArraySub); i++ {

		if arbol.AVDapArraySub[i] != -1 {
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
		if car != 49 {
			pos = i
			break
		}
		posIni++
	}

	file.Close()
	return int64(pos)
}

func insertaBitacora(ruta string, pos int64, bi bitacora, cantidad int64, tambit int64) {
	for i := 0; i < int(cantidad); i++ {
		bit := obtenerbitacora(ruta, pos)
		if bit.LOGtipo == 'x' {
			break
		} else {
			pos = pos + tambit
		}
	}
	escribirStructBitacora(ruta, pos, bi)
}

/*
func listadobitmap(ruta string, posIni int64, tam int) {

	file, err := os.OpenFile(ruta, os.O_RDWR, 0777)
	if err != nil {
		log.Fatal(err)
	}

	var car byte
	var size int = int(unsafe.Sizeof(car))

	for i := 0; i < tam; i++ {
		file.Seek(posIni, 0)

		data := obtenerBytesEbr(file, size)
		buffer := bytes.NewBuffer(data)

		err = binary.Read(buffer, binary.BigEndian, &car)
		if err != nil {
			log.Fatal("binary.Read failed", err)
		}

		fmt.Print("valor i " + strconv.Itoa(i) + " ")
		fmt.Println(car)
		posIni++
	}

	file.Close()

}*/
