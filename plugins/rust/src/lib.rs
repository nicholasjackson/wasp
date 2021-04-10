use std::ffi::{CStr, CString};
use std::mem;
use std::ptr;
use std::os::raw::{c_char, c_void};

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
pub extern fn reverse(raw: *mut u8) -> *mut c_void {
  unsafe {
    // get the length of the original vector
    let len = ptr::read(raw) as usize;
    let mut buffer: Vec<u8> = Vec::new();
    let original: Vec<u8> = Vec::from_raw_parts(raw,len+1,len+1);
    
    // set the size of the new buffer to the same as the old
    buffer.push(len as u8);
    //std::println!("len {}", original[0]);

    // reverse the array
    for b in (1..len+1).rev() {
      //std::println!("out {}", original[b]);
      buffer.push(original[b]);
    }
  
    let pointer = buffer.as_mut_ptr();
    mem::forget(buffer);

    return pointer as *mut c_void
  }
}

#[cfg(test)]
#[test]
fn it_works() {
  let mut buffer:Vec<u8> = Vec::with_capacity(2);
  buffer.push(1);
  buffer.push(42);

  let pointer = buffer.as_mut_ptr();
  mem::forget(buffer);

  reverse(pointer);
}