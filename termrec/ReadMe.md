# Ashirt Terminal Recorder

It records your terminal, then lets you upload the result.

## Overview / User's Guide

The terminal recorder can be started via the `termrec` binary. Once the application starts, a terminal recording file will be generated per your configuration. While this terminal is active, all _output_ events will be stored in this file, so there is no need to worry about passwords, etc being stored from keyboard input.
_To exit the terminal_: Use any standard way to exit a terminal: namely, `exit` or `ctrl+D` (on an empty prompt).

After the recording, you will be prompted to either upload the file, or exit. If you choose to upload, you will need to supply an operation, file, and description for the event. This will create a new piece of evidence inside the Ashirt application. During this interface, you can navigate with up and down arrows, and can exit some menus with `ctrl+c`. Additionally, many menus support searching for options by pressing `/` and then adding a search term.

### Configuration

This binary supports a few configuration options, and will attempt to load from each configuration level in order to come up with a complete view of how the interaction should be handled. The configuration levels are as follows: First, load from the config file, then replace with defined values from the env vars, then replace with command line switches.

Currently, the configuration file is expected to be found here: `$(HOME)/.config/ashirt/term-recorder.yaml`, however this may be changed in the future on a per-os level. If this file is not found, one will be generated with default (mostly empty) values. Ultimately, the code can be reviewed here for the current location. This default location can be found in `cmd/termrec/config/config.go`.

| Config File Parameter | Env Parameter                         | CLI flag         | Meaning                                                                     |
| --------------------- | ------------------------------------- | ---------------- | --------------------------------------------------------------------------- |
| outputDir             | ASHIRT_TERM_RECORDER_OUTPUT_DIR       | --output-dir     | Determines where to store recording files. Defaults to OS temp directory    |
| recordingShell        | ASHIRT_TERM_RECORDER_RECORDING_SHELL  | --shell       -s | Which shell to use when starting up (defaults to env's SHELL)               |
| operationID           | ASHIRT_TERM_RECORDER_OPERATION_ID     | --operation      | Which operation to upload to (by default -- can be selected during upload)  |
| apiURL                | ASHIRT_TERM_RECORDER_API_URL          | --svc            | Where the **backend** service is located.                                   |
| N/A                   | ASHIRT_TERM_RECORDER_OUTPUT_FILE_NAME | --output-file -o | What filename to use when writing the file locally (and remotely as well)   |
| accessKey             | ASHIRT_TERM_RECORDER_ACCESS_KEY       | N/A              | The Access Key needed to connect with the backend (created on the frontend) |
| secretKey             | ASHIRT_TERM_RECORDER_SECRET_KEY       | N/A              | The Secret Key needed to connect with the backend (created on the frontend) |

### Known Issues

1. It is not "natural" to upload a file that wasn't just recorded.
2. Not possible to re-name files on upload

## Development Overview

This program is actually two programs mascarading as one, and perhaps there's some value in splitting them up. The first sub-application is a somewhat complicated pty recorder. What this does is present a terminal for the user, then captures all _output_ generated from that pty session. The second sub-application is a pretty straight forward file uploader, with some somewhat complex cli interaction. Luckily, in both cases, a lot of the true complexity lies in the libraries being used. That said, there's a still a bit of wiring to go over.

### Requirements

This application has only a few requirements:

* Operating System is Linux or Mac OSX.
* Go 1.12

There's a soft requirement of Make as well, though this can be omitted if the make commands are run directly instead.

### Project Structure

The code is organized over 2 source directories: `root/cmd/termrec` and `root/termrec`. This directory (`root/termrec`) holds all of the core logic of the application, orangized into mini-libraries. The `root/cmd/termrec` contains the actual exectuable, in addition to configuration parsing code.

### Phase 1: Terminal recording

The big idea with this phase is simply starting up the pty console in such a way as to store the result. The pty itself is managed in `cmd/termrec/ptystate.go`. However, in order to capture results, we need to wire up the recorder with a set of `io.Writer`s (one for stdin, one for stdout). Once we have a pair of writers, we need a way to generate events, and something to record those events somewhere. Finally, we need something to intepret those events and convert the raw event into something an interpreter will later be able to parse. Within the code, this is broken into the following components:

| Component / Package     | Role                                                                   | Notes                                                          |
| ----------------------- | ---------------------------------------------------------------------- | -------------------------------------------------------------- |
| Eventer                 | Converts io.Writer bytes into a raw event                              | Can be customized with middleware for more customized behavior |
| Recorders               | Provides mechanism to control output sections / manage stream metadata |                                                                |
| Formatter               | Converts the raw event into a particular format                        | (e.g. asciinema format)                                        |
| Terminal Writer (write) | Manages output mechanism                                               | (e.g. writing to a file)                                       |

As mentioned above, the pty is configured with a pair of io.Writer instances. The first, and primary, is a muxed/multiplexed stdout and eventer. This essentially alows the pty to communicate stdout events both to the user (via stdout) and to the recording system (via the configured eventer). The second writer is a feed off of stdin, allowing the underlying system to react to input-related events as needed. Note, however, that these events are passed unfilterd into the pty, so it is not currently possible to ignore key events.

With this knowledge, the flow then, for output related events is as follows:

pty generates stdout-bound event -> Eventer sees this, and runs through various middleware to generate a raw event -> Event is passed to recorder, which in turn wirtes to the terminal writer -> Terminal writer passes to formatter conform the event -> Terminal writer writes to it's output stream

### Phase 2: Uploading

The upload section is really just a collection of CLI dialogs. The upload dialog itself can be found in `cmd/termrec/appdialogs/upload_prompt.go`, but a lot of the underlying code used there is actually located under `termrec/dialogs`. Once the upload is triggered, the upload action is deferred to `termrec/network`.

### Building

Executables can be created via `make build-termrec-linux`, `make build-termrec-osx` or `make build-termrec-all` (which will make both linux/osx binaries).

### Getting Started

The easiest way to get started is to run `make .env && make update-go`. From here, you can modify the `.env` file that gets created with appropriate values for your dev instance.

### Known Dev Issues
