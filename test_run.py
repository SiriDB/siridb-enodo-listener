import os

os.environ["ENODO_HUB_HOSTNAME"] = "localhost"
os.environ["ENODO_HUB_PORT"] = "9103"
os.environ["ENODO_TCP_PORT"] = "9104"
os.environ["ENODO_READY_PORT"] = "8082"
os.environ["ENODO_INTERNAL_SECURITY_TOKEN"] = ""

os.system("go run .")
