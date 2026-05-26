# Roamer UI/UX Plan

This plan tracks the frontend and interaction work for turning Roamer into a clearer live control deck.

## Completed

- [x] Redesign the main UI into a control-deck shell while keeping the large background image and central headline.
- [x] Replace plain page buttons with profile cards that show title and macro count.
- [x] Replace the macro button list with a responsive macro card grid.
- [x] Show the active macro state directly on its macro button.
- [x] Make a second tap on the active macro button abort that macro.
- [x] Add visible playing animation to active macro cards.
- [x] Add a distinct paused visual state for active macro cards.
- [x] Add a stable top control bar with connection state, current macro, pause/resume, and abort.
- [x] Load and use `Action.Icon` via `icons.css`.
- [x] Fix the `ico-kayak-forward` icon class typo.
- [x] Fix finite sequences publishing an idle state after completion.
- [x] Add `PageId`, `SequenceIndex`, and explicit sequencer `State` to sequence state messages.
- [x] Make `Abort()` release pressed keys through the sequencer adapter.
- [x] Fix abort-while-paused crash and add regression coverage.
- [x] Add macro metadata badges for loop, mouse usage, held keys, and estimated duration.
- [x] Move the attributions link into a cleaner fixed footer position.
- [x] Replace vague idle status text with context-specific labels.
- [x] Add a dedicated global "Release Inputs" action that releases all known pressed keyboard and mouse buttons without requiring a running macro.
- [x] Turn logs into an accessible diagnostics panel instead of hidden background output.
- [x] Make the log of last action, last error, and websocket state visible in the UI.
- [x] Add save-time script validation with readable error messages in the macro editor.
- [x] Add command help for `KD`, `KU`, `W`, `LD`, `LU`, `SM`, `MM`, `R`, and `L`.
- [x] Add a script format/normalize action using the existing parser/writer.
- [x] Add macro duplicate and reorder controls.
- [x] Add icon selection in the macro editor.
- [x] Add live script validation while editing.
- [x] Fix parser panic for incomplete commands during live validation.
- [x] Detect loop status server-side instead of only via frontend sequence heuristics.
- [x] Improve macro card metadata with backend-derived parsing information.
- [x] Add a visible error state for failed sequence parsing.
- [x] Add a distinct `looping` state for active loop macros.
- [x] Add a visible error state for execution errors.
- [x] Add page/profile theme support without global CSS collisions.
- [x] Convert game-specific CSS into controlled theme classes.
- [x] Redesign the macro editor as a panel/drawer over the current page context.
- [x] Improve the sound settings panel as a compact drawer or bottom sheet.
- [x] Review mobile/tablet layout after the core panel work.
- [x] Add more precise macro states: queued and stopped.
- [x] Add focused tests for sequence state transitions, abort behavior, and parser-derived metadata.

## Next

- [ ] Decide the next UX/workflow improvement package.

## Feedback Loop Notes

- Current checkpoint: control deck and macro card UX are in place and tested manually.
- Recent bug fixed: aborting while paused no longer panics and no longer blocks future sequences.
- Current work package completed: global safety controls plus diagnostics panel.
- Current work package completed: save-time script validation and clearer editor save feedback.
- Current work package completed: script format/normalize and command help.
- Current work package completed: macro duplicate/reorder controls.
- Current work package completed: icon selection in the macro editor.
- Current work package completed: live validation while editing.
- Current work package completed: backend-derived macro metadata and server-side loop detection.
- Current work package completed: visible macro card state for parser errors.
- Current work package completed: active loop macros now show a distinct looping state.
- Current work package completed: execution errors are reported from the sequencer to the UI and shown on the affected macro card.
- Current work package completed: page/profile theme data, legacy `CssFile` theme fallback, CSS variables, and basic theme editing in the Pages panel.
- Current work package completed: remaining game-specific CSS files are scoped to `.theme-*` presets and no longer style `body`, `button`, or `*` globally.
- Current work package completed: macro editor now opens as a right-side drawer over the current page context.
- Current work package completed: sound settings now opens as a compact bottom/right panel with horizontal sliders and visible percentage values.
- Current work package completed: mobile/tablet layout spacing for topbar, deck grids, macro drawer, sound panel, diagnostics, and page editor fields.
- Current work package completed: queued and stopped macro card states, with backend state split between active and queued macros.
- Current work package completed: focused tests for sequence state transitions, abort behavior, and parser-derived metadata.
- Current plan is complete; next step is choosing a new UX/workflow improvement package.
