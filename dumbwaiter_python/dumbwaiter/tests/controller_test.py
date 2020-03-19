from datetime import datetime, timedelta
from mock import (patch, call)
import time
import unittest

from dumbwaiter.common.rpi import (PiPin, RPiDevice)
from dumbwaiter.controller import (Controller, Direction, default_loop_frequency_sec)

class TestDumbwaiterController(unittest.TestCase):

  @patch('dumbwaiter.common.rpi.RPiDevice.send_signal')
  def test_send_up_from_stopped(self, mock_pi_send):
    dwc = self._setup_test(2, 2, Direction.Stopped)

    mock_pi_send.return_value=None
    dwc.set_requested_floor(3)
    self._wait_for_status(dwc, 2, Direction.Up, 3, 3)
    mock_pi_send.assert_called_once_with(PiPin.OpenerUp)

  @patch('dumbwaiter.common.rpi.RPiDevice.send_signal')
  def test_send_up_from_down(self, mock_pi_send):
    dwc = self._setup_test(2, 1, Direction.Down)

    mock_pi_send.return_value=None
    dwc.set_requested_floor(3)
    self._wait_for_status(dwc, 2, Direction.Up, 3, 3)
    mock_pi_send.assert_has_calls([call(PiPin.OpenerStop), call(PiPin.OpenerUp)])

  @patch('dumbwaiter.common.rpi.RPiDevice.send_signal')
  def test_send_down_from_stopped(self, mock_pi_send):
    dwc = self._setup_test(2, 2, Direction.Stopped)

    mock_pi_send.return_value=None
    dwc.set_requested_floor(1)
    self._wait_for_status(dwc, 2, Direction.Down, 1, 3)
    mock_pi_send.assert_called_once_with(PiPin.OpenerDown)

  @patch('dumbwaiter.common.rpi.RPiDevice.send_signal')
  def test_send_down_from_up(self, mock_pi_send):
    dwc = self._setup_test(2, 3, Direction.Up)

    mock_pi_send.return_value=None
    dwc.set_requested_floor(1)
    self._wait_for_status(dwc, 2, Direction.Down, 1, 3)
    mock_pi_send.assert_has_calls([call(PiPin.OpenerStop), call(PiPin.OpenerDown)])

  # create a controller in the specified state and start it running
  def _setup_test(self, last_seen_floor, requested_floor, direction):
    default_loop_frequency_sec = 0.1  # speed up tests
    dwc = Controller(3)
    dwc.set_last_seen_floor(last_seen_floor)
    dwc.set_requested_floor(requested_floor)
    dwc.set_moving_direction(direction)
    dwc.start_processing_loop()
    return dwc

  # _wait_for_status loops up to the timeout, trying to verify that the controller
  # reached the targeted state
  def _wait_for_status(self, dwc, last_seen_floor, direction, requested_floor, timeout_sec):
    end = datetime.now() + timedelta(seconds=timeout_sec)
    while datetime.now() < end:
      if dwc.get_last_seen_floor() == last_seen_floor and \
        dwc.get_moving_direction() == direction and \
        dwc.get_requested_floor() == requested_floor:
        return
      time.sleep(0.25)

    if dwc.get_last_seen_floor() == last_seen_floor and \
      dwc.get_moving_direction() == direction and \
      dwc.get_requested_floor() == requested_floor:
      return
