backend for shree


Environment Variables need - 

1. BACKEND_SERVER_ADDR = Contains addr of backend server
2. CA_SERVER_ADDR = Contains addr of ca server
3. REDIS_SERVER_ADDR = Contains address of local redis server

## Usage:

  -baddr string
    	(required)  Addrs to start backend server
  -caddr string
    	(required)  Addrs to CA  server
  -logdir string
    	  Directory for server logs (default "App Directory")
  -raddr string
    	(required)  Addrs to the redis server
