import sys

def print_stdout(string, newLine = True):
    """ Print string to stdout. """
    sys.stdout.write(
        "%s%s" % (string, "\n" if newLine else "")
    )
    sys.stdout.flush()

def log_stdout(string, level = 0, newLine = True):
    """
    Print string to stdout but prefix it based on level.
    Used to show different levels of sub tasks.
    """
    if level < 0: return
    levelStr = ">"
    if level > 0:
        levelStr = "-".rjust(level * 3)
    print_stdout(
        "%s %s" % (
            levelStr,
            string
        ),
        newLine
    )

def seperator_stdout():
    print_stdout("=======================================")