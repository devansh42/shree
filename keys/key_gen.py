#!/usr/bin/python3

"""
This script generates ssh keys for testing 
"""

import os
keys=["ca_host_key","ca_user_key","id_host","id_user"]

#gen_keys generates key 
def gen_keys(keys):
    for x in keys:
        os.system("ssh-keygen -t rsa -b 4096 -P '' %s "%x)


keyp=[("ca_host_key","id_host"),("ca_user_key","id_user")]
def gen_cert(keyp):
    i=0
    cmd="ssh-keygen -s %s -I %s %s -n %s %s"
    for x in keyp:
        signer,key=x
        if i==0:
            os.system(cmd%(signer,"localhost","-h","localhost",key+".pub"))
        else:
            os.system(cmd%(signer,"devansh42","","devansh42",key+".pub"))
        i+=1
    #Lets show this certificate
    for x in keyp:
        _,k=x
        os.system("ssh-keygen -L -f %s-cert.pub"%k) #displaying certificates
        
gen_cert(keyp)

