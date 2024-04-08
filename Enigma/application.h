#ifndef SRC_ENIGMA_APPLICATION_H_
#define SRC_ENIGMA_APPLICATION_H_
#define S21_INT_MAX 0x7FFFFFFF

#include "enigma.h"

namespace s21 {
class Application {
 public:
  Application() {}
  ~Application() {}
  void choose_start_menu();

 private:
  void start_menu();
  void job_menu();
  void choose_job_menu(s21::Enigma &enigma);
  void clear_cin();
  void erase_bracket(std::string &path_to_file);
  void something_went_wrong(std::exception const &e);
  s21::Enigma create_enigma();
  s21::Enigma get_enigma_from_file(std::string path_to_file);
  void save_cfg(s21::Enigma &enigma, std::string path_to_file);
  std::string encoded_to_decoded(std::string path_to_file);
  void to_encoded(s21::Enigma &enigma);
  void to_decoded(s21::Enigma &enigma);
};
}  // namespace s21
#endif  // SRC_ENIGMA_APPLICATION_H_
