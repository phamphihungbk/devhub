import os
import time


def run() -> None:
    queue_name = os.getenv("QUEUE_NAME", "[[ QUEUE_NAME ]]")
    while True:
        print(f"[[ SERVICE_NAME ]] polling {queue_name}")
        time.sleep(5)
