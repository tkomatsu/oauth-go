TARGET = oauth

all:
	go build -o $(TARGET)

clean:
	rm $(TARGET)

.PHONY: all clean
