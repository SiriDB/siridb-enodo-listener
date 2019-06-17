import asyncio
import datetime
import errno
import fcntl
import json
import os
import uuid

from lib.socket.package import *


class Client:

    def __init__(self, loop, hostname, port, heartbeat_interval=5):
        self.loop = loop
        self._hostname = hostname
        self._port = port
        self._heartbeat_interval = heartbeat_interval

        self._id = uuid.uuid4()
        self._messages = {}
        self._current_message_id = 1
        self._current_message_id_locked = False

        self._last_heartbeat_send = datetime.datetime.now()
        self._updates_on_heartbeat = []
        self._cbs = None
        self._sock = None

    async def setup(self, cbs=None):
        self._sock = socket.socket(socket.AF_INET, socket.SOCK_STREAM)

        self._sock.connect((self._hostname, self._port))
        fcntl.fcntl(self._sock, fcntl.F_SETFL, os.O_NONBLOCK)

        self._cbs = cbs
        if cbs is None:
            self._cbs = {}

        await self._handshake()

    async def run(self):
        while 1:
            if (datetime.datetime.now() - self._last_heartbeat_send).total_seconds() > int(
                    self._heartbeat_interval):
                await self._send_heartbeat()

            await self._read_from_socket()

    async def close(self):
        print('Close the socket')
        self._sock.close()

    async def _read_from_socket(self):
        try:
            header = self._sock.recv(PACKET_HEADER_LEN)
        except socket.error as e:
            err = e.args[0]
            if err == errno.EAGAIN or err == errno.EWOULDBLOCK:
                await asyncio.sleep(1)
                pass
            else:
                # a "real" error occurred
                print(e)
                # sys.exit(1)
        else:
            await self._read_message(header)

    async def _read_message(self, header):
        packet_type, packet_id, data = await read_packet(self._sock, header)

        if packet_type == 0:
            print("Connection lost, trying to reconnect")
            try:
                await self.setup(self._cbs)
            except Exception as e:
                print(e)
                await asyncio.sleep(5)
        elif packet_type == HANDSHAKE_OK:
            print(f'Hands shaked with hub')
        elif packet_type == HANDSHAKE_FAIL:
            print(f'Hub does not want to shake hands')
        elif packet_type == HEARTBEAT:
            print(f'Heartbeat back from hub')
        elif packet_type == REPONSE_OK:
            print(f'Hub received update correctly')
        elif packet_type == UNKNOW_CLIENT:
            print(f'Hub does not recognize us')
            await self._handshake()
        else:
            if packet_type in self._cbs.keys():
                await self._cbs.get(packet_type)(data)
            else:
                print(f'Message type not implemented: {packet_type}')

    async def _send_message(self, length, message_type, data):
        if self._current_message_id_locked:
            while self._current_message_id_locked:
                await asyncio.sleep(0.1)

        self._current_message_id_locked = True
        header = create_header(length, message_type, self._current_message_id)
        self._current_message_id += 1
        self._current_message_id_locked = False

        self._sock.send(header + data)

    async def send_message(self, body, message_type):
        await self._send_message(len(body), message_type, body)

    async def _handshake(self):
        data = json.dumps({'client_id': str(self._id), 'client_type': 'listener'}).encode('utf-8')
        await self._send_message(len(data), HANDSHAKE, data)
        self._last_heartbeat_send = datetime.datetime.now()

    async def _send_heartbeat(self):
        print('Sending heartbeat to hub')
        id_encoded = str(self._id).encode('utf-8')
        await self._send_message(len(id_encoded), HEARTBEAT, id_encoded)
        self._last_heartbeat_send = datetime.datetime.now()
