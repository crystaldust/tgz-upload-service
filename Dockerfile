FROM alpine
COPY ./tgz-upload-service /root
RUN chmod +x /root/tgz-upload-service

CMD ["/root/tgz-upload-service"]

