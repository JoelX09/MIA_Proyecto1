package main

import (
	"container/list"
	"fmt"
)

type estructDisco struct {
	estado uint8
	Ruta   string
	//Letra  byte
	discos [50]estructParticion
}

type estructParticion struct {
	estado    uint8
	Partname  [16]byte
	Partfit   byte
	Partstart int64
	Partsize  int64
	//Partnext int64
	//Parttype byte
}

func montarParticion(path string, name string) {
	//letra := byte('a')
	var tempcomp [16]byte
	copy(tempcomp[:], name)
	fmt.Println(tempcomp)

	listaP, _ := listaInicialPE(path)
	existeNombrePE, valoresExt, _ := imprimirListaPE(name, true, true, listaP)
	existeNombreL := false
	var discoL estructEBR

	if existeNombrePE == false {
		listaNL.Init()
		listaL := listaInicialL(path, valoresExt.inicioE, valoresExt.tamE, valoresExt.inicioE)
		existeNombreL, discoL = imprimirListaL(name, true, true, listaL)
	}

	if existeNombrePE == true {
		for ele := listaP.Front(); ele != nil; ele = ele.Next() {
			temp := ele.Value.(nodoPart)
			if temp.Partname == tempcomp {
				if temp.Partstatus == 0 {
					var insertDisco estructDisco
					var inserParticion estructParticion

					insertDisco.estado = 1
					insertDisco.Ruta = path
					//insertDisco.Letra = letra

					inserParticion.estado = 1
					inserParticion.Partname = temp.Partname
					inserParticion.Partfit = temp.Partfit
					inserParticion.Partstart = temp.Partstart
					inserParticion.Partsize = temp.Partsize

					insertarMount(path, insertDisco, inserParticion)

					temp.Partstatus = 1

					if ele.Prev() == nil {
						if ele.Next() == nil {
							listaP.Remove(ele)
							listaP.PushFront(temp)
						} else {
							tempSig := ele.Next()
							listaP.Remove(ele)
							listaP.InsertBefore(temp, tempSig)
						}
					} else if ele.Next() == nil {
						//tempAnt := ele.Prev()
						listaP.Remove(ele)
						listaP.PushBack(temp)
					} else {
						tempSig := ele.Next()
						listaP.Remove(ele)
						listaP.InsertBefore(temp, tempSig)
					}
					actualizarMBR(path, listaP)
				} else if temp.Partstatus == 1 {
					fmt.Println("La particion ya esta montada")
				}
				break
			}
		}
	} else if existeNombreL == true {
		if discoL.PartstatusL == 0 {
			var insertDisco estructDisco
			var inserParticion estructParticion

			insertDisco.estado = 1
			insertDisco.Ruta = path
			//insertDisco.Letra = letra

			inserParticion.estado = 1
			inserParticion.Partname = discoL.PartnameL
			inserParticion.Partfit = discoL.PartfitL
			inserParticion.Partstart = discoL.PartstartL
			inserParticion.Partsize = discoL.PartsizeL

			insertarMount(path, insertDisco, inserParticion)
			var listaTemp = list.New()
			discoL.PartstatusL = 1
			listaTemp.PushBack(discoL)
			escribirListaEbr(path, listaTemp)

		} else if discoL.PartstatusL == 1 {
			fmt.Println("La particion ya esta montada")
		}

	} else {
		fmt.Println("NO se encontro una particion para montar con el nombre: " + name)
	}
}

func insertarMount(path string, disco estructDisco, particion estructParticion) {
	existeRuta := false
	for i := 0; i < len(arregloMount); i++ {
		if arregloMount[i].Ruta == path {
			arregloDiscos := arregloMount[i].discos
			for j := 0; i < len(arregloDiscos); j++ {
				if arregloDiscos[j].estado == 0 {
					arregloMount[i].discos[j] = particion
					break
				}
			}
			existeRuta = true
			break
		}
	}

	if existeRuta == false {
		for i := 0; i < len(arregloMount); i++ {
			if arregloMount[i].estado == 0 {
				disco.discos[0] = particion
				arregloMount[i] = disco
				return
			}
		}

	}

}
