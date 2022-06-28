from dataclasses import dataclass


@dataclass(frozen=True)
class ProjectConfig:
    """
    ProjectConfig stores the configuration read from the passed dictionary, if using from_dict
    (this is intended to be os.environ). You can then access these values via the fields below.
    """
    dev_mode: bool
    backend_url: str
    access_key: str
    secret_key_b64: str
    port: str

    @classmethod
    def from_dict(cls, data: dict[str, str]):
        """
        from_dict attempts to get all of the configuration needs from the provided dictionary.
        If a field is not in the dictionary, then the default value is used instead.
        """
        dev_mode = data.get('ENABLE_DEV', 'false').lower() == 'true'
        backend_url = data.get('ASHIRT_BACKEND_URL', '')
        access_key = data.get('ASHIRT_ACCESS_KEY', '')
        secret_key_b64 = data.get('ASHIRT_SECRET_KEY', '')
        port = data.get('PORT', '5000')

        return cls(
            dev_mode=dev_mode,
            backend_url=backend_url,
            access_key=access_key,
            secret_key_b64=secret_key_b64,
            port=port,
        )
