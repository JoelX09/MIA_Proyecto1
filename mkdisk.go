package main

import "os"

func crearDisco(size int, path string, name string, unit byte) {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		err = os.MkdirAll(path, 0777)
		if err != nil {
			panic(err)
		}
	}

	f, err := os.Create(path + name)
	if err != nil {
		panic(err)
	}
	f.Close()

}
