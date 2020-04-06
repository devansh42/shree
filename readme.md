# Shree - An open soure Solution for Remote Tunneling

## Features - 

1. Local Port Forwarding Tunneling 
2. Remote Port Forwarding Tunneling

### **Local Port Forwarding Tunneling** - This feature provides port forwarding on local system, you can join ports. i.e. If a server is running at port X, then you can expose another port Y on the system and can redirect all the traffic from port Y to X.
> localhost:**Y** -> localhost:**X**

### **Remote Port Forwarding Tunneling** - This feature provides remote port forwarding using SSH Protocol. You can expose a local port to public port. Public port will be completely random due to server implementation. e.g. If you are running an application on local port X, you can expose that port to the remote server (The server on which tunneling server is running), let Z be the exposed port on remote machine then all the traffic on the remote port will be tunneled to port X.
> remote_server:**Z** -> localhost:**X**