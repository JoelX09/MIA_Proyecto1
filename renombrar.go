package main

import (
	"fmt"
	"strconv"
	"strings"
)

func renombrar(vd string, path string, nombreNuevo string) {
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
						if encontrado == false {
							fmt.Println("No existe: " + pathPart[i])
							break
						}
					}
					if encontrado == true {
						carpetaPadre := obtenerAVD(rutaDisco, posEncontrado)

						archivoEncontrado, posRef, posDD := buscarArchivoRM(rutaDisco, carpetaPadre.AVDapDetalleDir, nombre)

						if archivoEncontrado == true {
							nuevoDD := obtenerDD(rutaDisco, posDD)
							var vacio [20]byte
							for i := 0; i < 5; i++ {
								if nuevoDD.DDarrayFiles[i].DDfileApInodo == posRef {
									nuevoDD.DDarrayFiles[i].DDfileNombre = vacio
									copy(nuevoDD.DDarrayFiles[i].DDfileNombre[:], nombreNuevo)
									break
								}
							}
							escribirStructDD(rutaDisco, posDD, nuevoDD)
							fmt.Println("***************\nEl nombre del archivo se ha modificado\n***************")
						} else {
							fmt.Println("---------------\nEl archivo dado no existe\n---------------")
						}
					} else {
						fmt.Println("---------------\nCarpetas de la ruta inexistentes\n---------------")
					}
				}
			} else {
				encontrado := false
				posEncontrado := superBloque.SBapAVD
				for i := 0; i < len(pathPart); i++ {
					encontrado, posEncontrado = buscarDir(posEncontrado, pathPart[i], rutaDisco)
					if encontrado == false {
						fmt.Println("No existe: " + pathPart[i])
						break
					}
				}
				if encontrado == true {
					cambiarNombreCarpeta(rutaDisco, posEncontrado, nombreNuevo)
					fmt.Println("***************\nEl nombre de la carpeta se ha modificado\n***************")
				} else {
					fmt.Println("---------------\nCarpeta inexistentes en la ruta\n---------------")
				}
			}
		}
	}
}

func cambiarNombreCarpeta(ruta string, pos int64, nombre string) {
	nuevoAVD := obtenerAVD(ruta, pos)
	var tempNom [20]byte
	nuevoAVD.AVDnombreDirectorio = tempNom
	copy(nuevoAVD.AVDnombreDirectorio[:], nombre)
	escribirStructAVD(ruta, pos, nuevoAVD)

	if nuevoAVD.AVDapAVD != -1 {
		cambiarNombreCarpeta(ruta, nuevoAVD.AVDapAVD, nombre)
	}
}
