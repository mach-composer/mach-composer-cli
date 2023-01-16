# Autonomous component deployments

We want to provide the ability for components to deploy themselves independently from other components, through their own CI/CD pipelines. This will allow teams to operate more independently from other teams, while still building within a larger composable architecture.

!!! info "In design phase"

    This feature is currently being designed, so consider this page our current thinking around the implementation.

## How would it work?
Each component is in fact a Terraform module that can be executed individually when provided with the required input parameters. Our current thining is that MACH composer ☁️ will provide all input parameters (Terraform variables) for each of the sites that this component needs to be deployed to. And after the deployment is done, the output parameters will be stored in MACH composer ☁️ again.

<div class="mermaid">
sequenceDiagram
    participant D as Developer
    participant S as SCM
    participant CI as CI/CD
    participant MC as MACH composer
    participant MCC as MACH Composer☁️ 
    participant TF as Terraform

    D->>S: Merges new code
    S->>CI: run build and tests
    CI->>MC: deploy component action
    MC->>MCC: fetch required configuration
    MC->>TF: execute terraform apply for each site
    MC->>MCC: send outputs
</div>
