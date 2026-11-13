# SlideBolt Scheduler SDK

Go client contract for creating, updating, enabling, disabling, and deleting
durable SlideBolt timers over the messenger service.

Timers target a switch entity and publish ordinary `switch_turn_on` and
`switch_turn_off` commands at one-shot or cron edges.
