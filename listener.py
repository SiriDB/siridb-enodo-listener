import asyncio
import configparser
import datetime
import json
import os

from lib.siridb.pipeserver import PipeServer
from lib.socket.client import Client
from lib.socket.package import *


class Listener:

    def __init__(self, loop):
        self._loop = loop
        self._config = configparser.ConfigParser()
        self._config.read(os.path.join(os.path.dirname(__file__), 'listener.conf'))
        self._series_to_watch = ()
        self._serie_counter_updates = {}
        self._client = Client(loop, self._config['enodo']['hub_hostname'], int(self._config['enodo']['hub_port']),
                              heartbeat_interval=int(self._config['enodo']['heartbeat_interval']))
        self._client_run_task = None
        self._updater_task = None
        self._last_update = datetime.datetime.now()

    async def _start_siridb_pipeserver(self):
        pipe_server = PipeServer(self._config['enodo']['pipe_path'], self._on_data)
        await pipe_server.create()

    def _on_data(self, data):
        """
        Forwards incoming data to a async handler
        :param data:
        :return:
        """
        asyncio.ensure_future(self._handle_pipe_data(data))

    async def _handle_pipe_data(self, data):
        """
        Handles incoming data, when not relevant, it will be ignored
        :param data:
        :return:
        """
        print("INCOMMING DATA")
        for serie_name, values in data.items():
            if serie_name in self._series_to_watch:
                if serie_name in self._serie_counter_updates:
                    serie_counter = self._serie_counter_updates.get(serie_name)
                    serie_counter += len(values)
                else:
                    self._serie_counter_updates[serie_name] = len(values)

        print(self._serie_counter_updates)

    async def _updater(self):
        while 1:
            if (datetime.datetime.now() - self._last_update).total_seconds() > int(
                    self._config['enodo']['counter_update_interval']) and len(self._serie_counter_updates.keys()):
                await self._send_update()
                self._last_update = datetime.datetime.now()
            await asyncio.sleep(1)

    async def _send_update(self):
        update_encoded = json.dumps(self._serie_counter_updates).encode('utf-8')
        self._serie_counter_updates = {}
        await self._client.send_message(update_encoded, LISTENER_ADD_SERIE_COUNT)

    async def start_listener(self):
        await self._start_siridb_pipeserver()
        await self._client.setup(cbs={
            UPDATE_SERIES: self._handle_update_series
        })
        self._client_run_task = self._loop.create_task(self._client.run())
        self._updater_task = self._loop.create_task(self._updater())

    async def _handle_update_series(self, data):
        self._series_to_watch = set(json.loads(data))

    def close(self):
        self._client_run_task.cancel()
        self._updater_task.cancel()
        self._loop.run_until_complete(self._client_run_task)
        self._loop.run_until_complete(self._updater_task)
