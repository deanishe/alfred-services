#!/usr/bin/osascript -l JavaScript

/*
	Copyright (c) 2020 Dean Jackson <deanishe@deanishe.net>
	MIT Licence applies http://opensource.org/licenses/MIT
	Created on 2020-08-01

	Puts filepaths passed as ARGV onto general pasteboard.
*/

ObjC.import('Cocoa')

const pboard = $.NSPasteboard.generalPasteboard

function run(paths) {
	console.log('.')
	console.log('/--------- INPUT FILES ---------\\')
	let arr = $.NSMutableArray.alloc.init,
	paths.forEach(p => {
		let url = $.NSURL.fileURLWithPath(p)
		arr.addObject(url.absoluteURL)
		console.log(ObjC.unwrap(url.absoluteString))
	})
	console.log('\\--------- INPUT FILES ---------/')
	pboard.clearContents
	pboard.writeObjects(arr)
	return JSON.stringify({
		alfredworkflow: {
			variables: {
				PBOARD_TYPES: 'public.file-url',
				CLIPBOARD: paths.join('\n')
			}
		}
	})
}