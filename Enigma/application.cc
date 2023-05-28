#include "application.h"

void s21::Application::choose_start_menu() {
  std::string path_to_file;
  s21::Enigma enigma;
  size_t choose;
  do {
    start_menu();
    std::cout << "Введите цифру от 1 до 3: ";
    std::cin >> choose;
    std::cout << std::endl;
    clear_cin();
  } while (choose < 1 || choose > 3);
  switch (choose) {
    case 1:
      printf("%c\033[2J", 27);
      std::cout << "Bye Bye !!!" << std::endl;
      return;
    case 2:
      std::cout << "Введите имя желаемого файла настроек: ";
      std::cin >> path_to_file;
      erase_bracket(path_to_file);
      path_to_file += ".cfg";
      enigma = create_enigma();
      printf("%c\033[2J", 27);
      try {
        save_cfg(enigma, path_to_file);
      } catch (std::exception const &e) {
        something_went_wrong(e);
        break;
      }
      std::cout << "Создан файл настроек: '" << path_to_file << "'\n"
                << std::endl;
      choose_start_menu();
      break;
    case 3:
      std::cout << "Введите имя желаемого файла настроек: ";
      std::cin >> path_to_file;
      erase_bracket(path_to_file);
      printf("%c\033[2J", 27);
      try {
        enigma = get_enigma_from_file(path_to_file);
      } catch (std::exception const &e) {
        something_went_wrong(e);
        break;
      }
      std::cout << "Настройки применены: '" << path_to_file << "'\n\n";
      choose_job_menu(enigma);
      break;
  }
}

void s21::Application::start_menu() {
  std::cout << "\t ----------------" << std::endl;
  std::cout << "\t| Начальное меню |" << std::endl;
  std::cout << "\t ----------------\n" << std::endl;

  std::cout << " ------------------------------------" << std::endl;
  std::cout << "| Для выхода из программы нажмите: 1 |" << std::endl;
  std::cout << " --------------------------------------" << std::endl;
  std::cout << "| Создать новую шифровальную машину: 2 |" << std::endl;
  std::cout << " ------------------------------------------" << std::endl;
  std::cout << "| Получить шифровальную машину из файла: 3 |" << std::endl;
  std::cout << " ------------------------------------------" << std::endl;
}

void s21::Application::job_menu() {
  std::cout << "\t --------------------------" << std::endl;
  std::cout << "\t| Меню шифровальной машины |" << std::endl;
  std::cout << "\t --------------------------\n" << std::endl;

  std::cout << " ------------------------------------" << std::endl;
  std::cout << "| Чтобы закодировать файл нажмите: 1 |" << std::endl;
  std::cout << " -------------------------------------" << std::endl;
  std::cout << "| Чтобы раскодировать файл нажмите: 2 |" << std::endl;
  std::cout << " -------------------------------------------" << std::endl;
  std::cout << "| Для возврата в предыдущее меню нажмите: 3 |" << std::endl;
  std::cout << " -------------------------------------------\n" << std::endl;
}

void s21::Application::choose_job_menu(s21::Enigma &enigma) {
  size_t choose;
  do {
    job_menu();
    std::cout << "Введите цифру от 1 до 3: ";
    std::cin >> choose;
    clear_cin();
  } while (choose < 1 || choose > 3);
  switch (choose) {
    case 1:
      try {
        to_encoded(enigma);
      } catch (std::exception const &e) {
        something_went_wrong(e);
        break;
      }
      choose_start_menu();
      break;
    case 2:
      try {
        to_decoded(enigma);
      } catch (std::exception const &e) {
        something_went_wrong(e);
        break;
      }
      choose_start_menu();
      break;
    case 3:
      printf("%c\033[2J", 27);
      choose_start_menu();
      break;
  }
}

void s21::Application::clear_cin() {
  std::cin.clear();
  std::cin.ignore(S21_INT_MAX, '\n');
}

void s21::Application::erase_bracket(std::string &path_to_file) {
  if (path_to_file[0] == '\'') {
    path_to_file.erase(path_to_file.begin());
    path_to_file.erase(path_to_file.end() - 1);
  }
}

void s21::Application::something_went_wrong(std::exception const &e) {
  printf("%c\033[2J", 27);
  std::cout << "Что-то пошло не так: " << e.what() << std::endl;
  choose_start_menu();
}

s21::Enigma s21::Application::create_enigma() {
  size_t num_rotors;
  do {
    std::cout << "Введите количество роторов ( от 1 до 255): ";
    std::cin >> num_rotors;
    clear_cin();
  } while (num_rotors < 1 || num_rotors > 255);
  s21::Enigma enigma(num_rotors);
  return enigma;
}

s21::Enigma s21::Application::get_enigma_from_file(std::string path_to_file) {
  std::ifstream fin(path_to_file, std::ios::binary);
  if (!fin.is_open()) throw std::invalid_argument("'Path is not valid'\n");
  int num_rotors;
  fin >> num_rotors;
  s21::Enigma enigma(num_rotors);
  for (int j = 0; !fin.eof(); j++) {
    std::string str;
    fin >> str;
    if (j == 0) enigma.set_state(str);
    if (j > 0 && j <= num_rotors) enigma.get_rotors()[j - 1].set_rotor(str);
    if (j > num_rotors) enigma.get_reflector()->set_reflector(str);
  }
  fin.close();
  return enigma;
}

void s21::Application::save_cfg(s21::Enigma &enigma, std::string path_to_file) {
  std::ofstream fout(path_to_file, std::ios::binary);
  if (!fout.is_open())
    throw std::invalid_argument("'Failed to create configuring file'");
  fout << enigma.get_number_rotors() << std::endl;
  for (size_t i = 0; i < enigma.get_number_rotors(); i++)
    fout << enigma.get_state()[i];
  fout << std::endl;
  for (size_t i = 0; i < enigma.get_number_rotors(); i++) {
    for (size_t j = 0; j < enigma.get_reflector()->get_reflector().size();
         j++) {
      fout << enigma.get_rotors()[i].get_rotor()[j + *s21::alphabet.begin()];
    }
    fout << std::endl;
  }
  for (size_t i = 0; i < enigma.get_reflector()->get_reflector().size(); i++) {
    fout << enigma.get_reflector()->get_reflector()[i + *s21::alphabet.begin()];
  }
}

std::string s21::Application::encoded_to_decoded(std::string path_to_file) {
  auto responce = path_to_file.find("_encoded");
  if (responce != path_to_file.npos) path_to_file.erase(responce);
  path_to_file += "_decoded";
  return path_to_file;
}

void s21::Application::to_encoded(s21::Enigma &enigma) {
  std::string path_to_file;
  std::cout << "Введите путь до файла который требуется закодировать: ";
  std::cin >> path_to_file;
  erase_bracket(path_to_file);
  printf("%c\033[2J", 27);
  std::ifstream fin(path_to_file, std::ios::binary);
  if (!fin.is_open()) throw std::invalid_argument("'Path is not valid'\n");
  std::ofstream fout(path_to_file + "_encoded", std::ios::binary);
  if (!fout.is_open())
    throw std::invalid_argument("'Failed to create encoded file'\n");
  char ch;
  while (fin >> std::noskipws >> ch) {
    if (ch >= *s21::alphabet.begin() && ch <= *--s21::alphabet.end())
      ch = enigma.coder(ch);
    fout << ch;
  }
  fin.close();
  fout.close();
  std::cout << "Файл закодирован: '" << path_to_file + "_encoded'\n"
            << std::endl;
}

void s21::Application::to_decoded(s21::Enigma &enigma) {
  std::string path_to_file;
  std::cout << "Введите путь до файла который требуется раскодировать: ";
  std::cin >> path_to_file;
  erase_bracket(path_to_file);
  printf("%c\033[2J", 27);
  std::ifstream fin(path_to_file, std::ios::binary);
  std::string decoded_file = encoded_to_decoded(path_to_file);
  std::ofstream fout(decoded_file, std::ios::binary);
  if (!fin.is_open()) throw std::invalid_argument("'Path is not valid'\n");
  char ch;
  while (fin >> std::noskipws >> ch) {
    if (ch >= *s21::alphabet.begin() && ch <= *--s21::alphabet.end())
      ch = enigma.coder(ch);
    fout << ch;
  }
  std::cout << "Файл расшифрован: '" << decoded_file << "'\n" << std::endl;
  fin.close();
  fout.close();
}
