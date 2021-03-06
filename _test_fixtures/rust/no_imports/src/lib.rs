use std::ffi::{CStr, CString};
use std::fs::write;
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
pub extern fn workspace_write(dir: *mut c_char) {
  // fetch the string from the ptr passed to the function
  let in_param = string_from_ptr(dir);

  let mut in_file =  in_param.to_owned();
  in_file.push_str("/in.txt");
  
  match std::fs::read_to_string(in_file) {
    Ok(s) => print!("Read file {}\n",s),
    Err(e) => print!("No Read file {}\n",e),
  };

  let mut out_file =  in_param.to_owned();
  out_file.push_str("/hello.txt");

  print!("Writing file {}\n",out_file);

  match std::fs::write(in_param,"blah") {
    Ok(_) => print!("Written file\n"),
    Err(e) => print!("No Read file {}\n",e),
  }
}