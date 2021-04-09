use std::ffi::{CStr, CString};
use std::mem;
use std::os::raw::{c_char, c_void};

#[no_mangle]
pub extern fn allocate(size: usize) -> *mut c_void {
  let mut buffer = Vec::with_capacity(size);
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