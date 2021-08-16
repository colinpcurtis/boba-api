sudo docker run \
    -e PYTHONUNBUFFERED=1 \
    -e MONGO_URL=$MONGO_URL\
    -p 8000:8000 \
    --restart unless-stopped \
    bobaapi:latest
