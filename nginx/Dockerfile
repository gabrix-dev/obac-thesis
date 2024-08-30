FROM nginx:latest

# Copy the Nginx config
COPY default /etc/nginx/sites-available/default
COPY default.conf /etc/nginx/conf.d/default.conf
COPY certs/ /etc/nginx/ssl

# Run the Nginx server
CMD ["/usr/sbin/nginx", "-g", "daemon off;"]
