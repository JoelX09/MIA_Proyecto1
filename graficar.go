package main

import (
	"container/list"
	"fmt"
	"io/ioutil"
	"math"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"unsafe"
)

func graficar(path string, nombre string, id string, ruta string) {
	if nombre == "mbr" {
		dot := graficaMBR(id)
		pathDot, nombreDot := descomponer(path)
		generarDOT(dot, pathDot, nombreDot)
	} else if nombre == "disk" {
		dot := graficaDisco(id)
		pathDot, nombreDot := descomponer(path)
		generarDOT(dot, pathDot, nombreDot)
	}
}

func graficaMBR(vd string) string {
	var idDisco byte
	idDisco = vd[2]
	idDisco2 := idDisco - 97
	idP, _ := strconv.Atoi(vd[3:])
	idP--
	dot := ""

	if arregloMount[idDisco2].estado == 1 {
		if arregloMount[idDisco2].discos[idP].estado == 1 {
			rutaDisco := arregloMount[idDisco2].Ruta
			name := arregloMount[idDisco2].discos[idP].Partname
			nameSt := string(name[:])

			listaP, _ := listaInicialPE(rutaDisco)
			_, valoresExt, _ := imprimirListaPE(nameSt, false, true, listaP)
			listaNL.Init()
			listaL := listaInicialL(rutaDisco, valoresExt.inicioE, valoresExt.tamE, valoresExt.inicioE)

			dot += "digraph G {\n" +
				"\tnode [shape=plaintext]\n" +
				"\ta [label=<\n" +
				"\t<table border=\"0\" cellborder=\"1\" cellspacing=\"0\">\n\n"
			m := obtenerMbr(rutaDisco)
			dot += "		<tr><td><b>Nombre</b></td><td><b>Valor</b></td></tr>\n" +
				"		<tr><td><b>mbr_tamano</b></td><td>" + strconv.FormatInt(m.Mbrtam, 10) + "</td></tr>\n" +
				"		<tr><td><b>mbr_fecha_creacion</b></td><td>" + string(m.Mbrfecha[:]) + "</td></tr>\n" +
				"		<tr><td><b>mbr_disk_signature</b></td><td>" + strconv.Itoa(int(m.Mbrdisksig)) + "</td></tr>\n\n"

			pos := 1
			for ele := listaP.Front(); ele != nil; ele = ele.Next() {
				temp := ele.Value.(nodoPart)
				nombrePart := ""
				for i := 0; i < len(temp.Partname); i++ {
					if temp.Partname[i] != 0 {
						nombrePart += string(temp.Partname[i])
					}
				}
				dot += "		<tr><td><b>part_status_" + strconv.Itoa(pos) + "</b></td><td>" + strconv.Itoa(int(temp.Partstatus)) + "</td></tr>\n" +
					"		<tr><td><b>part_type_" + strconv.Itoa(pos) + "</b></td><td>" + string(temp.Parttype) + "</td></tr>\n" +
					"		<tr><td><b>part_fit_" + strconv.Itoa(pos) + "</b></td><td>" + string(temp.Partfit) + "</td></tr>\n" +
					"		<tr><td><b>part_start_" + strconv.Itoa(pos) + "</b></td><td>" + strconv.Itoa(int(temp.Partstart)) + "</td></tr>\n" +
					"		<tr><td><b>part_size_" + strconv.Itoa(pos) + "</b></td><td>" + strconv.Itoa(int(temp.Partsize)) + "</td></tr>\n" +
					"		<tr><td><b>part_name_" + strconv.Itoa(pos) + "</b></td><td>" + nombrePart + "</td></tr>\n\n"
				pos++
			}
			dot += "\t</table>>];\n"

			pos = 1
			for ele := listaL.Front(); ele != nil; ele = ele.Next() {
				temp := ele.Value.(estructEBR)
				nombrePart := ""
				for i := 0; i < len(temp.PartnameL); i++ {
					if temp.PartnameL[i] != 0 {
						nombrePart += string(temp.PartnameL[i])
					}
				}
				dot += "\n\tb" + strconv.Itoa(pos) + " [label=<\n" +
					"\t<table border=\"0\" cellborder=\"1\" cellspacing=\"0\">\n\n" +
					"\t<tr><td colspan=\"2\">EBR" + strconv.Itoa(pos) + "</td></tr>\n" +
					"		<tr><td><b>Nombre</b></td><td><b>Valor</b></td></tr>\n" +
					"		<tr><td><b>part_status_" + strconv.Itoa(pos) + "</b></td><td>" + strconv.Itoa(int(temp.PartstatusL)) + "</td></tr>\n" +
					"		<tr><td><b>part_fit_" + strconv.Itoa(pos) + "</b></td><td>" + string(temp.PartfitL) + "</td></tr>\n" +
					"		<tr><td><b>part_start_" + strconv.Itoa(pos) + "</b></td><td>" + strconv.Itoa(int(temp.PartstartL)) + "</td></tr>\n" +
					"		<tr><td><b>part_size_" + strconv.Itoa(pos) + "</b></td><td>" + strconv.Itoa(int(temp.PartsizeL)) + "</td></tr>\n" +
					"		<tr><td><b>part_next_" + strconv.Itoa(pos) + "</b></td><td>" + strconv.Itoa(int(temp.PartnextL)) + "</td></tr>\n" +
					"		<tr><td><b>part_name_" + strconv.Itoa(pos) + "</b></td><td>" + nombrePart + "</td></tr>\n\n" +
					"\t</table>>];\n"
				pos++
			}

			dot += "}"
		} else {
			fmt.Println("La particion indica no esta mon")
		}
	} else {
		fmt.Println("EL disco proporcionado no esta montado")
	}
	return dot
}

func graficaDisco(vd string) string {
	var idDisco byte
	idDisco = vd[2]
	idDisco2 := idDisco - 97
	idP, _ := strconv.Atoi(vd[3:])
	idP--
	dot := ""

	if arregloMount[idDisco2].estado == 1 {
		if arregloMount[idDisco2].discos[idP].estado == 1 {
			rutaDisco := arregloMount[idDisco2].Ruta
			name := arregloMount[idDisco2].discos[idP].Partname
			nameSt := string(name[:])

			m := obtenerMbr(rutaDisco)
			var sizeMBR int = int(unsafe.Sizeof(m))

			listaP, _ := listaInicialPE(rutaDisco)
			_, valoresExt, _ := imprimirListaPE(nameSt, false, true, listaP)
			var listaPtemp = list.New()
			listaPtemp.PushFrontList(listaP)
			listaP.Init()
			_, listaP = espaciosPEdisp(sizeMBR, m, listaPtemp)

			listaNL.Init()
			listaL := listaInicialL(rutaDisco, valoresExt.inicioE, valoresExt.tamE, valoresExt.inicioE)
			var listaLtemp = list.New()
			listaLtemp.PushFrontList(listaL)
			listaL.Init()
			listaL = espaciosLL(valoresExt.inicioE, valoresExt.tamE, listaLtemp)
			imprimirListaL("", false, false, listaL)

			dot += "digraph G {\n" +
				"\tnode [shape=plaintext]\n" +
				"\ta [label=<\n" +
				"\t<table border=\"1\" cellborder=\"1\" cellspacing=\"0\">\n\n"

			tamDisco := m.Mbrtam - int64(sizeMBR)
			//tamDisco64 := float64(tamDisco)

			dot += "		<tr>\n" +
				"		<td>MBR</td>\n"
			pos := 1
			for ele := listaP.Front(); ele != nil; ele = ele.Next() {
				temp := ele.Value.(nodoPart)
				nombrePart := ""
				for i := 0; i < len(temp.Partname); i++ {
					if temp.Partname[i] != 0 {
						nombrePart += string(temp.Partname[i])
					}
				}
				porcentaje1 := (float64(temp.Partsize) * 100) / float64(tamDisco)
				porcentaje := int(math.Round(porcentaje1))
				contenido := ""
				if temp.Estado == 0 {
					contenido = "Libre \n" + strconv.Itoa(porcentaje) + "%"
				} else if temp.Estado == 1 {
					if temp.Parttype == 'E' {
						a := 0
						b := 0
						for n := listaL.Front(); n != nil; n = n.Next() {
							a1 := n.Value.(estructEBR)
							if a1.EstadoL == 1 {
								a++
							} else if a1.EstadoL == 0 {
								b++
							}
						}
						contenido = " <table>\n" +
							"<tr><td colspan=\"" + strconv.Itoa(a*2+b) + "\">" + nombrePart + " (" + string(temp.Parttype) + ") " + strconv.Itoa(porcentaje) + "%</td></tr>" +
							"<tr>\n"
						for eleL := listaL.Front(); eleL != nil; eleL = eleL.Next() {
							//contenido += "<td>EBR</td><td>"
							tempL := eleL.Value.(estructEBR)
							nombrePartL := ""
							for i := 0; i < len(tempL.PartnameL); i++ {
								if tempL.PartnameL[i] != 0 {
									nombrePartL += string(tempL.PartnameL[i])
								}
							}
							porcentajeL1 := float64(tempL.PartsizeL) * 100 / float64(valoresExt.tamE)
							porcentajeL := int(math.Round(porcentajeL1))
							//fmt.Println(valoresExt.tamE)
							//fmt.Println(tempL.PartsizeL)
							//fmt.Println(math.Round(porcentajeL))
							if tempL.EstadoL == 0 {
								contenido += "<td>Libre \n" + strconv.Itoa(porcentajeL) + "%</td>"
							} else if tempL.EstadoL == 1 {
								contenido += "<td>EBR</td><td>" + nombrePartL + " \n" + strconv.Itoa(porcentajeL) + "%</td>"
							}
							//contenido += "</td>"
						}
						contenido += "</tr>" +
							"</table>"
					} else {
						contenido = nombrePart + " (" + string(temp.Parttype) + ") " + strconv.Itoa(porcentaje) + "%"
					}
				}
				dot += "		<td>" + contenido + "</td>\n"
				pos++
			}
			dot += "		</tr>\n"
			dot += "\t</table>>];\n"

			dot += "}"
		} else {
			fmt.Println("La particion indica no esta mon")
		}
	} else {
		fmt.Println("EL disco proporcionado no esta montado")
	}
	return dot
}

func generarDOT(dot string, path string, nombre string) {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		err = os.MkdirAll(path, 0777)
		if err != nil {
			panic(err)
		}
	}

	nombreDot := nombre
	nombreDot = strings.ReplaceAll(nombreDot, "png", "dot")
	f, err := os.Create(path + nombreDot)
	if err != nil {
		panic(err)
	}
	f.Close()

	err = ioutil.WriteFile(path+nombreDot, []byte(dot), 0777)
	if err != nil {
		panic(err)
	}
	pathdot := path + nombreDot
	pathimg := path + nombre

	cmd := exec.Command("dot", "-Tpng", pathdot, "-o", pathimg)
	cmd.Run()

}

func descomponer(path string) (string, string) {
	var carpeta, archivo string
	pathPart := strings.SplitAfter(path, "/")
	for i := 0; i < len(pathPart)-1; i++ {
		carpeta += pathPart[i]
	}
	archivo = pathPart[len(pathPart)-1]
	return carpeta, archivo
}
