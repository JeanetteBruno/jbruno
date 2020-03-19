package dumbwaiter.common;

import java.util.concurrent.TimeUnit;

/*
Sleeper simply sleeps the thread for the specified time (sleep(ms)).
 */
public class Sleeper implements TimingControl {
  public void sleep(long durationMs) {
    try {
      TimeUnit.MILLISECONDS.sleep(durationMs);
    } catch (Exception e) {
      System.out.println(String.format("exception %e", e));
    }
  }
}
