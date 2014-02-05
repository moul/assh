import socket


def create_socket(hostname, port):
    s = socket.socket(socket.AF_INET, socket.SOCK_STREAM)
    s.connect((hostname, port))
    logging.error('socket created')
    bytes = 0L
    last = 0
    # communicate
