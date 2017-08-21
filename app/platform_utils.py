import sys

def print_stdout(string, newLine = True):
    """ Print string to stdout. """
    sys.stdout.write(
        "%s%s" % (string, "\n" if newLine else "")
    )
    sys.stdout.flush()    