import logging
import os


def get_logger() -> logging.Logger:
    logging.basicConfig(format='[%(levelname)s %(asctime)s (%(name)s:%(lineno)d)]: %(message)s')

    level = 'DEBUG' if os.environ.get('DEBUG', 'false') == 'true' else 'INFO'
    logger = logging.getLogger()
    logger.setLevel(level)
    return logger
