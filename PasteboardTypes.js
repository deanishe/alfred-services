#!/usr/bin/osascript -l JavaScript

/*
	Copyright (c) 2020 Dean Jackson <deanishe@deanishe.net>
	MIT Licence applies http://opensource.org/licenses/MIT
	Created on 2020-08-01

	Returns type(s) of pasteboard items.
*/

ObjC.import('Cocoa')

const pboard = $.NSPasteboard.generalPasteboard

function pboardTypes() {
	let types = []
	ObjC.unwrap(pboard.pasteboardItems).forEach(it => {
		ObjC.unwrap(it.types).forEach(t => {
			let s = ObjC.unwrap(t)
			if (!types.includes(s)) {
				types.push(s)
				console.log(`[pboard] type=${s}`)
			}
		})
	})
	return types
}

function run() {
	return JSON.stringify(pboardTypes())
}