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
	} else if nombre == "disk" {
		dot := graficaDisco(id)
		generar(dot, path)
	} else if nombre == "sb" {
		dot := graficasb(id)
		generar(dot, path)
	} else if nombre == "directorio" {
		dot := graficarTreeDirectorio(id, false, true)
		generar(dot, path)
	} else if nombre == "tree_file" {
		dot := graficarTreeFile(id, true, ruta)
		generar(dot, path)
	} else if nombre == "tree_complete" {
		dot := graficarTreeDirectorio(id, true, true)
		generar(dot, path)
	} else if nombre == "tree_directorio" {
		dot := graficarTreeDirectorioUnico(id, true, false, ruta)
		generar(dot, path)
	} else if nombre == "bm_arbdir" {
		dot := graficarbitmap(id, true, false, false, false)
		pathtxt, nombretxt := descomponer(path)
		archivotxt(pathtxt, nombretxt, dot)
	} else if nombre == "bm_detdir" {
		dot := graficarbitmap(id, false, true, false, false)
		pathtxt, nombretxt := descomponer(path)
		archivotxt(pathtxt, nombretxt, dot)
	} else if nombre == "bm_inode" {
		dot := graficarbitmap(id, false, false, true, false)
		pathtxt, nombretxt := descomponer(path)
		archivotxt(pathtxt, nombretxt, dot)
	} else if nombre == "bm_block" {
		dot := graficarbitmap(id, false, false, false, true)
		pathtxt, nombretxt := descomponer(path)
		archivotxt(pathtxt, nombretxt, dot)
	} else if nombre == "bitacora" {
		dot := dotBitacora(id)
		generar(dot, path)
	} else {
		fmt.Println("El reportes solicitado es incorrecto")
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
			fecha1 := ""
			for i := 0; i < len(superbloque.SBdateCreacion); i++ {
				if superbloque.SBdateCreacion[i] != 0 {
					fecha1 += string(superbloque.SBdateCreacion[i])
				}
			}
			if fecha1 == "" {
				fecha1 = "Loss"
			}

			fecha2 := ""
			for i := 0; i < len(superbloque.SBdateLastMount); i++ {
				if superbloque.SBdateLastMount[i] != 0 {
					fecha2 += string(superbloque.SBdateLastMount[i])
				}
			}
			if fecha2 == "" {
				fecha2 = "Loss"
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
				"		<tr><td><b>sb_date_creacion</b></td><td><b>" + fecha1 + "</b></td></tr>\n" +
				"		<tr><td><b>sb_date_ultimo_montaje</b></td><td><b>" + fecha2 + "</b></td></tr>\n" +
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

func graficarTreeDirectorio(vd string, treeComplete bool, conInodo bool) string {
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

			dot += nodosTreeDirectorio(rutaDisco, superBloque.SBapAVD, treeComplete, conInodo)

			dot += "}"
		} else {
			fmt.Println("La particion indicada no esta montada")
		}
	} else {
		fmt.Println("EL disco proporcionado no esta montado")
	}

	return dot
}

func graficarTreeDirectorioUnico(vd string, treeComplete bool, conInodo bool, ruta string) string {
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

			//dot += nodosTreeDirectorio(rutaDisco, superBloque.SBapAVD, treeComplete)

			path1 := strings.TrimPrefix(ruta, "/")
			path2 := strings.TrimSuffix(path1, "/")
			pathPart := strings.Split(path2, "/")
			//fmt.Println(nombre)

			encontrado := false
			posEncontrado := superBloque.SBapAVD
			//listaCarpetas := list.New()

			type carpeta struct {
				nombreC string
				posC    int64
			}
			for i := 0; i < len(pathPart); i++ {
				fmt.Println(pathPart[i])
				encontrado, posEncontrado = buscarDir(posEncontrado, pathPart[i], rutaDisco)
				if encontrado == false {
					break
				}
			}

			if encontrado == true {
				dot += nodosTreeDirectorio(rutaDisco, posEncontrado, treeComplete, conInodo)
			} else {
				fmt.Println("Carpetas inexistentes en la ruta proporcionada")
			}

			dot += "}"
		} else {
			fmt.Println("La particion indicada no esta montada")
		}
	} else {
		fmt.Println("EL disco proporcionado no esta montado")
	}
	return dot
}

func graficarTreeFile(vd string, treeComplete bool, ruta string) string {
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

			//dot += nodosTreeDirectorio(rutaDisco, superBloque.SBapAVD, treeComplete)
			path, nombre := descomponer(ruta)
			path1 := strings.TrimPrefix(path, "/")
			path2 := strings.TrimSuffix(path1, "/")
			pathPart := strings.Split(path2, "/")
			//fmt.Println(nombre)

			encontrado := false
			posEncontrado := superBloque.SBapAVD
			listaCarpetas := list.New()

			type carpeta struct {
				nombreC string
				posC    int64
			}
			for i := 0; i < len(pathPart); i++ {
				fmt.Println(pathPart[i])
				encontrado, posEncontrado = buscarDir(posEncontrado, pathPart[i], rutaDisco)
				if encontrado == false {
					break
				} else {
					nuevaCarpeta := carpeta{}
					nuevaCarpeta.nombreC = pathPart[i]
					nuevaCarpeta.posC = posEncontrado
					listaCarpetas.PushBack(nuevaCarpeta)
				}
			}

			if encontrado == true {
				ultima := listaCarpetas.Back().Value.(carpeta)
				ultimaCarpeta := obtenerAVD(rutaDisco, ultima.posC)

				archivoEncontrado, archivoPos := buscarArchivo(rutaDisco, ultimaCarpeta.AVDapDetalleDir, nombre)

				if archivoEncontrado == true {

					arbolraiz := obtenerAVD(rutaDisco, superBloque.SBapAVD)
					nombrest := ""

					for i := 0; i < len(arbolraiz.AVDnombreDirectorio); i++ {
						if arbolraiz.AVDnombreDirectorio[i] != 0 {
							nombrest += string(arbolraiz.AVDnombreDirectorio[i])
						}
					}
					dot += "\tstruct" + strconv.FormatInt(superBloque.SBapAVD, 10) + " [label=\"{ " + nombrest + " |{"

					for i := 0; i < 6; i++ {
						dot += "<f" + strconv.Itoa(i) + ">|"
					}

					dot += "<f6>|<f7>}}\"];\n\n"

					hijoRaiz := listaCarpetas.Front().Value.(carpeta)

					dot += "\tstruct" + strconv.FormatInt(superBloque.SBapAVD, 10) + " -> " + "struct" + strconv.FormatInt(hijoRaiz.posC, 10)

					nombrest = ""
					for ele := listaCarpetas.Front(); ele != nil; ele = ele.Next() {
						carpetaGrap := ele.Value.(carpeta)
						arbol := obtenerAVD(rutaDisco, carpetaGrap.posC)
						nombrest = ""

						for i := 0; i < len(arbol.AVDnombreDirectorio); i++ {
							if arbol.AVDnombreDirectorio[i] != 0 {
								nombrest += string(arbol.AVDnombreDirectorio[i])
							}
						}
						dot += "\tstruct" + strconv.FormatInt(carpetaGrap.posC, 10) + " [label=\"{ " + nombrest + " |{"
						for i := 0; i < len(arbol.AVDapArraySub); i++ {
							dot += "<f" + strconv.Itoa(i) + ">|"
						}

						dot += "<f6>|<f7>}}\"];\n\n"

						if ele.Next() != nil {
							proximo := ele.Next().Value.(carpeta)
							dot += "\tstruct" + strconv.FormatInt(carpetaGrap.posC, 10) + " -> " + "struct" + strconv.FormatInt(proximo.posC, 10)
						}
					}
					padreInodo := listaCarpetas.Back().Value.(carpeta)
					dot += "\tstructarchivo [label=\"{ " + nombre + "}\"];\n\n"
					dot += "\tstruct" + strconv.FormatInt(padreInodo.posC, 10) + " -> structarchivo"
					dot += "\tstructarchivo -> struct" + strconv.FormatInt(archivoPos, 10)
					dot += graficarInodo(archivoPos, rutaDisco)
				} else {
					fmt.Println("El archivo no existe")
				}

			} else {
				fmt.Println("Carpetas inexistentes en la ruta proporcionada")
			}

			dot += "}"
		} else {
			fmt.Println("La particion indicada no esta montada")
		}
	} else {
		fmt.Println("EL disco proporcionado no esta montado")
	}

	return dot
}

func buscarArchivo(rutaDisco string, pos int64, nombre string) (bool, int64) {
	nuevoDD := obtenerDD(rutaDisco, pos)
	var nombreb [20]byte
	copy(nombreb[:], nombre)
	var posInodo int64
	encontrado := false
	for i := 0; i < 5; i++ {
		if nuevoDD.DDarrayFiles[i].DDfileNombre == nombreb {
			posInodo = nuevoDD.DDarrayFiles[i].DDfileApInodo
			encontrado = true
			break
		}
	}

	if encontrado == false {
		if nuevoDD.DDapDD != -1 {
			encontrado, posInodo = buscarArchivo(rutaDisco, nuevoDD.DDapDD, nombre)
		}
	}

	return encontrado, posInodo
}

func nodosTreeDirectorio(rutaDisco string, pos int64, treeComplete bool, conInodo bool) string {
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
	nombre = ""
	dot += "<f6>|<f7>}}\"];\n\n"

	if conInodo == true {

		for i := 0; i < len(arbol.AVDapArraySub); i++ {
			if arbol.AVDapArraySub[i] != -1 {
				dot += "\tstruct" + strconv.FormatInt(pos, 10) + ":f" + strconv.Itoa(i) + " -> " +
					"struct" + strconv.FormatInt(arbol.AVDapArraySub[i], 10) + ";\n"
			}
		}
	}
	if treeComplete == true {
		if arbol.AVDapDetalleDir != -1 {
			dot += "\tstruct" + strconv.FormatInt(pos, 10) + ":f6 -> " +
				"struct" + strconv.FormatInt(arbol.AVDapDetalleDir, 10) + ";\n"
		}
	}

	if arbol.AVDapAVD != -1 {
		dot += "\n\tstruct" + strconv.FormatInt(pos, 10) + ":f7 -> " +
			"struct" + strconv.FormatInt(arbol.AVDapAVD, 10) + "\n\n\n"
	}

	if conInodo == true {
		for i := 0; i < len(arbol.AVDapArraySub); i++ {
			if arbol.AVDapArraySub[i] != -1 {
				dot += nodosTreeDirectorio(rutaDisco, arbol.AVDapArraySub[i], treeComplete, conInodo)
			}
		}
	}

	if treeComplete == true {
		if arbol.AVDapDetalleDir != -1 {
			dot += graficarDD(arbol.AVDapDetalleDir, rutaDisco, conInodo)
		}
	}

	if arbol.AVDapAVD != -1 {
		dot += nodosTreeDirectorio(rutaDisco, arbol.AVDapAVD, treeComplete, conInodo)
	}

	return dot
}

func graficarDD(posDD int64, rutaDisco string, conInodo bool) string {
	directorio := obtenerDD(rutaDisco, posDD)

	dot := ""
	dot += "\tstruct" + strconv.FormatInt(posDD, 10) + " [label=\"{"
	for i := 0; i < 5; i++ {
		dot += "<f" + strconv.Itoa(i) + "> "
		if directorio.DDarrayFiles[i].DDfileApInodo != -1 {
			nombre := ""
			for j := 0; j < len(directorio.DDarrayFiles[i].DDfileNombre); j++ {
				if directorio.DDarrayFiles[i].DDfileNombre[j] != 0 {
					nombre += string(directorio.DDarrayFiles[i].DDfileNombre[j])
				}
			}
			dot += nombre + " |"
		} else {
			dot += " |"
		}

	}
	dot += "<f5> Indi }\"];\n\n"

	if conInodo == true {
		for i := 0; i < 5; i++ {
			if directorio.DDarrayFiles[i].DDfileApInodo != -1 {
				dot += "\tstruct" + strconv.FormatInt(posDD, 10) + ":f" + strconv.Itoa(i) + " -> struct" + strconv.FormatInt(directorio.DDarrayFiles[i].DDfileApInodo, 10) + ";\n\n"
			}
		}
	}

	if directorio.DDapDD != -1 {
		dot += "\tstruct" + strconv.FormatInt(posDD, 10) + ":f5 -> struct" + strconv.FormatInt(directorio.DDapDD, 10) + ";\n\n"
	}

	if conInodo == true {
		for i := 0; i < 5; i++ {
			if directorio.DDarrayFiles[i].DDfileApInodo != -1 {
				dot += graficarInodo(directorio.DDarrayFiles[i].DDfileApInodo, rutaDisco)
			}
		}
	}

	if directorio.DDapDD != -1 {
		dot += graficarDD(directorio.DDapDD, rutaDisco, conInodo)
	}

	return dot
}

func graficarInodo(posInodo int64, rutaDisco string) string {

	nuevoInodo := obtenerINODO(rutaDisco, posInodo)
	dot := ""
	dot += "\n\tstruct" + strconv.FormatInt(posInodo, 10) + " [label=\"{" +
		"<fo> " + strconv.FormatInt(nuevoInodo.IcountInodo, 10) + " |" +
		"<f1> " + strconv.FormatInt(nuevoInodo.IsizeArchivo, 10) + " |" +
		"<f2> " + strconv.FormatInt(nuevoInodo.IcountBloquesAsignados, 10) + " |" +
		"<f3> B1 |" +
		"<f4> B2 |" +
		"<f5> B3 |" +
		"<f6> B4 |" +
		"<f7> Inodo indirecto}\"];\n"

	for i := 0; i < len(nuevoInodo.IarrayBloques); i++ {
		if nuevoInodo.IarrayBloques[i] != -1 {
			dot += "\n\tstruct" + strconv.FormatInt(posInodo, 10) + ":f" + strconv.Itoa(i+3) + " -> struct" + strconv.FormatInt(nuevoInodo.IarrayBloques[i], 10)
		}
	}

	if nuevoInodo.IapIndirecto != -1 {
		dot += "\n\tstruct" + strconv.FormatInt(posInodo, 10) + ":f7 -> struct" + strconv.FormatInt(nuevoInodo.IapIndirecto, 10)
	}

	for i := 0; i < len(nuevoInodo.IarrayBloques); i++ {
		if nuevoInodo.IarrayBloques[i] != -1 {
			nuevoBloque := obtenerBLOQUE(rutaDisco, nuevoInodo.IarrayBloques[i])
			cont := ""
			for j := 0; j < 25; j++ {
				cont += string(nuevoBloque.DBdata[j])
			}
			dot += "\n\tstruct" + strconv.FormatInt(nuevoInodo.IarrayBloques[i], 10) + " [label=\" " + cont + " \"];"
		}
	}

	if nuevoInodo.IapIndirecto != -1 {
		dot += graficarInodo(nuevoInodo.IapIndirecto, rutaDisco)
	}

	return dot
}

func graficarbitmap(vd string, bavd bool, bdd bool, binodo bool, bbloque bool) string {
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
			superBloque := obtenerSB(rutaDisco, inicioPart)

			var listadoBits []byte

			if bavd == true {
				listadoBits = obtenerBitmap(rutaDisco, superBloque.SBapBAVD, superBloque.SBavdCount)
			} else if bdd == true {
				listadoBits = obtenerBitmap(rutaDisco, superBloque.SBapBDD, superBloque.SBddCount)
			} else if binodo == true {
				listadoBits = obtenerBitmap(rutaDisco, superBloque.SBapBINODO, superBloque.SBinodosCount)
			} else if bbloque == true {
				listadoBits = obtenerBitmap(rutaDisco, superBloque.SBapBBLOQUE, superBloque.SBbloquesCount)
			}

			for i := 0; i < len(listadoBits); i++ {
				if listadoBits[i] == '1' {
					dot += "1 |"
				} else {
					dot += "0 |"
				}

				if (i+1)%20 == 0 {
					dot += "\n"
				}
			}

		} else {
			fmt.Println("La particion indicada no esta montada")
		}
	} else {
		fmt.Println("EL disco proporcionado no esta montado")
	}

	return dot
}

func dotBitacora(vd string) string {
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
			superBloque := obtenerSB(rutaDisco, inicioPart)
			var sizeBitacora int64 = int64(unsafe.Sizeof(bitacora{}))
			dot += graficarBitacora(rutaDisco, superBloque.SBapLOG, superBloque.SBavdCount, sizeBitacora)
		} else {
			fmt.Println("La particion indicada no esta montada")
		}
	} else {
		fmt.Println("EL disco proporcionado no esta montado")
	}

	return dot

}

func graficarBitacora(ruta string, pos int64, cantidad int64, tambit int64) string {
	dot := ""
	dot += "digraph treedd {\n" +
		"\tnode [shape=record];\n\n"
	indice := 0
	for i := 0; i < int(cantidad); i++ {
		bit := obtenerbitacora(ruta, pos)
		if bit.LOGtipo != 'x' {
			indice++
			tipoOp := ""
			for a := 0; a < 10; a++ {
				if bit.LOGtipoOperacion[a] != 0 {
					tipoOp += string(bit.LOGtipoOperacion[a])
				}
			}
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
			fecha := ""
			for a := 0; a < 19; a++ {
				if bit.LOGfecha[a] != 0 {
					fecha += string(bit.LOGfecha[a])
				}
			}
			dot += "\tstruct" + strconv.FormatInt(pos, 10) + " [label=\"{ Log " + strconv.Itoa(indice) + " |{" +
				"Tipo Operacion | Tipo | Path | Contenido | Fecha Log | Size}|{" + tipoOp + " | " + string(bit.LOGtipo) +
				" | " + path + " | " + cont + " | " + fecha + " | " + strconv.Itoa(int(bit.LOGsize)) +
				"}}\"];\n\n"

		}
		pos = pos + tambit

	}
	dot += "}"

	return dot
}

func archivotxt(path string, nombre string, datos string) {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		err = os.MkdirAll(path, 0777)
		if err != nil {
			panic(err)
		}
	}

	f, err := os.Create(path + nombre)
	if err != nil {
		panic(err)
	}
	f.Close()

	err = ioutil.WriteFile(path+nombre, []byte(datos), 0777)
	if err != nil {
		panic(err)
	}
}

func generarDOT(dot string, path string, nombre string) {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		err = os.MkdirAll(path, 0777)
		if err != nil {
			panic(err)
		}
	}

	nombreDot := nombre
	ext := strings.Split(nombre, ".")
	nombreDot = strings.ReplaceAll(nombreDot, ext[1], "dot")
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

	param := "-T" + ext[1]

	cmd := exec.Command("dot", param, pathdot, "-o", pathimg)
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
