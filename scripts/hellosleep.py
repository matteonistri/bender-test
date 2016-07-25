#!/usr/bin/env python
import sys
import time
for i in range(2):
    sys.stdout.write("%s, writing on stdout:%s\n" % (time.time(),i))
    sys.stderr.write("writing on stderr:%s\n" % i)
    sys.stdout.flush()
    sys.stderr.flush()
    time.sleep(1)
