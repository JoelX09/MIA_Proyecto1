package main

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"log"
	"os"
	"strconv"
	"time"
	"unsafe"
)

/*
obtenes el numero de estructuras n con la formula, luego tendrias que tu cantidad
de avd va a ser igual a n, DD = n, Inodos = 5*n, bloques = 20*n y bitacoras = n
*/

type sb struct {
	SBnombreHd           [16]byte
	SBavdCount           int64
	SBddCount            int64
	SBinodosCount        int64
	SBbloquesCount       int64
	SBavdFree            int64
	SBddFree             int64
	SBinodosFree         int64
	SBbloquesFree        int64
	SBdateCreacion       [19]byte
	SBdateLastMount      [19]byte
	SBmontajesCount      int64
	SBapBAVD             int64
	SBapAVD              int64
	SBapBDD              int64
	SBapDD               int64
	SBapBINODO           int64
	SBapINODO            int64
	SBapBBLOQUE          int64
	SBapBLOQUE           int64
	SBapLOG              int64
	SBsizeStructAVD      int64
	SBsizeStructDD       int64
	SBsizeStructINODO    int64
	SBsizeStructBLOQUE   int64
	SBfirstFreeBitAVD    int64
	SBfirstFreeBitDD     int64
	SBfirstFreeBitINODO  int64
	SBfirstFreeBitBLOQUE int64
	SBmagicNum           int64
}

type avd struct {
	AVDfechaCreacion    [19]byte
	AVDnombreDirectorio [20]byte
	AVDapArraySub       [6]int64
	AVDapDetalleDir     int64
	AVDapAVD            int64
	AVDproper           int8
	AVDgid              int8
	AVDperm             [3]byte
}

type dd struct {
	DDarrayFiles [5]arregloArchivos // <-----------------------------------Preguntar
	DDapDD       int64
}

type arregloArchivos struct {
	DDfileNombre           [20]byte
	DDfileApInodo          int64
	DDfileDateCreacion     [19]byte
	DDfileDateModificacion [19]byte
}

type inodo struct {
	IcountInodo            int64
	IsizeArchivo           int64
	IcountBloquesAsignados int64
	IarrayBloques          [4]int64
	IapIndirecto           int64
	IidProper              int8
	AVDgid                 int8
	AVDperm                [3]byte
}

type bloque struct {
	DBdata [25]byte
}

type bitacora struct {
	LOGtipoOperacion [10]byte
	LOGtipo          byte
	LOGnombre        [100]byte
	LOGcontenido     [100]byte
	LOGfecha         [19]byte
}

func formatearPart(vd string, tipo string, add int64, unit byte) {
	var idDisco byte
	idDisco = vd[2]
	idDisco2 := idDisco - 97
	idP, _ := strconv.Atoi(vd[3:])
	idP--

	if arregloMount[idDisco2].estado == 1 {
		if arregloMount[idDisco2].discos[idP].estado == 1 {
			inicioPart := arregloMount[idDisco2].discos[idP].Partstart
			tamPart := arregloMount[idDisco2].discos[idP].Partsize
			rutaDisco := arregloMount[idDisco2].Ruta

			deleteFull(rutaDisco, inicioPart, tamPart)

			superBloque := sb{}
			var sizeSB int64 = int64(unsafe.Sizeof(superBloque))
			//fmt.Println("Tamano SUPERBLOQUE")
			//fmt.Println(sizeSB)
			var sizeAVD int64 = int64(unsafe.Sizeof(avd{}))
			//fmt.Println("Tamano AVD")
			//fmt.Println(sizeAVD)
			var sizeDD int64 = int64(unsafe.Sizeof(dd{}))
			//fmt.Println("Tamano DD")
			//fmt.Println(sizeDD)
			var sizeInodo int64 = int64(unsafe.Sizeof(inodo{}))
			//fmt.Println("Tamano INODO")
			//fmt.Println(sizeInodo)
			var sizeBloque int64 = int64(unsafe.Sizeof(bloque{}))
			//fmt.Println("Tamano BLOQUE")
			//fmt.Println(sizeBloque)
			var sizeBitacora int64 = int64(unsafe.Sizeof(bitacora{}))
			//fmt.Println("Tamano BITACORA")
			//fmt.Println(sizeBitacora)

			nEstructuras := (tamPart - (2 * sizeSB)) / (27 + sizeAVD + sizeDD + (5*sizeInodo + (20 * sizeBloque) + sizeBitacora))

			cantidadAVD := nEstructuras
			cantidadDD := nEstructuras
			cantidadInodos := 5 * nEstructuras
			cantidadBloques := 4 * cantidadInodos //20*nEstructuras
			//cantidadBitacoras := nEstructuras

			_, nombreDisco := descomponer(rutaDisco)

			// ----- LLENANDO EL SUPER BLOQUE -------------------------------------
			copy(superBloque.SBnombreHd[:], nombreDisco)
			superBloque.SBavdCount = cantidadAVD
			superBloque.SBddCount = cantidadDD
			superBloque.SBinodosCount = cantidadInodos
			superBloque.SBbloquesCount = cantidadBloques
			superBloque.SBavdFree = cantidadAVD
			superBloque.SBddFree = cantidadDD
			superBloque.SBinodosFree = cantidadInodos
			superBloque.SBbloquesFree = cantidadBloques

			fecha := time.Now().Format("2006-01-02 15:04:05")
			copy(superBloque.SBdateCreacion[:], fecha)
			fecha = time.Now().Format("2006-01-02 15:04:05")
			copy(superBloque.SBdateLastMount[:], fecha)

			superBloque.SBmontajesCount = 1

			superBloque.SBapBAVD = inicioPart + sizeSB
			superBloque.SBapAVD = superBloque.SBapBAVD + cantidadAVD
			superBloque.SBapBDD = superBloque.SBapAVD + (cantidadAVD * sizeAVD)
			superBloque.SBapDD = superBloque.SBapBDD + cantidadDD
			superBloque.SBapBINODO = superBloque.SBapDD + (cantidadDD * sizeDD)
			superBloque.SBapINODO = superBloque.SBapBINODO + (cantidadInodos)
			superBloque.SBapBBLOQUE = superBloque.SBapINODO + (cantidadInodos * sizeInodo)
			superBloque.SBapBLOQUE = superBloque.SBapBBLOQUE + (cantidadBloques)
			superBloque.SBapLOG = superBloque.SBapBLOQUE + (cantidadBloques * sizeBloque)

			superBloque.SBsizeStructAVD = sizeAVD
			superBloque.SBsizeStructDD = sizeDD
			superBloque.SBsizeStructINODO = sizeInodo
			superBloque.SBsizeStructBLOQUE = sizeBloque
			superBloque.SBfirstFreeBitAVD = 0
			superBloque.SBfirstFreeBitDD = 0
			superBloque.SBfirstFreeBitINODO = 0
			superBloque.SBfirstFreeBitBLOQUE = 0
			superBloque.SBmagicNum = 201403975
			// --------------------------------------------------------------------

			// ----- ABRO ARCHIVO PARA ESCRIBIR EL SB -----------------------------
			escribirSuperBloque(rutaDisco, inicioPart, superBloque) // <-------------- Comprobar
			// --------------------------------------------------------------------
			escribirStructInicial(rutaDisco, superBloque) // <----------------VERIFICAR
			// --------------------------------------------------------------------
			crearCarpeta(vd, "/", true)
			crearArvhi(vd, "/users.txt", true, 0, "")

		} else {
			fmt.Println("La particion indica no esta mon")
		}
	} else {
		fmt.Println("EL disco proporcionado no esta montado")
	}

}

func escribirSuperBloque(path string, inicioPart int64, super sb) {
	superBloque := super
	s := &superBloque
	file, err := os.OpenFile(path, os.O_RDWR, 0777)

	if err != nil {
		log.Fatal(err)
	}
	file.Seek(inicioPart, 0)
	var binario2 bytes.Buffer
	err1 := binary.Write(&binario2, binary.BigEndian, s)
	if err1 != nil {
		fmt.Println("SuperBloque- binary error ", err1)
	}
	escribirBytes(file, binario2.Bytes())
	file.Close()
}

func escribirStructAVD(path string, posAVD int64, arbol avd) {
	ar := arbol
	a := &ar
	file, err := os.OpenFile(path, os.O_RDWR, 0777)

	if err != nil {
		log.Fatal(err)
	}
	file.Seek(posAVD, 0)
	var binario2 bytes.Buffer
	err1 := binary.Write(&binario2, binary.BigEndian, a)
	if err1 != nil {
		fmt.Println("SuperBloque- binary error ", err1)
	}
	escribirBytes(file, binario2.Bytes())
	file.Close()
}

func escribirStructDD(path string, posDD int64, detalle dd) {
	dir := detalle
	d := &dir
	file, err := os.OpenFile(path, os.O_RDWR, 0777)

	if err != nil {
		log.Fatal(err)
	}
	file.Seek(posDD, 0)
	var binario2 bytes.Buffer
	err1 := binary.Write(&binario2, binary.BigEndian, d)
	if err1 != nil {
		fmt.Println("SuperBloque- binary error ", err1)
	}
	escribirBytes(file, binario2.Bytes())
	file.Close()
}

func escribirStructINODO(path string, posINO int64, ino inodo) {
	inod := ino
	i := &inod
	file, err := os.OpenFile(path, os.O_RDWR, 0777)

	if err != nil {
		log.Fatal(err)
	}
	file.Seek(posINO, 0)
	var binario2 bytes.Buffer
	err1 := binary.Write(&binario2, binary.BigEndian, i)
	if err1 != nil {
		fmt.Println("SuperBloque- binary error ", err1)
	}
	escribirBytes(file, binario2.Bytes())
	file.Close()
}

func escribirStructBLOQUE(path string, posBLOQUE int64, bloquedato bloque) {
	bd := bloquedato
	b := &bd
	file, err := os.OpenFile(path, os.O_RDWR, 0777)

	if err != nil {
		log.Fatal(err)
	}
	file.Seek(posBLOQUE, 0)
	var binario2 bytes.Buffer
	err1 := binary.Write(&binario2, binary.BigEndian, b)
	if err1 != nil {
		fmt.Println("SuperBloque- binary error ", err1)
	}
	escribirBytes(file, binario2.Bytes())
	file.Close()
}

func actualizarValorBitmap(path string, pos int64, valor byte) {
	file, err := os.OpenFile(path, os.O_RDWR, 0777)
	val := valor
	/*fmt.Println("Valor a meter a bitmap")
	fmt.Println(val)
	fmt.Println("\nEjecuacion pausada... Presione enter para continuar")
	fmt.Scanln()*/
	//v := &val
	if err != nil {
		log.Fatal(err)
	}
	/*fmt.Println("Nuevo valor bitmap")
	fmt.Println(val)
	fmt.Println("EN la posicion")
	fmt.Println(pos)*/
	file.Seek(pos, 0)
	var binario2 bytes.Buffer
	err1 := binary.Write(&binario2, binary.BigEndian, val)
	if err1 != nil {
		fmt.Println("ValorBit- binary error ", err1)
	}
	escribirBytes(file, binario2.Bytes())
	file.Close()
}

func escribirStructInicial(path string, superbloque sb) {
	file, err := os.OpenFile(path, os.O_RDWR, 0777)

	if err != nil {
		log.Fatal(err)
	}

	var sizeAVD int64 = int64(unsafe.Sizeof(avd{}))
	var sizeDD int64 = int64(unsafe.Sizeof(dd{}))
	var sizeInodo int64 = int64(unsafe.Sizeof(inodo{}))
	var sizeBloque int64 = int64(unsafe.Sizeof(bloque{}))
	var sizeBitacora int64 = int64(unsafe.Sizeof(bitacora{}))

	posSeek := superbloque.SBapAVD
	a := avd{}
	for i := 0; i < 6; i++ {
		a.AVDapArraySub[i] = -1
	}
	a.AVDapDetalleDir = -1
	a.AVDapAVD = -1
	a.AVDproper = -1

	for i := 0; i < int(superbloque.SBavdCount); i++ {
		file.Seek(posSeek, 0)
		var binario2 bytes.Buffer
		err1 := binary.Write(&binario2, binary.BigEndian, a)
		if err1 != nil {
			fmt.Println("SuperBloqueVacio- binary error ", err1)
		}
		escribirBytes(file, binario2.Bytes())
		posSeek = posSeek + sizeAVD
	}

	posSeek = superbloque.SBapDD
	d := dd{}
	for i := 0; i < 5; i++ {
		d.DDarrayFiles[i].DDfileApInodo = -1
	}
	d.DDapDD = -1
	for i := 0; i < int(superbloque.SBddCount); i++ {
		file.Seek(posSeek, 0)
		var binario2 bytes.Buffer
		err1 := binary.Write(&binario2, binary.BigEndian, d)
		if err1 != nil {
			fmt.Println("SuperBloqueVacio- binary error ", err1)
		}
		escribirBytes(file, binario2.Bytes())
		posSeek = posSeek + sizeDD

	}

	posSeek = superbloque.SBapINODO
	ino := inodo{}
	ino.IcountInodo = -1
	ino.IcountBloquesAsignados = -1
	for i := 0; i < 4; i++ {
		ino.IarrayBloques[i] = -1
	}
	ino.IapIndirecto = -1
	ino.IidProper = -1
	for i := 0; i < int(superbloque.SBinodosCount); i++ {
		file.Seek(posSeek, 0)
		var binario2 bytes.Buffer
		err1 := binary.Write(&binario2, binary.BigEndian, ino)
		if err1 != nil {
			fmt.Println("SuperBloqueVacio- binary error ", err1)
		}
		escribirBytes(file, binario2.Bytes())
		posSeek = posSeek + sizeInodo

	}

	posSeek = superbloque.SBapBLOQUE
	b := bloque{}
	for i := 0; i < int(superbloque.SBbloquesCount); i++ {
		file.Seek(posSeek, 0)
		var binario2 bytes.Buffer
		err1 := binary.Write(&binario2, binary.BigEndian, b)
		if err1 != nil {
			fmt.Println("SuperBloqueVacio- binary error ", err1)
		}
		escribirBytes(file, binario2.Bytes())
		posSeek = posSeek + sizeBloque

	}

	posSeek = superbloque.SBapLOG
	bi := bitacora{}
	for i := 0; i < int(superbloque.SBavdCount); i++ {
		file.Seek(posSeek, 0)
		var binario2 bytes.Buffer
		err1 := binary.Write(&binario2, binary.BigEndian, bi)
		if err1 != nil {
			fmt.Println("SuperBloqueVacio- binary error ", err1)
		}
		escribirBytes(file, binario2.Bytes())
		posSeek = posSeek + sizeBitacora
	}
	file.Close()
}

func obtenerSB(path string, pos int64) sb {
	//fallo := false
	file, err := os.Open(path)
	defer file.Close()
	if err != nil {
		//fallo = true
		panic(err)
	}

	s := sb{}
	var size int = int(unsafe.Sizeof(s))

	file.Seek(pos, 0)
	data := obtenerBytes(file, size)
	buffer := bytes.NewBuffer(data)

	//fmt.Println(data)

	err = binary.Read(buffer, binary.BigEndian, &s)
	if err != nil {
		log.Fatal("super bloque binary.Read failed", err)
	}

	//fmt.Println(m)
	return s //, fallo
}

func obtenerAVD(path string, pos int64) avd {
	//fallo := false
	file, err := os.Open(path)
	defer file.Close()
	if err != nil {
		//fallo = true
		panic(err)
	}

	s := avd{}
	var size int = int(unsafe.Sizeof(s))

	file.Seek(pos, 0)
	data := obtenerBytes(file, size)
	buffer := bytes.NewBuffer(data)

	//fmt.Println(data)

	err = binary.Read(buffer, binary.BigEndian, &s)
	if err != nil {
		log.Fatal("avd binary.Read failed", err)
	}

	//fmt.Println(m)
	return s //, fallo
}

func obtenerDD(path string, pos int64) dd {
	//fallo := false
	file, err := os.Open(path)
	defer file.Close()
	if err != nil {
		//fallo = true
		panic(err)
	}

	s := dd{}
	var size int = int(unsafe.Sizeof(s))

	file.Seek(pos, 0)
	data := obtenerBytes(file, size)
	buffer := bytes.NewBuffer(data)

	//fmt.Println(data)

	err = binary.Read(buffer, binary.BigEndian, &s)
	if err != nil {
		log.Fatal("avd binary.Read failed", err)
	}

	//fmt.Println(m)
	return s //, fallo
}

func obtenerINODO(path string, pos int64) inodo {
	//fallo := false
	file, err := os.Open(path)
	defer file.Close()
	if err != nil {
		//fallo = true
		panic(err)
	}

	s := inodo{}
	var size int = int(unsafe.Sizeof(s))

	file.Seek(pos, 0)
	data := obtenerBytes(file, size)
	buffer := bytes.NewBuffer(data)

	//fmt.Println(data)

	err = binary.Read(buffer, binary.BigEndian, &s)
	if err != nil {
		log.Fatal("avd binary.Read failed", err)
	}

	//fmt.Println(m)
	return s //, fallo
}

func obtenerBLOQUE(path string, pos int64) bloque {
	//fallo := false
	file, err := os.Open(path)
	defer file.Close()
	if err != nil {
		//fallo = true
		panic(err)
	}

	s := bloque{}
	var size int = int(unsafe.Sizeof(s))

	file.Seek(pos, 0)
	data := obtenerBytes(file, size)
	buffer := bytes.NewBuffer(data)

	//fmt.Println(data)

	err = binary.Read(buffer, binary.BigEndian, &s)
	if err != nil {
		log.Fatal("avd binary.Read failed", err)
	}

	//fmt.Println(m)
	return s //, fallo
}

/*func escribirBytesS(file *os.File, bytes []byte) {
	_, err := file.Write(bytes)
	if err != nil {
		log.Fatal(err)
	}
}*/
