package main

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"os"
	"strconv"
	"strings"
)

type mbr struct {
	Mbrtam     int64
	Mbrfecha   [20]byte
	Mbrdisksig int8
	Prt        [4]partition
}

type partition struct {
	Partstatus byte
	Parttype   byte
	Partfit    byte
	Partstart  int64
	Partsize   int64
	Partname   [16]byte
}

type ebr struct {
	Partstatus byte
	Partfit    byte
	Partstart  int64
	Partsize   int64
	Partnext   int64
	Partname   [16]byte
}

type datoDisco struct {
	path    string
	size    int
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
	/*if err := os.Chmod("/home", 0777); err != nil {
		log.Fatal(err)
	}*/
	/*cmd := exec.Command("sudo chmod 777 /home")
	err := cmd.Run()
	if err != nil {
		panic(err)
	}*/
	//log.Printf("Command finished with error: %v", err)

	fmt.Println("Joel Obdulio Xicara Rios \n201403975")
	fmt.Println("\nSistema de Archivos LWH")
	fmt.Println("\nIngrese un comando:")

	reader := bufio.NewReader(os.Stdin)
	lectura, _ := reader.ReadString('|')
	comando := strings.TrimRight(lectura, "|") //Elimino el caracter | que acepta la cadena
	comando1 := strings.Split(comando, "\n")   //Hago split con salto si es un comando divido

	temp := comando1[0]
	if temp[0] == '#' { //Verifico si la primera linea es comentario
		fmt.Println("\nComentario: " + temp[1:])
	} else if len(comando1) == 2 {
		if temp[len(temp)-2:] == "\\*" { //verifico si valida que el comando continua
			temp = strings.ReplaceAll(temp, "\\*", "") + comando1[1] //Elmino \* y concateno el comando
			analizador(temp)
		} else {
			fmt.Println("Revisar el comando ingresado")
		}
	} else {
		temp1 := strings.Split(temp, "#")
		if len(temp1) > 1 {
			fmt.Println("\nComentario: " + temp1[1])
		}
		analizador(temp1[0])
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
				declaracionComando := strings.SplitN(temp2, " ", 2)
				tipo := strings.ToLower(declaracionComando[0])
				if len(declaracionComando) > 1 {
					parametro = declaracionComando[1]
					comandoUnico = false
				}

				fmt.Println("\nComando a ejecutar: \"" + tipo + "\"")
				if comandoUnico == false {
					if parametro[len(parametro)-2:] == "\\*" {
						i++
						parametro = strings.ReplaceAll(parametro, "\\*", "") + listado[i]
					}
					fmt.Println("Contenido: " + parametro)
					dato = datoDisco{"", 0, "", 0, 0, "", "", 0, ""}
					flagP = banderaParam{false, false, false, false, false, false, false, false, false}
					dato = analizadorParametros(parametro, i+1)
					//fmt.Println(dato)
				}

				switch tipo {
				case "pause":
					fmt.Println("\nEjecuacion pausada... Presione enter para continuar")
					fmt.Scanln()
				case "exec":
					//fmt.Println("ruta:" + dato.path)

					bytesLeidos, err := ioutil.ReadFile(dato.path)
					if err != nil {
						fmt.Printf("Error leyendo archivo: %v", err)
					}

					contenido := string(bytesLeidos)
					analizador(contenido)
				case "mkdisk":
					//fmt.Println("Ruta para crear el disco: " + dato.path)
					if flagP.sizeY == true && flagP.pathY == true && flagP.nameY == true {
						fmt.Printf("Se crear el disco en la ruta: %s de tamano: %d con nombre: %s", dato.path, dato.size, dato.name)
						fmt.Println("")
						crearDisco(dato.size, dato.path, dato.name, dato.unit)
						readFile(dato.path + dato.name)
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
					adminParticion()

				case "mount":
					//
				case "unmount":
					//

				}
			}
			fmt.Println("Presione enter para continuar")
			fmt.Scanln()
		}

	}

}

func analizadorParametros(cadena string, linea int) datoDisco {
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
				almacenarValor(parametro, contParam)
				estado = 0
			} else if cadena[i] == '#' {
				//fmt.Println("Valor del Parametro: \"" + contParam + "\"")
				almacenarValor(parametro, contParam)
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
				almacenarValor(parametro, contParam)
				estado = 0
			} else if cadena[i] == '#' {
				//fmt.Println("Valor del Parametro: \"" + contParam + "\"")
				almacenarValor(parametro, contParam)
				//fmt.Println("Parametros analizados")
			}

		}
	}
	return dato
}

func almacenarValor(parametro string, contParam string) {
	switch valor := strings.ToLower(parametro); valor {
	case "path":
		flagP.pathY = true
		dato.path = contParam
	case "size":
		flagP.sizeY = true
		val, _ := strconv.Atoi(contParam)
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
	default:
		flagP.idY = true
		dato.idn = contParam
	}
}
