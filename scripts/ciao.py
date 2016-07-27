#!/usr/bin/env python
import time

for i in range(5):
	prova = open("prova.txt", "a+w")
	prova.write(str(i))
	prova.close()
	time.sleep(1)