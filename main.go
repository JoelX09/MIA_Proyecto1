package main

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"os"
	"regexp"
	"strconv"
	"strings"
)

type mbr struct {
	Mbrtam     int64
	Mbrfecha   [19]byte
	Mbrdisksig int32
	Prt        [4]partition
}

type partition struct {
	Partstatus int8
	Parttype   byte
	Partfit    [2]byte
	Partstart  int64
	Partsize   int64
	Partname   [16]byte
}

type ebr struct {
	Partstatus int8
	Partfit    [2]byte
	Partstart  int64
	Partsize   int64
	Partnext   int64
	Partname   [16]byte
}

type datoDisco struct {
	path    string
	size    int64
	name    string
	unit    byte
	typeP   byte
	fit     string
	deleteP string
	add     int
	idn     string
}

type banderaParam struct {
	pathY   bool
	sizeY   bool
	nameY   bool
	unitY   bool
	typeY   bool
	fitY    bool
	deleteY bool
	addY    bool
	idY     bool
}

var dato datoDisco
var flagP banderaParam

func main() {
	fmt.Println("Joel Obdulio Xicara Rios \n201403975")
	fmt.Println("\nSistema de Archivos LWH")
	fmt.Println("\nIngrese un comando o ingrese 'e' para salir:")
	comando := pedirComando()
	finalizar := 0

	if comando == "e" {
		finalizar = 1
	} else if comando == "x" {
		finalizar = 0
	} else {
		analizador(comando)
		finalizar = 0
	}

	for finalizar != 1 {
		fmt.Println("\nIngrese un comando o ingrese 'e' para salir:")
		comando = pedirComando()
		if comando == "e" {
			finalizar = 1
		} else if comando == "x" {
			finalizar = 0
		} else {
			analizador(comando)
			finalizar = 0
		}
	}
}

func pedirComando() string {
	reader := bufio.NewReader(os.Stdin)
	lectura, _ := reader.ReadString('\n')
	comando := strings.TrimRight(lectura, "\n") //Elimino el caracter | que acepta la cadena
	if comando == "e" {                         //Verifico si la solicitud es de salida
		fmt.Println("\nSaliendo...")
		return comando
	} else if comando[0] == '#' { //Verifico si la primera linea es comentario
		fmt.Println("\nComentario: " + comando[1:])
		return "x"
	} else if comando[len(comando)-2:] == "\\*" { //verifico si valida que el comando continua
		comando = strings.ReplaceAll(comando, "\\*", "") //Elmino \*
		temp := pedirComando()                           //Solicito la continuacion del comando
		return comando + temp
	} else {
		temp1 := strings.Split(comando, "#") //Verifico si el comando trae un comentario
		if len(temp1) > 1 {
			fmt.Println("\nComentario: " + temp1[1])
			temp1[0] = strings.TrimRight(temp1[0], " ")
		}
		return temp1[0]
	}

}

func analizador(cadena string) {
	listado := strings.Split(cadena, "\n")

	for i := 0; i < len(listado); i++ {
		comandoUnico := true
		var parametro string
		comando := listado[i]

		if comando != "" {
			if comando[0] == '#' { //Verifico si la primera linea es comentario
				fmt.Println("\nComentario: " + comando[1:])
			} else {
				temp1 := strings.Split(comando, "#")
				if len(temp1) > 1 {
					fmt.Println("\nComentario: " + temp1[1])
				}
				temp2 := temp1[0]
				temp2 = strings.TrimSuffix(temp2, " ")
				temp2 = strings.TrimSuffix(temp2, "\r")
				declaracionComando := strings.SplitN(temp2, " ", 2)
				tipo := strings.ToLower(declaracionComando[0])
				if len(declaracionComando) > 1 {
					parametro = declaracionComando[1]
					comandoUnico = false
				}

				fmt.Println("\nComando a ejecutar: \"" + tipo + "\"")
				if comandoUnico == false {
					flag := 0

					for flag != 1 {
						if parametro[len(parametro)-2:] == "\\*" {
							i++
							parametro = strings.ReplaceAll(parametro, "\\*", "") + listado[i]
							flag = 0
						} else {
							flag = 1
						}
					}

					fmt.Println("Contenido: " + parametro)
					dato = datoDisco{"", 0, "", 0, 0, "", "", 0, ""}
					flagP = banderaParam{false, false, false, false, false, false, false, false, false}
					/*dato = */ analizadorParametros(parametro, i+1)
					//fmt.Println(dato)
				}

				switch tipo {
				case "pause":
					fmt.Println("\nEjecuacion pausada... Presione enter para continuar")
					fmt.Scanln()
				case "exec":
					archivo, err := ioutil.ReadFile(dato.path)
					if err != nil {
						fmt.Printf("Error leyendo archivo: %v", err)
					}
					contenido := string(archivo)
					analizador(contenido)
				case "mkdisk":
					fmt.Println("Ruta para crear el disco: " + dato.path)
					if flagP.sizeY == true && flagP.pathY == true && flagP.nameY == true {
						fmt.Printf("Se crear el disco en la ruta: %s de tamano: %d con nombre: %s", dato.path, dato.size, dato.name)
						fmt.Println("")
						crearDisco(dato.size, dato.path, dato.name, dato.unit)
					}

				case "rmdisk":
					fmt.Println("Desea remover el diso: " + dato.path + " [y/n]")
					reader := bufio.NewReader(os.Stdin)
					lectura, _ := reader.ReadString('\n')
					eleccion := strings.TrimRight(lectura, "\n")
					if eleccion == "y" {
						err := os.Remove(dato.path)
						if err != nil {
							fmt.Printf("Error eliminando archivo: %v\n", err)
						} else {
							fmt.Println("Eliminado correctamente")
						}
					} else if eleccion == "n" {
						fmt.Println("No se eliminara el archivo")
					} else {
						fmt.Println("Confirmacion invalida.")
					}

				case "fdisk":
					fmt.Println("Ruta del disco a utilizar: " + dato.path)
					if flagP.sizeY == true && flagP.pathY == true && flagP.nameY == true {
						fmt.Printf("Se creara la particion en la ruta: %s de tamano: %d con nombre: %s", dato.path, dato.size, dato.name)
						fmt.Println("")
						adminParticion(dato, flagP)
					}

				case "mount":
					//
				case "unmount":
					//
				default:
					fmt.Println("El comando " + tipo + " no es valido. Linea: " + strconv.Itoa(i+1))

				}
			}
			fmt.Println("Presione enter para continuar")
			fmt.Scanln()
		}

	}

}

func analizadorParametros(cadena string, linea int) /*datoDisco */ {
	cadena += "#"
	estado := 0
	var parametro, contParam string
	for i := 0; i < len(cadena); i++ {
		switch estado {
		case 0:
			if cadena[i] == '-' {
				estado = 1
				parametro = ""
			} else {
				fmt.Println("Error en la linea: " + strconv.Itoa(linea) + "columna: " + strconv.Itoa(i+1))
			}
		case 1:
			if cadena[i] == '-' {
				//fmt.Println("\nParametro: \"" + parametro + "\"")
				estado = 2
			} else if cadena[i] == 32 {
				//fmt.Println("\nParametro \"" + parametro + "\" suelto en la linea: " + strconv.Itoa(linea))
				estado = 0
			} else if cadena[i] == '#' {
				//fmt.Println("\nParametro \"" + parametro + "\" suelto en la linea: " + strconv.Itoa(linea))
				//fmt.Println("Parametros analizados")
			} else {
				parametro += string(cadena[i])
				estado = 1
			}
		case 2:
			if cadena[i] == '>' {
				estado = 3
			} else {
				fmt.Println("Error en la linea: " + strconv.Itoa(linea) + "columna: " + strconv.Itoa(i+1))
			}
		case 3:
			if cadena[i] == '"' {
				contParam = ""
				estado = 5
			} else {
				contParam = ""
				contParam += string(cadena[i])
				estado = 4
			}
		case 4:
			if cadena[i] == 32 {
				//fmt.Println("Valor del Parametro: \"" + contParam + "\"")
				almacenarValor(parametro, contParam, linea)
				estado = 0
			} else if cadena[i] == '#' {
				//fmt.Println("Valor del Parametro: \"" + contParam + "\"")
				almacenarValor(parametro, contParam, linea)
				//fmt.Println("Parametros analizados")
			} else {
				contParam += string(cadena[i])
				estado = 4
			}

		case 5:
			if cadena[i] == '"' {
				estado = 6
			} else {
				contParam += string(cadena[i])
				estado = 5
			}
		case 6:
			if cadena[i] == 32 {
				//fmt.Println("Valor del Parametro: \"" + contParam + "\"")
				almacenarValor(parametro, contParam, linea)
				estado = 0
			} else if cadena[i] == '#' {
				//fmt.Println("Valor del Parametro: \"" + contParam + "\"")
				almacenarValor(parametro, contParam, linea)
				//fmt.Println("Parametros analizados")
			}

		}
	}
	//return dato
}

func almacenarValor(parametro string, contParam string, linea int) {
	valor := strings.ToLower(parametro)
	match, _ := regexp.MatchString("^id[0-9]", valor)
	if match == true {
		valor = "id"
	}
	switch valor {
	case "path":
		flagP.pathY = true
		dato.path = contParam
	case "size":
		flagP.sizeY = true
		val, _ := strconv.ParseInt(contParam, 10, 64)
		dato.size = val
	case "name":
		flagP.nameY = true
		dato.name = contParam
	case "unit":
		flagP.unitY = true
		dato.unit = contParam[0]
	case "type":
		flagP.typeY = true
		dato.typeP = contParam[0]
	case "fit":
		flagP.fitY = true
		dato.fit = contParam
	case "delete":
		flagP.deleteY = true
		dato.deleteP = contParam
	case "add":
		flagP.addY = true
		val, _ := strconv.Atoi(contParam)
		dato.add = val
	case "id":
		flagP.idY = true
		dato.idn = contParam
	default:
		fmt.Println("El parametro: " + valor + " no es valido. Linea: " + strconv.Itoa(linea))
	}
}
