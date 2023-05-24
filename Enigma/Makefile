CC		=	g++
FLAGS	=	-Wall -Wextra -Werror -std=c++17 -pedantic -fsanitize=address -g
SRC		=	application.cc rotor.cc reflector.cc enigma.cc main.cc

all:	enigma app

app:	enigma
		./enigma

enigma:
		$(CC) $(FLAGS) $(SRC) -o enigma

test:

clean:
		rm -rf  enigma.dSYM *_encoded *_decoded
		rm -f ./a.exe  ./a.out ./enigma *.cfg .*.cfg

re:		clean enigma app