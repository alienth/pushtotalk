# pushtotalk

Global hotkeys in X are typically achieved using XGrabKey(). When a grabbed key
is pressed, a focus event is triggered. The result of this is that whenver you
press your hotkey, whatever app you're focused on loses focus until the hotkey
is released. This is incredibly annoying when you want a global hotkey for
push-to-talk.

I get around this by polling X for the keyboard state. When a key is pressed,
polling will detect the press and tell PulseAudio to unmute an input. A little
crude but it works. Incidentally, this is the same way Discord on Linux happens
to get their push-to-talk to work.
