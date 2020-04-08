#!/bin/sh
#This files build docker images and push them to docker hub

docker build  -t devansh42/shree-backend ./back
docker build  -t devansh42/shree-ca ./ca
docker build  -t devansh42/shree-serv ./exe/serv
 #Asking for login
 docker login

#docker push devansh42/shree-backend
#docker push devansh42/shree-ca
#docker push devansh42/shree-serv

#Now let's compile shree-service image
docker-compose build 

