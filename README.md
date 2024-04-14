# processor-framework
A lightweight framework for extract-transform-load (ETL) based processes that interfaces with the Cluster.tools platform.

This library is still in its early days, more documentation will be added as it's used for internal projects.

## writting your own processors
The framework is responsible for goroutine instantiation and scaling based on configurations passed to it. Processors are given full control of
the instantiated goroutine unless they are defined as the extract process in Stream mode. An extract processes can be terminated by the framework upon
receiving an SIGINT or SIGTERM signal from the kernel so that blocking code can be placed within it. Blocking code places in a transform or load process **will block the framework from terminating** to ensure data beyond the extract stage cannot be lost.

## example
A quick example on how to get started with the framework can be found [here](https://github.com/GabeCordo/processor-template)

### disclosure

This repository is not related to the contributing members (of the repository) to the organizations they currently belong, the work they have, currently, or will perform at such organizations. All work completed within this repository pre-dates these organizations. All work completed withon this repository shall not be through company resources. Where "company resources" includes but is not limited to working hours, intellectual property, and electronic devices.
