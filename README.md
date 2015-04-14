# get-zk-leader

---

Simple cli app that gets the zk leader node as long as it follows the standard zk sequence numbering scheme.

flags:


|  Flag |                Usage                |                               Description                              |
|:-----:|:-----------------------------------:|:----------------------------------------------------------------------:|
|   zk  | -zk zooserver1:2181,zooserver2:2181 | Comma delimited list of zookeeper servers                              |
|  path |    -path /chronos/state/candidate   | zookeeper path to location containing sequence znodes                  |
| mesos |             -mesos true             | [optional] defaults to false, set to true is querying mesos leader     |
| debug |             -debug true             | [optional] defaults to false, set to true to see panic(err) of errors. |
