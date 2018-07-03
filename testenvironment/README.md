# Test Environment

This project sets up a test environment containing:

   - [x] 6 Skycoin nodes
   - [x] A load balancer ([Traefik](https://traefik.io/))
   - [x] A set of CockroachDB
   - [x] Metrics
   - [x] Logging system (working through swarm)
   - [x] Update system (using swarm)
   - [x] The coin nodes communicate with each other
   - [ ] Autoupdate
   - [ ] Gui dashboard for users (traefik dashboard, kibana dashboard, cocroach dashboard)
   - [x] Node registration (through swarm)


This environment is orchestrated by [Docker Swarm](https://docs.docker.com/engine/swarm/), so it will
facilitate service discovery, handle system failures, handle system
scalability and automate updates.

# Installation

### Create a swarm

Firstable we need to create a [Docker Swarm](https://docs.docker.com/engine/swarm/) environment. Docker Swarm
comes packaged within [Docker CE](https://www.docker.com/get-docker).

You will need to install Docker on every machine that is part of
your test environment cluster. Then follow the [documentation](https://docs.docker.com/engine/swarm/swarm-tutorial/create-swarm/) to
create a swarm.

### Create a registry

We need to have an image repository that will make our images public to
every other node in the swarm. You can:
1. Use a public image repository, like https://hub.docker.com/.
2. Host a private repository in the Swarm cluster.

If you choose to hold your own image repository locally:
```bash
    # We will deploy the registry into a single node.
    # First tag the node with a label
    docker node update --label-add registry=true $node_name

    # Then pull the image and create the service, which will
    # be deployed on the tagged node
    docker pull registry:2
    docker service create --name registry --constraint 'node.labels.registry==true' registry:2
```

To make it available to other hosts you need to follow the steps defined in the [docs](https://docs.docker.com/registry/deploying/#run-an-externally-accessible-registry)
in the section called "Get a certificate".

### Create a Docker image for a new coin

Now we need to create a new fiber coin that allows us to host our
own network for that coin.

Follow the instructions to download and install
[Skycoin](https://github.com/skycoin/skycoin).

Then, inside the project package, move to the development branch:

    git checkout develop

Follow the instructions [here](https://github.com/skycoin/skycoin/tree/develop/cmd/newcoin) to install the
newcoin binary.

Now we can use it to create a new coin with the configuration
file provided in this project (skynode/fiber.toml).

From the skycoin project root, run the following:
```bash
newcoin createcoin --coin skycoin \
    --template-dir template \
    --coin-template-file coin.template \
    --visor-template-file visor.template \
    --config-dir $path_to_testenvironment_dir \
    --config-file skyfiber.toml
```

This will generate a new coin (also called skycoin) with our
configuration. Now, following the steps defined [here](https://github.com/skycoin/skycoin/blob/develop/docker/images/mainnet/README.md)
we will create a new Docker image for this coin. This Docker image
will later be used to build another Docker image to deploy in the swarm.

Also from the skycoin project root, run the following:

    docker build -f docker/images/mainnet/Dockerfile -t skycoin:new .

### Create a Docker image to deploy in the swarm

Now, get into testenvironment project root and run the following:

    docker build -f sky-node/Dockerfile -t registry/sky-node:1 sky-node

Registry here refers to the one you are using, either hosted externally or not. If you
host it locally you may need to append the port where the service is listening, for example
registry:5000/sky-node:1.

Push the image as normal, so its available by the other services in the swarm.

### Deploy the services

Now, also from this project root run:

    docker stack deploy --compose-file ./docker-compose.yaml skystack

With this command swarm starts deploying the services, now you can perform
commands on both docker and swarm to get information from the services,
read the logs, update the services, etc.

All the services deployed are preceded by the name of the stack, in this case
called skystack.

If for example we want to see the logs:

    docker service logs -f skystack_skycoin-node

If we want remove all the services on the stack:

    docker stack rm skystack

If we want to update the services we can update the docker-compose.yaml
specifiying new versions of the images, then use the command to deploy
again and swarm will handle the update of the services with no downtime.

We can also update the services manually and one by one.

You can read more about docker swarm commands in the official documentation.

### Access the coin nodes from outside

The coin nodes cannot be accessed directly from outside. The way
to communicate to them is through the load-balancer that is running
in front of them. There is one instance of the load-balancer running
on each manager node.

You can access a load-balancer through the ip of each manager node.

These load-balancers exposes two endpoints to contact the coin services:
1. http://{manager_node_ip}/api/... to access the coin json API that
   usually runs on the port 6420.
2. http://{manager_node_ip}/node/ to connect to the one of the nodes.

Additionally, you can access http://{manager_ip}:8080 to see the
load-balancer dashboard.

### Security notes

If security is a concern and the environment is going to be publicly
available both traefik and the image registry should enable TLS options
and make use of valid certificates issued by a CA.

You can find instructions to do so in these links:
1. https://docs.docker.com/registry/deploying/#run-an-externally-accessible-registry)
2. https://docs.traefik.io/user-guide/docker-and-lets-encrypt/
3. https://docs.traefik.io/user-guide/cluster-docker-consul/

### Accessing logs and metrics
In order to visualize whats happening inside the hosts and containers you
can additionally deploy a second stack.

You can find how [here](./LOG_STACK.md)