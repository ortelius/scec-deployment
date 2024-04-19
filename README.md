# Ortelius v11 deployment Microservice

> Version 11.0.0

RestAPI for the Deployment Object
![Release](https://img.shields.io/github/v/release/ortelius/scec-deployment?sort=semver)
![license](https://img.shields.io/github/license/ortelius/scec-deployment)

![Build](https://img.shields.io/github/actions/workflow/status/ortelius/scec-deployment/build-push-chart.yml)
[![MegaLinter](https://github.com/ortelius/scec-deployment/workflows/MegaLinter/badge.svg?branch=main)](https://github.com/ortelius/scec-deployment/actions?query=workflow%3AMegaLinter+branch%3Amain)
![CodeQL](https://github.com/ortelius/scec-deployment/workflows/CodeQL/badge.svg)
[![OpenSSF-Scorecard](https://api.securityscorecards.dev/projects/github.com/ortelius/scec-deployment/badge)](https://api.securityscorecards.dev/projects/github.com/ortelius/scec-deployment)

![Discord](https://img.shields.io/discord/722468819091849316)

## Path Table

| Method | Path | Description |
| --- | --- | --- |
| GET | [/msapi/deployment](#getmsapideployment) | Get a List of Deployments |
| POST | [/msapi/deployment](#postmsapideployment) | Create a Deployment |
| GET | [/msapi/deployment/:key](#getmsapideploymentkey) | Get a Deployment |

## Reference Table

| Name | Path | Description |
| --- | --- | --- |

## Path Details

***

### [GET]/msapi/deployment

- Summary  
Get a List of Deployments

- Description  
Get a list of deploymentss.

#### Responses

- 200 OK

***

### [POST]/msapi/deployment

- Summary  
Create a Deployment

- Description  
Create a new Deployment and persist it

#### Responses

- 200 OK

***

### [GET]/msapi/deployment/:key

- Summary  
Get a Deployment

- Description  
Get a deployment based on the _key or name.

#### Responses

- 200 OK

## References
