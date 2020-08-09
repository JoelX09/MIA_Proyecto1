package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

func main() {
	fmt.Println("Joel Obdulio Xicara Rios \n201403975")
	fmt.Println("\nSistema de Archivos LWH")
	fmt.Println("\nIngrese un comando:")

	reader := bufio.NewReader(os.Stdin)
	lectura, _ := reader.ReadString('|')
	comando := strings.TrimRight(lectura, "|") //Elimino el caracter | que acepta la cadena

	comando1 := strings.Split(comando, "\n") //Hago split con salto si es un comando divido

	if temp := comando1[0]; temp[0] == '#' { //Verifico si la primera linea es comentario
		fmt.Println("Comentario: " + temp[1:])
	} else if len(comando1) == 2 {
		if temp := comando1[0]; temp[len(temp)-2:] == "\\*" { //verifico si valida que el comando continua
			temp = strings.ReplaceAll(temp, "\\*", "")
			temp1 := temp + comando1[1]
			fmt.Println("Comando:")
			fmt.Println(temp1)
		} else {
			fmt.Println("Revisar el comando ingresado")
		}
	} else {
		fmt.Println("Revisar el comando ingresado")
	}

	//comando = strings.ReplaceAll(comando, "\\*", "")

	//comando := strings.TrimRight(comando1, "*")

	//fmt.Println("Ejecutando...")
	//fmt.Println(comando)

	//fmt.Println(comando[len(comando)-2:]) ///<-------- para el simbolo \*

	//fmt.Println(b)
}
