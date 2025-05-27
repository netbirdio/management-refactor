# Pocketbase Hooks

This package original authoring is by Gani Georgiev the creator and mantainer of Pocketbase, I've merely ported and adapted it by removing dependencies. You can read more about it and its use on Pocketbase [here](https://pocketbase.io/docs/go-event-hooks/).

This is a very powerful event/hook system, especially this:

> All hook handler functions share the same func(e T) error signature and expect the user to call e.Next() if they want to proceed with the execution chain.

I've included it in our cloud repository and implemented the concept with a sub-set of billing events as a PoC of how we could re-architect Management for a cleaner developer experience and faster development by decoupling server bootstrapping (dependencies), serve (http and gRPC) and actual features (billing, idp, event streaming, etc).

