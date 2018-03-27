class StateError(Exception):
    """
    Raised when current state prevents code execution.
    (i.e. trying to shell into container when it's not running) 
    """