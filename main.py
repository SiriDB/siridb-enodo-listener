import asyncio

from listener import Listener

loop = asyncio.get_event_loop()
listener = Listener(loop)
loop.run_until_complete(listener.start_listener())
loop.run_forever()
loop.close()