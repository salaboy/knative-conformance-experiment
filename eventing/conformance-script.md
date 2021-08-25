# Knative Eventing Conformance Test Plan

This document descibe a test plan for testing Knative Eventing Conformance based on the specs that can be found here: https://github.com/knative/specs/blob/main/specs/eventing

The specs are splitted into Control Plane and Data Plane tests, this document follows the same approach and further divide the tests into further subsections. 

# Control Plane

https://github.com/knative/specs/blob/main/specs/eventing/control-plane.md


## Requirements: 

If you want to test conformance (**MUST, MUST NOT, REQUIRED**) you need: 
- **Prerequisites**: 
    - Knative Eventing Installed. 
    - `kubectl` access to the cluster as defined in the spec: https://github.com/knative/specs/blob/main/specs/eventing/control-plane.md#rbac-profile
    - `jq` installed
- A Kubernetes Service that can be addresable to receive and count cloudevents that arrive
  - Clone `https://github.com/salaboy/knative-conformance-experiment` and `cd` to `events-counter` and then run `ko apply -f config/` 
- `curl` to send CloudEvents

## Test Plan for Control Plane

- [Broker Lifecycle Conformance](broker-lifecycle-conformance.md)
- [Trigger Lifecycle Conformance](trigger-lifecycle-conformance.md)
- [Channel Lifecycle Conformance](broker-lifecycle-conformance.md)
- [Subscription Lifecycle Conformance](subscription-lifecycle-conformance.md)


## Other Notes about Control Plane Spec

- The [**Resource Lifecycle Section**](https://github.com/knative/specs/blob/main/specs/eventing/control-plane.md#resource-lifecycle) (Broker, Trigger, Channel, Subscription and Destination Resolution) can be automated by creating resources and running commands as described in the previous section.

- The [**Event Routing Section**](https://github.com/knative/specs/blob/main/specs/eventing/control-plane.md#event-routing) describes internal Broker behaviours, which needs to be observed and inferred based on results. In the **Topology Based Routing** section, the sentence `Before acknowledging an event, the Channel MUST durably enqueue the event (be able to deliver with retry without receiving the event again).` implies that we can check this from outside the Channel. Conformance should check that events are delivered, if they are `durably enqueued` is loosly defined here, and impossible to check from the outside. This is also a Data Plane concern not a control plane one. 

- The [**Detailed Resources Section**](https://github.com/knative/specs/blob/main/specs/eventing/control-plane.md#detailed-resources) can be tested by creating different resources with the REQUIRED fields and see if they work. There are tons of optionals, which we shouldn't be covering at this stage, so automating this and creating the resources shouldn't take much. 


# References

- Feature Language already defined in reconciler-tests: https://github.com/knative/eventing/blob/main/test/rekt/features/broker/control_plane.go#L95-L120  

- `Event Library` for Data Plane tests: https://github.com/knative/eventing/blob/main/test/test_images/event-library/main.go

- `CloudEvents Conformance` CLI (listen, and invoke CE using Events Library format): https://github.com/cloudevents/conformance

- `EventsHub` for testing events delivery (look for it inside `reconciler-tests`)

# Data Plane

https://github.com/knative/specs/blob/main/specs/eventing/data-plane.md

## Test Plan for Data Plane

- [Event Ack and Delivery Retry](event-ack-and-retry.md)


# Emit Events

```
curl -X POST -H "Content-Type: application/json" \
  -H "ce-specversion: 1.0" \
  -H "ce-source: curl-command" \
  -H "ce-type: ConformanceTested" \
  -H "ce-id: 123-abc" \
  -d '{"name":"Salaboy testing conformance"}' \
  http://broker-ingress.knative-eventing.127.0.0.1.nip.io/default/conformance-broker 
```











