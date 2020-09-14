package main

import (
	"fmt"
	"strconv"
	"strings"
)

func eliminarFileDir(vd string, path string, rf bool) {

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

			path, nombre, isFile := descomponer(path)
			path1 := strings.TrimPrefix(path, "/")
			path2 := strings.TrimSuffix(path1, "/")
			pathPart := strings.Split(path2, "/")

			if isFile == true {
				fmt.Println("**********************************************")
				fmt.Println("Carpetas donde esta el archivo")
				fmt.Println(pathPart)
				fmt.Println("Archivo a eliminar")
				fmt.Println(nombre)
				fmt.Println("\nEjecuacion pausada... Presione enter para continuar")
				fmt.Scanln()

				raiz := false
				if len(pathPart) == 1 {
					if pathPart[0] == "" {
						raiz = true
					}
				}
				if raiz == false {
					encontrado := false
					posEncontrado := superBloque.SBapAVD

					for i := 0; i < len(pathPart); i++ {
						encontrado, posEncontrado = buscarDir(posEncontrado, pathPart[i], rutaDisco)
					}

					if encontrado == true {
						carpetaPadre := obtenerAVD(rutaDisco, posEncontrado)

						archivoEncontrado, archivoPosInodo, posDD := buscarArchivoRM(rutaDisco, carpetaPadre.AVDapDetalleDir, nombre)

						if archivoEncontrado == true {
							borrarInodo(rutaDisco, archivoPosInodo, inicioPart)
							tempDD := obtenerDD(rutaDisco, posDD)
							for i := 0; i < 5; i++ {
								if tempDD.DDarrayFiles[i].DDfileApInodo == archivoPosInodo {
									nuevo := arregloArchivos{}
									nuevo.DDfileApInodo = -1
									tempDD.DDarrayFiles[i] = nuevo
									escribirStructDD(rutaDisco, posDD, tempDD)
								}
							}
						} else {
							fmt.Println("EL archivo no existe")
						}
					} else {
						fmt.Println("Carpetas inexistentes en la ruta proporcionada")
					}

				} else {
					avdraiz := obtenerAVD(rutaDisco, superBloque.SBapAVD)
					archivoEncontrado, archivoPosInodo, posDD := buscarArchivoRM(rutaDisco, avdraiz.AVDapDetalleDir, nombre)

					if archivoEncontrado == true {
						borrarInodo(rutaDisco, archivoPosInodo, inicioPart)
						tempDD := obtenerDD(rutaDisco, posDD)
						for i := 0; i < 5; i++ {
							if tempDD.DDarrayFiles[i].DDfileApInodo == archivoPosInodo {
								nuevo := arregloArchivos{}
								nuevo.DDfileApInodo = -1
								tempDD.DDarrayFiles[i] = nuevo
								escribirStructDD(rutaDisco, posDD, tempDD)
							}
						}
					} else {
						fmt.Println("EL archivo no existe")
					}
				}
			} else {
				fmt.Println("**********************************************")
				fmt.Println("Carpeta a eliminar - ES LA ULTIMA ")
				fmt.Println(pathPart)
				fmt.Println("\nEjecuacion pausada... Presione enter para continuar")
				fmt.Scanln()
				encontrado := false
				posEncontrado := superBloque.SBapAVD
				var posPadre int64
				for i := 0; i < len(pathPart); i++ {
					encontrado, posEncontrado, posPadre = buscarDirRM(posEncontrado, pathPart[i], rutaDisco)
				}
				if encontrado == true {
					borraAVD(rutaDisco, posEncontrado, inicioPart)
					avdraiz := obtenerAVD(rutaDisco, posPadre)
					fmt.Println("**********************************************")
					fmt.Println("Posicion de la carpeta a eliminar ")
					fmt.Println(posEncontrado)
					fmt.Println("Posicion de la carpeta padre ")
					fmt.Println(posPadre)
					fmt.Println("\nEjecuacion pausada... Presione enter para continuar")
					fmt.Scanln()
					for i := 0; i < 6; i++ {
						if avdraiz.AVDapArraySub[i] == posEncontrado {
							avdraiz.AVDapArraySub[i] = -1
						}
					}
					escribirStructAVD(rutaDisco, posPadre, avdraiz)
				} else {
					fmt.Println("Carpetas inexistentes en la ruta proporcionada")
				}
			}

		} else {
			fmt.Println("La particion indicada no esta montada")
		}
	} else {
		fmt.Println("EL disco proporcionado no esta montado")
	}
}

func borraAVD(ruta string, pos int64, iniciopart int64) {
	avdtemp := obtenerAVD(ruta, pos)
	superBloque := obtenerSB(ruta, iniciopart)

	nuevo := avd{}
	for i := 0; i < 6; i++ {
		nuevo.AVDapArraySub[i] = -1
	}
	nuevo.AVDapAVD = -1
	nuevo.AVDapDetalleDir = -1
	escribirStructAVD(ruta, pos, nuevo)
	valortemp := (pos - superBloque.SBapAVD) / superBloque.SBsizeStructAVD
	actualizarValorBitmap(ruta, (superBloque.SBapBAVD + valortemp), '0')
	superBloque.SBavdFree++
	superBloque.SBfirstFreeBitAVD = obtenerFirstFreeBit(ruta, superBloque.SBapBAVD, int(superBloque.SBddCount))
	escribirSuperBloque(ruta, iniciopart, superBloque)

	for i := 0; i < 6; i++ {
		if avdtemp.AVDapArraySub[i] != -1 {
			borraAVD(ruta, avdtemp.AVDapArraySub[i], iniciopart)
		}
	}
	if avdtemp.AVDapDetalleDir != -1 {
		borrarDD(ruta, avdtemp.AVDapDetalleDir, iniciopart)
	}
	if avdtemp.AVDapAVD != -1 {
		borraAVD(ruta, avdtemp.AVDapAVD, iniciopart)
	}

}

func borrarDD(ruta string, pos int64, iniciopart int64) {
	ddtemp := obtenerDD(ruta, pos)
	superBloque := obtenerSB(ruta, iniciopart)

	for i := 0; i < 5; i++ {
		if ddtemp.DDarrayFiles[i].DDfileApInodo != -1 {
			borrarInodo(ruta, ddtemp.DDarrayFiles[i].DDfileApInodo, iniciopart)
		}
	}

	nuevo := dd{}
	for i := 0; i < 5; i++ {
		nuevo.DDarrayFiles[i].DDfileApInodo = -1
	}
	nuevo.DDapDD = -1
	escribirStructDD(ruta, pos, nuevo)
	valortemp := (pos - superBloque.SBapDD) / superBloque.SBsizeStructDD
	actualizarValorBitmap(ruta, (superBloque.SBapBDD + valortemp), '0')
	superBloque.SBddFree++
	superBloque.SBfirstFreeBitDD = obtenerFirstFreeBit(ruta, superBloque.SBapBDD, int(superBloque.SBddCount))
	escribirSuperBloque(ruta, iniciopart, superBloque)

	if ddtemp.DDapDD != -1 {
		borrarDD(ruta, ddtemp.DDapDD, iniciopart)
	}

}

func borrarInodo(ruta string, pos int64, inicioPart int64) {
	inodoTemp := obtenerINODO(ruta, pos)
	superBloque := obtenerSB(ruta, inicioPart)

	for i := 0; i < 4; i++ {
		if inodoTemp.IarrayBloques[i] != -1 {
			nuevo := bloque{}
			escribirStructBLOQUE(ruta, inodoTemp.IarrayBloques[i], nuevo)
			valortemp := (inodoTemp.IarrayBloques[i] - superBloque.SBapBLOQUE) / superBloque.SBsizeStructBLOQUE
			actualizarValorBitmap(ruta, (superBloque.SBapBBLOQUE + valortemp), '0')
			superBloque.SBbloquesFree++
		}
	}

	siguiente := inodoTemp.IapIndirecto

	nuevoInodo := inodo{}

	for i := 0; i < 4; i++ {
		nuevoInodo.IarrayBloques[i] = -1
	}
	nuevoInodo.IapIndirecto = -1
	escribirStructINODO(ruta, pos, nuevoInodo)
	valortemp := (pos - superBloque.SBapINODO) / superBloque.SBsizeStructINODO
	actualizarValorBitmap(ruta, (superBloque.SBapBINODO + valortemp), '0')
	superBloque.SBinodosFree++

	superBloque.SBfirstFreeBitBLOQUE = obtenerFirstFreeBit(ruta, superBloque.SBapBBLOQUE, int(superBloque.SBbloquesCount))
	superBloque.SBfirstFreeBitINODO = obtenerFirstFreeBit(ruta, superBloque.SBapBINODO, int(superBloque.SBinodosCount))
	escribirSuperBloque(ruta, inicioPart, superBloque)

	if siguiente != -1 {
		borrarInodo(ruta, siguiente, inicioPart)
	}
}

func buscarArchivoRM(rutaDisco string, pos int64, nombre string) (bool, int64, int64) {
	nuevoDD := obtenerDD(rutaDisco, pos)
	var nombreb [20]byte
	copy(nombreb[:], nombre)
	var posInodo int64
	encontrado := false
	for i := 0; i < 5; i++ {
		if nuevoDD.DDarrayFiles[i].DDfileNombre == nombreb {
			posInodo = nuevoDD.DDarrayFiles[i].DDfileApInodo
			encontrado = true
			break
		}
	}

	if encontrado == false {
		if nuevoDD.DDapDD != -1 {
			encontrado, posInodo, pos = buscarArchivoRM(rutaDisco, nuevoDD.DDapDD, nombre)
		}
	}

	return encontrado, posInodo, pos
}

func buscarDirRM(pos int64, nombre string, ruta string) (bool, int64, int64) {
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
				//fmt.Println("Se encontro el directorio: " + nombre)
				break
			}
		}
	}
	if encontrado == false {
		if arbol.AVDapAVD != -1 {
			encontrado, posEncontrado, pos = buscarDirRM(arbol.AVDapAVD, nombre, ruta)
		}
	}
	return encontrado, posEncontrado, pos
}
