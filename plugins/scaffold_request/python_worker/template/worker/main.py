import time


def main() -> None:
    while True:
        print("[[SERVICE_NAME]] processing background work", flush=True)
        time.sleep(30)


if __name__ == "__main__":
    main()
