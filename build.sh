#!/bin/sh
#This files build docker images and push them to docker hub


docker build --no-cache -t devansh42/shree-backend ./back
docker build --no-cache	 -t devansh42/shree-ca ./ca
docker build  --no-cahce  -t devansh42/shree-serv ./exe/serv
 #Asking for login

#docker push devansh42/shree-backend
#docker push devansh42/shree-ca
#docker push devansh42/shree-serv

#Now let's compile shree-service image


