#!/usr/bin/python3

"""
This file build script builds docker images ca,server and backend container
This file also compiles the cli
"""
import os

try:
    import docker
except ModuleNotFoundError as m:
    print("Couldn't found docker module, installing ")
    status=os.system("pip3 install docker")    
    if status!=0:
        print("Couldn't install docker module")
        exit(1)
    else:
        import docker

import docker

rootDir=os.path.abspath(os.curdir)
cli=docker.from_env()

"""
Builds necssary images
"""
def build_images():
    dirs=[("ca","shree-ca"),("back","shree-back"),("exe/serv","shree-serv")]
    builds=[build_img(p,t) for (p,t) in dirs]

def build_img(dockerfile,tag):
    print( "Building %s with tag %s ..."%(dockerfile,tag))
    return cli.images.build(path=dockerfile,tag=tag)

