#include <emscripten/emscripten.h>

EMSCRIPTEN_KEEPALIVE int sum(int a, int b) {
  return a + b;
}