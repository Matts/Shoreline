FROM ubuntu

EXPOSE 443 # Caddy external webhook port
EXPOSE 8800 # Internal admin ui port

RUN make caddy