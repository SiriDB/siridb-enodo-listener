import os

os.environ["ENODO_HUB_HOSTNAME"] = "localhost"
os.environ["ENODO_HUB_PORT"] = "9103"
os.environ["ENODO_PIPE_PATH"] = "/tmp/test.sock"
os.environ["ENODO_INTERNAL_SECURITY_TOKEN"] = ""

os.system("go run .")