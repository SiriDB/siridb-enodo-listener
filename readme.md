
<p align="center"><img src="https://github.com/siridb/siridb-enodo-hub/raw/development/assets/logo_full.png" alt="Enodo"></p>

# Enodo

### Listener

The enode listener listens to pipe socket with siridb server. It sums up the totals of added datapoints to each serie. 
It periodically sends an update to the enode hub. The listener only keeps track of series that are registered via an ADD_SERIE of the UPDATE_SERIE message. The listener is seperated from the enode hub, so that it can be placed close to the siridb server, so it can locally access the pipe socket.
Every interval for heartbeat and update can be configured with the listener.conf file next to the main.py


## Getting started

To get the Enodo Listener setup you need to following the following steps:

### Locally

1. Install dependencies via `pip3 install -r requirements.txt`
2. Setup a .conf file file `python3 main.py --create_config` There will be made a `default.conf` next to the main.py.
3. Fill in the `default.conf` file
4. Call `python3 main.py --config=default.conf` to start the hub.
5. You can also setup the config by environment variables. These names are identical to those in the default.conf file, except all uppercase.