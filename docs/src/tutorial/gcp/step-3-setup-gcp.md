# Step 3. Prepare your GCP environment

## Create shared/services project

## Create project per MACH stack
### Enable services

```bash
gcloud services enable compute.googleapis.com
...
```

**Custom domains**

Custom domain names are not supported by API Gateway. 
For custom domains, we need to create a load balancer and direct requests to the `gateway.dev` domain of the deployed API.

- [Setting up load balancing 'the hard way'](https://cloud.google.com/blog/topics/developers-practitioners/serverless-load-balancing-terraform-hard-way)
- [API gateway behind load balancer](https://medium.com/swlh/google-api-gateway-and-load-balancer-cdn-9692b7a976df)
- [Load balancing with Terraform](https://cloud.google.com/community/tutorials/modular-load-balancing-with-terraform)
