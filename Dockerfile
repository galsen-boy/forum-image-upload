FROM golang:latest
# cree un repertoire de travail
WORKDIR /app
# copy des dependance dans /app
COPY . ./
#Bulder pour obtenir un executable 
RUN go build -o forum
#Le port a ecouter
EXPOSE 8080
#Route vers l'executation ici main.go
CMD [ "./forum" ]