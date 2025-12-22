import logging
import sys


def get_logger(correlation_id: str) -> logging.Logger:
    logger = logging.getLogger("idp-action")
    logger.setLevel(logging.INFO)

    handler = logging.StreamHandler(sys.stdout)
    formatter = logging.Formatter(
        fmt='{"level":"%(levelname)s","msg":"%(message)s","correlation_id":"%s"}'
        % correlation_id
    )
    handler.setFormatter(formatter)

    if not logger.handlers:
        logger.addHandler(handler)

    return logger