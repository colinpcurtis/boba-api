sudo docker run \
    -e MONGO_URL=$MONGO_URL\
    -p 8000:8000 \
    --restart unless-stopped \
    boba-api:latest
