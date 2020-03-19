package dumbwaiter;

import java.util.ArrayList;

import org.junit.Assert;

import dumbwaiter.common.PiPin;
import dumbwaiter.common.RPi;
import dumbwaiter.common.RPiDeviceException;

public class FakeRPi implements RPi {
  ArrayList<PiPin> expectedCalls;
  int callCount;

  public FakeRPi(ArrayList<PiPin> expectedCallsIn) {
    expectedCalls = expectedCallsIn;
    callCount = 0;
  }

  @Override
  public void SendSignal(PiPin pin) throws RPiDeviceException {
    Assert.assertEquals(this.expectedCalls.get(callCount), pin);
  }

  @Override
  public boolean GetSignal(PiPin pin) throws RPiDeviceException {
    return false;
  }

}
