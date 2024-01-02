# Talks

# Sub0 2024 Conference

Talk Ideas

- Privacy Preserving Node Observability
- One Click Parachain Deployment

# Privacy Preserving Web3 Observability with Telescope

[https://www.notion.so](https://www.notion.so)

Observability is a major challenge in the Polkadot ecosystem, builders, node operators and parachains do not have proper tooling for having a 360º observability insight into their nodes and dApps performance. 

Developers currently depends on traditional monitoring tools like Prometheus and Grafana, which fail to capture chain-specific data and require extensive technical knowledge and effort to set up. This reliance leads to delayed error detection and negatively impacts network health.

Also, while handling logs and metrics is straightforward with these tools, tracing is more challenging. Effective tracing demands that each component is in a request's path forwards tracking information, a complex task in decentralized blockchain systems compared to centralized platforms. 

Currently, most chains directly apply Web2 observability tools and methods to Web3 dApps without taking note of the nuances we have with web3 dApps. There is no standardized approach for instrumenting, monitoring, and reporting in the realm of dApps and Web3 protocols.

To ensure the reliability of parachains and nodes, it's essential to develop a robust Web3 observability infrastructure that addresses these challenges.

**Let's ask ourselves the following Questions:** 

- What is the ease of access for developers to view  events log for a chain without needing to set up their own node?
- Is there a reliable and verifiable  source for obtaining metrics, logs, and traces from blockchain nodes? Similar to how Subscan provides detailed on-chain transactions, how can we monitor and evaluate the performance of the constituent nodes of a blockchain in a decentralized manner?

A good observable system should be able to i) capture data with context, ii) query data to get full insight iii) proactively report on the insight found and use that in either debugging an issue or preventing a major issue from happening.

**How do we achieve this?** 

This talk aims to help solve that by presenting a one-solution for all observability challenges, where netowkrk can easily monitor node operator running their clients and node operator, builder and validator runners also have deep insight into the runing of their clients that even if something happens they immediately know how to fix it with or without any technical knowledge

This talk is broken down into 3 catgories:

- Highlighting the problems stakers, builders and node operators faces when running nodes and building dApps
- Review the challenges faced by current setup popularly utilized and propose improvements
- Describe how to setup and monitor a self manage robust observablity infrastructure using Telescope
- Explain the design and architecture of privacy preserving observability and how telescope achieves that

# **Observability Challenges in Web3**

- Web3 products complexity makes it difficult to quickly identify and resolve issues in Nodes and dApps.
- Measure Everything, Learn Nothing
- Technical difficulties for builders in implementing comprehensive monitoring of their nodes.
- Blockchain networks and chains find it difficult to evaluate the effectiveness of Validator Runners and Node Operators without end user pushing metrics and logs to them
- Blockchain Network would rather focus on their main business logic rather helping developers and node operators troubleshoot issues
- Major debugging challenges due to inadequate insights and log visualization for node operator activities.
- No real-time alerting systems to proactively notify Builders, Node Operators and Parachains about performance issues or problems in relay, parachains and dApps built on them.

# What does an Ideal Solution looks like?

- Generally available yet verifiable Node Observability:
    - Parachains gain comprehensive insight into all Node Operators running their clients.
    - Node Operators receive a ready-to-use monitoring solution for their entire node fleet.
- Ready-to-Use Monitoring:
    - Early warning alerts for node performance issues.
    - Web3-focused insights specifically tailored for Polkadot Parachain operations.
- Deep insight & analytics into all Node Operators running their Charon clients
- Node Operators & Enterprise Customers should have deep out of the box Monitoring solution for all their node fleets.
- Multiple Channels Alert subscription for Node Operators receive customized alerts.
- Integrated SLO for both Node Operator and Networks
- Intelligent Log Analysis

# **Introducing Telescope**

1. Central Node Observability:
    - Node Operators receive a ready-to-use monitoring solution for their entire node fleet.

- Parachains gain comprehensive insight into all Node Operators running their clients.

Telescope is a one-stop observability tool for web3 observability, telescope helps builders and networks collect, store and visualize all monitoring data such as metrics, logs, traces and runtime events thereby saving time and money by reducing operational overhead, improving security and helping deliver a consistent, seamless observability expereince.

**Core Features**

- Easy to Setup
- Privacy Preserving Observability Setup
- Zero Setup Costs & Fixed Monthly Fee
- Handoff observability & focus on core business logic
- Supports Multiple Observability Signals(logs, traces, metrics)
- AI powered log ag and event analysis: threat and anomaly detection, event correlation and automated retrospective
- Multi-Tenancy setup for both networks and builders
- Simple Dashboard that only displays what matters
- Built-in Runbook
- Observability API that allows client to automate their entire observability pipeline alongside their application deployments

# **How Telescope works**

**Observability Architecture** 

![blockops web3 observability(beta).drawio.png](./Talks-sub0/blockops_web3_observability(beta).drawio.png)

## How Telescope Preserve Privacy

![telescope-privacy-preserving-architecture.png](./Talks-sub0/telescope-privacy-preserving-architecture.png)

- Dropping all labels containing sensitive information like IP Address, Waller addresses from source
- Encrypting sensitive information with a source private key and only displaying the public key as the unique identifier at destination

## Benefits

- Proactive Incident Management: Improve developer experience by allowing your team to prevent issues before they happen.
- Performance Efficiency:  Understand and increase the performance and efficiency of your app.
- Full Application Insight: Gain insights about user behavior and growth trends that are currently unavailable.
- Seamless Integration: Telescope has a cli and sdk that makes it easy for developers to instruments observability into their dApps and also Node Operators to easily install a one in all observability tooling

