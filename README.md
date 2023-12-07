# ASHIRT

Adversary Simulators High-Fidelity Intelligence and Reporting Toolkit. This repo contains only the server side and frontend code for ASHIRT. You can find the ASHIRT client [here](https://www.github.com/ashirt-ops/ashirt) and aterm [here](https://www.github.com/ashirt-ops/aterm).

## Table of Contents

- [Background](#background)
- [Install](#install)
- [Configuration](#configuration)
- [Usage](#usage)
- [Contribute](#contribute)
- [License](#license)

## Background

Documenting and reporting is a key part of our jobs and generally the part we all look forward to the least. Compared to the rest of the work we do it's not the most fun and by the time we get around to it, it's not always clear exactly what happened or we don't have the evidence to prove it. Teams generally solve this with ad hoc solutions for note taking, recording and sharing screenshots, and collecting other evidence but these solutions rarely scale, are not always easily shared, and typically require manual steps to manage. Having to dig through a pile of evidence after an operation to find the one screenshot you need, if you even have it, can be cumbersome especially as evidence starts to span multiple operators and computers. ASHIRT attempts to solve this by serving as a non-intrusive, automatic when possible, way to capture, index, and provide search over a centralized synchronization point of high fidelity data from all your evidence sources during an operation. You can learn more in our 2020 presentation linked below. 

[![2020 presentation](https://img.youtube.com/vi/NNPsfTbDPG0/0.jpg)](https://youtu.be/NNPsfTbDPG0)

## Install
If you just want to try ashirt out locally, check out the [quickstart docs](https://github.com/ashirt-ops/ashirt-server/wiki/quickstart). If you're looking for a production-ready deployment reference, check out [ashirt-deployments](https://github.com/ashirt-ops/ashirt-deployments/tree/main/ecs).

Detailed instructions for building and installation are available for the [frontend](frontend/Readme.md) and [backend](backend/Readme.md).

## Configuration

All configuration options for the backend are described [here](backend/Readme.md). These cover the various components and configuration options necessary for deployment and outlines how the components interact.

## Contribute

Please refer to [the contributing.md file](Contributing.md) for information about how to get involved. We welcome issues, questions, and pull requests.

## Maintainers

- John Kennedy: john.kennedy@yahooinc.com
- Joe Rozner: joe.rozner@yahooinc.com

## License

This project is licensed under the terms of the [MIT](LICENSE-MIT) open source license. Please refer to [LICENSE](LICENSE) for the full terms.
