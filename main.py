import argparse
import asyncio

from listener import Listener

parser = argparse.ArgumentParser(description='Process config')
parser.add_argument('--config', help='Config path', required=True)

loop = asyncio.get_event_loop()
listener = Listener(loop, parser.parse_args().config)
loop.run_until_complete(listener.start_listener())
loop.run_forever()
loop.close()