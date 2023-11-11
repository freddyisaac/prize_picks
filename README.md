# Prize Picks Assessment

## Download and Install

Clone this repository.
To ensure that the correct packages are available use
```
make setup
```
This should ensure that the correct packages are downloaded and installed
You may additionally install ``mockgen`` and ``golangci-lint`` should you wish to run lint checking or recreate the generated mock files for testing

## Modes
There are 3 basic modes of operation
1. Stand alone server with an external postgres database
2. Stand alone server with the provided postgres database
3. server and postgres db running in a siongle docker compose

## Data Model

## Rest Api definitions



## Testing
Although the provided testing is woefully inadequate it does at least demonstrate the use of ``mock`` and ``httptest``. In order to run the testing use
```
make test
```

## Scripts

## Mocks
Should you wish to recreate the mock files, you will need to have mockgen installed on your system. Running
```
make mocks
```
should recreate them and place the generated output in the ``mocks/`` directory
