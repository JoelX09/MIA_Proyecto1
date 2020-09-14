package main

import (
	"fmt"
	"strconv"
	"strings"
)

func cat(vd string, path string) {
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

			for ele := listaFile.Front(); ele != nil; ele = ele.Next() {
				path = ele.Value.(string)
				path, nombre, _ := descomponer(path)
				path1 := strings.TrimPrefix(path, "/")
				path2 := strings.TrimSuffix(path1, "/")
				pathPart := strings.Split(path2, "/")

				raiz := false
				if len(pathPart) == 1 {
					if pathPart[0] == "" {
						raiz = true
					}
				}
				if raiz == false {
					encontrado := false
					posEncontrado := superBloque.SBapAVD
					/*fmt.Println("EMpezre analisis de carpetas")
					fmt.Println("ELmentos")
					fmt.Println(pathPart)
					fmt.Println("longitud")
					fmt.Println(len(pathPart))
					fmt.Println("\nEjecuacion pausada... Presione enter para continuar")
					fmt.Scanln()*/
					for i := 0; i < len(pathPart); i++ {
						encontrado, posEncontrado = buscarDir(posEncontrado, pathPart[i], rutaDisco)
						if encontrado == false {
							fmt.Println("No existe: " + pathPart[i])
							break
						}
					}
					/*fmt.Println("termine")
					fmt.Println("\nEjecuacion pausada... Presione enter para continuar")
					fmt.Scanln()*/
					if encontrado == true {
						carpetaPadre := obtenerAVD(rutaDisco, posEncontrado)

						archivoEncontrado, posRef, _ := buscarArchivoRM(rutaDisco, carpetaPadre.AVDapDetalleDir, nombre)

						if archivoEncontrado == true {
							imprimirArchivo(posRef, rutaDisco)
							fmt.Println("")
						} else {
							fmt.Println("---------------\nEl archivo dado no existe\n---------------")
						}
					} else {
						fmt.Println("---------------\nCarpetas de la ruta inexistentes\n---------------")
					}
				}
			}

		}
	}
}

func imprimirArchivo(posInodo int64, ruta string) {
	nuevo := obtenerINODO(ruta, posInodo)
	for i := 0; i < 4; i++ {
		if nuevo.IarrayBloques[i] != -1 {
			nb := obtenerBLOQUE(ruta, nuevo.IarrayBloques[i])
			for j := 0; j < 25; j++ {
				if nb.DBdata[j] != 0 {
					fmt.Print(string(nb.DBdata[j]))
				}
			}
		}
	}
	if nuevo.IapIndirecto != -1 {
		imprimirArchivo(nuevo.IapIndirecto, ruta)
	}
}
