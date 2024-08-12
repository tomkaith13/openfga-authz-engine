# OpenFGA based Authz Engine
## What is OpenFGA?
> [OpenFGA](https://openfga.dev/) is an open-source authorization solution that allows developers to build granular access control using an easy-to-read modeling language and friendly APIs.

This PoC hosts a webapp along that connects to an OpenFGA server and adds Authz Models, along with tuple.
It later exposes an endpoint that performs a `check` 