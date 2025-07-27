```yaml
version: '2.4'

services:

  fiber-app:
    image: ghcr.io/anotherdanger/khaira-admin:2.3
    container_name: fiber-app
    env_file:
      - .env
    networks:
      - web
      - admin-networks
    volumes:
      - product-image:/app/uploads
    depends_on:
      - admin-data
      - elasticsearch
    mem_limit: 256m
    cpus: 0.5
    restart: on-failure

  khaira-user:
    image: ghcr.io/anotherdanger/khaira-user:dev16
    container_name: khaira-user
    env_file:
      - .env
    networks:
      - web
      - admin-networks
    volumes:
      - product-image:/app/uploads
    depends_on:
      - admin-data
      - elasticsearch
    mem_limit: 256m
    cpus: 0.5
    restart: on-failure

  jwt-auth:
    image: ghcr.io/anotherdanger/khaira-jwt-auth:2.4
    container_name: jwt-auth
    env_file:
      - .env
    networks:
      - web
    mem_limit: 256m
    cpus: 0.5
    restart: on-failure

  admin-data:
    image: mysql:latest
    container_name: admin-data
    environment:
      MYSQL_DATABASE: catering
      MYSQL_ROOT_PASSWORD: ${DB_PASS}
      MYSQL_PASSWORD: ${DB_PASS}
    volumes:
      - admin-data:/var/lib/mysql
    networks:
      - admin-networks
    healthcheck:
      test: ["CMD", "mysqladmin", "ping", "-h", "localhost", "-u", "root", "-p${DB_PASS}"]
      interval: 10s
      timeout: 5s
      retries: 5
    restart: always

  elasticsearch:
    image: docker.elastic.co/elasticsearch/elasticsearch:8.13.4
    container_name: elasticsearch
    environment:
      - discovery.type=single-node
      - bootstrap.memory_lock=true
      - xpack.security.enabled=false
      - xpack.security.transport.ssl.enabled=false
      - ES_JAVA_OPTS=-Xms1g -Xmx1g
    ulimits:
      memlock:
        soft: -1
        hard: -1
    ports:
      - "9200:9200"
      - "9300:9300"
    volumes:
      - esdata:/usr/share/elasticsearch/data
    networks:
      - admin-networks
    healthcheck:
      test: ["CMD-SHELL", "curl --silent --fail http://localhost:9200/_cluster/health || exit 1"]
      interval: 10s
      timeout: 5s
      retries: 5
    restart: always

networks:
  admin-networks:
  web:
    external: true

volumes:
  admin-data:
  product-image:
  esdata:
