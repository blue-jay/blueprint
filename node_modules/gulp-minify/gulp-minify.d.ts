/// <reference path="../typings/main.d.ts" />

declare namespace GulpMinify {
	interface MinifyOptions {
		ext?: {
			src?: string;
			min?: string;
		},
		exclude?: Array<string>,
		ignoreFiles?: Array<string>
	}
}

declare module "gulp-minify" {
	import { Duplex } from 'stream';

	function minify(options?: GulpMinify.MinifyOptions): Duplex; 

	export = minify;
}