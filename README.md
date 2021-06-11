# DIPLOY

A utility to manage your linux deployments.


## How does it work

- You start a server using `diploy server`
- You create a `diploy.yml` with below format file and run `diploy add` in the same directory.
- **diploy** exposes a http server for instructions [take care of *firewall* !!]
- Now you get endpoits for each configuration you add with the below format.
- Use these as webhooks with Github or just manually.
- Or just use the CLI.

## diploy.yml

```
name: <Application Name>    // These should be unique across installation
update:                     // Specify command to update
                            // codebase (optional [default `git pull`])
  command: git pull         // Command/Script to run to update
  type: command             // Whether a script (run using /bin/sh) or command
build:                      // To build stuff
  command: echo This is build
  type: command
run:                        // To start the application
  command: echo This is run
  type: command
```

## endpoints
All requests are POST requests.

#### Start processes
start:
  - update codebase: `/start/update/{name}`
  - build application: `/start/build/{name}`
  - run application: `/start/run/{name}`

<!-- #### Add an application [IDK why this is there]
<!-- TODO: fix this -->
<!-- add:
    `/add` <br/>
    Body: diploy.yml file (stringified ["Content-Type": "text/plain"]) -->

#### Stop processes for a given application
stop:
    `/stop/{name}`
  

### Caveats
Processes started with **diploy** will also stop if diploy is stopped.

### Todo
See the dedicated file.


## Stuff Used

- golang
- gorm + sqlite
- gorilla/mux
- cobra
- go-yaml