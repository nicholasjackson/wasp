export function allocate(size: i32): ArrayBuffer {
  return new ArrayBuffer(size);
}

export function get_string_size(b: ArrayBuffer): i32 {
  return b.byteLength;
}

export function sum(a: i32, b: i32): i32 {
  return a + b;
}

export function hello(name: ArrayBuffer): ArrayBuffer {
  let inParam = String.UTF8.decode(name,true)

  return String.UTF8.encode("Hello " + inParam, true)
}