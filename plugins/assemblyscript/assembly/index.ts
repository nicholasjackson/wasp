import { call_me } from "./plugins";

function bytes_from_buffer(raw: ArrayBuffer): Int8Array {
  // the length of the data is stored in the buffer in the first 4 bytes we can discard this
  return Int8Array.wrap(raw.slice(4));
}

function buffer_from_bytes(data: Int8Array): ArrayBuffer {
  // the output buffer is 4 bytes longer to contain the length
  let buffer = new ArrayBuffer(data.byteLength+4);
  let view = new DataView(buffer);

  // set the length as a uint 32
  view.setUint32(0,data.byteLength,true);

  // copy the remaining data
  let out = Int8Array.wrap(buffer);
  out.set(data,4);

  return buffer;
} 

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

  //console.log("writing" + inParam)
  return String.UTF8.encode("Hello " + inParam, true)
}

export function reverse(inRaw: ArrayBuffer) : ArrayBuffer {
  let inData = bytes_from_buffer(inRaw);
  let outData = inData.reverse();

  return buffer_from_bytes(outData);
}

export function callback(): ArrayBuffer {
  let inParam = call_me(String.UTF8.encode("Nic"));

  return inParam;
}