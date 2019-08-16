# 'Break Sawtooth History' Demo
This repo demonstrates a suspected bug in Sawtooth 1.1 where historical state can seemingly "disappear" after number of
writes/blocks have passed. The bug in question is triggered by the following steps:

1) Write a piece of state to an address, read it back to verify and notate the current head.
2) Do some number of further operations, including (for clarity, but not necessary), over-writing the state written in step 1.
3) Attempt to read back the original address at the previously notated head.

If enough has happened in step 2 to trigger the bug, the REST API will return an Error 75 (State Not Found).

For any questions about this demo, please contact me at `bill` `at` `taekion.com`.

## Implementation
This repo provides a simple transaction processor that writes blocks of bytes to Sawtooth at a given address referenced by UUID, as well
as a client library to interface with the TP and a CLI program that uses the client library to invoke the bug.

## Requirements
- Docker
- Docker Compose
- Git
- Make/Gmake

## Instructions

Clone the repo into a convenient location:
    
    git clone https://github.com/taekion-org/break-sawtooth-history
    
Build the Docker containers and start the compose environment (assumes you have `make/gmake` installed):

    cd break-sawtooth-history
    make docker
    docker-compose -f docker/break-sawtooth-history.yaml up

The compose script starts up a container that is ready to execute the test inside. To see the
help message for the test:

    docker exec -ti break-sawtooth-history-cli break-sawtooth-history-cli --help

You can execute the test with default settings as follows.

    docker exec -ti break-sawtooth-history-cli break-sawtooth-history-cli --url http://rest-api:8008

For each successful payload submission, the test will print a `.` If a submission is refused due to throttling, the
test will print a `(rXXX)` and will repeat the `XXXth` submission.

If the bug shows itself, you should see output similar to the following:

    Final Check:
    Requesting address 5139f58349ee7641ede01a288738dfbfb1dccebeaba039dc82b8c4b45edd85e4b9d363 at head 498eec2db2daf64cc01d674f6e1a55b8c09c9b6b08228080cd689947968032547ca3d827f3e4614aaad9ab13f62e7a2a516b80a6fa6c6631410fbb5c02972fe5
    Sawtooth Error: 75 -- Transport Error: Sawtooth REST API Error: method=GET, status=404, code=75, title=State Not Found

You can further verify this yourself by sending an HTTP request to the REST API with the parameters shown (values are from the previous example):

    # To try using httpie    
    http get http://localhost:8008/state/5139f58349ee7641ede01a288738dfbfb1dccebeaba039dc82b8c4b45edd85e4b9d363?head=498eec2db2daf64cc01d674f6e1a55b8c09c9b6b08228080cd689947968032547ca3d827f3e4614aaad9ab13f62e7a2a516b80a6fa6c6631410fbb5c02972fe5
