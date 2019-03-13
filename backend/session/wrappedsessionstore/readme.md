# Local Session Store

_A fork of [MySQLStore](https://github.com/srinathgs/mysqlstore)_ with identifable sessions

## Overview

This subproject attempts to be a minimal fork of MySQLStore which adds the ability to identify which
users own which sessions. The goal is to be able to find all sessions for a particular user so that
they may be invalidated.

### Changes

1. The `sessions` table now includes a `user_id` field (which acts as a foreign key to the `users` table)
2. The `sessions` table must be created _prior_ to using this store
3. the `_on` suffix on dates is now `_at`, to better correspond to the Ashirt proejct's usage

In addition to the above changes, some code related to the above has been changed to conform to the new usage.

## Original Copyright

The original copyright can be found in local_session_store.go as the header

## License

The original project is licensed as MIT, and a copy of the license can be found in the LICENSE file

## Contributors

A list of contributors to MySQLStore can be found in the CONTRIBUTORS file
