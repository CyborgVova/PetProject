#ifndef SRC_ENIGMA_REFLECTOR_H_
#define SRC_ENIGMA_REFLECTOR_H_

#include "rotor.h"

namespace s21 {
static const std::set<char> alphabet_set = {s21::alphabet.begin(),
                                            s21::alphabet.end()};

class Reflector {
 public:
  Reflector();
  Reflector(s21::Reflector const &other);
  Reflector(s21::Reflector &&other);
  ~Reflector();
  void operator=(s21::Reflector const &other);
  void operator=(s21::Reflector &&other);
  std::map<char, char> get_reflector();
  void set_reflector(std::string str);

 private:
  void make_reflector();

  std::map<char, char> reflector_;
};
}  // namespace s21
#endif  // SRC_ENIGMA_REFLECTOR_H_
