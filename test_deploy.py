"""
This file deploys shree for integration testing puporses 
"""

import os
import docker

cli=docker.from_env()
allContainers=map(lambda c:c.name,cli.containers.list(all=True))
net_backend_services="shree-backend-services"
net_server="shree-server-network"

def isDuplicate(name):
    if name in allContainers:
        print("Container",name," already exists")
        return True
    return False


volume_value={
    "bind":"/",
    "mode":"rw"
}



"""
name - name of the container
"""
def deploy_redis(name):
    
    if isDuplicate(name):return
    vol=dict()
    vol["shree-redis-vol"]=volume_value
    cli.containers.run(name=name,detach=True,image="redis:alpine",volume=vol,network=net_backend_services)


def deploy_ca(name):
    if isDuplicate(name):return
    vol=dict()
    vol["shree-ca-vol"]=volume_value
    cli.containers.run(name=name,detach=True,volume=vol,image="devansh42/shree-ca",network=net_backend_services)


def deploy_backend(name):
    if isDuplicate(name):return
    vol=dict()
    vol["shree-backend-vol"]=volume_value
    cli.containers.run(name=name,detach=True,volume=vol ,image="devansh42/shree-backend",network=net_backend_services)
   

def deploy_remote_server(name):
    if isDuplicate(name):return
    vol=dict()
    vol["shree-remote-vol"]=volume_value
    cli.containers.run(name=name,detach=True,volume=vol,image="devansh42/shree-serv",network=net_server)
   



def make_vols():
    names=["shree-backend-vol","shree-remote-vol","shree-ca-vol","shree-redis-vol"]
    l=[x.name for x in cli.volumes.list()]

    for x in names:
        if x in l:
            print("Volume",x," already exists")
            continue
        cli.volumes.create(name=x)

"""
It tries to make 2 networks one for backend services and other for ssh server
"""
def make_network():

    nets=[net_backend_services, #For backend,redis and ca
    net_server] #For ssh server and backend
    l=map(lambda n:n.name,cli.networks.list())
    for x in nets:
        if x in l:
            print("Network",x," already exists")
            continue #network previously exists

        cli.networks.create(name=x,attachable=True)