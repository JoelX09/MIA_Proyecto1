package main

import (
	"container/list"
	"fmt"
	"regexp"
	"strconv"
)

func desmontar() {
	for ele := listaID.Front(); ele != nil; ele = ele.Next() {
		valorVD := ele.Value.(string)
		match, _ := regexp.MatchString("^vd[a-z][0-9]+", valorVD)
		if match == true {
			eliminarVD(valorVD)
		} else {
			fmt.Println("La estructura del ID es incorrecta")
		}
	}
	listaID.Init()
	fmt.Println("===========================\nLista de desmontar finalizo")
}

func eliminarVD(vd string) {
	var idDisco byte
	idDisco = vd[2]
	idDisco2 := idDisco - 97
	idP, _ := strconv.Atoi(vd[3:])
	idP--

	fmt.Println("vd" + string(idDisco) + strconv.Itoa(idP+1))

	if arregloMount[idDisco2].estado == 1 {
		if arregloMount[idDisco2].discos[idP].estado == 1 {
			name := arregloMount[idDisco2].discos[idP].Partname
			nameSt := string(name[:])
			path := arregloMount[idDisco2].Ruta

			listaP, _ := listaInicialPE(path)
			existeNombrePE, valoresExt, _ := imprimirListaPE(nameSt, false, true, listaP)
			existeNombreL := false
			var discoL estructEBR

			if existeNombrePE == false {
				if valoresExt.inicioE != 0 {
					listaNL.Init()
					listaL := listaInicialL(path, valoresExt.inicioE, valoresExt.tamE, valoresExt.inicioE)
					existeNombreL, discoL = imprimirListaL(nameSt, false, true, listaL)
				}
			}

			if existeNombrePE == true {
				for ele := listaP.Front(); ele != nil; ele = ele.Next() {
					temp := ele.Value.(nodoPart)
					if temp.Partname == name {
						temp.Partstatus = 0

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
							listaP.Remove(ele)
							listaP.PushBack(temp)
						} else {
							tempSig := ele.Next()
							listaP.Remove(ele)
							listaP.InsertBefore(temp, tempSig)
						}
						actualizarMBR(path, listaP)
						break
					}
				}
			} else if existeNombreL == true {
				var listaTemp = list.New()
				discoL.PartstatusL = 0
				listaTemp.PushBack(discoL)
				escribirListaEbr(path, listaTemp)
			}
			var nueva estructParticion
			arregloMount[idDisco2].discos[idP] = nueva
			tam := 0
			for j := 0; j < len(arregloMount[idDisco2].discos); j++ {
				if arregloMount[idDisco2].discos[j].estado == 1 {
					tam++
				}
			}
			if tam == 0 {
				var nuevo estructDisco
				arregloMount[idDisco2] = nuevo
			}
			fmt.Println("Particion desmontada")
		}
	} else {
		fmt.Println("EL disco solicitado no esta montado")
	}

}
