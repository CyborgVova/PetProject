CC			=	g++
FLAGS		=	-Wall -Wextra -Werror -std=c++17 -pedantic
SRC			=	application.cc rotor.cc reflector.cc enigma.cc main.cc

all:		enigma app

app:		enigma
			./enigma

enigma:
			$(CC) $(FLAGS) $(SRC) -o enigma

clean:
			rm -rf ./enigma *_encoded *_decoded *.cfg

re:			clean enigma app
