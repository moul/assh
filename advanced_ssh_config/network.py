import fcntl
import os
import select
import sys
import socket
import logging


class Socket(object):
    def __init__(self, hostname, port, bufsize=4096, stdin=None, stdout=None):
        self.hostname = hostname
        self.port = port
        self.socket = socket.socket(socket.AF_INET, socket.SOCK_STREAM)
        self.bufsize = bufsize

        if not stdin:
            stdin = sys.stdin
        if not stdout:
            stdout = sys.stdout
        self.stdin = stdin
        self.stdout = stdout

    def run(self):
        self.socket.connect((self.hostname, self.port))
        self.communicate()

    def communicate(self):
        self.socket.setblocking(0)

        fd = sys.stdin.fileno()
        fl = fcntl.fcntl(fd, fcntl.F_GETFL)
        fcntl.fcntl(fd, fcntl.F_SETFL, fl | os.O_NONBLOCK)

        inputs = [self.socket, self.stdin]

        while True:
            read_ready, write_ready, in_error = select.select(inputs, [], [])
            for sock in read_ready:
                if sock == self.stdin:
                    buffer = self.stdin.read(self.bufsize)
                    try:
                        while buffer != '':
                            self.socket.send(buffer)
                            buffer = self.stdin.read(self.bufsize)
                    except:
                        pass
                else:
                    try:
                        buffer = self.socket.recv(self.bufsize)
                        while buffer != '':
                            self.stdout.write(buffer)
                            self.stdout.flush()
                            buffer = self.socket.recv(self.bufsize)
                        if buffer == '':
                            logging.warn('Server disconnected')
                            return
                    except socket.error:
                        pass
