# Setting up GCP

!!! todo
    GCP is not yet supported by MACH composer, unfortunately.
    We are happy to accept contributions to add this! And think with you if you plan to work on it. Feel free to reach out to us via [opensource@labdigital.nl](opensource@labdigital.nl).

## What needs to happen to support GCP?

Luckily, adding GCP should not be a lot of work. Most of the implementation boils down to adding the right terraform code in MACH composer, to generate the nessesary resources. This is primarily about setting up an API gateway that ties together many services, through [GCP's Terraform support](https://registry.terraform.io/providers/hashicorp/google/latest/docs).

- [ ] Add support in MACH yaml to support GCP cloud
- [ ] Add terraform templates to MACH composer, to support GCP resources
- [ ] Reference implementation of a CloudRun serverless function
- [ ] Implement routing using Google API Gateway
- [ ] Base Google Cloud account setup where the architecture runs in
- [ ] Ideally: extend component bootstrapper and component cookiecutter to include Google Cloud setup