# Trigger Lifecycle 

From: https://github.com/knative/specs/blob/main/specs/eventing/control-plane.md#trigger-lifecycle

> A Trigger MAY be created before the referenced Broker indicated by its spec.broker field; if the Broker does not currently exist or the Broker's Ready condition is not true, then the Trigger's Ready condition MUST be false, and the reason SHOULD indicate that the corresponding Broker is missing or not ready.

> The Trigger's controller MUST also set the status.subscriberUri field based on resolving the spec.subscriber field before setting the Ready condition to true. If the spec.subscriber.ref field points to a resource which does not exist or cannot be resolved via Destination resolution, the Trigger MUST set the Ready condition to false, and at least one condition MUST indicate the reason for the error. The Trigger MUST also set status.subscriberUri to the empty string if the spec.subscriber.ref cannot be resolved.

> If the Trigger's spec.delivery.deadLetterSink field it set, it MUST be resolved to a URI and reported in status.deadLetterSinkUri in the same manner as the spec.subscriber field before setting the Ready condition to true.

> Once created, the Trigger's spec.broker MUST NOT permit updates; to change the spec.broker, the Trigger can instead be deleted and re-created. This pattern is chosen to make it clear that changing the spec.broker is not an atomic operation, as it could span multiple storage systems. Changes to spec.subscriber, spec.filter and other fields SHOULD be permitted, as these could occur within a single storage system.

> When a Trigger becomes associated with a Broker (either due to creating the Trigger or the Broker), the Trigger MUST only set the Ready condition to true after the Broker has been configured to send all future events matching the spec.filter to the Trigger's spec.subscriber. The Broker MAY send some events to the Trigger's spec.subscriber prior to the Trigger's Readycondition being set to true. When a Trigger is deleted, the Broker MAY send some additional events to the Trigger's spec.subscriber after the deletion.

# Testing Trigger Lifecycle Conformance: 

We are going to be testing the previous paragraphs coming from the Knative Eventing Spec. To do this we will be creating a trigger  checking its Ready Status and then creating a Trigger that links to it by making a reference. We will also checking the Trigger Status, as it depends on the Broker to be ready to work correclty. We will be also checking that the broker is addresable by looking at the status conditions fields. Because this is a Control Plane test, we are not going to be sending Events to these components. 

You can find the resources for running these tests inside the `control-plane/trigger-lifecycle/` directory. 
- A trigger resource that reference a non-existent broker: `control-plane/trigger-lifecycle/trigger-no-broker.yaml`

## [Pre] Creating a Trigger with a reference to a non-existent Broker 

```
kubectl apply -f control-plane/trigger-lifecycle/trigger-no-broker.yaml
```

## [Test] Trigger Non Readyness if no Broker is available

Check that the trigger is not Ready, as the Broker doesn't exist

```
 kubectl get trigger conformance-trigger-no-broker -ojson | jq '.status.conditions[] |select(.type == "Ready")' 
```

### [Output]

```
{
  "test": "control-plane/trigger-lifecycle/trigger-not-ready-no-broker"
  "output": {
	"expectedType": "Ready",
	"expectedStatus": "False"
    "expectedReason": "BrokerDoesNotExist"
  }
}
```

## [Pre] Creating a Trigger with a reference to a non-existent Broker 

```
kubectl apply -f control-plane/trigger-lifecycle/trigger-no-subscriber-ref.yaml
kubectl apply -f control-plane/trigger-lifecycle/broker.yaml
```


## [Test] Trigger Non Readyness if no subscriber ref resolvable

Check that the trigger is not Ready, as the subscriber ref cannot be resolved

```
 kubectl get trigger conformance-trigger-no-subscriber-ref -ojson | jq '.status.conditions[] |select(.type == "Ready")' 
```

# Clean up & Congrats

Make sure that you clean up all resources created in these tests by running: 

```
kubectl delete -f control-plane/trigger-lifecycle/
```

Congratulations you have tested the **Trigger Lifecycle Conformance** :metal: !