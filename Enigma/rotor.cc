#include "rotor.h"

s21::Rotor::Rotor() { make_rotor(); }

s21::Rotor::Rotor(s21::Rotor const &other) { rotor_ = other.rotor_; }

s21::Rotor::Rotor(s21::Rotor &&other) { rotor_ = std::move(other.rotor_); }

s21::Rotor::~Rotor() { rotor_.clear(); }

void s21::Rotor::operator=(s21::Rotor const &other) { rotor_ = other.rotor_; }

void s21::Rotor::operator=(s21::Rotor &&other) {
  std::swap(rotor_, other.rotor_);
}

std::map<char, char> s21::Rotor::get_rotor() { return rotor_; }

void s21::Rotor::set_rotor(std::string str) {
  for (size_t i = 0; i < rotor_.size(); i++) rotor_[i + 'A'] = str[i];
}

char s21::Rotor::get_out_char(char current_ch) { return rotor_.at(current_ch); }

char s21::Rotor::get_key(char ch) {
  for (auto it = rotor_.begin(); it != rotor_.end(); ++it) {
    if ((*it).second == ch) {
      ch = (*it).first;
      break;
    }
  }
  return ch;
}

void s21::Rotor::make_rotor() {
  std::vector<char> tmp(alphabet);
  for (int i = 0; i < (int)alphabet.size(); i++) {
    int random = rand() % tmp.size();
    rotor_.emplace(tmp[random], alphabet[i]);
    tmp.erase(tmp.begin() + random);
  }
}
