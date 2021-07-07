#!/usr/bin/osascript -l JavaScript

/*
	Copyright (c) 2021 Dean Jackson <deanishe@deanishe.net>
	MIT Licence applies http://opensource.org/licenses/MIT
	Created on 2021-07-07

	Trigger ⌘C and wait for pasteboard to change.
*/

ObjC.import('Cocoa')
ObjC.import('Carbon')


function commandC() {
	let source = $.CGEventSourceCreate($.kCGEventSourceStateCombinedSessionState);

	let copyCommandDown = $.CGEventCreateKeyboardEvent(source, $.kVK_ANSI_C, true);
	$.CGEventSetFlags(copyCommandDown, $.kCGEventFlagMaskCommand);
	let copyCommandUp = $.CGEventCreateKeyboardEvent(source, $.kVK_ANSI_C, false);

	$.CGEventPost($.kCGAnnotatedSessionEventTap, copyCommandDown);
	$.CGEventPost($.kCGAnnotatedSessionEventTap, copyCommandUp);
}


function run() {
	const pboard = $.NSPasteboard.generalPasteboard,
		value = 'net.deanishe.alfred.macos-services',
		type = 'public.utf8-plain-text'

	// put a sentinel value on the clipboard
	pboard.clearContents
	pboard.setStringForType('', 'org.nspasteboard.ConcealedType')
	pboard.setStringForType(value, type)

	commandC()

	// wait up to 2 secs for sentinel value to be replaced with new
	// clipboard contents from ⌘C
	for (let i = 0; i < 200; i++) {
		let v = ObjC.unwrap(pboard.stringForType(type))
		if (v != value) {
			return
		}
		delay(0.01)
	}
	console.log('pasteboard was not populated within 2 seconds')
	return JSON.stringify({'alfredworkflow': {'variables': {'copy_failed': true}}})
}