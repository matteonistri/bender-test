package main

import ("io/ioutil")

func Read(name, id string, c,z int) []byte{
	buf, _ :=ioutil.ReadFile("prova.txt")
	return buf
}