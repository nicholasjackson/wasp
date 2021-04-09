package io.nicholasjackson;

import java.io.BufferedReader;
import java.io.IOException;
import java.io.InputStreamReader;
import java.net.URL;

import org.teavm.interop.Async;
import org.teavm.interop.Export;
import org.teavm.platform.async.AsyncCallback;

public class WasmFunctions {

    @Export(name = "sum")
    public static int sum(int a, int  b){
        return a + b;
    }

    @Async
    public static native String getURL(String url) throws IOException;

    @Export(name = "get")
    public static void getURL(String url, AsyncCallback<String> callback) {
      try {
        StringBuffer out = new StringBuffer();
        URL oracle;
        oracle = new URL("http://www.oracle.com/");
        BufferedReader in;
        in = new BufferedReader(
        new InputStreamReader(oracle.openStream()));

        String inputLine;

        while ((inputLine = in.readLine()) != null) {
          out.append(inputLine);
        }
    
        in.close();
      } catch (IOException e) {
        e.printStackTrace();
      }
    }

    public static void main(String[] args) {

    }
}