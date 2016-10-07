#!/usr/bin/env python
import commands, sys
import webbrowser

def open_page(ip, port):
	print ip, port
	webbrowser.open("http://" + ip + ":" + port)

out = commands.getoutput('avahi-browse -prt _http._tcp | grep Bender-test')
outlist = out.split("\n")

for item in outlist:
    fields = item.split(";")

    if fields[2] == "IPv4" and len(fields) >= 7:
        open_page(fields[7], fields[8])
        break
