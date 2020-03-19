package dumbwaiter.common;


// RPi and interface to the RPi device
public interface RPi {
  void SendSignal(PiPin pin) throws RPiDeviceException; // send a signal on the target pin
  boolean GetSignal(PiPin pin) throws RPiDeviceException;   //  get a signal from a target pin, true when signal is on, false when off
}

