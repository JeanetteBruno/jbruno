package dumbwaiter.common;

/*
Timing Control interface allows tests to inject control over the interactions
between test and controller threads
 */
public interface TimingControl {
  public void sleep(long durationMs);
}
