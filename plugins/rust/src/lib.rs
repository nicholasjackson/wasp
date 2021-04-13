use std::ffi::{CStr, CString};
use std::mem;
use std::os::raw::{c_char, c_void};

#[no_mangle]
fn bytes_from_ptr(raw: *mut c_void) -> Vec<u8> {
  unsafe {
    let data = raw as *mut u8;
    // get the length of the original vector
    let len_vec: Vec<u8> = Vec::from_raw_parts(data, 4, 4);
    let mut len_array:[u8; 4] = [0; 4];

    for i in 0..4 {
      len_array[i] = len_vec[i];
    }

    let len = u32::from_le_bytes(len_array) as usize;
    let mut original: Vec<u8> = Vec::from_raw_parts(data,len+4,len+4);

    return original.split_off(4);
  }
}

fn ptr_from_bytes(data: Vec<u8>) -> *mut c_void {
  let len = data.len() as u32; // cast len to u32 should never exceed this size
  let len_bytes = len.to_le_bytes();

  // add the length to the beginning of the buffer
  let mut buffer:Vec<u8> = data.clone();
  buffer.splice(0..0, len_bytes.iter().cloned());

  // get a pointer for the buffer and stop rust from managing the memory
  let pointer = buffer.as_mut_ptr();
  mem::forget(buffer);
    
  pointer as *mut c_void
}

fn ptr_from_string(data: String) -> *mut c_char {
  let str = CString::new(data).expect("Expected CString to be created from string");
  let ptr = str.into_raw();

  return ptr
}

fn string_from_ptr(ptr: *mut c_char) -> String {
  let str = unsafe { CString::from_raw(ptr) }.into_string().expect("Expected CString to be created from ptr");
  return str;
}

#[no_mangle]
pub extern fn allocate(size: usize) -> *mut c_void {
  let mut buffer:Vec<u8> = Vec::with_capacity(size);
  let pointer = buffer.as_mut_ptr();
  mem::forget(buffer);

  pointer as *mut c_void
}

#[no_mangle]
pub extern fn deallocate(pointer: *mut c_void, capacity: usize) {
  unsafe {
      let _ = Vec::from_raw_parts(pointer, 0, capacity);
  }
}

#[no_mangle]
pub extern fn get_string_size(raw: *mut c_char) -> usize {
  unsafe {
    // use CStr as this borrows the reference
    // CString will reclaim the reference
    let c_string = CStr::from_ptr(raw);
    c_string.to_bytes().len()
  }
}

#[no_mangle]
pub extern fn sum(x: i32, y: i32) -> i32 {
  x + y
}

#[link(wasm_import_module = "plugin")]
extern "C" {
    fn call_me(name: *mut c_char) -> *mut c_char;
}

#[no_mangle]
pub extern fn hello(name: *mut c_char) -> *mut c_char {
  // fetch the string from the ptr passed to the function
  let in_param = string_from_ptr(name);

  // append the name
  let mut output = "Hello ".to_owned();
  output.push_str(&in_param);
  
  // create a pointer to a c_str to return to the caller
  return ptr_from_string(output);
}

#[no_mangle]
pub extern fn reverse(raw: *mut c_void) -> *mut c_void {
  let data_in = bytes_from_ptr(raw);
  return ptr_from_bytes(data_in.into_iter().rev().collect());
}

#[no_mangle]
pub extern fn callback() -> *mut c_char {
  // uses the Wasi interface for printing on the host
  println!("Hello, world!");

  let name = ptr_from_string("World".to_owned());
  let result = unsafe { call_me(name) };
  
  return result;
}

#[cfg(test)]
#[test]
fn it_works() {
  let mut buffer:Vec<u8> = Vec::with_capacity(2);
  buffer.push(1);
  buffer.push(42);

  let pointer = ptr_from_bytes(buffer);

  reverse(pointer);
}