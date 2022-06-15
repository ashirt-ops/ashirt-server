class RequestState(object):
    """
    RequestState captures the memory needs of an in-flight request. If you need to store data
    temporarily (for the lifetime of a request), you can stick it here
    """
    def __init__(self, request_logger):
        self.req_log = request_logger
