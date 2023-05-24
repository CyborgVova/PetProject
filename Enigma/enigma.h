#ifndef SRC_ENIGMA_ENIGMA_H_
#define SRC_ENIGMA_ENIGMA_H_
#include <iostream>

#include "reflector.h"
#include "rotor.h"

namespace s21 {
class Enigma {
 public:
  Enigma(size_t number_rotors = 0);
  Enigma(s21::Enigma const &other);
  Enigma(s21::Enigma &&other);
  ~Enigma();
  void operator=(s21::Enigma const &other);
  void operator=(s21::Enigma &&other);
  void initial_state();
  s21::Rotor *get_rotors();
  s21::Reflector *get_reflector();
  std::vector<char> &get_state();
  void set_state(std::string str);
  char coder(char &ch);

 private:
  void clear_enigma();
  void move_enigma(s21::Enigma &other);
  void first_step(char &ch);
  void rotor_to_reflector(char &ch);
  void after_reflector(char &ch);
  void rotor_back_and_out(char &ch);
  bool check_state(char ch);
  void rotors_rotation();
  char counter_how_add(int number);

  s21::Rotor *rotors_;
  s21::Reflector *reflector_;
  std::vector<char> state_;
  size_t number_rotors_;
};
}  // namespace s21
#endif //SRC_ENIGMA_ENIGMA_H_
