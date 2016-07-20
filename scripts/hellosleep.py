#!/usr/bin/env python
import sys
import time

for i in range(10):
    print "hello" + str(i)
    sys.stderr.write("from the other side\n")
    sys.stdout.flush()
    time.sleep(1)