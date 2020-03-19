package dumbwaiter;

import dumbwaiter.common.TimingControl;

/*
FakeTicker is used by the tests to inject a timing controller that allows the
tests to synchronize actions across the main controller thread and the test thread.

This timing controller causes the main thread to pause till the test thread
sets 'tick' to true.  The main thread will consume that tick, ending its 'sleep'
then wait for the next 'tick' on the next 'sleep'.
 */
public class FakeTicker implements TimingControl {
  boolean tick;

  public FakeTicker() {
    this.consumeTick();  // set tick to false
  }

  /*
  sleep till we get a tick, then consume the tick
   */
  @Override
  public void sleep(long durationMs) {
    while (true) {
      // keep sleeping till we get a tick from the test control
      try {
        Thread.sleep(100);
      } catch (InterruptedException e) {
        e.printStackTrace();
      }
      if (this.isTick()) {
        this.consumeTick();
        return;
      }
    }
  }

  /*
  single tick: if we already have tick, wait for it to be consumed,
  then set tick.
   */
  public void tick() {
    while (true) {  // wait for the current tick to be consumed
      if (! this.tick) {
        this.setTick();
        return;
      }
      try {
        Thread.sleep(100);
      } catch (InterruptedException e) {
        e.printStackTrace();
      }
    }
  }

  public synchronized boolean isTick() {
    return tick;
  }

  public synchronized void consumeTick() {
    this.tick = false;
  }

  public synchronized void setTick() {
    this.tick = true;
  }

}
