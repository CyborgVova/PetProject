#ifndef SRC_ENIGMA_ROTOR_H_
#define SRC_ENIGMA_ROTOR_H_
#include <ctime>
#include <fstream>
#include <iostream>
#include <map>
#include <set>
#include <thread>
#include <vector>

namespace s21 {
static const std::vector<char> alphabet = {
    '!', '"', '#', '$',  '%', '&', '\'', '(', ')', '*', '+', ',', '-', '.',
    '/', '0', '1', '2',  '3', '4', '5',  '6', '7', '8', '9', ':', ';', '<',
    '=', '>', '?', '@',  'A', 'B', 'C',  'D', 'E', 'F', 'G', 'H', 'I', 'J',
    'K', 'L', 'M', 'N',  'O', 'P', 'Q',  'R', 'S', 'T', 'U', 'V', 'W', 'X',
    'Y', 'Z', '[', '\\', ']', '^', '_',  '`', 'a', 'b', 'c', 'd', 'e', 'f',
    'g', 'h', 'i', 'j',  'k', 'l', 'm',  'n', 'o', 'p', 'q', 'r', 's', 't',
    'u', 'v', 'w', 'x',  'y', 'z', '{',  '|', '}', '~'};

class Rotor {
 public:
  Rotor();
  Rotor(s21::Rotor const &other);
  Rotor(s21::Rotor &&other);
  ~Rotor();
  void operator=(s21::Rotor const &other);
  void operator=(s21::Rotor &&other);
  std::map<char, char> get_rotor();
  void set_rotor(std::string str);
  char get_out_char(char current_ch);
  char get_key(char ch);

 private:
  void make_rotor();

  std::map<char, char> rotor_;
};
}  // namespace s21
#endif  // SRC_ENIGMA_ROTOR_H_
