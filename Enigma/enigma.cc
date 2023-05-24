#include "enigma.h"

s21::Enigma::Enigma(size_t number_rotors) : number_rotors_(number_rotors) {
  if (number_rotors_) {
    rotors_ = new s21::Rotor[number_rotors];
    reflector_ = new s21::Reflector();
    initial_state();
  }
}

s21::Enigma::Enigma(s21::Enigma const &other) { *this = other; }

s21::Enigma::Enigma(s21::Enigma &&other) { move_enigma(other); }

void s21::Enigma::operator=(s21::Enigma const &other) {
  rotors_ = new s21::Rotor[other.number_rotors_];
  for (size_t i = 0; i < other.number_rotors_; i++)
    rotors_[i] = other.rotors_[i];
  reflector_ = new s21::Reflector(*other.reflector_);
  state_ = other.state_;
  number_rotors_ = other.number_rotors_;
}

void s21::Enigma::operator=(s21::Enigma &&other) {
  clear_enigma();
  move_enigma(other);
}

s21::Enigma::~Enigma() { clear_enigma(); }

void s21::Enigma::clear_enigma() {
  if (number_rotors_) {
    delete[] rotors_;
    delete reflector_;
    state_.clear();
    number_rotors_ = 0;
  }
}

void s21::Enigma::move_enigma(s21::Enigma &other) {
  std::swap(rotors_, other.rotors_);
  std::swap(reflector_, other.reflector_);
  std::swap(state_, other.state_);
  std::swap(number_rotors_, other.number_rotors_);
  other.clear_enigma();
}

void s21::Enigma::initial_state() {
  size_t count = number_rotors_;
  while (count--) state_.push_back((rand() % 26) + 'A');
}

s21::Rotor *s21::Enigma::get_rotors() { return rotors_; }

s21::Reflector *s21::Enigma::get_reflector() { return reflector_; }

std::vector<char> &s21::Enigma::get_state() { return state_; }

void s21::Enigma::set_state(std::string str) {
  for (size_t i = 0; i < str.size(); ++i) state_[i] = str[i];
}

[[nodiscard]] bool s21::Enigma::check_state(char ch) { return ch == 'Z' ? true : false; }

char s21::Enigma::counter_how_add(int number) {
  return number < 0 ? 26 - abs(number) % 26 : number % 26;
}

void s21::Enigma::rotors_rotation() {
  for(size_t i = 0; i < number_rotors_; ++i) {
    if(check_state(state_[i])) {
      state_[i] = 'A';
    } else {
      state_[i] += 1;
      break;
    }
  }
}

void s21::Enigma::first_step(char &ch) {
  rotors_rotation();
  ch = counter_how_add((state_[0] - 'A') + (ch - 'A')) + 'A';
}

void s21::Enigma::rotor_to_reflector(char &ch) {
  ch = counter_how_add((ch - 'A') + (state_[1] - state_[0])) + 'A';
  for (size_t i = 1; i < state_.size(); i++) {
    ch = rotors_[i].get_rotor()[ch];
    if (i + 1 < state_.size()) {
      ch = counter_how_add((ch - 'A') + (state_[i + 1] - state_[i])) + 'A';
    }
  }
}

void s21::Enigma::after_reflector(char &ch) {
  ch = counter_how_add(ch - state_[state_.size() - 1]) + 'A';
  ch = reflector_->get_reflector()[ch];
}

void s21::Enigma::rotor_back_and_out(char &ch) {
  int rotor_index = state_.size() - 1;
  s21::Rotor r = rotors_[rotor_index];
  ch = counter_how_add((ch - 'A') + (state_[rotor_index] - 'A')) + 'A';
  ch = r.get_key(ch);
  if (rotor_index) {
    for (int i = rotor_index - 1; i >= 0; --i) {
      r = rotors_[i];
      ch = counter_how_add((ch - 'A') - (state_[i + 1] - state_[i])) + 'A';
      ch = r.get_key(ch);
    }
  }
  ch = counter_how_add(ch - state_[0]) + 'A';
}

char s21::Enigma::coder(char &ch) {
  first_step(ch);
  ch = rotors_[0].get_out_char(ch);
  if (state_.size() > 1) rotor_to_reflector(ch);
  after_reflector(ch);
  rotor_back_and_out(ch);
  return ch;
}
