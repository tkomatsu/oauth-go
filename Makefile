TARGET = oauth

all:
	go build -o $(TARGET)

clean:
	rm $(TARGET)

env:
	echo 'CLIENT_ID=""' >> .env
	echo 'SECRET=""' >> .env

.PHONY: all clean
