# Overview of primary flows

## Pushing a component to the component registry

MACH composer is primarily used to deploy components in an individual site, and provide it with the right context (settings, endpoints, etc).

Components themselves are responsible for publishing their own artifacts into the component registry (which can be a simple S3 bucket). And usually that is implemented by configuring a CI/CD pipeline that manages that automatically.

An artifact is usually a zip file containing a serverless function, i.e. `my-component-vXYZ.zip`. These can be conveniently generated using the serverless framework, using `sls package`, but can also be built using a different process.

!!! tip "A component is a Terraform module"
    Next to the component publishing its artifact into a component registry, it should also provide the necessary terraform resources for the component. Usually, components have a `/terraform` directory in the root of the component, containing all those resources. Effectively, this makes a component a terraform module. And MACH composer in turn leverages [Terraforms 'modules sources'](https://www.terraform.io/docs/language/modules/sources.html) functionality to 'pull together' different modules from different Git repositories.

<div class="mermaid">
sequenceDiagram
    participant D as Developer
    participant S as SCM
    participant CI as CI/CD
    participant C as Component Registry
    
    D->>S: Merges new code
    S->>CI: run build and tests
    opt if master branch
        CI->>C: Push artifact to component registry
    end
</div>


## MACH composer deployment

MACH composer itself is primarily a code generator that generates the required Terraform code per site in the YAML configuration. Optionally (though recommended) MACH composer decrypts the YAML file through SOPS.

<div class="mermaid">
sequenceDiagram
    participant D as Developer
    participant M as MACH composer
    participant S as SOPS
    participant T as Terraform
    participant SC as component SCM
    participant C as Component Registry
    participant MACH as MACH services
    participant Cloud as Cloud (AWS/Azure/GCP)
    
    D->>M: executes mach apply
    opt if encrypted with SOPS
        S->>Cloud: fetch encryption key
        S->>M: decrypt & return yaml configuration
    end
    loop per site configuration
        M->>M: generates Terraform code based on YAML
        M->>T: execute Terraform code
        loop per component
            SC->>T: fetch terraform module from SCM
            C->> T: fetch component
        end
        T->>MACH: create/configure MACH resources
        T->>Cloud: create/configure cloud resources
        T->>Cloud: deploy components (i.e. Lambda zipfiles)
    end
</div>


!!! tip "Running MACH composer in CI/CD is a best practice"
    For production deployments we recommend running MACH composer in a CI/CD pipeline. This because running it requires access to sensitive resources and should be secured properly, as well as providing a good audit-trail about who deployed what at which time.