import argparse
import asyncio
import os

from listener import Listener

from lib.config import create_standard_config_file

parser = argparse.ArgumentParser(description='Process config')
parser.add_argument('--config', help='Config path', required=False)
parser.add_argument('--create_config', help='Create standard config file', action='store_true', default=False)

if parser.parse_args().create_config:
    create_standard_config_file(os.path.join(os.path.dirname(os.path.realpath(__file__)), 'default.conf'))
    exit()

loop = asyncio.get_event_loop()
listener = Listener(loop, parser.parse_args().config)
loop.run_until_complete(listener.start_listener())
loop.run_forever()
loop.close()