# Detecting changes

One of the most difficult elements of a deployment is detecting changes. Mach
Composer uses a combination of Terraform and comparing configurations to detect
what components have changes.

In basic terms Mach Composer stores a hash of the configuration of each deployed
element in the state. When a new configuration is applied, the hash of the new
configuration is compared to the hash of the old configuration. If they are
different, the component is marked as changed and will be updated

This hash is based on the definition of the component, declared dependencies and
the variables and secrets. If any of these have changed the hash will also
change, telling Mach Composer that the component should be updated.

[//]: <> (@formatter:off)
!!! warning "Changes in variables"
    Due to current technical limitations, Mach Composer cannot detect changes 
    in actual values of variables or secrets, but only in their shape. This 
    means that if only a password is changed Mach Composer will not detect a 
    change. 

    If this is the case you can force a change by changing the order of the 
    variables as a quick-fix. This will be fixed in a future release. 
[//]: <> (@formatter:on)

