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

export function reverse(inRaw: ArrayBuffer) : ArrayBuffer {
  let inData = Int8Array.wrap(inRaw)
  let outRaw = new ArrayBuffer(4);
  let outData = Int8Array.wrap(outRaw)

  outData[0] = 3; // size of the array

  // read the inData and reverse
  // length is always position 1
  var pos = 1;
  for (let i = inData[0]; i > 0; --i) {
    outData[pos] = inData[i];
    pos++
  }

  return outRaw;
}