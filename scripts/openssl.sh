

# check if key.pem and cert.pem exist
# if not, create them

if [ ! -f key.pem ] || [ ! -f cert.pem ]; then
  echo "Creating key.pem and cert.pem"
  openssl req -x509 -newkey rsa:4096 -keyout key.pem -out cert.pem -days 365 -nodes -subj "/C=DE/ST=Bavaria/L=Munich/O=JeschekDev/OU=Software/CN=jeschek.dev"
fi

