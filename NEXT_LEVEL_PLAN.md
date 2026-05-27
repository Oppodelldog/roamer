# Roamer Next Level Plan

This plan tracks larger follow-up features after the first UI/UX control-deck package.

## Goals

- Make macro creation easier for non-trivial game workflows.
- Make playback more observable while a macro is running.
- Improve testing safety by allowing input execution to be stubbed.
- Improve remote/second-device operation.

## Next

- [x] Add an input execution mode toggle: real WinAPI input vs stubbed dry-run mode.
- [x] Show live playback command feedback as transient action messages.
- [x] Add an execution timeline with current step, next step, progress, and estimated remaining time.
- [x] Add a macro recorder for keyboard and mouse input.
- [x] ~~Add a visual macro builder alongside the text editor.~~ Not needed after recorder/editor/validation workflow.
- [x] Add device pairing / remote mode for second-device control.

## Feature Packages

### Input Execution Mode

- Add a visible toggle for input mode: `Live` and `Dry Run`.
- `Live` keeps current behavior and sends inputs through WinAPI.
- `Dry Run` parses, queues, and plays sequences without executing real keyboard or mouse input.
- Dry-run playback should still update macro state, queue state, logs, diagnostics, command feedback, and timeline.
- The current mode should be visible in the top control area.
- Safety default needs a decision before implementation: keep current behavior as default, or default to dry-run on startup.

### Playback Command Feedback

- Show transient command messages during playback, similar to toast messages.
- Each executed command should appear briefly and move/fade upward.
- Wait commands should use their wait duration as part of the animation timing.
- Repeated commands from `R n [...]` should be visible as individual runtime actions.
- Feedback should be useful but not noisy enough to hide the macro deck.
- Dry-run mode should use the same command feedback path.

### Execution Timeline

- Show current command, next command, progress, and estimated remaining time while a macro runs.
- For finite macros, progress should move from start to completion.
- For loop macros, timeline should show the current loop cycle instead of pretending there is a final end.
- For queued macros, show active macro plus queued macro.
- Timeline should integrate with pause, resume, abort, stopped, execution error, and dry-run mode.

### Macro Recorder

- Add a recording mode for keyboard and mouse input.
- Capture key down/up, mouse button down/up, mouse movement or position, and waits between actions.
- Allow saving the recorded input as a new macro or replacing an existing macro.
- Provide a cleanup step after recording: normalize waits, optionally remove tiny mouse movements, and format the script.
- Recording should be clearly armed/stopped so accidental capture is unlikely.
- Dry-run mode should not block recording, but recording must not execute inputs by itself.

### Visual Macro Builder

- Skipped deliberately after testing the first slice.
- Recorder, formatter, validation, and text editor cover the workflow with less UI and less code complexity.

### Device Pairing / Remote Mode

- Add a remote-friendly control mode for second-device usage.
- Provide a simple pairing entry point, ideally a QR code or short URL.
- Remote view should focus on profile selection, macro deck, active status, pause/resume, release inputs, and abort.
- Editing features can be hidden or secondary in remote mode.
- Remote view should include a simplified macro recorder: `New Macro`, fullscreen recording overlay, large pulsing record button, name input, save, and discard.
- Consider a lightweight access token so another device on the network cannot control Roamer accidentally.
- Keep local-first operation unchanged.

## Feedback Loop Notes

- Current checkpoint: the first UI/UX control-deck plan is complete and pushed on `ui-ux-control-deck`.
- Current work package completed: Live/Dry-Run input execution mode toggle backed by a swappable input executor.
- Current work package completed: live playback command feedback with transient action messages, runtime command events, wait countdowns, and wait-based animation timing.
- Current work package completed: execution timeline in the control area with current command, next command, step progress, loop cycle, and estimated remaining time.
- Current work package completed: macro recorder in the macro editor with Start/Stop, recorded script preview, live recorder toasts, and per-macro `Use Recording`.
- Current work package skipped: visual macro builder, because recorder plus script editor is the better workflow.
- Current work package completed: `/remote` view with reduced controls, local remote URLs, and simplified remote macro recorder.
- Current package status: complete. No open feature package remains in this plan.
