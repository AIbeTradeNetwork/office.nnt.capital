# AIbeTrade

Platform for trading using automatic mechanics

## Service description

API service for customer registration and authorisation, and control of trading mechanics

## Project structure

### Business layer:

- `/domain`: data models, used in business logic
- `/service`: business logic, separated by logical use cases

### Technical layers:

- `/config`: configuration of application
- `/repository`: database storages
- `/provider`: external service providers
- `/transport`: interfaces for interaction with application

## Development & Test

Scripts for dev tooling and debugging can be found in `/Makefile`