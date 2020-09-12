package main

import (
	"container/list"
	"fmt"
	"strconv"
	"unsafe"
)

type datosLoss struct {
	posLog     int64
	posCopiaSB int64
}

var listaLoss list.List
var structLoss datosLoss

func loss(vd string) {
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

			structLoss.posLog = superBloque.SBapLOG
			var sizeBitacora int64 = int64(unsafe.Sizeof(bitacora{}))
			posCopiaSB := superBloque.SBapLOG + (superBloque.SBavdCount * sizeBitacora)
			structLoss.posCopiaSB = posCopiaSB
			deleteFull(rutaDisco, inicioPart, superBloque.SBapLOG-superBloque.SBapAVD)

		} else {
			fmt.Println("La particion indicada no esta montada")
		}
	} else {
		fmt.Println("EL disco proporcionado no esta montado")
	}
}

func recorery(vd string) {
	var idDisco byte
	idDisco = vd[2]
	idDisco2 := idDisco - 97
	idP, _ := strconv.Atoi(vd[3:])
	idP--

	if arregloMount[idDisco2].estado == 1 {
		if arregloMount[idDisco2].discos[idP].estado == 1 {
			inicioPart := arregloMount[idDisco2].discos[idP].Partstart
			rutaDisco := arregloMount[idDisco2].Ruta
			superBloque := obtenerSB(rutaDisco, structLoss.posCopiaSB)

			superBloque.SBavdFree = superBloque.SBavdCount
			superBloque.SBddFree = superBloque.SBddCount
			superBloque.SBinodosFree = superBloque.SBinodosCount
			superBloque.SBbloquesFree = superBloque.SBbloquesCount

			superBloque.SBfirstFreeBitAVD = 0
			superBloque.SBfirstFreeBitDD = 0
			superBloque.SBfirstFreeBitINODO = 0
			superBloque.SBfirstFreeBitBLOQUE = 0
			escribirSuperBloque(rutaDisco, inicioPart, superBloque)

			var sizeBitacora int64 = int64(unsafe.Sizeof(bitacora{}))
			pos := structLoss.posLog
			for i := 0; i < int(superBloque.SBavdCount); i++ {
				bit := obtenerbitacora(rutaDisco, pos)
				if bit.LOGtipo != 'x' {
					cont := ""
					for a := 0; a < 100; a++ {
						if bit.LOGcontenido[a] != 0 {
							cont += string(bit.LOGcontenido[a])
						}
					}
					path := ""
					for a := 0; a < 100; a++ {
						if bit.LOGnombre[a] != 0 {
							path += string(bit.LOGnombre[a])
						}
					}

					if bit.LOGtipo == '1' {
						crearArvhi(vd, path, true, int64(int(bit.LOGsize)), cont, false)
					} else if bit.LOGtipo == '0' {
						crearCarpeta(vd, path, true, false)
					}

				}
				pos = pos + sizeBitacora

			}

		} else {
			fmt.Println("La particion indicada no esta montada")
		}
	} else {
		fmt.Println("EL disco proporcionado no esta montado")
	}
}
