#!/usr/bin/env python
import sys
import time

for i in range(5):
    print "hello", i
    sys.stdout.flush()
    time.sleep(1)

sys.exit(3)