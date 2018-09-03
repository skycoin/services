# LOGGING ENVIRONMENT
At this point we have a whole system running services in different hosts,
and our only way to have visibility over them is using docker commands on
a swarm manager. But this is not providing that much information over whats
going on inside them, and also will only allow people with credentials to use
the manager nodes to see that information.

In order to allow other people to have access to this information without allowing
them to manage the swarm we want a logging system, that reports logs and statistics
to a dashboard.

### System
1. [logspout-logstash](https://hub.docker.com/r/bekt/logspout-logstash/)
2. [metricbeat](https://www.elastic.co/products/beats/metricbeat)
3. [logstash](https://www.elastic.co/products/logstash)
4. [elasticsearch](https://www.elastic.co/products/elasticsearch)
5. [kibana](https://www.elastic.co/products/kibana)

Logspout and Metricbeat run on every host, the first one retrieving logs, while the second retrieves
metrics.

Logspout sends the data to logstash, that formats that data on the fly, and then sends them elasticsearch.
Metricbeat sends the data to elasticserach directly. Metricbeat also connect to Kibana to tell Kibana
how to retrieve its data and how to display it.

Elasticsearch is a powerful database that allows for very fast text queries, Kibana is a dashboard
that displays data retrieved from elasticseach.

### Deploying the logging stack
Firstable, the 3 heviest services are deployed in tagged hosts, so every time the stack
is redeployed the don't need to download a big docker image, however this can also hurt
the availability of the system, so if its a concern to have high availability for these
services you can tag more than one node, so they are eligible to deploy the services on them.

```bash
docker node update --label-add logstash=true $list_of_nodes
docker node update --label-add node.labels.kibana=true $list_of_nodes
docker node udpate --label-add node.labels.es1=true $node
```
Elasticsearch is deployed as a stateful service using a local volume, so its is
convinient, in case of deploying a cluster of them, to use different tags for every replica.

To deploy the stack, simply:
```bash
docker stack deploy --compose-file ./logger-docker-compose.yaml logstack
```

You can access the Kibana dashboard in the port 8082 of the host that you
deployed Kibana in.

However, there are other considerations to have in mind about this stack:

Both logspout and metricbeat are very lightweight services that won't overload the host, even
in low-specs hardware.

But logstash, elasticsearch and kibana are much more resource-hungry services, and if the stack is
going to be deployed in low-specs hardware those should be hosted outside, while configuring logspout
and metricbeat to connect to the correct addresses.

To change where logstash writes the data just change the value of its ROUTE_URIS environment variable
in the logger-docker-compose.yml.

To change where metricbeat writes the data you need to change the values in logger/metricbeat/metricbeat.yml.
The variables that you need to change are `setup.kibana.host` and `setup.elasticsearch.hosts` to point
at the correct values.

### Security Notes
You may want to secure the access to logstash, elasticsearch and kibana.
You will find information regarding this topic on the following links:
[elastic stack](https://www.elastic.co/guide/en/elastic-stack-overview/6.3/xpack-security.html)
[logstash](https://www.elastic.co/guide/en/logstash/current/keystore.html)
[metricbeat](https://www.elastic.co/guide/en/beats/metricbeat/current/securing-communication-elasticsearch.html)
[elasticsearch](https://www.elastic.co/guide/en/elasticsearch/reference/current/configuring-security.html)
[kibana](https://www.elastic.co/guide/en/kibana/current/using-kibana-with-security.html)

### Configure Kibana
In order to visualize the data with Kibana you need to set indexes to tell it how to
retrieve it from elasticsearch. When Metricbeat connects to Kibana it tells Kibana how to
access its data, but Logstash does not, so in order to visualize Logspout's data you need to
add the next index to it: "logstash-*".

However, Kibana allows the user to define its own visualizations and dashboards, you can read
more about that on the [docs](https://www.elastic.co/guide/en/kibana/current/using-kibana-with-security.html)