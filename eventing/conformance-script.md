# Conformance Steps to manually test Knative Eventing


## Broker Lifecycle 

From: https://github.com/knative/specs/blob/main/specs/eventing/control-plane.md#broker-lifecycle

```
A Broker represents an Addressable endpoint (i.e. it MUST have a status.address.url field) which can receive, store, and forward events to multiple recipients based on a set of attribute filters (Triggers). 

Triggers MUST be associated with a Broker based on the spec.broker field on the Trigger; it is expected that the controller for a Broker will also control the associated Triggers. 

When the Broker's Ready condition is true, the Broker MUST provide a status.address.url which accepts all valid CloudEvents and MUST attempt to forward the received events for filtering to each associated Trigger whose Ready condition is true. As described in the Trigger Lifecycle section, a Broker MAY forward events to an associated Trigger destination which does not currently have a true Ready condition, including events received by the Broker before the Trigger was created.

The annotation eventing.knative.dev/broker.class SHOULD be used to select a particular implementation of a Broker, if multiple implementations are available. It is RECOMMENDED to default the eventing.knative.dev/broker.class field on creation if it is unpopulated. Once created, the eventing.knative.dev/broker.class annotation and the spec.config field MUST be immutable; the Broker MUST be deleted and re-created to change the implementation class or spec.config. This pattern is chosen to make it clear that changing the implementation class or spec.config is not an atomic operation and that any implementation would be likely to result in event loss during the transition.
```



### Requirements: 

If I want to test conformance (MUST, MUST NOT, REQUIRED) for the previous two paragraphs I need: 
- **Prerequisites**: 
    - Knative Eventing Installed. 
    - `kubectl` access to the cluster as defined in the spec: https://github.com/knative/specs/blob/main/specs/eventing/control-plane.md#rbac-profile
    - `yq` installed
- A broker resource: broker.yaml
- A trigger resource that reference the broker: trigger.yaml 
- A trigger resource that doesn't reference the broker: trigger-no-broker.yaml
- A Kubernetes Service that can be addresable to receive and count cloudevents that arrive
  - Clone `https://github.com/salaboy/knative-conformance-experiment` and `cd` to `events-counter` and then run `ko apply -f config/` 
- `curl` to send CloudEvents


### Testing for Conformance: 


Create a broker to test conformance

```
kubectl apply -f broker.yaml
```

Check for default annotations, this should return the name of the selected implementation: 

```
kubectl get broker conformance-broker -ojson | jq '.metadata.annotations["eventing.knative.dev/broker.class"]'
```

Try to patch the annotation: `eventing.knative.dev/broker.class` to see if the resource mutates: 

```
kubectl patch broker conformance-broker --type merge -p '{"metadata":{"annotations":{"eventing.knative.dev/broker.class":"mutable"}}}'
```

You should get the following error: 
```
Error from server (BadRequest): admission webhook "validation.webhook.eventing.knative.dev" denied the request: validation failed: Immutable fields changed (-old +new): annotations
{string}:
	-: "MTChannelBasedBroker"
	+: "mutable"
```

Try to mutate the `.spec.config` to see if the resource mutates: 

```
kubectl patch broker conformance-broker --type merge -p '{"spec":{"config":{"apiVersion":"v1"}}}'
```

**@TODO**: check why this is not returning an error, it seems that a validation webhook is missing


Check for condition type `Ready` with status `True`: 

```
 kubectl get broker conformance-broker -ojson | jq '.status.conditions[] |select(.type == "Ready")' 
```

Running the following command should return a URL

```
kubectl get broker conformance-broker -ojson | jq .status.address.url
```


Create a trigger that points to the broker:

```
kubectl apply -f trigger.yaml
```

Check that the `Trigger` is making a reference to the `Broker`

```
kubectl get trigger conformance-trigger -ojson | jq '.spec.broker'
```

Check for condition type `Ready` with status `True`: 

```
kubectl get trigger conformance-trigger -ojson | jq '.status.conditions[] |select(.type == "Ready")'
```

Congratulations you have tested the **Broker Lifecycle Conformance**!



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











