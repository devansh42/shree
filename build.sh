#!/bin/sh
#This files build docker images and push them to docker hub
docker build --no-cache -t devansh42/shree-backend ./back
docker build --no-cache	 -t devansh42/shree-ca ./ca
docker build  --no-cahce  -t devansh42/shree-serv ./exe/serv
