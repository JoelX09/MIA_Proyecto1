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

	temp := comando1[0]
	if temp[0] == '#' { //Verifico si la primera linea es comentario
		fmt.Println("\nComentario: " + temp[1:])
	} else if len(comando1) == 2 {
		if temp[len(temp)-2:] == "\\*" { //verifico si valida que el comando continua
			temp = strings.ReplaceAll(temp, "\\*", "") //Elimino /*
			temp = temp + comando1[1]                  //Concateno el comando
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
	comandoAnalizar := strings.Split(cadena, "\n")

	for i := 0; i < len(comandoAnalizar); i++ {
		comando := comandoAnalizar[i]
		declaracionComando := strings.SplitN(comando, " ", 2)
		fmt.Println("\nComando a ejecutar: \"" + declaracionComando[0] + "\"")
		fmt.Println("Contenido: " + declaracionComando[1])
	}
}
