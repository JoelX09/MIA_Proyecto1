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
		generar(dot, path)
		/*pathDot, nombreDot := descomponer(path)
		generarDOT(dot, pathDot, nombreDot)*/
	} else if nombre == "disk" {
		dot := graficaDisco(id)
		generar(dot, path)
		/*pathDot, nombreDot := descomponer(path)
		generarDOT(dot, pathDot, nombreDot)*/
	} else if nombre == "sb" {
		dot := graficasb(id)
		generar(dot, path)
		/*pathDot, nombreDot := descomponer(path)
		generarDOT(dot, pathDot, nombreDot)*/
	} else if nombre == "tree_directorio" {
		dot := graficarTreeDirectorio(id)
		generar(dot, path)
		/*pathDot, nombreDot := descomponer(path)
		generarDOT(dot, pathDot, nombreDot)*/
	}
}

func generar(dot string, path string) {
	pathDot, nombreDot := descomponer(path)
	generarDOT(dot, pathDot, nombreDot)
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

			listaP, valoresPE := listaInicialPE(rutaDisco)
			_, valoresExt, _ := imprimirListaPE(nameSt, false, true, listaP)
			listaNL.Init()
			var listaL = list.New()

			if valoresPE[1] == 1 {
				listaL = listaInicialL(rutaDisco, valoresExt.inicioE, valoresExt.tamE, valoresExt.inicioE)
			}

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

			if valoresPE[1] == 1 {
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

			listaP, valoresPE := listaInicialPE(rutaDisco)
			_, valoresExt, _ := imprimirListaPE(nameSt, false, true, listaP)
			var listaPtemp = list.New()
			listaPtemp.PushFrontList(listaP)
			listaP.Init()
			_, listaP = espaciosPEdisp(sizeMBR, m, listaPtemp)

			listaNL.Init()
			var listaL = list.New()

			if valoresPE[1] == 1 {
				listaL = listaInicialL(rutaDisco, valoresExt.inicioE, valoresExt.tamE, valoresExt.inicioE)
				var listaLtemp = list.New()
				listaLtemp.PushFrontList(listaL)
				listaL.Init()
				listaL = espaciosLL(valoresExt.inicioE, valoresExt.tamE, listaLtemp)
				imprimirListaL("", false, false, listaL)
			}

			dot += "digraph G {\n" +
				"\tnode [shape=plaintext]\n" +
				"\ta [label=<\n" +
				"\t<table border=\"1\" cellborder=\"1\" cellspacing=\"0\">\n\n"

			tamDisco := m.Mbrtam - int64(sizeMBR)

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
							tempL := eleL.Value.(estructEBR)
							nombrePartL := ""
							for i := 0; i < len(tempL.PartnameL); i++ {
								if tempL.PartnameL[i] != 0 {
									nombrePartL += string(tempL.PartnameL[i])
								}
							}
							porcentajeL1 := float64(tempL.PartsizeL) * 100 / float64(valoresExt.tamE)
							porcentajeL := int(math.Round(porcentajeL1))
							if tempL.EstadoL == 0 {
								contenido += "<td>Libre \n" + strconv.Itoa(porcentajeL) + "%</td>"
							} else if tempL.EstadoL == 1 {
								contenido += "<td>EBR</td><td>" + nombrePartL + " \n" + strconv.Itoa(porcentajeL) + "%</td>"
							}
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

func graficasb(vd string) string {
	var idDisco byte
	idDisco = vd[2]
	idDisco2 := idDisco - 97
	idP, _ := strconv.Atoi(vd[3:])
	idP--
	dot := ""

	if arregloMount[idDisco2].estado == 1 {
		if arregloMount[idDisco2].discos[idP].estado == 1 {
			inicioPart := arregloMount[idDisco2].discos[idP].Partstart
			rutaDisco := arregloMount[idDisco2].Ruta
			superbloque := obtenerSB(rutaDisco, inicioPart)

			nombrePart := ""
			for i := 0; i < len(superbloque.SBnombreHd); i++ {
				if superbloque.SBnombreHd[i] != 0 {
					nombrePart += string(superbloque.SBnombreHd[i])
				}
			}

			dot += "digraph G {\n" +
				"\tnode [shape=plaintext]\n" +
				"\ta [label=<\n" +
				"\t<table border=\"1\" cellborder=\"1\" cellspacing=\"0\">\n\n" +
				"		<tr><td><b>Nombre</b></td><td><b>Valor</b></td></tr>\n" +
				"		<tr><td><b>sb_nombre_hd</b></td><td>" + nombrePart + "</td></tr>\n" +
				"		<tr><td><b>sb_arbol_virtual_count</b></td><td><b>" + strconv.FormatInt(superbloque.SBavdCount, 10) + "</b></td></tr>\n" +
				"		<tr><td><b>sb_detalle_directorio_count</b></td><td><b>" + strconv.FormatInt(superbloque.SBddCount, 10) + "</b></td></tr>\n" +
				"		<tr><td><b>sb_inodos_count</b></td><td><b>" + strconv.FormatInt(superbloque.SBinodosCount, 10) + "</b></td></tr>\n" +
				"		<tr><td><b>sb_bloques_count</b></td><td><b>" + strconv.FormatInt(superbloque.SBbloquesCount, 10) + "</b></td></tr>\n" +
				"		<tr><td><b>sb_arbol_virtual_free</b></td><td><b>" + strconv.FormatInt(superbloque.SBavdFree, 10) + "</b></td></tr>\n" +
				"		<tr><td><b>sb_detalle_directorio_free</b></td><td><b>" + strconv.FormatInt(superbloque.SBddFree, 10) + "</b></td></tr>\n" +
				"		<tr><td><b>sb_inodos_free</b></td><td><b>" + strconv.FormatInt(superbloque.SBinodosFree, 10) + "</b></td></tr>\n" +
				"		<tr><td><b>sb_bloques_free</b></td><td><b>" + strconv.FormatInt(superbloque.SBbloquesFree, 10) + "</b></td></tr>\n" +
				"		<tr><td><b>sb_date_creacion</b></td><td><b>" + string(superbloque.SBdateCreacion[:]) + "</b></td></tr>\n" +
				"		<tr><td><b>sb_date_ultimo_montaje</b></td><td><b>" + string(superbloque.SBdateLastMount[:]) + "</b></td></tr>\n" +
				"		<tr><td><b>sb_montajes_count</b></td><td><b>" + strconv.FormatInt(superbloque.SBmontajesCount, 10) + "</b></td></tr>\n" +
				"		<tr><td><b>sb_ap_bitmap_arbol_directorio</b></td><td><b>" + strconv.FormatInt(superbloque.SBapBAVD, 10) + "</b></td></tr>\n" +
				"		<tr><td><b>sb_ap_arbol_directorio</b></td><td><b>" + strconv.FormatInt(superbloque.SBapAVD, 10) + "</b></td></tr>\n" +
				"		<tr><td><b>sb_ap_bitmap_detalle_directorio</b></td><td><b>" + strconv.FormatInt(superbloque.SBapBDD, 10) + "</b></td></tr>\n" +
				"		<tr><td><b>sb_ap_detalle_directorio</b></td><td><b>" + strconv.FormatInt(superbloque.SBapDD, 10) + "</b></td></tr>\n" +
				"		<tr><td><b>sb_ap_bitmap_tabla_inodo</b></td><td><b>" + strconv.FormatInt(superbloque.SBapBINODO, 10) + "</b></td></tr>\n" +
				"		<tr><td><b>sb_ap_tabla_inodo</b></td><td><b>" + strconv.FormatInt(superbloque.SBapINODO, 10) + "</b></td></tr>\n" +
				"		<tr><td><b>sb_ap_bitmap_bloques</b></td><td><b>" + strconv.FormatInt(superbloque.SBapBBLOQUE, 10) + "</b></td></tr>\n" +
				"		<tr><td><b>sb_ap_bloques</b></td><td><b>" + strconv.FormatInt(superbloque.SBapBLOQUE, 10) + "</b></td></tr>\n" +
				"		<tr><td><b>sb_ap_log</b></td><td><b>" + strconv.FormatInt(superbloque.SBapLOG, 10) + "</b></td></tr>\n" +
				"		<tr><td><b>sb_size_struct_arbol_directorio</b></td><td><b>" + strconv.FormatInt(superbloque.SBsizeStructAVD, 10) + "</b></td></tr>\n" +
				"		<tr><td><b>sb_size_struct_detalle_directorio</b></td><td><b>" + strconv.FormatInt(superbloque.SBsizeStructDD, 10) + "</b></td></tr>\n" +
				"		<tr><td><b>sb_size_struct_inodo</b></td><td><b>" + strconv.FormatInt(superbloque.SBsizeStructINODO, 10) + "</b></td></tr>\n" +
				"		<tr><td><b>sb_size_struct_bloque</b></td><td><b>" + strconv.FormatInt(superbloque.SBsizeStructBLOQUE, 10) + "</b></td></tr>\n" +
				"		<tr><td><b>sb_first_free_bit_arbol_directorio</b></td><td><b>" + strconv.FormatInt(superbloque.SBfirstFreeBitAVD, 10) + "</b></td></tr>\n" +
				"		<tr><td><b>sb_first_free_bit_detalle_directorio</b></td><td><b>" + strconv.FormatInt(superbloque.SBfirstFreeBitDD, 10) + "</b></td></tr>\n" +
				"		<tr><td><b>sb_first_free_bit_tabla_inodo</b></td><td><b>" + strconv.FormatInt(superbloque.SBfirstFreeBitINODO, 10) + "</b></td></tr>\n" +
				"		<tr><td><b>sb_first_free_bit_bloques</b></td><td><b>" + strconv.FormatInt(superbloque.SBfirstFreeBitBLOQUE, 10) + "</b></td></tr>\n" +
				"		<tr><td><b>sb_magic_num</b></td><td><b>" + strconv.FormatInt(superbloque.SBmagicNum, 10) + "</b></td></tr>\n" +
				"\t</table>>];\n" +
				"}"

		} else {
			fmt.Println("La particion indica no esta mon")
		}
	} else {
		fmt.Println("EL disco proporcionado no esta montado")
	}
	return dot
}

func graficarTreeDirectorio(vd string) string {
	var idDisco byte
	idDisco = vd[2]
	idDisco2 := idDisco - 97
	idP, _ := strconv.Atoi(vd[3:])
	idP--

	dot := ""

	dot += "digraph treedd {\n" +
		"\tnode [shape=record];\n\n"

	if arregloMount[idDisco2].estado == 1 {
		if arregloMount[idDisco2].discos[idP].estado == 1 {

			inicioPart := arregloMount[idDisco2].discos[idP].Partstart
			rutaDisco := arregloMount[idDisco2].Ruta
			superBloque := obtenerSB(rutaDisco, inicioPart)

			dot += nodosTreeDirectorio(rutaDisco, superBloque.SBapAVD)

			dot += "}"
		} else {
			fmt.Println("La particion indicada no esta montada")
		}
	} else {
		fmt.Println("EL disco proporcionado no esta montado")
	}

	return dot
}

func nodosTreeDirectorio(rutaDisco string, pos int64) string {
	arbol := obtenerAVD(rutaDisco, pos)
	nombre := ""
	dot := ""
	for i := 0; i < len(arbol.AVDnombreDirectorio); i++ {
		if arbol.AVDnombreDirectorio[i] != 0 {
			nombre += string(arbol.AVDnombreDirectorio[i])
		}
	}
	dot += "\tstruct" + strconv.FormatInt(pos, 10) + " [label=\"{ " + nombre + " |{"
	for i := 0; i < len(arbol.AVDapArraySub); i++ {
		dot += "<f" + strconv.Itoa(i) + ">|"
	}
	dot += "<f6>|<f7>}}\"];\n\n"

	for i := 0; i < len(arbol.AVDapArraySub); i++ {
		if arbol.AVDapArraySub[i] != -1 {
			dot += "\tstruct" + strconv.FormatInt(pos, 10) + ":f" + strconv.Itoa(i) + " -> " +
				"struct" + strconv.FormatInt(arbol.AVDapArraySub[i], 10) + ";\n"
		}
	}

	if arbol.AVDapAVD != -1 {
		dot += "\n\tstruct" + strconv.FormatInt(pos, 10) + ":f7 -> " +
			"struct" + strconv.FormatInt(arbol.AVDapAVD, 10) + "\n\n\n"
	}

	for i := 0; i < len(arbol.AVDapArraySub); i++ {
		if arbol.AVDapArraySub[i] != -1 {
			dot += nodosTreeDirectorio(rutaDisco, arbol.AVDapArraySub[i])
		}
	}

	if arbol.AVDapAVD != -1 {
		dot += nodosTreeDirectorio(rutaDisco, arbol.AVDapAVD)
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
	/*partar := strings.Split(path, "/")
	for i := 0; i < len(partar); i++ {
		fmt.Println(partar[i])
	}*/
	for i := 0; i < len(pathPart)-1; i++ {
		carpeta += pathPart[i]
	}
	archivo = pathPart[len(pathPart)-1]
	return carpeta, archivo
}
