#include <emscripten/emscripten.h>
#include<stdlib.h>
#include<stdint.h>
#include<string.h>
#include<stdio.h>

// allocate is used by the host to allocate memory that complex types
// can be written to
EMSCRIPTEN_KEEPALIVE void * allocate(int size) {
  return malloc(size);
}

// get the size of a byte array at the pointer
// byte arrays allways have a uint32 size in the first
// 4 bytes
int size_of_bytes(void * ptr) {
  uint8_t * bytes = (uint8_t *)ptr;

  union
  {
    unsigned int integer;
    unsigned char byte[4];
  } len;

  len.byte[0] = bytes[0];
  len.byte[1] = bytes[1];
  len.byte[2] = bytes[2];
  len.byte[3] = bytes[3];

  return len.integer;
}

// takes a c array and returns another c array with the 
// length appended to the begining
// the ptr that is returned is not managed internally
void * ptr_from_bytes(uint8_t bytes[], int size) {
  void * ptr = allocate(size + 4);
  uint8_t * data = (uint8_t *)ptr;

  // convert the length to bytes and write it to the first 4 bytes
  union
  {
    unsigned int integer;
    unsigned char byte[4];
  } len;

  len.integer = size;
  data[0] = len.byte[0];
  data[1] = len.byte[1];
  data[2] = len.byte[2];
  data[3] = len.byte[3];

  // copy the data
  for(int i = 0; i < size; ++i) {
    data[i+4] = bytes[i];
  }

  return data;
}

// returns a byte array from a pointer
void bytes_from_ptr(uint8_t bytes[], int size, void * ptr) {
  uint8_t * data = (uint8_t *)ptr;

  // copy the data
  for(int i = 0; i < size; ++i) {
    bytes[i] = data[i+4];
  }
}


EMSCRIPTEN_KEEPALIVE int get_string_size(char * ptr) {
  return strlen(ptr);
}

EMSCRIPTEN_KEEPALIVE int sum(int a, int b) {
  return a + b;
}

EMSCRIPTEN_KEEPALIVE char * hello(char * name) {
  const char * message = "Hello";
  return message;
}

EMSCRIPTEN_KEEPALIVE void * reverse(void * bytes) {
  // convert the input to a byte array
  int size = size_of_bytes(bytes);
  uint8_t data[size];
  bytes_from_ptr(data,size,bytes);

  // reverse the input
  uint8_t out[size];
  int n = 0;
  for(int i = size-1; i >= 0; --i) {
    out[n] = data[i];
    ++n;
  }

  uint8_t * ptr = ptr_from_bytes(out, 3);
  return ptr;
}

extern char * call_me(char * name);

EMSCRIPTEN_KEEPALIVE char * callback() {
  char * result = call_me("Nic");
  return result;
}

int main() {}