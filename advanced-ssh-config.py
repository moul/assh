#!/usr/bin/env python

import sys, os, ConfigParser, re, socket, subprocess, optparse

class advanced_ssh_config():
    def __init__(self, hostname = None, port = 22, configfile = None, verbose = False, update_sshconfig = False):
        self.verbose, self.hostname, self.port = verbose, hostname, port
        self.configfiles = [ '/etc/ssh/config.advanced', os.path.expanduser("~/.ssh/config.advanced") ]
        if configfile:
            self.configfiles += configfile
        self.parser = ConfigParser.ConfigParser()

        self.parser.read(self.configfiles)
        includes = self.conf_get('includes', 'default', '').strip()
        for include in includes.split():
            include = os.path.expanduser(include)
            if not include in self.configfiles and os.path.exists(include):
                self.parser.read(include)

        if update_sshconfig:
            self._update_sshconfig()

        self.debug()
        self.debug("================")

    def debug(self, str = None, force = False):
        if self.verbose or force:
            if str:
                sys.stderr.write("Debug: %s\n" % str)
            else:
                sys.stderr.write("\n")

    def conf_get(self, key, host, default = None):
        for section in self.parser.sections():
            if re.match(section, host):
                if self.parser.has_option(section, key):
                    return self.parser.get(section, key)
        if self.parser.has_option('default', key):
            return self.parser.get('default', key)
        return default

    def connect(self):
        mkdir_path = os.path.dirname(os.path.join(os.path.dirname(os.path.expanduser(self.conf_get('controlpath', 'default', '/tmp'))), self.hostname))
        try:
            os.makedirs(mkdir_path)
        except:
            pass
        path = self.hostname.split('/')

        args = {}
        options = {'p': 'Port',
                   'u': 'User',
                   'h': 'Hostname',
                   'i': 'IdentifyFile'}
        for key in options:
            value = self.conf_get(options[key], path[0])
            self.debug("get (-%-1s) %-12s : %s" % (key, options[key], value))
            if value:
                args[key] = value
        if not 'h' in args:
            args['h'] = path[0]
        self.debug('args: %s' % args)
        self.debug()

        self.debug("hostname    : %s" % self.hostname)
        self.debug("path        : %s" % path)
        self.debug("path[0]     : %s" % path[0])
        self.debug("path[1:]    : %s" % path[1:])
        self.debug("args        : %s" % args)
        self.debug("configfiles : %s" % self.configfiles)
        self.debug("port        : %s" % self.port)

        self.debug()
        gateways = self.conf_get('Gateways', path[-1], 'direct').strip().split(' ')
        reallocalcommand = self.conf_get('RealLocalCommand', path[-1], '').strip().split(' ')
        for gateway in gateways:
            right_path = path[1:]
            if gateway != 'direct':
                right_path += [gateway]
            cmd = []
            if len(right_path):
                cmd += ['ssh', '/'.join(right_path)]

            if len(cmd):
                cmd += ['nc', args['h'], args['p']]
            else:
                cmd += ['nc', args['h'], args['p']]

            self.debug("cmd         : %s" % cmd)
            self.debug("================")
            self.debug()
            ssh_process = subprocess.Popen(cmd)
            if len(reallocalcommand[0]):
                reallocalcommand_process = subprocess.Popen(reallocalcommand)
            if ssh_process.wait() != 0:
                self.debug("There were some errors")
            if len(reallocalcommand[0]):
                reallocalcommand_process.kill()

    def _update_sshconfig(self, write = True):
        config = []

        for section in self.parser.sections():
            if section != 'default':
                host = section
                host = re.sub('\.\*', '*', host)
                host = re.sub('\\\.', '.', host)
                config += ["Host %s" % host]
                for key, value in self.parser.items(section):
                    if key not in ['hostname', 'gateways', 'reallocalcommand', 'remotecommand']:
                        if key == 'alias':
                            key = 'hostname'
                        config += ["  %s %s" % (key, value)]
                config += ['']

        config += ['Host *']
        for key, value in self.parser.items('default'):
            if key not in ['hostname', 'gateways', 'includes']:
                config += ["  %s %s" % (key, value)]

        if write:
            file = open(os.path.expanduser("~/.ssh/config"), 'w+')
            file.write('\n'.join(config))
            file.close()
        else:
            print '\n'.join(config)

if __name__ == "__main__":
    parser = optparse.OptionParser(usage = "%prog [-v] -h hostname -p port", version = "%prog 1.0")
    parser.add_option("-H", "--hostname", dest = "hostname", help = "Host")
    parser.add_option("-p", "--port", dest = "port")
    parser.add_option("-v", "--verbose", dest = "verbose", action="store_true")
    parser.add_option("-u", "--update-sshconfig", dest = "update_sshconfig", action="store_true")
    (options, args) = parser.parse_args()
    ssh = advanced_ssh_config(hostname = options.hostname, port = options.port, verbose = options.verbose, update_sshconfig = options.update_sshconfig)
    ssh.connect()
    #sys.stderr.write("\n")
