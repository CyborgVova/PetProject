#ifndef SRC_ENIGMA_APPLICATION_H_
#define SRC_ENIGMA_APPLICATION_H_
#include <climits>
#include <fstream>
#include <iostream>
#include <thread>

#include "enigma.h"

namespace s21{
class Application {
 public:
  Application() {}
  ~Application() {}

    void choose_start_menu();
 private:
    void to_encoded(s21::Enigma &enigma);
    void start_menu();
    void job_menu();
    void choose_job_menu(s21::Enigma &enigma);
    void to_decoded(s21::Enigma &enigma);
    std::string encoded_to_decoded(std::string path_to_file);
    void save_cfg(s21::Enigma &enigma, std::string path_to_file);
    s21::Enigma create_enigma();
    void erase_bracket(std::string &path_to_file);
    void clear_cin();
    s21::Enigma get_enigma_from_file(std::string path_to_file);
    void something_went_wrong(std::exception const &e);
};
} // namespace s21
#endif // SRC_ENIGMA_APPLICATION_H_
