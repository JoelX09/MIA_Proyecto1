package main

import (
	"fmt"
	"strconv"
	"strings"
	"time"
)

func crearArvhi(vd string, path string, p bool, size int64, cont string) {
	var idDisco byte
	idDisco = vd[2]
	idDisco2 := idDisco - 97
	idP, _ := strconv.Atoi(vd[3:])
	idP--

	if arregloMount[idDisco2].estado == 1 {
		if arregloMount[idDisco2].discos[idP].estado == 1 {
			rutaAchivo, nombreArchivo := descomponer(path)
			/*fmt.Println("NOmbre obtenido al descomponer")
			fmt.Println(nombreArchivo)
			fmt.Println("\nEjecuacion pausada... Presione enter para continuar")
			fmt.Scanln()*/
			inicioPart := arregloMount[idDisco2].discos[idP].Partstart
			rutaDisco := arregloMount[idDisco2].Ruta
			superBloque := obtenerSB(rutaDisco, inicioPart)

			if path == "/users.txt" {

				//raiz := obtenerAVD(rutaDisco, superBloque.SBapAVD)

				//****************************************** Contenido **************************************************//
				fmt.Println("Se va a crear users.txt")
				contenido := "1,G,root      \n1,U,root      ,root      ,201403975 \n"
				contBytes := []byte(contenido)
				tamContenido := len(contBytes)

				cantidadBloques := tamContenido / 25
				cantidadBloqueD := tamContenido % 25
				if cantidadBloqueD != 0 {
					cantidadBloques++
				}

				//******************************************** DE AQUI **************************************************//
				codigoRepetido(rutaDisco, superBloque.SBapAVD, contenido, int64(tamContenido), inicioPart, "users.txt")

				/*posDD := raiz.AVDapDetalleDir
				nuevoDD := obtenerDD(rutaDisco, posDD)
				//nuevoDD := dd{}
				/*for i := 0; i < 5; i++ {
					nuevoDD.DDarrayFiles[i].DDfileApInodo = -1
				}*/
				//nuevoDD.DDapDD = -1

				//posDD := superBloque.SBapDD + superBloque.SBfirstFreeBitDD*superBloque.SBsizeStructDD

				//superBloque.SBddFree--
				//actualizarValorBitmap(rutaDisco, superBloque.SBapBDD+superBloque.SBfirstFreeBitDD, '1')
				//nuevoFFBDD := obtenerFirstFreeBit(rutaDisco, superBloque.SBapBDD, int(superBloque.SBddCount))

				//superBloque.SBfirstFreeBitDD = nuevoFFBDD

				/*nuevoArreglogArchivo := arregloArchivos{}

				nuevoInodo := inodo{}
				posNuevoInodo := superBloque.SBapINODO + superBloque.SBfirstFreeBitINODO*superBloque.SBsizeStructINODO

				superBloque.SBinodosFree--
				actualizarValorBitmap(rutaDisco, superBloque.SBapBINODO+superBloque.SBfirstFreeBitINODO, '1')
				nuevoFFBINODO := obtenerFirstFreeBit(rutaDisco, superBloque.SBapBINODO, int(superBloque.SBinodosCount))

				bInodoAnt := superBloque.SBfirstFreeBitINODO
				superBloque.SBfirstFreeBitINODO = nuevoFFBINODO

				copy(nuevoArreglogArchivo.DDfileNombre[:], "users.txt")
				nuevoArreglogArchivo.DDfileApInodo = posNuevoInodo
				fecha := time.Now().Format("2006-01-02 15:04:05")
				copy(nuevoArreglogArchivo.DDfileDateCreacion[:], fecha)
				fecha = time.Now().Format("2006-01-02 15:04:05")
				copy(nuevoArreglogArchivo.DDfileDateModificacion[:], fecha)
				for i := 0; i < 5; i++ {
					if nuevoDD.DDarrayFiles[i].DDfileApInodo == -1 {
						nuevoDD.DDarrayFiles[i] = nuevoArreglogArchivo
						break
					}
				}

				escribirStructDD(rutaDisco, posDD, nuevoDD)

				nuevoInodo.IcountInodo = bInodoAnt + 1
				nuevoInodo.IsizeArchivo = int64(tamContenido)
				nuevoInodo.IcountBloquesAsignados = int64(cantidadBloques)
				nuevoInodo.IidProper = 1
				nuevoInodo.IapIndirecto = -1
				for i := 0; i < 4; i++ {
					nuevoInodo.IarrayBloques[i] = -1
				}
				escribirStructINODO(rutaDisco, posNuevoInodo, nuevoInodo)

				//raiz.AVDapDetalleDir = posDD

				//escribirStructAVD(rutaDisco, superBloque.SBapAVD, raiz)

				escribirSuperBloque(rutaDisco, inicioPart, superBloque)
				superBloque = obtenerSB(rutaDisco, inicioPart)

				llenarNuevoArchivo(contenido, posNuevoInodo, rutaDisco, superBloque, inicioPart, true)*/

			} else {

				path1 := strings.TrimPrefix(rutaAchivo, "/")
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

				if encontrado == true {
					fmt.Println("******************************************************")
					fmt.Println("Todas las carpetas existen")
					fmt.Println("\nEjecuacion pausada... Presione enter para continuar")
					/*fmt.Scanln()
					fmt.Println(posEncontrado)*/

					raiz := obtenerAVD(rutaDisco, posEncontrado)
					posDDr := raiz.AVDapDetalleDir
					nuevoDD, posDD := verficarNuevoDD(rutaDisco, posDDr, inicioPart)

					nuevoArreglogArchivo := arregloArchivos{}

					nuevoInodo := inodo{}
					posNuevoInodo := superBloque.SBapINODO + superBloque.SBfirstFreeBitINODO*superBloque.SBsizeStructINODO

					superBloque.SBinodosFree--
					actualizarValorBitmap(rutaDisco, superBloque.SBapBINODO+superBloque.SBfirstFreeBitINODO, '1')
					nuevoFFBINODO := obtenerFirstFreeBit(rutaDisco, superBloque.SBapBINODO, int(superBloque.SBinodosCount))

					bInodoAnt := superBloque.SBfirstFreeBitINODO
					superBloque.SBfirstFreeBitINODO = nuevoFFBINODO

					/*fmt.Println("NOmbre a insertar obtenido al descomponer")
					fmt.Println(nombreArchivo)
					fmt.Println("\nEjecuacion pausada... Presione enter para continuar")
					fmt.Scanln()*/
					copy(nuevoArreglogArchivo.DDfileNombre[:], nombreArchivo)
					/*fmt.Println("NOmbre inserta")
					fmt.Println(string(nuevoArreglogArchivo.DDfileNombre[:]))
					fmt.Println("\nEjecuacion pausada... Presione enter para continuar")
					fmt.Scanln()*/
					nuevoArreglogArchivo.DDfileApInodo = posNuevoInodo
					fecha := time.Now().Format("2006-01-02 15:04:05")
					copy(nuevoArreglogArchivo.DDfileDateCreacion[:], fecha)
					fecha = time.Now().Format("2006-01-02 15:04:05")
					copy(nuevoArreglogArchivo.DDfileDateModificacion[:], fecha)
					fmt.Println("Pos arreglo y si esta libre")
					for i := 0; i < 5; i++ {
						if nuevoDD.DDarrayFiles[i].DDfileApInodo == -1 {
							fmt.Println("La posicion esta libre")
							fmt.Println(i)
							fmt.Println("Valor")
							fmt.Println(nuevoDD.DDarrayFiles[i].DDfileApInodo)
						} else {
							fmt.Println("La posicion esta ocupada")
							fmt.Println(i)
							fmt.Println("Valor")
							fmt.Println(nuevoDD.DDarrayFiles[i].DDfileApInodo)
						}
					}
					fmt.Println("=======================================================================================")
					fmt.Println("\nEjecuacion pausada... Presione enter para continuar")
					fmt.Scanln()
					for i := 0; i < 5; i++ {
						if nuevoDD.DDarrayFiles[i].DDfileApInodo == -1 {
							fmt.Println("Voy a insertar el archivo")
							fmt.Println(nombreArchivo)
							fmt.Println("EN la pos del arreglo de archivos")
							fmt.Println(i)
							fmt.Println("Caber tantos arreglos")
							fmt.Println(len(nuevoDD.DDarrayFiles))
							fmt.Println("\nEjecuacion pausada... Presione enter para continuar")
							fmt.Scanln()
							fmt.Println("\nEjecuacion pausada... Presione enter para continuar")
							fmt.Scanln()
							nuevoDD.DDarrayFiles[i] = nuevoArreglogArchivo
							break
						}
					}

					escribirStructDD(rutaDisco, posDD, nuevoDD)

					nuevoInodo.IcountInodo = bInodoAnt + 1
					nuevoInodo.IsizeArchivo = int64(size)
					cantidadBloques := size / 25
					cantidadBloqueD := size % 25
					if cantidadBloqueD != 0 {
						cantidadBloques++
					}
					nuevoInodo.IcountBloquesAsignados = int64(cantidadBloques)
					nuevoInodo.IidProper = 1
					nuevoInodo.IapIndirecto = -1
					for i := 0; i < 4; i++ {
						nuevoInodo.IarrayBloques[i] = -1
					}
					escribirStructINODO(rutaDisco, posNuevoInodo, nuevoInodo)

					//raiz.AVDapDetalleDir = posDD

					//escribirStructAVD(rutaDisco, superBloque.SBapAVD, raiz)

					escribirSuperBloque(rutaDisco, inicioPart, superBloque)
					superBloque = obtenerSB(rutaDisco, inicioPart)

					llenarNuevoArchivo(cont, posNuevoInodo, rutaDisco, superBloque, inicioPart, true)

				} else {
					fmt.Println("******************************************************")
					fmt.Println("Faltan carpetas")
					fmt.Println("\nEjecuacion pausada... Presione enter para continuar")
					fmt.Scanln()
					var posUltimaCarpeta int64
					if t == len(pathPart)-1 {
						posUltimaCarpeta = crearDir(rutaDisco, superBloque, path2, inicioPart)
						codigoRepetido(rutaDisco, posUltimaCarpeta, cont, size, inicioPart, nombreArchivo)
						/*raiz := obtenerAVD(rutaDisco, posUltimaCarpeta)
						posDD := raiz.AVDapDetalleDir
						nuevoDD := obtenerDD(rutaDisco, posDD)
						nuevoArreglogArchivo := arregloArchivos{}

						nuevoInodo := inodo{}
						posNuevoInodo := superBloque.SBapINODO + superBloque.SBfirstFreeBitINODO*superBloque.SBsizeStructINODO

						superBloque.SBinodosFree--
						actualizarValorBitmap(rutaDisco, superBloque.SBapBINODO+superBloque.SBfirstFreeBitINODO, '1')
						nuevoFFBINODO := obtenerFirstFreeBit(rutaDisco, superBloque.SBapBINODO, int(superBloque.SBinodosCount))

						bInodoAnt := superBloque.SBfirstFreeBitINODO
						superBloque.SBfirstFreeBitINODO = nuevoFFBINODO

						copy(nuevoArreglogArchivo.DDfileNombre[:], nombreArchivo)
						nuevoArreglogArchivo.DDfileApInodo = posNuevoInodo
						fecha := time.Now().Format("2006-01-02 15:04:05")
						copy(nuevoArreglogArchivo.DDfileDateCreacion[:], fecha)
						fecha = time.Now().Format("2006-01-02 15:04:05")
						copy(nuevoArreglogArchivo.DDfileDateModificacion[:], fecha)
						for i := 0; i < 5; i++ {
							if nuevoDD.DDarrayFiles[i].DDfileApInodo == -1 {
								nuevoDD.DDarrayFiles[i] = nuevoArreglogArchivo
								break
							}
						}

						escribirStructDD(rutaDisco, posDD, nuevoDD)

						nuevoInodo.IcountInodo = bInodoAnt + 1
						nuevoInodo.IsizeArchivo = int64(size)
						cantidadBloques := size / 25
						cantidadBloqueD := size % 25
						if cantidadBloqueD != 0 {
							cantidadBloques++
						}
						nuevoInodo.IcountBloquesAsignados = int64(cantidadBloques)
						nuevoInodo.IidProper = 1
						nuevoInodo.IapIndirecto = -1
						for i := 0; i < 4; i++ {
							nuevoInodo.IarrayBloques[i] = -1
						}
						escribirStructINODO(rutaDisco, posNuevoInodo, nuevoInodo)

						//raiz.AVDapDetalleDir = posDD

						//escribirStructAVD(rutaDisco, superBloque.SBapAVD, raiz)

						escribirSuperBloque(rutaDisco, inicioPart, superBloque)
						superBloque = obtenerSB(rutaDisco, inicioPart)

						llenarNuevoArchivo(cont, posNuevoInodo, rutaDisco, superBloque, inicioPart)*/

					} else {
						if p == true {

							for i := 0; i < len(pathPart); i++ {
								superBloque = obtenerSB(rutaDisco, inicioPart)
								posUltimaCarpeta = crearDir(rutaDisco, superBloque, path2, inicioPart)
								fmt.Println("Nuevo archivo en dd")
								fmt.Println(posUltimaCarpeta)
							}
							codigoRepetido(rutaDisco, posUltimaCarpeta, cont, size, inicioPart, nombreArchivo)
							/*raiz := obtenerAVD(rutaDisco, posUltimaCarpeta)
							posDD := raiz.AVDapDetalleDir
							nuevoDD := obtenerDD(rutaDisco, posDD)
							nuevoArreglogArchivo := arregloArchivos{}

							nuevoInodo := inodo{}
							posNuevoInodo := superBloque.SBapINODO + superBloque.SBfirstFreeBitINODO*superBloque.SBsizeStructINODO

							superBloque.SBinodosFree--
							actualizarValorBitmap(rutaDisco, superBloque.SBapBINODO+superBloque.SBfirstFreeBitINODO, '1')
							nuevoFFBINODO := obtenerFirstFreeBit(rutaDisco, superBloque.SBapBINODO, int(superBloque.SBinodosCount))

							bInodoAnt := superBloque.SBfirstFreeBitINODO
							superBloque.SBfirstFreeBitINODO = nuevoFFBINODO

							copy(nuevoArreglogArchivo.DDfileNombre[:], nombreArchivo)
							nuevoArreglogArchivo.DDfileApInodo = posNuevoInodo
							fecha := time.Now().Format("2006-01-02 15:04:05")
							copy(nuevoArreglogArchivo.DDfileDateCreacion[:], fecha)
							fecha = time.Now().Format("2006-01-02 15:04:05")
							copy(nuevoArreglogArchivo.DDfileDateModificacion[:], fecha)
							for i := 0; i < 5; i++ {
								if nuevoDD.DDarrayFiles[i].DDfileApInodo == -1 {
									nuevoDD.DDarrayFiles[i] = nuevoArreglogArchivo
									break
								}
							}

							escribirStructDD(rutaDisco, posDD, nuevoDD)

							nuevoInodo.IcountInodo = bInodoAnt + 1
							nuevoInodo.IsizeArchivo = int64(size)
							cantidadBloques := size / 25
							cantidadBloqueD := size % 25
							if cantidadBloqueD != 0 {
								cantidadBloques++
							}
							nuevoInodo.IcountBloquesAsignados = int64(cantidadBloques)
							nuevoInodo.IidProper = 1
							nuevoInodo.IapIndirecto = -1
							for i := 0; i < 4; i++ {
								nuevoInodo.IarrayBloques[i] = -1
							}
							escribirStructINODO(rutaDisco, posNuevoInodo, nuevoInodo)

							//raiz.AVDapDetalleDir = posDD

							//escribirStructAVD(rutaDisco, superBloque.SBapAVD, raiz)

							escribirSuperBloque(rutaDisco, inicioPart, superBloque)
							superBloque = obtenerSB(rutaDisco, inicioPart)

							llenarNuevoArchivo(cont, posNuevoInodo, rutaDisco, superBloque, inicioPart)*/

						} else {
							fmt.Println("No se pueden crear las carpetas padres, falta de parametro de permiso")
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

func verficarNuevoDD(rutaDisco string, posDD int64, pos int64) (dd, int64) {
	nuevoDD := obtenerDD(rutaDisco, posDD)
	cantidad := 0

	for i := 0; i < len(nuevoDD.DDarrayFiles); i++ {
		if nuevoDD.DDarrayFiles[i].DDfileApInodo != -1 {
			cantidad++
		}
	}

	if cantidad == 5 {
		if nuevoDD.DDapDD != -1 {
			nuevoDD, posDD = verficarNuevoDD(rutaDisco, nuevoDD.DDapDD, pos)
		} else {
			superBloque := obtenerSB(rutaDisco, pos)
			nuevo := dd{}

			for i := 0; i < 5; i++ {
				nuevo.DDarrayFiles[i].DDfileApInodo = -1
			}
			nuevo.DDapDD = -1

			posNuevo := superBloque.SBapDD + superBloque.SBfirstFreeBitDD*superBloque.SBsizeStructDD

			superBloque.SBddFree--
			actualizarValorBitmap(rutaDisco, superBloque.SBapBDD+superBloque.SBfirstFreeBitDD, '1')
			nuevoFFBDD := obtenerFirstFreeBit(rutaDisco, superBloque.SBapBDD, int(superBloque.SBddCount))

			superBloque.SBfirstFreeBitDD = nuevoFFBDD

			escribirStructDD(rutaDisco, posNuevo, nuevo)

			nuevoDD = nuevo
			posDD = posNuevo
		}

	}
	return nuevoDD, posDD
}

func verficarNuevoI(rutaDisco string, posI int64, pos int64) (inodo, int64) {
	nuevoI := obtenerINODO(rutaDisco, posI)
	cantidad := 0

	for i := 0; i < len(nuevoI.IarrayBloques); i++ {
		if nuevoI.IarrayBloques[i] != -1 {
			cantidad++
		}
	}

	if cantidad == 3 {
		if nuevoI.IapIndirecto != -1 {
			nuevoI, posI = verficarNuevoI(rutaDisco, nuevoI.IapIndirecto, pos)
		} else {
			superBloque := obtenerSB(rutaDisco, pos)
			nuevo := inodo{}

			for i := 0; i < 4; i++ {
				nuevo.IarrayBloques[i] = -1
			}
			nuevo.IapIndirecto = -1

			posNuevo := superBloque.SBapINODO + superBloque.SBfirstFreeBitINODO*superBloque.SBsizeStructINODO

			superBloque.SBinodosFree--
			actualizarValorBitmap(rutaDisco, superBloque.SBapBINODO+superBloque.SBfirstFreeBitINODO, '1')
			nuevoFFBINODO := obtenerFirstFreeBit(rutaDisco, superBloque.SBapBINODO, int(superBloque.SBinodosCount))

			superBloque.SBfirstFreeBitDD = nuevoFFBINODO

			escribirStructINODO(rutaDisco, posNuevo, nuevo)

			nuevoI = nuevo
			posI = posNuevo
		}

	}
	return inodo{}, posI
}

func llenarNuevoArchivo(cont string, posInodo int64, rutaDisco string, superBloque sb, inicioPart int64, existeContenido bool) {
	contB := []byte(cont)
	termino := false
	//****************************** AQUI PTM**************
	//nuevoInodo, _ := verficarNuevoI(rutaDisco, posInodo, inicioPart)
	nuevoInodo := obtenerINODO(rutaDisco, posInodo)

	pos := 0
	for i := 0; i < 4; i++ {
		var letra byte = 'A'
		nuevoBloque := bloque{}
		for j := 0; j < 25; j++ {
			if j < len(contB) {
				nuevoBloque.DBdata[j] = contB[j]
			} else {
				termino = true
				nuevoBloque.DBdata[j] = letra
				letra++
			}
		}

		if termino == true {

		}

		posNuevoBloque := superBloque.SBapBLOQUE + superBloque.SBfirstFreeBitBLOQUE*superBloque.SBsizeStructBLOQUE

		escribirStructBLOQUE(rutaDisco, posNuevoBloque, nuevoBloque)
		superBloque.SBbloquesFree--
		actualizarValorBitmap(rutaDisco, superBloque.SBapBBLOQUE+superBloque.SBfirstFreeBitBLOQUE, '1')
		nuevoFFBBLOQUE := obtenerFirstFreeBit(rutaDisco, superBloque.SBapBBLOQUE, int(superBloque.SBbloquesCount))

		superBloque.SBfirstFreeBitBLOQUE = nuevoFFBBLOQUE
		escribirSuperBloque(rutaDisco, inicioPart, superBloque)
		superBloque = obtenerSB(rutaDisco, inicioPart)

		nuevoInodo.IarrayBloques[i] = posNuevoBloque

		pos = i
		if termino == true {
			escribirStructINODO(rutaDisco, posInodo, nuevoInodo)
			superBloque = obtenerSB(rutaDisco, inicioPart)
			break
		} else if len(contB) == 25 {
			termino = true
			escribirStructINODO(rutaDisco, posInodo, nuevoInodo)
			superBloque = obtenerSB(rutaDisco, inicioPart)
			break
		} else {
			contB = contB[25:]
		}

	}

	if termino == false {
		if pos == 3 {
			nuevoInodoIndirecto := inodo{}
			posNuevoInodoIndirecto := superBloque.SBapINODO + superBloque.SBfirstFreeBitINODO*superBloque.SBsizeStructINODO

			superBloque.SBinodosFree--
			actualizarValorBitmap(rutaDisco, superBloque.SBapBINODO+superBloque.SBfirstFreeBitINODO, '1')
			nuevoFFBINODO := obtenerFirstFreeBit(rutaDisco, superBloque.SBapBINODO, int(superBloque.SBinodosCount)) //Error por el count
			bInodoAnt := superBloque.SBfirstFreeBitINODO
			superBloque.SBfirstFreeBitINODO = nuevoFFBINODO

			nuevoInodoIndirecto.IcountInodo = bInodoAnt + 1
			nuevoInodoIndirecto.IsizeArchivo = nuevoInodo.IsizeArchivo
			nuevoInodoIndirecto.IcountBloquesAsignados = nuevoInodo.IcountBloquesAsignados
			nuevoInodoIndirecto.IidProper = 1
			nuevoInodoIndirecto.IapIndirecto = -1
			for i := 0; i < 4; i++ {
				nuevoInodoIndirecto.IarrayBloques[i] = -1
			}

			nuevoInodo.IapIndirecto = posNuevoInodoIndirecto

			escribirStructINODO(rutaDisco, posInodo, nuevoInodo)

			escribirStructINODO(rutaDisco, posNuevoInodoIndirecto, nuevoInodoIndirecto)
			escribirSuperBloque(rutaDisco, inicioPart, superBloque)
			superBloque = obtenerSB(rutaDisco, inicioPart)

			contres := ""
			for i := 0; i < len(contB); i++ {
				contres += string(contB[i])
			}

			llenarNuevoArchivo(contres, posNuevoInodoIndirecto, rutaDisco, superBloque, inicioPart, true)
		}
	}
}

func codigoRepetido(rutaDisco string, posUltimaCarpeta int64, cont string, size int64, inicioPart int64, nombreArchivo string) {
	superBloque := obtenerSB(rutaDisco, inicioPart)
	raiz := obtenerAVD(rutaDisco, posUltimaCarpeta)
	posDD := raiz.AVDapDetalleDir
	nuevoDD := obtenerDD(rutaDisco, posDD)
	nuevoArreglogArchivo := arregloArchivos{}

	nuevoInodo := inodo{}
	posNuevoInodo := superBloque.SBapINODO + superBloque.SBfirstFreeBitINODO*superBloque.SBsizeStructINODO

	superBloque.SBinodosFree--
	actualizarValorBitmap(rutaDisco, superBloque.SBapBINODO+superBloque.SBfirstFreeBitINODO, '1')
	nuevoFFBINODO := obtenerFirstFreeBit(rutaDisco, superBloque.SBapBINODO, int(superBloque.SBinodosCount))

	bInodoAnt := superBloque.SBfirstFreeBitINODO
	superBloque.SBfirstFreeBitINODO = nuevoFFBINODO

	copy(nuevoArreglogArchivo.DDfileNombre[:], nombreArchivo)
	nuevoArreglogArchivo.DDfileApInodo = posNuevoInodo
	fecha := time.Now().Format("2006-01-02 15:04:05")
	copy(nuevoArreglogArchivo.DDfileDateCreacion[:], fecha)
	fecha = time.Now().Format("2006-01-02 15:04:05")
	copy(nuevoArreglogArchivo.DDfileDateModificacion[:], fecha)
	for i := 0; i < 5; i++ {
		if nuevoDD.DDarrayFiles[i].DDfileApInodo == -1 {
			fmt.Println("Voy a insertar el archivo")
			fmt.Println(nombreArchivo)
			fmt.Println("EN la pos del arreglo de archivos")
			fmt.Println(i)
			fmt.Println("Caber tantos arreglos")
			fmt.Println(len(nuevoDD.DDarrayFiles))
			fmt.Println("\nEjecuacion pausada... Presione enter para continuar")
			fmt.Scanln()
			fmt.Println("\nEjecuacion pausada... Presione enter para continuar")
			fmt.Scanln()
			nuevoDD.DDarrayFiles[i] = nuevoArreglogArchivo
			break
		}
	}

	//agregarDDIndirecto()

	escribirStructDD(rutaDisco, posDD, nuevoDD)

	nuevoInodo.IcountInodo = bInodoAnt + 1
	nuevoInodo.IsizeArchivo = size
	cantidadBloques := size / 25
	cantidadBloqueD := size % 25
	if cantidadBloqueD != 0 {
		cantidadBloques++
	}
	nuevoInodo.IcountBloquesAsignados = int64(cantidadBloques)
	nuevoInodo.IidProper = 1
	nuevoInodo.IapIndirecto = -1
	for i := 0; i < 4; i++ {
		nuevoInodo.IarrayBloques[i] = -1
	}
	escribirStructINODO(rutaDisco, posNuevoInodo, nuevoInodo)

	//raiz.AVDapDetalleDir = posDD

	//escribirStructAVD(rutaDisco, superBloque.SBapAVD, raiz)

	escribirSuperBloque(rutaDisco, inicioPart, superBloque)
	superBloque = obtenerSB(rutaDisco, inicioPart)

	llenarNuevoArchivo(cont, posNuevoInodo, rutaDisco, superBloque, inicioPart, true)

}

func agregarDDIndirecto(rutaDisco string, posDD int64) {
	nuevoDD := obtenerDD(rutaDisco, posDD)
	cantidad := 0
	for i := 0; i < 5; i++ {
		if nuevoDD.DDarrayFiles[i].DDfileApInodo != -1 {
			cantidad++
		}
	}
	//if cantidad==5
}
