*golang-evdev*
--------------

This package provides Go language bindings to the generic input event
interface in Linux. The *evdev* interface serves the purpose of
passing events generated in the kernel directly to userspace through
character devices that are typically located in `/dev/input/`.

Documentation:
    http://godoc.org/github.com/gvalkov/golang-evdev

Development:
    https://github.com/gvalkov/golang-evdev

*Fork Notes:*
-------------

The only changes between this and the original project (thus far) is 
the keycode which is problematic. It works if you assume a particular
layout (maybe). The change is similar to some code in evdev-python.

One of two things needs to happen:
* The keycode is looked up in the keymap. I have no idea how to do
  this.
* Keymap is stored as an array of possible codes requiring a compare
  function. This is problematic in that any scancode can 
  theorhetically produce any keycode depending on layout and really
  this just then becomes about the meaningful probability of a key
  presenting a particular character/keycode.
