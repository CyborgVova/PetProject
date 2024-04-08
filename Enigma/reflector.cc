#include "reflector.h"

s21::Reflector::Reflector() { make_reflector(); }

s21::Reflector::Reflector(s21::Reflector const &other) {
  reflector_ = other.reflector_;
}

s21::Reflector::Reflector(s21::Reflector &&other) {
  reflector_ = std::move(other.reflector_);
}

s21::Reflector::~Reflector() { reflector_.clear(); }

void s21::Reflector::operator=(s21::Reflector const &other) {
  reflector_ = other.reflector_;
}

void s21::Reflector::operator=(s21::Reflector &&other) {
  std::swap(reflector_, other.reflector_);
}

std::map<char, char> s21::Reflector::get_reflector() { return reflector_; }

void s21::Reflector::set_reflector(std::string str) {
  for (size_t i = 0; i < reflector_.size(); i++)
    reflector_[i + *s21::alphabet_set.begin()] = str[i];
}

void s21::Reflector::make_reflector() {
  srand(time(nullptr));
  std::set<char> tmp(s21::alphabet_set);
  for (size_t i = 0; i < s21::alphabet_set.size() / 2; i++) {
    auto it = tmp.begin();
    char key = *it;
    tmp.erase(it);
    int random = rand() % tmp.size();
    it = tmp.begin();
    for (int i = 0; i < random; ++i) it++;
    reflector_.emplace(key, *it);
    reflector_.emplace(*it, key);
    tmp.erase(it);
  }
}
