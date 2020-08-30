
<div align="center">
    <img width="128" height="128" src="https://raw.githubusercontent.com/deanishe/alfred-services/master/icons/icon-large.png" alt="grey cog in grey ring" title="workflow icon">
</div>

macOS Services for Alfred
=========================

Run macOS services via Alfred 4+.

This workflow can execute macOS services on the clipboard contents, current selection or files (via an Alfred File Action).


Installation
------------

Download the latest version of the workflow from the [releases page][releases], then double-click the `macOS-Services-X.X.X.alfredworkflow` file to install.


Usage
-----

When run via keyword (`services` by default), you can choose from services that are applicable to the current contents of the general pasteboard. There is also a Hotkey to run the workflow with the current pasteboard.

Alternatively, you can call the workflow via its File Action (called "macOS Services") to run a service with the selected files, or use the second Hotkey to run the workflow on the current macOS selection.

Finally, you can assign your own Hotkeys to specific services (though this only works with the pasteboard contents). See the red EXAMPLE element in the workflow (which calls the "New TextEdit Window Containing Selection" service).


Configuration
-------------

There is one variable in the workflow configuration sheet: `delay_after_copy`. When using the Hotkey to use the current macOS selection, the workflow simulates a `⌘C` keypress to put the selection on the clipboard (where it can read it). `delay_after_copy` specifies how long (in seconds) the workflow should wait for the clipboard to populate after the simulated keypress. The default delay is 0.3s. Increase this if you find the workflow is being run with the previous clipboard contents.


Licensing & thanks
------------------

- The workflow is released under the [MIT licence][mit].
- It is heavily based on the [AwGo][awgo] library, also MIT licensed.
- The workflow icons are based on [Font Awesome][awesome], released under the [Creative Commons Attribution 4.0 licence][ccby40].

[releases]: https://github.com/deanishe/alfred-services/releases/latest
[mit]: LICENCE.txt
[awesome]: https://github.com/FortAwesome/Font-Awesome
[ccby40]: https://creativecommons.org/licenses/by/4.0/legalcode
[awgo]: https://github.com/deanishe/awgo
