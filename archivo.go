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

func crearArvhi(vd string, path string, p bool, size int64, cont string, registro bool) {
	var idDisco byte
	idDisco = vd[2]
	idDisco2 := idDisco - 97
	idP, _ := strconv.Atoi(vd[3:])
	idP--

	if arregloMount[idDisco2].estado == 1 {
		if arregloMount[idDisco2].discos[idP].estado == 1 {
			rutaAchivo, nombreArchivo, _ := descomponer(path)
			inicioPart := arregloMount[idDisco2].discos[idP].Partstart
			rutaDisco := arregloMount[idDisco2].Ruta
			superBloque := obtenerSB(rutaDisco, inicioPart)

			if path == "/users.txt" {

				//****************************************** Contenido **************************************************//
				fmt.Println("Se va a crear users.txt")
				contenido := "1,G,root      \n1,U,root      ,root      ,201403975 \n"
				cont = contenido
				contBytes := []byte(contenido)
				tamContenido := len(contBytes)
				size = int64(tamContenido)

				cantidadBloques := tamContenido / 25
				cantidadBloqueD := tamContenido % 25
				if cantidadBloqueD != 0 {
					cantidadBloques++
				}

				//******************************************** DE AQUI **************************************************//
				codigoRepetido(rutaDisco, superBloque.SBapAVD, contenido, int64(tamContenido), inicioPart, "users.txt")

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
					/*fmt.Println("******************************************************")
					fmt.Println("Todas las carpetas existen")
					fmt.Println("\nEjecuacion pausada... Presione enter para continuar")
					fmt.Scanln()*/
					//fmt.Println(posEncontrado)

					raiz := obtenerAVD(rutaDisco, posEncontrado)
					posDDr := raiz.AVDapDetalleDir
					nuevoDD, posDD := verficarNuevoDD(rutaDisco, posDDr, inicioPart)
					superBloque = obtenerSB(rutaDisco, inicioPart)

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

					escribirSuperBloque(rutaDisco, inicioPart, superBloque)
					var sizeBitacora int64 = int64(unsafe.Sizeof(bitacora{}))
					posCopiaSB := superBloque.SBapLOG + (superBloque.SBavdCount * sizeBitacora)
					escribirSuperBloque(rutaDisco, posCopiaSB, superBloque)
					superBloque = obtenerSB(rutaDisco, inicioPart)

					llenarNuevoArchivo(cont, posNuevoInodo, rutaDisco, superBloque, inicioPart, true, size)

				} else {
					var posUltimaCarpeta int64
					if t == len(pathPart)-1 {
						posUltimaCarpeta = crearDir(rutaDisco, superBloque, path2, inicioPart)
						codigoRepetido(rutaDisco, posUltimaCarpeta, cont, size, inicioPart, nombreArchivo)

					} else {
						if p == true {

							for i := 0; i < len(pathPart); i++ {
								superBloque = obtenerSB(rutaDisco, inicioPart)
								posUltimaCarpeta = crearDir(rutaDisco, superBloque, path2, inicioPart)
								fmt.Println("Nuevo archivo en dd")
								fmt.Println(posUltimaCarpeta)
							}
							codigoRepetido(rutaDisco, posUltimaCarpeta, cont, size, inicioPart, nombreArchivo)

						} else {
							fmt.Println("No se pueden crear las carpetas padres, falta de parametro de permiso")
						}
					}

				}
			}
			if registro == true {
				bit := bitacora{}
				copy(bit.LOGtipoOperacion[:], "mkfile")
				bit.LOGtipo = '1'
				copy(bit.LOGnombre[:], path)
				copy(bit.LOGcontenido[:], cont)
				fecha := time.Now().Format("2006-01-02 15:04:05")
				copy(bit.LOGfecha[:], fecha)
				bit.LOGsize = size
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
			nuevoDD.DDapDD = posNuevo
			escribirStructDD(rutaDisco, posDD, nuevoDD)
			escribirStructDD(rutaDisco, posNuevo, nuevo)
			escribirSuperBloque(rutaDisco, pos, superBloque)
			var sizeBitacora int64 = int64(unsafe.Sizeof(bitacora{}))
			posCopiaSB := superBloque.SBapLOG + (superBloque.SBavdCount * sizeBitacora)
			escribirSuperBloque(rutaDisco, posCopiaSB, superBloque)

			nuevoDD = nuevo
			posDD = posNuevo
		}

	}
	return nuevoDD, posDD
}

func llenarNuevoArchivo(cont string, posInodo int64, rutaDisco string, superBloque sb, inicioPart int64, existeContenido bool, size int64) {
	contB := []byte(cont)
	//tamContenido := len(contB)
	/*fmt.Println("Size recibido")
	fmt.Println(size)
	fmt.Println("**************************************")
	fmt.Println("\nEjecuacion pausada... Presione enter para continuar")
	fmt.Scanln()*/
	cantidadBloques := int(size) / 25
	cantidadBloqueD := int(size) % 25
	if cantidadBloqueD != 0 {
		cantidadBloques++
	}
	/*fmt.Println("Cantidad de bloque a llenar")
	fmt.Println(cantidadBloques)
	fmt.Println("**************************************")
	fmt.Println("\nEjecuacion pausada... Presione enter para continuar")
	fmt.Scanln()*/
	//termino := false
	//****************************** AQUI PTM**************
	//nuevoInodo, _ := verficarNuevoI(rutaDisco, posInodo, inicioPart)
	nuevoInodo := obtenerINODO(rutaDisco, posInodo)
	//for a := 0; a < cantidadBloques; a++ {
	pos := 0
	for i := 0; i < 4; i++ {
		var letra byte = 'A'
		nuevoBloque := bloque{}
		for j := 0; j < 25; j++ {
			if existeContenido == true {
				if j < len(contB) {
					nuevoBloque.DBdata[j] = contB[j]
				} else {
					existeContenido = false
					//termino = true
					nuevoBloque.DBdata[j] = letra
					letra++
				}
			} else {
				nuevoBloque.DBdata[j] = letra
				letra++
			}

		}
		posNuevoBloque := superBloque.SBapBLOQUE + superBloque.SBfirstFreeBitBLOQUE*superBloque.SBsizeStructBLOQUE

		escribirStructBLOQUE(rutaDisco, posNuevoBloque, nuevoBloque)
		superBloque.SBbloquesFree--
		actualizarValorBitmap(rutaDisco, superBloque.SBapBBLOQUE+superBloque.SBfirstFreeBitBLOQUE, '1')
		nuevoFFBBLOQUE := obtenerFirstFreeBit(rutaDisco, superBloque.SBapBBLOQUE, int(superBloque.SBbloquesCount))

		superBloque.SBfirstFreeBitBLOQUE = nuevoFFBBLOQUE
		escribirSuperBloque(rutaDisco, inicioPart, superBloque)
		var sizeBitacora int64 = int64(unsafe.Sizeof(bitacora{}))
		posCopiaSB := superBloque.SBapLOG + (superBloque.SBavdCount * sizeBitacora)
		escribirSuperBloque(rutaDisco, posCopiaSB, superBloque)
		superBloque = obtenerSB(rutaDisco, inicioPart)

		nuevoInodo.IarrayBloques[i] = posNuevoBloque
		cantidadBloques--
		size -= 25
		pos = i
		//if termino == true {
		if cantidadBloques == 0 {
			nuevoInodo.IcountBloquesAsignados = int64(i + 1)
			escribirStructINODO(rutaDisco, posInodo, nuevoInodo)
			superBloque = obtenerSB(rutaDisco, inicioPart)
			break
			/*} else if len(contB) == 25 {
			termino = true
			escribirStructINODO(rutaDisco, posInodo, nuevoInodo)
			superBloque = obtenerSB(rutaDisco, inicioPart)
			break*/
		} else {
			if existeContenido == true {
				contB = contB[25:]
			}
		}

	}

	if cantidadBloques > 0 {
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
			nuevoInodo.IcountBloquesAsignados = 4
			escribirStructINODO(rutaDisco, posInodo, nuevoInodo)

			escribirStructINODO(rutaDisco, posNuevoInodoIndirecto, nuevoInodoIndirecto)
			escribirSuperBloque(rutaDisco, inicioPart, superBloque)
			var sizeBitacora int64 = int64(unsafe.Sizeof(bitacora{}))
			posCopiaSB := superBloque.SBapLOG + (superBloque.SBavdCount * sizeBitacora)
			escribirSuperBloque(rutaDisco, posCopiaSB, superBloque)
			superBloque = obtenerSB(rutaDisco, inicioPart)

			contres := ""
			if existeContenido == true {
				for i := 0; i < len(contB); i++ {
					contres += string(contB[i])
				}
			}

			llenarNuevoArchivo(contres, posNuevoInodoIndirecto, rutaDisco, superBloque, inicioPart, existeContenido, size)
		}
	}
	//}
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
			nuevoDD.DDarrayFiles[i] = nuevoArreglogArchivo
			break
		}
	}

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

	escribirSuperBloque(rutaDisco, inicioPart, superBloque)
	var sizeBitacora int64 = int64(unsafe.Sizeof(bitacora{}))
	posCopiaSB := superBloque.SBapLOG + (superBloque.SBavdCount * sizeBitacora)
	escribirSuperBloque(rutaDisco, posCopiaSB, superBloque)
	superBloque = obtenerSB(rutaDisco, inicioPart)

	llenarNuevoArchivo(cont, posNuevoInodo, rutaDisco, superBloque, inicioPart, true, size)

}

func obtenerBitmap(path string, inicio int64, tam int64) []byte {

	f, err := os.Open(path)
	//defer f.Close()
	if err != nil {
		log.Fatal(err)
	}
	temptam := make([]byte, tam)
	var size int = len(temptam)
	f.Seek(inicio, 0)

	data := obtenerBytes(f, size)
	buffer := bytes.NewBuffer(data)

	err = binary.Read(buffer, binary.BigEndian, &temptam)
	if err != nil {
		log.Fatal("bitmap.Read failed", err)
	}
	f.Close()
	return temptam
}
