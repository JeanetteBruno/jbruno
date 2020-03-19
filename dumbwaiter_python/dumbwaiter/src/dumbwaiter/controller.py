from enum import Enum
import logging as log
from threading import (Lock, Thread)
import sys
import time

from dumbwaiter.common.rpi import (PiPin, RPiDevice)

default_loop_frequency_sec = 0.5

class Direction(Enum):
  Up = 1
  Down = 2
  Stopped = 3

class Controller:
  def __init__(self, top_floor):
    self._top_floor = top_floor
    self._pi_device = RPiDevice()
    self._requested_floor_lock = Lock()
    self._direction_lock = Lock()
    self._last_seen_floor_lock = Lock()
    self._loop_frequency = default_loop_frequency_sec

    self._requested_floor = -1
    self._last_seen_floor = -1
    self._direction = Direction.Stopped

  def start_processing_loop(self):
    t = Thread(target=self._processing_loop, args=(), daemon=True)
    t.start()

  # _processing_loop translates signals from the dumbwaiter platform and user controls
  # to up/down/stop commands to the garage door opener
  def _processing_loop(self):

    while True:
      # if the car is stationary and another floor is requested, start it moving in the requested direction
      # if the car is moving and a floor in the opposite direction has been requested stop the car                                                                                                                                                                       # (let the next iteration start it moving)
      if self.get_requested_floor() > self.get_last_seen_floor():
        if self.get_moving_direction() == Direction.Stopped:
          self._send_up()
        elif self.get_moving_direction() == Direction.Down:
          self._stop() # stop the machine, it start moving up on next iteration
        # else do nothing it is already moving up
      elif self.get_requested_floor() < self.get_last_seen_floor():
        if self.get_moving_direction() == Direction.Stopped:
          self._send_down()
        elif self.get_moving_direction() == Direction.Up:
          self._stop() # stop the machine, it will start moving down on next iteration
      elif self.get_moving_direction() != Direction.Stopped:
        self._stop()

      time.sleep(self._loop_frequency)

  def _send_up(self):
    self._pi_device.send_signal(PiPin.OpenerUp)
    self.set_moving_direction(Direction.Up)

  def _send_down(self):
    self._pi_device.send_signal(PiPin.OpenerDown)
    self.set_moving_direction(Direction.Down)

  def _stop(self):
    self._pi_device.send_signal(PiPin.OpenerStop)
    self.set_moving_direction(Direction.Stopped)

  def get_requested_floor(self):
    with self._requested_floor_lock:
      return self._requested_floor

  def set_requested_floor(self, floor):
    with self._requested_floor_lock:
      self._requested_floor = floor

  def get_last_seen_floor(self):
    with self._last_seen_floor_lock:
      return self._last_seen_floor

  def set_last_seen_floor(self, floor):
    with self._last_seen_floor_lock:
      self._last_seen_floor = floor

  def get_moving_direction(self):
    with self._direction_lock:
      return self._direction

  def set_moving_direction(self, direction):
    with self._direction_lock:
      self._direction = direction
