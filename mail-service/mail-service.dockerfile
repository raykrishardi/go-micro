FROM alpine:latest

RUN mkdir /app

COPY mailApp /app
# Copy to /templates because the path is templateToRender := "./templates/mail.html.gohtml" and the default work dir is /
COPY templates /templates

CMD ["/app/mailApp"]