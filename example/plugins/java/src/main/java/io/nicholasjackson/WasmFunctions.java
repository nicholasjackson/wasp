package io.nicholasjackson;

import org.teavm.interop.Export;

public class WasmFunctions {

    @Export(name = "sum")
    public static int sum(int a, int  b){
        return a + b;
    }

    public static void main(String[] args) {

    }
}