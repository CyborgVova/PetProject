#ifndef SRC_ENIGMA_REFLECTOR_H_
#define SRC_ENIGMA_REFLECTOR_H_
#include <iostream>
#include <map>
#include <set>
#include <time.h>

namespace s21 {
static const std::set<char> alphabet_set = {
    'A', 'B', 'C', 'D', 'E', 'F', 'G', 'H', 'I', 'J', 'K', 'L', 'M',
    'N', 'O', 'P', 'Q', 'R', 'S', 'T', 'U', 'V', 'W', 'X', 'Y', 'Z'};

class Reflector {
 public:
  Reflector();
  Reflector(s21::Reflector const &other);
  Reflector(s21::Reflector &&other);
  ~Reflector();
  void operator=(s21::Reflector const &other);
  void operator=(s21::Reflector &&other);
  void make_reflector();
  std::map<char, char> get_reflector();
  void set_reflector(std::string str);

 private:
  std::map<char, char> reflector_;
};
}  // namespace s21
#endif  // SRC_ENIGMA_REFLECTOR_H_