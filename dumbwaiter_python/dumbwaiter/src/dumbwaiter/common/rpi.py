from enum import Enum

class PiPin(Enum):
  OpenerUp = 1
  OpenerDown = 2
  OpenerStop = 3
  Floor3Requested = 4
  Floor2Requested = 5
  Floor1Requested = 6
  StopRequested = 7
  AtFloor = 8


class RPi:
  def send_signal(self, pin):
    """send a signal on the specified pin"""
    pass

  def get_signal(self, pin):
    """ get the value from the specified pin """
    pass


class RPiDevice(RPi):
  def send_signal(self, pin):
    # TODO implement
    pass

  def get_signal(self, pin):
    # TODO implement
    return None
