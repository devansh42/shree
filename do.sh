apt update
apt install -y docker.io
git clone -b test https://github.com/devansh42/shree.git
cd shree
curl -L "https://github.com/docker/compose/releases/download/1.25.4/docker-compose-$(uname -s)-$(uname -m)" -o /usr/local/bin/docker-compose
chmod +x /usr/local/bin/docker-compose
