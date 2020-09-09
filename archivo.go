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
			rutaAchivo, _ := descomponer(path)

			inicioPart := arregloMount[idDisco2].discos[idP].Partstart
			//tamPart := arregloMount[idDisco2].discos[idP].Partsize
			rutaDisco := arregloMount[idDisco2].Ruta
			superBloque := obtenerSB(rutaDisco, inicioPart)

			if path == "/users.txt" {
				/*fmt.Println("Creando el archivo USERS.TXT")
				fmt.Println("\nEjecuacion pausada... Presione enter para continuar")
				fmt.Scanln()*/
				raiz := obtenerAVD(rutaDisco, superBloque.SBapAVD)
				/*fmt.Println("En el avd raiz")
				fmt.Println(raiz)
				fmt.Println("\nEjecuacion pausada... Presione enter para continuar")
				fmt.Scanln()*/
				//****************************************** Contenido **************************************************//
				fmt.Println("Se va a crear users.txt")
				contenido := "1,G,root      \n1,U,root      ,root      ,201403975 \n"
				contBytes := []byte(contenido)
				tamContenido := len(contBytes)

				/*fmt.Println("Contenido que va a tener")
				fmt.Println(contenido)
				fmt.Println("Tamano")
				fmt.Println(tamContenido)
				fmt.Println("\nEjecuacion pausada... Presione enter para continuar")
				fmt.Scanln()
				fmt.Println("contenido BYTES")
				fmt.Println(contBytes)
				fmt.Println("Contenido desde bytes")
				fmt.Println(string(contBytes[:]))
				fmt.Println("\nEjecuacion pausada... Presione enter para continuar")
				fmt.Scanln()*/

				cantidadBloques := tamContenido / 25
				cantidadBloqueD := tamContenido % 25
				if cantidadBloqueD != 0 {
					cantidadBloques++
				}
				/*fmt.Println("Cantidad de bloques")
				fmt.Println(cantidadBloques)
				fmt.Println("\nEjecuacion pausada... Presione enter para continuar")
				fmt.Scanln()*/

				//******************************************** DE AQUI **************************************************//
				nuevoDD := dd{}
				for i := 0; i < 5; i++ {
					nuevoDD.DDarrayFiles[i].DDfileApInodo = -1
				}
				nuevoDD.DDapDD = -1
				/*fmt.Println("COntenido en el DD NUEVO")
				fmt.Println(nuevoDD)*/
				posDD := superBloque.SBapDD + superBloque.SBfirstFreeBitDD*superBloque.SBsizeStructDD
				/*fmt.Println("Posicion para insertar el nuevo DD")
				fmt.Println(posDD)
				fmt.Println("\nEjecuacion pausada... Presione enter para continuar")
				fmt.Scanln()*/
				superBloque.SBddFree--
				actualizarValorBitmap(rutaDisco, superBloque.SBapBDD+superBloque.SBfirstFreeBitDD, '1')
				nuevoFFBDD := obtenerFirstFreeBit(rutaDisco, superBloque.SBapBDD, int(superBloque.SBddCount))
				/*listadobitmap(rutaDisco, superBloque.SBapBDD, int(superBloque.SBsizeStructDD)) // ESTA MALO <-------------
				fmt.Println("\nEjecuacion pausada... Presione enter para continuar")
				fmt.Scanln()*/
				superBloque.SBfirstFreeBitDD = nuevoFFBDD
				/*fmt.Println("Nuevo bitmap libre DD")
				fmt.Println(nuevoFFBDD)
				fmt.Println("\nEjecuacion pausada... Presione enter para continuar")
				fmt.Scanln()*/
				nuevoArreglogArchivo := arregloArchivos{}

				nuevoInodo := inodo{}
				posNuevoInodo := superBloque.SBapINODO + superBloque.SBfirstFreeBitINODO*superBloque.SBsizeStructINODO
				/*fmt.Println("Posicion para insertar el nuevo Inodo")
				fmt.Println(posNuevoInodo)
				fmt.Println("\nEjecuacion pausada... Presione enter para continuar")
				fmt.Scanln()*/
				superBloque.SBinodosFree--
				actualizarValorBitmap(rutaDisco, superBloque.SBapBINODO+superBloque.SBfirstFreeBitINODO, '1')
				nuevoFFBINODO := obtenerFirstFreeBit(rutaDisco, superBloque.SBapBINODO, int(superBloque.SBinodosCount))
				/*listadobitmap(rutaDisco, superBloque.SBapBINODO, int(superBloque.SBsizeStructINODO))
				fmt.Println("\nEjecuacion pausada... Presione enter para continuar")
				fmt.Scanln()*/
				bInodoAnt := superBloque.SBfirstFreeBitINODO
				superBloque.SBfirstFreeBitINODO = nuevoFFBINODO
				/*fmt.Println("Nuevo bitmap libre Inodo")
				fmt.Println(nuevoFFBINODO)
				fmt.Println("\nEjecuacion pausada... Presione enter para continuar")
				fmt.Scanln()*/

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
				/*fmt.Println("COntenido en el DD a insertar")
				fmt.Println(nuevoDD)
				fmt.Println("En la pos")
				fmt.Println(posDD)
				fmt.Println("\nEjecuacion pausada... Presione enter para continuar")
				fmt.Scanln()*/

				escribirStructDD(rutaDisco, posDD, nuevoDD)
				/*fmt.Println("COntenido en el INODO NUEVO")
				fmt.Println(nuevoInodo)*/

				nuevoInodo.IcountInodo = bInodoAnt + 1
				nuevoInodo.IsizeArchivo = int64(tamContenido)
				nuevoInodo.IcountBloquesAsignados = int64(cantidadBloques)
				nuevoInodo.IidProper = 1
				nuevoInodo.IapIndirecto = -1
				for i := 0; i < 4; i++ {
					nuevoInodo.IarrayBloques[i] = -1
				}

				/*fmt.Println("COntenido en el Inodo a insertar")
				fmt.Println(nuevoInodo)
				fmt.Println("En la pos")
				fmt.Println(posNuevoInodo)
				fmt.Println("\nEjecuacion pausada... Presione enter para continuar")
				fmt.Scanln()*/
				escribirStructINODO(rutaDisco, posNuevoInodo, nuevoInodo)

				raiz.AVDapDetalleDir = posDD
				/*fmt.Println("La raiz su apuntado a un DD vale")
				fmt.Println(posDD)
				fmt.Println(raiz)
				fmt.Println("\nEjecuacion pausada... Presione enter para continuar")
				fmt.Scanln()*/
				escribirStructAVD(rutaDisco, superBloque.SBapAVD, raiz)

				escribirSuperBloque(rutaDisco, inicioPart, superBloque)
				superBloque = obtenerSB(rutaDisco, inicioPart)

				/*fmt.Println("Posicion del nuevo inodo")
				fmt.Println(posNuevoInodo)
				fmt.Println("Sb")
				fmt.Println(superBloque)
				fmt.Println("\nEjecuacion pausada... Presione enter para continuar")
				fmt.Scanln()*/
				llenarNuevoArchivo(contenido, posNuevoInodo, rutaDisco, superBloque, inicioPart)
				/*fmt.Println("/***************************************************************************************************\\")
				fmt.Println("Finalizo Insercion")
				fmt.Println("\nEjecuacion pausada... Presione enter para continuar")
				fmt.Scanln()*/

				/*raiz = obtenerAVD(rutaDisco, superBloque.SBapAVD)
				fmt.Println("***************************** RAIZ ***************************************")
				fmt.Println("********************************************************************")
				fmt.Println("********************************************************************")
				fmt.Println(raiz)
				fmt.Println("\nEjecuacion pausada... Presione enter para continuar")
				fmt.Scanln()*/

			} else {

				path1 := strings.TrimPrefix(rutaAchivo, "/")
				path2 := strings.TrimSuffix(path1, "/")
				pathPart := strings.Split(path2, "/")

				encontrado := false
				posEncontrado := superBloque.SBapAVD
				//t := 0

				for i := 0; i < len(pathPart); i++ {
					//t = i
					fmt.Println(pathPart[i])
					encontrado, posEncontrado = buscarDir(posEncontrado, pathPart[i], rutaDisco)
					if encontrado == false {
						break
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

func llenarNuevoArchivo(cont string, posInodo int64, rutaDisco string, superBloque sb, inicioPart int64) {
	contB := []byte(cont)
	termino := false
	/*fmt.Println("Pos inodo obtenido")
	fmt.Println(posInodo)
	fmt.Println("Inodo obtenido")*/
	nuevoInodo := obtenerINODO(rutaDisco, posInodo)
	/*fmt.Println(nuevoInodo)
	fmt.Println("CONTENIDO")
	fmt.Println(contB)
	fmt.Println("\nEjecuacion pausada... Presione enter para continuar")
	fmt.Scanln()*/

	//var letra byte = 'A'
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
		/*fmt.Println("BLOQUE -- " + strconv.Itoa(i))
		fmt.Println(nuevoBloque)
		fmt.Println("\nEjecuacion pausada... Presione enter para continuar")
		fmt.Scanln()*/
		if termino == true {
			/*fmt.Println("SE supone que ingreso las 3")
			fmt.Println("\nEjecuacion pausada... Presione enter para continuar")
			fmt.Scanln()*/
		}

		posNuevoBloque := superBloque.SBapBLOQUE + superBloque.SBfirstFreeBitBLOQUE*superBloque.SBsizeStructBLOQUE
		/*fmt.Println("POsicion en la que voy a insertar el bloque")
		fmt.Println(posNuevoBloque)
		fmt.Println("\nEjecuacion pausada... Presione enter para continuar")
		fmt.Scanln()*/

		escribirStructBLOQUE(rutaDisco, posNuevoBloque, nuevoBloque)
		superBloque.SBbloquesFree--
		actualizarValorBitmap(rutaDisco, superBloque.SBapBBLOQUE+superBloque.SBfirstFreeBitBLOQUE, '1')
		nuevoFFBBLOQUE := obtenerFirstFreeBit(rutaDisco, superBloque.SBapBBLOQUE, int(superBloque.SBbloquesCount))
		/*fmt.Println("Nuevo bit libre de bloque")
		fmt.Println(nuevoFFBBLOQUE)
		fmt.Println("\nEjecuacion pausada... Presione enter para continuar")
		fmt.Scanln()*/
		superBloque.SBfirstFreeBitBLOQUE = nuevoFFBBLOQUE
		escribirSuperBloque(rutaDisco, inicioPart, superBloque)
		superBloque = obtenerSB(rutaDisco, inicioPart)

		nuevoInodo.IarrayBloques[i] = posNuevoBloque
		/*fmt.Println("INodo actualizado")
		fmt.Println(nuevoInodo)
		fmt.Println("\nEjecuacion pausada... Presione enter para continuar")
		fmt.Scanln()*/
		pos = i
		if termino == true {
			/*fmt.Println("Termino igual a true")
			fmt.Println("Insertar en pos ")
			fmt.Println(posInodo)
			fmt.Println(nuevoInodo)
			fmt.Println("\nEjecuacion pausada... Presione enter para continuar")
			fmt.Scanln()*/
			escribirStructINODO(rutaDisco, posInodo, nuevoInodo)
			superBloque = obtenerSB(rutaDisco, inicioPart)
			break
		} else if len(contB) == 25 {
			termino = true
			/*fmt.Println("Termino porque el ultimo tenia 25")
			fmt.Println("Insertar en pos ")
			fmt.Println(posInodo)
			fmt.Println(nuevoInodo)
			fmt.Println("\nEjecuacion pausada... Presione enter para continuar")
			fmt.Scanln()*/
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
			/*fmt.Println("Posicion para insertar el nuevo Inodo indirecto")
			fmt.Println(posNuevoInodoIndirecto)
			fmt.Println("\nEjecuacion pausada... Presione enter para continuar")
			fmt.Scanln()*/
			superBloque.SBinodosFree--
			actualizarValorBitmap(rutaDisco, superBloque.SBapBINODO+superBloque.SBfirstFreeBitINODO, '1')
			nuevoFFBINODO := obtenerFirstFreeBit(rutaDisco, superBloque.SBapBINODO, int(superBloque.SBinodosCount)) //Error por el count
			bInodoAnt := superBloque.SBfirstFreeBitINODO
			superBloque.SBfirstFreeBitINODO = nuevoFFBINODO

			/*fmt.Println("Nuevo bitmap libre Inodo")
			fmt.Println(nuevoFFBINODO)
			fmt.Println("\nEjecuacion pausada... Presione enter para continuar")
			fmt.Scanln()*/

			nuevoInodoIndirecto.IcountInodo = bInodoAnt + 1
			nuevoInodoIndirecto.IsizeArchivo = nuevoInodo.IsizeArchivo
			nuevoInodoIndirecto.IcountBloquesAsignados = nuevoInodo.IcountBloquesAsignados
			nuevoInodoIndirecto.IidProper = 1
			nuevoInodoIndirecto.IapIndirecto = -1
			for i := 0; i < 4; i++ {
				nuevoInodoIndirecto.IarrayBloques[i] = -1
			}

			nuevoInodo.IapIndirecto = posNuevoInodoIndirecto
			/*fmt.Println("Actualizando inodo con indirecto")
			fmt.Println("actualizar inodo en pos")
			fmt.Println(posInodo)
			fmt.Println("cONTENIDO INODOD")
			fmt.Println(nuevoInodo)
			fmt.Println("\nEjecuacion pausada... Presione enter para continuar")
			fmt.Scanln()*/
			escribirStructINODO(rutaDisco, posInodo, nuevoInodo)

			/*fmt.Println("insertar inodo indirecto")
			fmt.Println("inodo en pos")
			fmt.Println(posNuevoInodoIndirecto)
			fmt.Println("cONTENIDO INODOD")
			fmt.Println(nuevoInodoIndirecto)
			fmt.Println("\nEjecuacion pausada... Presione enter para continuar")
			fmt.Scanln()*/
			escribirStructINODO(rutaDisco, posNuevoInodoIndirecto, nuevoInodoIndirecto)
			escribirSuperBloque(rutaDisco, inicioPart, superBloque)
			superBloque = obtenerSB(rutaDisco, inicioPart)
			//
			//fmt.Println("Contenido restante")
			contres := ""
			for i := 0; i < len(contB); i++ {
				contres += string(contB[i])
			}
			/*fmt.Println(contres)
			fmt.Println(contB)
			fmt.Println("\nEjecuacion pausada... Presione enter para continuar")
			fmt.Scanln()*/
			llenarNuevoArchivo(contres, posNuevoInodoIndirecto, rutaDisco, superBloque, inicioPart)
		}
	}

}
