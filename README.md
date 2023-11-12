# Prize Picks Assessment

# Overview
The implementation is a gorilla mux based service with a postgres backend.
The intent is to provide a simple inventory system for a Jurassic park style management system of cages and dinosaurs.

# Requirements
The application was built and tested using the centos flavour of linux. It may require docker and/or a postgres database depending on which version is tested.

Some test scripts require that ``jq`` be installed for formatting output.

## Download and Install

Clone this repository.
To ensure that the correct packages are available use
```
make setup
```
This should ensure that the correct packages are downloaded and installed.
You may additionally install ``mockgen`` and ``golangci-lint`` should you wish to run lint checking or recreate the generated mock files for testing

## Modes
There are 3 basic modes of operation
1. Stand alone server with an external postgres database
2. Stand alone server with the provided postgres database
3. Server and postgres db running in a single docker compose

The server does require a reference file which lists the permitted dinosaurs and their diet. This is provided in the ``species.json`` file.

If the server is not executed in the same directory as this file then the full path can be specified using the ``-sf`` option on the command line

### Stand alone server mode
To build the server use

```make svr```

several environment variables must be set for connecting to a postgres database and can be found in the configuration file

```env.sh```

If the default values here are correct then the environment may be set using

```source env.sh```

This mode requires that the database tables are created manually and the schema can be found in

```dataset/init.sql```

Once the database is read simply start the server using

```./svr```

If all is well then it should report that is listening on the configure endpoint.

### Stand alone with docker compose database

In order to start the provided test database just issue

```docker compose -f docker-compose-db.yml```

This should start and will self initialize the tabls required for the test
Once the database is up and ready the create then in another terminal create the server executable

```make svr```

in a shell set the environment veriables

```source env.sh```

and start the server

```./svr```

If all is well the server should start and repost a list of valid species along with the port it is listening on.

### Docker compose startup (possibly easiest)

Before starting this version please ensuer that the server image is created using
```make docker-image```
** do not try to use this image as it does not have and entrypoint defined. That is defined in the docker compose file.
This will create an image ``dino_svr:latest`` for use in the docker compose
To start issue
```docker compose -f docker-compose.yml up```
All being well both the database and server will start and be available on port ``8000``

_NOTE_ I did notice that despite the specified dependency of the server on the postgres database that on occasion the server would start before the postgres endpoint was ready to accept connections. To avoid this a small sleep was added prior to starting the app server. It is possible that this may not be sufficient depending upon system performance and needs increasing to ensure that correct boot sequence is followed.


## Data Model
The data model is quite simple and self explanatory. It is composed of only two tables and can be found in the file ``dataset/init.sql``

## Rest Api definitions

Typically the rest api would not be detailed here, but using swagger or some other means, but in the interest of brevity th following is a short description of the available rest sdk

```GET /dino/list```

Will return a full list of all saved dinosaurs in json. It does not paginate and so may provide a length list. An alternate for is available for filtering on species
```/dino/list?species=<species name>```

As above but a species name is provided as a parameter. However as above this does not paginate and so could return a lengthy reply

```GET /cages```

This lists all cage information in json. As above as this may be quite extensive and so the paramater for may be used

```GET /cages?status=<ACTIVE|DOWN>```

This list cages in json for the provided status. Status must be either ACTIVE or DOWN. If an invalid status is given the server will return an error.

```POST /cage/{diet}/add```

This will create a new cage for the given dietary requirements of the species to be placed therein. It takes no payload. The diet must be either _H_ or _C_ or an error will be returned. There is no payload for this and it will be ignored if passed. The reply upon success will be the numerical identifier of the cage in json format. This api call takes an optional ``cap=`` parameter that will specify the dinosaur capacity.

```GET /cage/{cageid}/list_dinosaurs```

Returns a json response of the dinosaurs in a given cage.

```POST /cage/{cageid}/status/{status}```

Will set the given cage to the status _ACTIVE|DOWN_. If the cage is occupied the cage may not be powered down and will return an error

```POST /cage/{cageid}/add_dino```

Will add a dinosaur to the json provided dinosaur to the specified cage. Upon success code _200_ is returned and an error if the cage is full or not of the correct dietary requirements. For reference the payload may be seen in the file ``scripts/dino_c.h``

```POST /species/add```

Will add a new species to the in memory reference lookup. This is not persisted and so will not be available when the server is restarted. Additionally it is scoped to the receiving server instance and as such will not propagate in a clustered environment. An example payload for this can be found in  ``scripts/add_species.sh`` script.

```GET /species/list```

returns a json formatted list of available species. The result is not paginated and is returned from the memory cache.

```POST /cage/{cageid}/add_dino```

Places a dinosaur from the provided payload into the given cage. Examples of json payload for a dinosaur can be seen in ``scripts/dino_c.json```

```POST /dino```

This is a general purpose add dinosaur functionality. It takes a json dinosaur payload of which an example can be seen in ``scripts/dino_c.json``.

If a cage with capacity or of the required type does not exist one is created.



## Testing
Although the provided testing is woefully short of full coverage it does at least demonstrate the use of ``mock`` and ``httptest``. In order to run the testing use
```
make test
```

## Scripts
In order to assist testing some scripts have been provided in the ``scripts/`` directory. Most are quite self explanatory and use ``curl`` so the api invoked can be checked with the above documentation

```list_cages.sh``` - lists available cages

```list_dinos.sh```

lists the dinosaurs. If a command line paramater then the dinosaurs for a given cage are returned

```list_species.sh``` - lists the available species and their dietary designation

```cage_status.sh -id <cage id> -s <ACTIVE|DOWN>```

Set the status of a given cage to the provided status

```add_cage.sh -diet <H|C> -cap <capacity>```

Create a new cage of th specified dietary requirements. It takes a single paramater specifying the dietary requirements.

```put_in_cage.sh -id <cage id> -f <json file with dino info>```

Places a dinosaur in a given cage. If the cage is at capacity of not of the required dietary requirements then the it should fail.

## Mocks
Should you wish to recreate the mock files, you will need to have mockgen installed on your system. Running
```
make mocks
```
should recreate them and place the generated output in the ``mocks/`` directory

## Improvements/Shortcomings
1. Paginate large query responses
2. Provided more extensive and granular filtering
3. Provide referential integrity
4. Implement the species list as a reference table
5. Relax condition that a cage must be created for an explicit diet
6. The cage capacity should be configurable
7. Improve error messages from the data access layer
8. Improve the documentation for th rest api
9. code comments 
