package main

type estructDisco struct {
	Ruta   string
	Letra  byte
	discos [50]estructParticion
}

type estructParticion struct {
	PartnameL  [16]byte
	estado     uint8
	PartfitL   byte
	PartstartL int64
	PartsizeL  int64
}

func montarParticion(path string, name string) {
	/*//letra := byte('a')
	//num := 1
	//pos := 0

	//if listaMount.Len() != 0 {
	//for eled := listaMount.Front(); eled != nil; eled = eled.Next() {
	m := obtenerMbr(path)
	crearListaMBR(m)

	fmt.Println("-----------------------------------")
	fmt.Println("Particiones P y E existentes")
	existeNombrePE, valoresExt := imprimirListaPE(name, true, true)
	fmt.Println("-----------------------------------")
	encontrado := false
	for ele := listaP.Front(); ele != nil; ele = ele.Next() {
		temp := ele.Value.(nodoPart)
		if temp.Estado == 1 {

			var tempcomp [16]byte
			copy(tempcomp[:], name)

			if temp.Parttype == 'E' && existeNombrePE == false {
				fmt.Println("Recorrer Logicas para ver si es la que se monta")

				listaL.Init()
				listaLogica(path, valoresExt.inicioE, valoresExt.tamE, name, valoresExt.inicioE)

				fmt.Println("-----------------------------------")
				fmt.Println("Particiones L existentes")
				/*logicaEncontrada :=*/ /*mostrarListaLogica(name, true, true)
	fmt.Println("-----------------------------------")

	/*if logicaEncontrada == true {
		for eleL := listaL.Front(); eleL != nil; eleL = eleL.Next() {
			tempL := eleL.Value.(estructEBR)
			//if tempL.EstadoL == 1 {
			if tempL.PartnameL == tempcomp {
				if listaMount.Len() != 0 {
					for eled := listaMount.Front(); eled != nil; eled = eled.Next() {
						tempM := eled.Value.(estructDisco)
						letra = tempM.Letra

						if strings.Compare(tempM.Ruta, path) == 0 {
							for i := 1; i < len(tempM.discos); i++ {
								if tempM.discos[i].estado == 0 {

									tempM.discos[i].PartnameL = tempL.PartnameL
									tempM.discos[i].PartfitL = tempL.PartfitL
									tempM.discos[i].PartstartL = tempL.PartstartL
									tempM.discos[i].PartsizeL = tempL.PartsizeL
									break
								}
							}
						} else {
							if eled.Next() == nil {
								var temporalDisco estructDisco
								temporalDisco.Ruta = path
								temporalDisco.Letra = letra + 1
								temporalDisco.discos[1].estado = 1
								temporalDisco.discos[1].PartnameL = tempL.PartnameL
								temporalDisco.discos[1].PartfitL = tempL.PartfitL
								temporalDisco.discos[1].PartstartL = tempL.PartstartL
								temporalDisco.discos[1].PartsizeL = tempL.PartsizeL
								break
							}
						}
					}
				} else {
					var temporalDisco estructDisco
					temporalDisco.Ruta = path
					temporalDisco.Letra = letra
					temporalDisco.discos[1].estado = 1
					temporalDisco.discos[1].PartnameL = tempL.PartnameL
					temporalDisco.discos[1].PartfitL = tempL.PartfitL
					temporalDisco.discos[1].PartstartL = tempL.PartstartL
					temporalDisco.discos[1].PartsizeL = tempL.PartsizeL
					break
				}

			}
			//}
		}
	}*/
	//}

	/*if encontrado == false && temp.Partname == tempcomp {

			}

		}
	}
	if encontrado == false {
		fmt.Println("NO se encontro ninguna particion con ese nombre")
	}
	//}
	/*	} else {
		var dataDisco estructDisco

		var dataPartition estructParticion
		dataPartition.PartnameL =
		dataDisco.Ruta = path
		dataDisco.Letra = letra
		dataDisco.discos[pos] =

	}*/

}

/*func varr() list.List {
	var l1 list.List

	l1.PushFront(1)

	return l1
}*/
