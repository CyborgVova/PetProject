#ifndef SRC_ENIGMA_ROTOR_H_
#define SRC_ENIGMA_ROTOR_H_
#include <iostream>
#include <map>
#include <vector>

namespace s21 {
static const std::vector<char> alphabet = {
    'A', 'B', 'C', 'D', 'E', 'F', 'G', 'H', 'I', 'J', 'K', 'L', 'M',
    'N', 'O', 'P', 'Q', 'R', 'S', 'T', 'U', 'V', 'W', 'X', 'Y', 'Z'};

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
