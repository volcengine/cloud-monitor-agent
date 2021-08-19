### Cloud Monitor Agent

- Cloud monitor agent is an application that runs inside each instance. the application will collect the monitoring
  metrics within the instance and report them After the user authorizes.

### Supported Metrics

- Cloud monitor service will provide CPU, Memory, Disk, Network and other monitoring indicators to meet the basic
  monitoring operation and maintenance needs of the server.

### Supported systems

- Linux
    - CentOS6.9、CentOS7.6、CentOS7.7 and CentOS8.3
- Windows
    - Not support now

### Performance

- In ECS: CPU<2%, Memory<30M

- Although under normal situation, the agent consumes very little resources.we still use Cgroup v1 to limit the
  resources of agent used:
    - CPU: use cpu.share to limit agent will not use of a lot of slice when cpu shortage.
    - Memory: use memory.limit_in_bytes to limit agent will not exceed 50M
