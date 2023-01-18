# Step 4. Create and deploy component

{%
   include-markdown "../../overhaul_warning.md"
%}

## Create a MACH component

Create a MACH component which we can use in our MACH stack later on.

```bash
mach-composer bootstrap component
```

And follow the wizard to create a new component.

This component can now be pushed to a Git repository.

## Deploy your component

It depends on what component you have, but if you've created a component
containing a serverless function, that function needs to be built, packaged and
uploaded to a **artifact repository**.

For a component created with the `mach-composer bootstrap` command or one of the
provided example components, these commands will be enough:

```bash
./build.sh package
./build.sh upload
```

!!! info "Component deployment"
    The deployment process of a component can vary.<br>
    [Read more](../../topics/deployment/components.md) about component deployments.
