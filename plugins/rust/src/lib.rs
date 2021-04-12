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
    let c_string = CString::from_raw(raw);
    c_string.as_bytes().len()
  }
}

#[no_mangle]
pub extern fn sum(x: i32, y: i32) -> i32 {
  x + y
}

#[no_mangle]
pub extern fn hello(raw: *mut c_char) -> *mut c_char {
  // fetch the string from memory
  let in_param = unsafe{  CStr::from_ptr(raw).to_bytes().to_vec()  };

  // combine the input and output
  let mut output = b"Hello ".to_vec();
  output.extend(&in_param);

  // create a pointer to a c_str to return to the caller
  unsafe { CString::from_vec_unchecked(output) }.into_raw()
}

#[no_mangle]
pub extern fn reverse(raw: *mut c_void) -> *mut c_void {
  let data_in = bytes_from_ptr(raw);
  return ptr_from_bytes(data_in.into_iter().rev().collect());
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