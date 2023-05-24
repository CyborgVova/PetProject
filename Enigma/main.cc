#include "application.h"

int main() {
  s21::Application app;
  printf("%c\033[2J", 27);
  app.choose_start_menu();
  return 0;
}
