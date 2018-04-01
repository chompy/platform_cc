class ContainerCommandError(Exception):
    """
    Raised when command issued to container returns
    with unexpected exit code.
    """